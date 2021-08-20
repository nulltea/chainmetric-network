package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-contracts/shared/core"
	model "github.com/timoth-y/chainmetric-contracts/shared/model/user"
)

// GenerateJWT generates JWT token for given `user`.
func GenerateJWT(user *model.User) (string, *time.Time, error) {
	var (
		token        = jwt.New(jwt.SigningMethodRS512)
		org          = viper.GetString("organization")
		expiresAfter = viper.GetDuration("api.jwt_expiration")
		expiresAt    *time.Time

		claims = &jwt.StandardClaims{
			Id:       user.ID,
			Issuer:   fmt.Sprintf("identity.%s.org", org),
			IssuedAt: time.Now().Unix(),
		}
	)

	if expiresAfter != 0 {
		*expiresAt = time.Now().Add(expiresAfter)
		claims.ExpiresAt = expiresAt.Unix()
	}

	token.Claims = claims

	jwtToken, err := token.SignedString(jwtSigningPrivateKey())

	return jwtToken, expiresAt, err
}

// VerifyJWT performs verification of a given JWT `token`.
func VerifyJWT(token string) (string, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
				return jwtSigningPublicKey(), nil
			}
			return nil, errors.Errorf("unexpected signing method: %q", token.Header["alg"])
		},
	); if err != nil {
		return "", fmt.Errorf("access token is invalid: %w", err)
	}

	if claims, ok := jwtToken.Claims.(*jwt.StandardClaims); ok {
		return claims.Id, nil
	}

	return "", errors.New("invalid token claims")
}

func jwtSigningPrivateKey() *rsa.PrivateKey {
	var path = viper.GetString("api.jwt_signing_key")

	keyBytes, err := os.ReadFile(path)
	if err != nil {
		core.Logrus.WithField("path", path).Fatal("failed to read key")
		return nil
	}

	block, _ := pem.Decode(keyBytes)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes); if err != nil {
		core.Logrus.WithField("path", path).
			Fatal("failed to parse rsa private key")
		return nil
	}

	return key.(*rsa.PrivateKey)
}

func jwtSigningPublicKey() *rsa.PublicKey {
	var path = viper.GetString("api.jwt_signing_key")

	keyBytes, err := os.ReadFile(path)
	if err != nil {
		core.Logrus.WithField("path", path).Fatal("failed to read key")
		return nil
	}

	block, _ := pem.Decode(keyBytes)
	key, err := x509.ParsePKIXPublicKey(block.Bytes); if err != nil {
		core.Logrus.WithField("path", path).
			Fatal("failed to parse rsa public key")
		return nil
	}

	return key.(*rsa.PublicKey)
}
