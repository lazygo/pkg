package unipay

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/models"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment"
	notifyRequest "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/notify/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/request"
	wxresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/response"
)

var _ unipay = (*Wxpay)(nil)

type WxpayConfig struct {
	Appid                      string `json:"appid" toml:"appid"`
	MchID                      string `json:"mch_id" toml:"mch_id"`
	MchCertificateSerialNumber string `json:"mch_certificate_serial_number" toml:"mch_certificate_serial_number"`
	MchAPIv3Key                string `json:"mch_api_v3_key" toml:"mch_api_v3_key"`
	MchPrivateKey              string `json:"mch_private_key" toml:"mch_private_key"`
	NotifyURL                  string `json:"notify_url" toml:"notify_url"`
}

type Wxpay struct {
	config *WxpayConfig
	client *payment.Payment
}

var WxpayClient *Wxpay

func InitWxPay(conf *WxpayConfig) error {
	var err error
	WxpayClient, err = NewWxpay(conf)
	return err
}

func NewWxpay(conf *WxpayConfig) (*Wxpay, error) {
	if len(conf.MchAPIv3Key) != 32 {
		return nil, errors.New("wechatpay: mch api v3 key must be 32 bytes")
	}
	privateKey, err := loadRSAPrivateKey(conf.MchPrivateKey)
	if err != nil {
		return nil, err
	}
	// PowerWeChat 的签名器（Base、Order 等每个服务客户端各持有一个实例）只支持 PKCS#8 私钥，
	// 这里统一转成 PKCS#8 后经 KeyPath 传入；KeyPath 不是磁盘已有文件时，SDK 会把它当私钥内容直接解析
	pkcs8, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: pkcs8})

	app, err := payment.NewPayment(&payment.UserConfig{
		AppID:       conf.Appid,
		MchID:       conf.MchID,
		SerialNo:    conf.MchCertificateSerialNumber,
		MchApiV3Key: conf.MchAPIv3Key,
		KeyPath:     string(keyPEM),
		NotifyURL:   conf.NotifyURL,
	})
	if err != nil {
		return nil, err
	}

	return &Wxpay{
		config: conf,
		client: app,
	}, nil
}

// TradePreCreate Native 扫码支付下单，返回支付二维码链接
// amount 单位为分
// https://pay.weixin.qq.com/doc/v3/merchant/4012791863
func (wx *Wxpay) TradePreCreate(outTradeNo string, subject string, amount uint) (any, string, error) {
	res, err := wx.client.Order.TransactionNative(context.Background(), &request.RequestNativePrepay{
		PrepayBase: request.PrepayBase{
			AppID:     wx.config.Appid,
			MchID:     wx.config.MchID,
			NotifyUrl: wx.config.NotifyURL,
		},
		Description: subject,
		OutTradeNo:  outTradeNo,
		Amount: &request.NativeAmount{
			Total:    int(amount),
			Currency: "CNY",
		},
	})
	if err != nil {
		return nil, "", err
	}
	return res, res.CodeURL, nil
}

// TradeQuery 按商户订单号查询订单
// https://pay.weixin.qq.com/doc/v3/merchant/4012791908
func (wx *Wxpay) TradeQuery(outTradeNo string) (*Trade, error) {
	res, err := wx.client.Order.QueryByOutTradeNumber(context.Background(), outTradeNo)
	if err != nil {
		return nil, err
	}
	return &Trade{Transaction: res}, nil
}

// AckNotification 应答微信支付回调
// https://pay.weixin.qq.com/doc/v3/merchant/4012791906
func (wx *Wxpay) AckNotification(w http.ResponseWriter) error {
	// 验签通过：HTTP 应答状态码返回 200 或 204，无需返回应答报文
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(""))
	return err
}

// Client 返回底层 PowerWeChat 支付实例，可用于调用未封装的接口
func (wx *Wxpay) Client() *payment.Payment {
	return wx.client
}

// Notify 解析微信支付回调：解密回调报文并返回订单
// 解密失败（AES-GCM 校验不过）会返回错误，此时应返回 5xx 让微信重推
// https://pay.weixin.qq.com/doc/v3/merchant/4012791907
func (wx *Wxpay) Notify(req *http.Request) (*Trade, string, error) {
	var (
		txn        *models.Transaction
		outTradeNo string
	)
	_, err := wx.client.HandlePaidNotify(req, func(message *notifyRequest.RequestNotify, transaction *models.Transaction, fail func(message string)) interface{} {
		txn = transaction
		outTradeNo = transaction.OutTradeNo
		return nil
	})
	if err != nil {
		return nil, "", err
	}
	if txn == nil {
		return nil, "", errors.New("wechatpay: empty notify transaction")
	}
	return &Trade{Transaction: convertTransaction(txn)}, outTradeNo, nil
}

// convertTransaction 回调通知的 Transaction 转为查询订单的 ResponseOrder，统一 Trade.Transaction 类型
func convertTransaction(txn *models.Transaction) *wxresponse.ResponseOrder {
	if txn == nil {
		return nil
	}
	return &wxresponse.ResponseOrder{
		Amount:         txn.Amount,
		AppID:          txn.AppID,
		Attach:         txn.Attach,
		BankType:       txn.BankType,
		MchID:          txn.MchID,
		OutTradeNo:     txn.OutTradeNo,
		Payer:          txn.Payer,
		SuccessTime:    txn.SuccessTime,
		TradeState:     txn.TradeState,
		TradeStateDesc: txn.TradeStateDesc,
		TradeType:      txn.TradeType,
		TransactionID:  txn.TransactionID,
		SceneInfo:      txn.SceneInfo,
	}
}

func loadRSAPrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("wechatpay: invalid private key PEM")
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("wechatpay: private key is not RSA")
		}
		return rsaKey, nil
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
