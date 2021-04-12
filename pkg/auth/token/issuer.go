package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/spf13/viper"
	"time"

	"golang.org/x/oauth2"
)

// Issuer interface
type Issuer interface {
	Issue(*VerifiedClaims) (*oauth2.Token, error)
}

// NewIssuer returns a new instance of the issuer
func NewIssuer(generator Generator, duration time.Duration) Issuer {
	return &issuer{
		Generator:   generator,
		Duration:    duration,
	}
}

type issuer struct {
	Generator   Generator
	Duration    time.Duration
}

// Issue returns a token for the specified claims
func (i *issuer) Issue(verifiedClaims *VerifiedClaims) (*oauth2.Token, error) {
	if verifiedClaims.UserID == "" {
		return nil, errors.New("missing user id in verified claims")
	}

	if verifiedClaims.ConnectorID == "" {
		return nil, errors.New("missing connector id in verified claims")
	}

	sub := verifiedClaims.Sub
	email := verifiedClaims.Email
	name := verifiedClaims.Name
	userID := verifiedClaims.UserID
	userName := verifiedClaims.UserName
	claimGroups := verifiedClaims.Groups
	
	isAdmin := false
	for _, group := range claimGroups {
		if group == viper.GetString("dex_admin_group") {
			isAdmin = true
		}
	}

	teams := claimGroups

	if len(teams) == 0 {
		return nil, errors.New("user doesn't belong to any team")
	}

	return i.Generator.Generate(map[string]interface{}{
		"sub":       sub,
		"email":     email,
		"name":      name,
		"user_id":   userID,
		"user_name": userName,
		"teams":     teams,
		"is_admin":  isAdmin,
		"exp":       time.Now().Add(i.Duration).Unix(),
		"csrf":      RandomString(),
	})
}

// RandomString returns a random string
func RandomString() string {
	s := make([]byte, 32)
	rand.Read(s)
	return hex.EncodeToString(s)
}