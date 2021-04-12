package token

import (
	"context"
	"errors"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// VerifiedClaims structure defines a verified claim
type VerifiedClaims struct {
	Sub         string
	Email       string
	Name        string
	UserID      string
	UserName    string
	ConnectorID string
	Groups      []string
}

// NewVerifier returns an instance of a Verifier from the given clientID and issuerURL
func NewVerifier(clientID, issuerURL string) *Verifier {
	return &Verifier{
		ClientID:  clientID,
		IssuerURL: issuerURL,
	}
}

// Verifier struct
type Verifier struct {
	ClientID  string
	IssuerURL string
}

// Verify checks that the token is valid and returns verified claims
func (v *Verifier) Verify(ctx context.Context, token *oauth2.Token) (*VerifiedClaims, error) {
	if v.ClientID == "" {
		return nil, errors.New("Missing client id")
	}

	if v.IssuerURL == "" {
		return nil, errors.New("Missing issuer")
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("Missing id_token")
	}

	provider, err := oidc.NewProvider(ctx, v.IssuerURL)
	if err != nil {
		return nil, err
	}

	providerVerifier := provider.Verifier(&oidc.Config{
		ClientID: v.ClientID,
	})

	verifiedToken, err := providerVerifier.Verify(ctx, idToken)
	if err != nil {
		return nil, err
	}

	type Federated struct {
		ConnectorID string `json:"connector_id"`
		UserID      string `json:"user_id"`
		UserName    string `json:"user_name"`
	}

	type Claims struct {
		Sub       string    `json:"sub"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		Groups    []string  `json:"groups"`
		Federated Federated `json:"federated_claims"`
	}

	var claims Claims
	verifiedToken.Claims(&claims)

	return &VerifiedClaims{
		Sub:         claims.Sub,
		Email:       claims.Email,
		Name:        claims.Name,
		UserID:      claims.Federated.UserID,
		UserName:    claims.Federated.UserName,
		ConnectorID: claims.Federated.ConnectorID,
		Groups:      claims.Groups,
	}, nil
}