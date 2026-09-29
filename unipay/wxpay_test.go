package unipay

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"testing"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/support"
)

func testPrivateKeyPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))
}

func testWxpayConfig(t *testing.T) *WxpayConfig {
	return &WxpayConfig{
		Appid:                      "appid",
		MchID:                      "1234567890",
		MchCertificateSerialNumber: "MCHSERIAL",
		MchAPIv3Key:                "0123456789abcdef0123456789abcdef",
		MchPrivateKey:              testPrivateKeyPEM(t),
		NotifyURL:                  "https://example.com/notify",
	}
}

func TestNewWxpay(t *testing.T) {
	wx, err := NewWxpay(testWxpayConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	if wx.Client() == nil {
		t.Fatal("expected non-nil payment client")
	}
}

func TestWxpayOrderSigner(t *testing.T) {
	wx, err := NewWxpay(testWxpayConfig(t))
	if err != nil {
		t.Fatal(err)
	}
	// 回归测试：下单客户端的签名器必须能直接签名。
	// 此前只注入了 Base 的签名器，Order 等服务客户端各自的签名器为空，
	// 首次请求时报 "you must set privatekey to use SHA256WithRSASigner"
	_, err = wx.Client().Order.BaseClient.Signer.GenerateRequestSign(&support.RequestSignChain{
		Method:       http.MethodGet,
		CanonicalURL: "/v3/certificates",
		SignBody:     "",
	})
	if err != nil {
		t.Fatalf("order signer fail: %v", err)
	}
}

func TestNewWxpayBadConfig(t *testing.T) {
	conf := testWxpayConfig(t)
	conf.MchAPIv3Key = "short"
	if _, err := NewWxpay(conf); err == nil {
		t.Error("expect error for short api v3 key")
	}

	conf = testWxpayConfig(t)
	conf.MchPrivateKey = "not a pem"
	if _, err := NewWxpay(conf); err == nil {
		t.Error("expect error for invalid private key")
	}
}
