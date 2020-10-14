package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/concourse/concourse/go-concourse/concourse"
	"golang.org/x/oauth2"
)

type ConcourseClientConfig struct {
	ATCAddr            string
	Username           string
	Password           string
	TeamName           string
	InsecureSkipVerify bool
	EnableTracing      bool
}

func NewClientFromEnv(team string) (concourse.Client, error) {
	cfg := ConcourseClientConfig{
		ATCAddr:            os.Getenv("CONCOURSE_ATC_ADDR"),
		Username:           os.Getenv("CONCOURSE_USERNAME"),
		Password:           os.Getenv("CONCOURSE_PASSWORD"),
		TeamName:           team,
		InsecureSkipVerify: os.Getenv("CONCOURSE_INSECURE_SKIP_VERIFY") == "true",
		EnableTracing:      os.Getenv("CONCOURSE_ENABLE_TRACING") == "true",
	}
	return newClient(cfg)
}

func newClient(cfg ConcourseClientConfig) (concourse.Client, error) {
	tokenClient := concourse.NewClient(cfg.ATCAddr, &http.Client{
		Transport: ConcourseAuthTransport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
		},
	}, cfg.EnableTracing)
	// get bearer token from user/pass
	oauth2Config := oauth2.Config{
		ClientID:     "fly",
		ClientSecret: "Zmx5",
		Endpoint:     oauth2.Endpoint{TokenURL: tokenClient.URL() + "/sky/issuer/token"},
		Scopes:       []string{"openid", "profile", "email", "federated:id", "groups"},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, tokenClient.HTTPClient())
	token, err := oauth2Config.PasswordCredentialsToken(ctx, cfg.Username, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("resource: couldn't obtain auth token: %s", err)
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("failed to find id_token extra in oauth2 token")
	}
	// create a concourse client
	client := concourse.NewClient(cfg.ATCAddr, &http.Client{
		Transport: ConcourseAuthTransport{
			AccessToken:     idToken,
			TokenType:       token.TokenType,
			TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify},
		},
	}, cfg.EnableTracing)

	return client, nil
}

type ConcourseAuthTransport struct {
	AccessToken     string
	TokenType       string
	TLSClientConfig *tls.Config
}

// RoundTrip modifies the behaviour of the http.DefaultTransport by injecting Authorization header.
func (bat ConcourseAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if bat.TokenType != "" && bat.AccessToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", bat.TokenType, bat.AccessToken))
	}
	return http.DefaultTransport.RoundTrip(req)
}
