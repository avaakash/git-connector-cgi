package gitclient

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v64/github"
	"github.com/harness/git-connector-cgi/common"
	"golang.org/x/oauth2"
)

func GetTokenForGithubApp(config *common.GithubApp) (string, error) {

	privateKey, err := loadPrivateKey(config.PrivateKey)
	if err != nil {
		return "Error loading RSA private key", err
	}

	jwtToken, err := createJWTToken(config.AppId, privateKey)
	if err != nil {
		return "Failed to create JWT token", err
	}

	return getInstallationAccessToken(jwtToken, config.AppInstallationId, config.GithubUrl)
}

// loadPrivateKey loads the private key from a PEM file
func loadPrivateKey(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return privKey, nil
}

// createJWTToken creates a JWT token for the GitHub App using the private key
func createJWTToken(appID string, privateKey *rsa.PrivateKey) (string, error) {
	// To protect against clock drift, set the issuance time 60 seconds in the past.
	now := time.Now().Add(-60 * time.Second)
	expirationTime := now.Add(5 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.RegisteredClaims{
		Issuer:    appID,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(now),
	})

	// Sign the JWT with the private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

// getInstallationAccessToken exchanges the JWT for an installation access token
func getInstallationAccessToken(jwtToken string, installationID string, githubUrl string) (string, error) {
	var err error
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: jwtToken}),
		},
	}

	githubClient := github.NewClient(client)
	if !strings.Contains(githubUrl, "github.com") {
		if githubClient, err = githubClient.WithEnterpriseURLs(githubUrl, ""); err != nil {
			return "", fmt.Errorf("failed to create GitHub Entreprise client for URL: %v due to %w", githubUrl, err)
		}
	}
	// Use GitHub's API to exchange JWT for the installation access token
	installationIDInt, err := strconv.ParseInt(installationID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse installation ID: %w", err)
	}
	accessToken, _, err := githubClient.Apps.CreateInstallationToken(context.Background(), installationIDInt, nil)
	if err != nil {
		return "", fmt.Errorf("failed to exchange JWT for installation access token: %w", err)
	}

	return *accessToken.Token, nil
}
