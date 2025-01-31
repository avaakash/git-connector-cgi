// Copyright 2021 Harness Inc. All rights reserved.
// Use of this source code is governed by the PolyForm Free Trial 1.0.0 license
// that can be found in the licenses directory at the root of this repository, also available at
// https://polyformproject.org/wp-content/uploads/2020/05/PolyForm-Free-Trial-1.0.0.txt.

package gitclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
	"github.com/drone/go-scm/scm/transport"
	"github.com/harness/git-connector-cgi/common"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/drone/go-scm/scm/transport/oauth2"
)

const (
	SkipSSLVerify       bool = true
	AdditionalCertsPath      = ""
)

func oauthTransport(token string, skip bool, additionalCertsPath string, proxy string) http.RoundTripper {
	return &oauth2.Transport{
		Base: defaultTransport(skip, additionalCertsPath, proxy),
		Source: oauth2.StaticTokenSource(
			&scm.Token{
				Token: token,
			},
		),
	}
}
func privateTokenTransport(token string, skip bool, additionalCertsPath string, proxy string) http.RoundTripper {
	return &transport.PrivateToken{
		Base:  defaultTransport(skip, additionalCertsPath, proxy),
		Token: token,
	}
}

func tlsConfig(skip bool, additionalCertsPath string) *tls.Config {
	config := tls.Config{
		InsecureSkipVerify: skip,
	}
	if skip || additionalCertsPath == "" {
		return &config
	}
	// Try to read 	additional certs and add them to the root CAs
	// Create TLS config using cert PEM
	rootPem, err := os.ReadFile(additionalCertsPath)
	if err != nil {
		logrus.Warnf("could not read certificate file (%s), error: %s", additionalCertsPath, err.Error())
		return &config
	}

	// Use the system certs if possible
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	ok := rootCAs.AppendCertsFromPEM(rootPem)
	if !ok {
		logrus.Errorf("error adding cert (%s) to pool", additionalCertsPath)
		return &config
	}
	config.RootCAs = rootCAs
	return &config
}

// defaultTransport provides a default http.Transport.
// If skip verify is true, the transport will skip ssl verification.
// Otherwise, it will append all the certs from the provided path.
func defaultTransport(skip bool, additionalCertsPath string, proxy string) http.RoundTripper {
	if proxy == "" {
		return &http.Transport{
			Proxy:           http.ProxyFromEnvironment,
			TLSClientConfig: tlsConfig(skip, additionalCertsPath),
		}
	}

	proxyURL, _ := url.Parse(proxy)

	return &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),
		TLSClientConfig: tlsConfig(skip, additionalCertsPath),
	}
}
func GetGitClient(provider common.Provider, config *common.APIAccess) (client *scm.Client, err error) { //nolint:gocyclo,funlen
	switch provider {
	case common.Github:
		if config.Endpoint == "" {
			client = github.NewDefault()
		} else {
			client, err = github.New(config.Endpoint)
			if err != nil {
				logrus.Errorln("GetGitClient failure Github", "endpoint", config.Endpoint, zap.Error(err))
				return nil, err
			}
		}
		isGithubAnonymous := IsGithubAnonymous()
		if isGithubAnonymous {
			client.Client = &http.Client{
				Transport: defaultTransport(SkipSSLVerify, AdditionalCertsPath, config.ProxyURL),
			}
		} else {
			var token string
			switch config.AccessType {
			case common.APIAccessToken:
				token = config.Token
			case common.APIAccessGithubApp:
				token, err = GetTokenForGithubApp(config.GithubApp)
				if err != nil {
					return nil, err
				}
			default:
				return nil, status.Errorf(codes.Unimplemented, "Github Application not implemented yet")
			}
			client.Client = &http.Client{
				Transport: oauthTransport(token, SkipSSLVerify, AdditionalCertsPath, config.ProxyURL),
			}
		}
	default:
		logrus.Errorln("GetGitClient unsupported git provider", "endpoint", config.Endpoint)
		return nil, status.Errorf(codes.InvalidArgument, "Unsupported git provider")
	}
	return client, nil
}

// Finds out if provider is Github Anonymous
func IsGithubAnonymous() (out bool) {
	return false
}

type HarnessAuth struct {
	Username, Password string
}

func (h *HarnessAuth) String() string {
	return "HarnessAuth"
}

func (h *HarnessAuth) Name() string {
	masked := "*******"
	if h.Password == "" {
		masked = "<empty>"
	}

	return fmt.Sprintf("%s - %s:%s", h.Name(), h.Username, masked)
}

func (h *HarnessAuth) SetAuth(r *http.Request) {
	if r == nil {
		return
	}

	r.Header.Set("Authorization", fmt.Sprintf("CIManager %s", h.Password))
}
