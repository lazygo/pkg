package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

// openssl genrsa -out rsa_private_key.pem 1024
// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem

const PrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDl6cvp9HR+XlicKE1uhIqpePhkP8l1jXG4aQeaMfOWDW2RsSpf
QfilBQC15bCaHmcbIVrgy69/Y4ebs+gGjdN9elzsYfHytCJWSMEG/kT7l2HdOrBr
jXDs1SSQynCL9hZWIdw50hW2FhF6E0ZFeFNt7Hf5OU7cwpdopPhJ/J6Q+wIDAQAB
AoGAA9vkvEyKGATlX9mdUxmOakHJiYU4kGyLWBkLM59bA02+ZQ+gMnEdB0gKNwNf
73ZLLL1mlRdWHsFA6XAfmNyQjCZ4B/jdb5wFpokNsxXylvGm6IPKTIJRoeXwJTNK
8mzi0Y7XVYLyM1l0JYbdE89d+mzUl+q/3yQcgfB8zUt8jKECQQD4+P5V0A+0Negz
zHTL60cJIeD/UaSeL5qYzNTwQzNiAUOx7Xmd5sYQRae19w9nBCfKb595Ro6BBz6E
76XI+xgfAkEA7GcV3jCrvJrYgdEqBFiRH/i7XwVXFOyO2EKF/YRznRoWtJy/rFZu
tOykDKUp3MM4clG37MPoyn0RsJ4aXCtbpQJAFvrcds0ydd635PgNG7lGoDgpTUea
2yLnsQzO5rI9LuGQ/v49SG7Bf0T+mtQH7uk6RvwQiyARDSW/BoQcGDXc3wJBALne
s0rnaZ/4/5HSKv8Pw8snferQAA/rjsRqSX9yzJQRFxkaxXly28hU5wcqNSfmNlNr
/PijcD0E6Qu8w20EiiECQDh0nmpEICyMxAnwYHEYJ84XGL7GZ+SKj64IK6CxvL8N
1prWh+d6IzfWvgpCuBZk6QRSB265iodhbRrZ2TSdH0U=
-----END RSA PRIVATE KEY-----
`

func TestJwt2(t *testing.T) {
	enc, err := NewJwtEncoder([]byte(PrivateKey))
	if err != nil {
		t.Error(err)
		return
	}

	token, err := enc.Encode(&Claims{Subject: "5"})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(token)
}

func TestJwt(t *testing.T) {

	// genRsaKey(1024)

	privateKey := []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDl6cvp9HR+XlicKE1uhIqpePhkP8l1jXG4aQeaMfOWDW2RsSpf
QfilBQC15bCaHmcbIVrgy69/Y4ebs+gGjdN9elzsYfHytCJWSMEG/kT7l2HdOrBr
jXDs1SSQynCL9hZWIdw50hW2FhF6E0ZFeFNt7Hf5OU7cwpdopPhJ/J6Q+wIDAQAB
AoGAA9vkvEyKGATlX9mdUxmOakHJiYU4kGyLWBkLM59bA02+ZQ+gMnEdB0gKNwNf
73ZLLL1mlRdWHsFA6XAfmNyQjCZ4B/jdb5wFpokNsxXylvGm6IPKTIJRoeXwJTNK
8mzi0Y7XVYLyM1l0JYbdE89d+mzUl+q/3yQcgfB8zUt8jKECQQD4+P5V0A+0Negz
zHTL60cJIeD/UaSeL5qYzNTwQzNiAUOx7Xmd5sYQRae19w9nBCfKb595Ro6BBz6E
76XI+xgfAkEA7GcV3jCrvJrYgdEqBFiRH/i7XwVXFOyO2EKF/YRznRoWtJy/rFZu
tOykDKUp3MM4clG37MPoyn0RsJ4aXCtbpQJAFvrcds0ydd635PgNG7lGoDgpTUea
2yLnsQzO5rI9LuGQ/v49SG7Bf0T+mtQH7uk6RvwQiyARDSW/BoQcGDXc3wJBALne
s0rnaZ/4/5HSKv8Pw8snferQAA/rjsRqSX9yzJQRFxkaxXly28hU5wcqNSfmNlNr
/PijcD0E6Qu8w20EiiECQDh0nmpEICyMxAnwYHEYJ84XGL7GZ+SKj64IK6CxvL8N
1prWh+d6IzfWvgpCuBZk6QRSB265iodhbRrZ2TSdH0U=
-----END RSA PRIVATE KEY-----
	`)

	publicKey := []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDl6cvp9HR+XlicKE1uhIqpePhk
P8l1jXG4aQeaMfOWDW2RsSpfQfilBQC15bCaHmcbIVrgy69/Y4ebs+gGjdN9elzs
YfHytCJWSMEG/kT7l2HdOrBrjXDs1SSQynCL9hZWIdw50hW2FhF6E0ZFeFNt7Hf5
OU7cwpdopPhJ/J6Q+wIDAQAB
-----END PUBLIC KEY-----
	`)

	encoder, err := NewJwtEncoder(privateKey)
	if err != nil {
		t.Error(err)
	}
	decoder, err := NewJwtDecoder(publicKey)
	if err != nil {
		t.Error(err)
	}

	id := "1"
	token, err := encoder.Encode(&Claims{Subject: id})
	if err != nil {
		t.Error(err)
	}

	t.Log(token)

	claims, err := decoder.Decode(token)
	if err != nil {
		t.Error(err)
	}

	if id != claims.Subject {
		t.Errorf("id: %s, decode id: %s", id, claims.Subject)
	}

	id = "2"
	issuer := "p2link.cn"
	audience := "p"
	token, err = encoder.Encode(&Claims{Subject: id, Issuer: issuer, Audience: []string{audience}})
	if err != nil {
		t.Error(err)
	}

	t.Log(token)

	decid, session, err := ParseToken(decoder, token, VerifyIssuer("p2link.cn"), VerifyAudience("p"))
	if err != nil {
		t.Error(err)
		return
	}
	if id != strconv.FormatUint(decid, 10) {
		t.Errorf("id: %s, decode id: %s", id, claims.ID)
	}
	if session != 0 {
		t.Error("session error")
		return
	}

	_, _, err = ParseToken(decoder, token, VerifyIssuer("p2link.cn"), VerifyAudience("n"))
	if !strings.Contains(err.Error(), "invalid audience") {
		t.Error(err)
		return
	}

	_, _, err = ParseToken(decoder, token, VerifyIssuer("p2link"), VerifyAudience("p"))
	if !strings.Contains(err.Error(), "invalid issuer") {
		t.Error(err)
		return
	}

}

func genRsaKey(bits int) error {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "私钥",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "公钥",
		Bytes: derPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}
