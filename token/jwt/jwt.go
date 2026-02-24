package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/cristalhq/jwt/v3"
	"github.com/lazygo/pkg/token"
)

type Claims jwt.StandardClaims
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

	// 8. also you can parse and verify in 1 operation
	token, err := jwt.ParseAndVerifyString(str, verifier)
	if err != nil {
		return nil, err
	}

	// 9. get standard claims
	var claims Claims
	err = json.Unmarshal(token.RawClaims(), &claims)
	if err != nil {
		return nil, err
	}

	// 10. verify claims
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
		if claims.Issuer != issuer {
			return fmt.Errorf("invalid issuer: want: %s has: %s", claims.Issuer, issuer)
		}
		return nil
	}
}

func VerifyAudience(audience string) Validator {
	return func(claims *Claims) error {
		if slices.Contains(claims.Audience, audience) {
			return nil
		}
		return fmt.Errorf("invalid audience: want: %v has: %s", claims.Audience, audience)
	}
}
