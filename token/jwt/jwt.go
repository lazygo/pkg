package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strconv"

	"github.com/cristalhq/jwt/v5"
	"github.com/lazygo/pkg/token"
)

type Claims = jwt.RegisteredClaims
type Audience = jwt.Audience
type Validator func(*Claims) error

type JwtEncoder struct {
	privateKey *rsa.PrivateKey
}

func NewJwtEncoder(privateKey []byte) (*JwtEncoder, error) {
	je := &JwtEncoder{}
	var err error
	je.privateKey, err = parsePrivateKey(privateKey)
	return je, err
}

func (je *JwtEncoder) Encode(claims *Claims) (string, error) {
	// 1. create a signer & a verifier
	signer, err := jwt.NewSignerRS(jwt.RS256, je.privateKey)
	if err != nil {
		return "", err
	}

	// 3. create a builder
	builder := jwt.NewBuilder(signer)

	// 4. and build a token
	token, err := builder.Build(claims)
	if err != nil {
		return "", err
	}
	return token.String(), nil
}

type JwtDncoder struct {
	publicKey *rsa.PublicKey
}

func NewJwtDecoder(publicKey []byte) (*JwtDncoder, error) {
	jd := &JwtDncoder{}
	var err error
	jd.publicKey, err = parsePublicKey(publicKey)
	return jd, err
}

func (jd *JwtDncoder) Decode(str string) (*Claims, error) {
	// 1. create a signer & a verifier
	verifier, err := jwt.NewVerifierRS(jwt.RS256, jd.publicKey)
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := jwt.ParseClaims([]byte(str), verifier, &claims); err != nil {
		return nil, err
	}
	return &claims, nil
}

var (
	ErrKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1")
	ErrInvalidRSAKey       = errors.New("Key is not a valid RSA key")
)

func parsePrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func parsePublicKey(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pkey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidRSAKey
	}
	return pkey, nil
}

func ParseToken(dec *JwtDncoder, str string, validator ...Validator) (uint64, int64, error) {
	if dec == nil {
		return 0, 0, errors.New("jwt decoder fail")
	}
	token, session := token.UnwrapToken(str)
	claims, err := dec.Decode(token)
	if err != nil {
		return 0, 0, fmt.Errorf("decode token fail: %w", err)
	}

	for _, valid := range validator {
		if err := valid(claims); err != nil {
			return 0, 0, fmt.Errorf("invalid token: %w", err)
		}
	}
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return id, session, nil
}

func VerifyIssuer(issuer string) Validator {
	return func(claims *Claims) error {
		if claims.IsIssuer(issuer) {
			return nil
		}
		return fmt.Errorf("invalid issuer: want: %s has: %s", claims.Issuer, issuer)
	}
}

func VerifyAudience(audience string) Validator {
	return func(claims *Claims) error {
		if claims.IsForAudience(audience) {
			return nil
		}
		return fmt.Errorf("invalid audience: want: %v has: %s", claims.Audience, audience)
	}
}
