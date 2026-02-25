package unipay

import (
	"context"
	"fmt"
	"net/http"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	wxutils "github.com/wechatpay-apiv3/wechatpay-go/utils"
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
	client *core.Client
}

var WxpayClient *Wxpay

func InitWxPay(conf *WxpayConfig) error {
	var err error
	WxpayClient, err = NewWxpay(conf)
	return err
}

func NewWxpay(conf *WxpayConfig) (*Wxpay, error) {
	wx := &Wxpay{
		config: conf,
	}

	mchPrivateKey, err := wxutils.LoadPrivateKey(conf.MchPrivateKey)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(wx.config.MchID, wx.config.MchCertificateSerialNumber, mchPrivateKey, wx.config.MchAPIv3Key),
	}
	wx.client, err = core.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return wx, nil
}

// TradePreCreate
// amount 单位为分
func (wx *Wxpay) TradePreCreate(outTradeNo string, subject string, amount uint) (any, string, error) {

	trade := native.PrepayRequest{}
	trade.Appid = core.String(wx.config.Appid)
	trade.Mchid = core.String(wx.config.MchID)
	// 支付宝回调地址（需要在支付宝后台配置）
	// 支付成功后，支付宝会发送一个POST消息到该地址
	trade.NotifyUrl = core.String(wx.config.NotifyURL)
	// 支付成功之后，浏览器将会重定向到该 URL
	// trade.ReturnURL = "http://localhost:8088/return"
	// 支付标题
	trade.Description = core.String(subject)
	// 订单号，一个订单号只能支付一次
	trade.OutTradeNo = core.String(outTradeNo)
	// 销售产品码，与支付宝签约的产品码名称,目前仅支持FACE_TO_FACE_PAYMENT
	// trade.ProductCode = "FACE_TO_FACE_PAYMENT"
	// 金额
	trade.Amount = &native.Amount{
		Total: core.Int64(int64(amount)),
	}

	// pay.GoodsDetail = []*alipay.GoodsDetailItem{}
	svc := native.NativeApiService{Client: wx.client}
	// 得到prepay_id，以及调起支付所需的参数和签名
	res, _, err := svc.Prepay(context.Background(), trade)
	if err != nil {
		return nil, "", err
	}
	return res, *res.CodeUrl, err
}

func (wx *Wxpay) TradeQuery(outTradeNo string) (*Trade, error) {
	trade := native.QueryOrderByOutTradeNoRequest{}
	trade.Mchid = core.String(wx.config.MchID)
	trade.OutTradeNo = core.String(outTradeNo)

	svc := native.NativeApiService{Client: wx.client}
	res, _, err := svc.QueryOrderByOutTradeNo(context.Background(), trade)
	return &Trade{Transaction: res}, err
}

func (wx *Wxpay) AckNotification(w http.ResponseWriter) error {
	// 验签通过：商户需告知微信支付接收回调成功，HTTP应答状态码需返回200或204，无需返回应答报文。
	// 验签不通过：商户需告知微信支付接收回调失败，HTTP应答状态码需返回5XX或4XX，同时需返回以下应答报文：
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(""))
	return err
}

func (wx *Wxpay) Client() *core.Client {
	return wx.client
}

func (wx *Wxpay) Notify(req *http.Request) (*Trade, string, error) {

	mchPrivateKey, err := wxutils.LoadPrivateKey(wx.config.MchPrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("load private key fail: %w", err)
	}
	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	err = downloader.MgrInstance().RegisterDownloaderWithPrivateKey(context.Background(), mchPrivateKey, wx.config.MchCertificateSerialNumber, wx.config.MchID, wx.config.MchAPIv3Key)
	if err != nil {
		return nil, "", fmt.Errorf("register download with private key fail: %w", err)
	}
	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(wx.config.MchID)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(wx.config.MchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		return nil, "", fmt.Errorf("new notify handler fail: %w", err)
	}

	transaction := new(payments.Transaction)
	_, err = handler.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		// 如果验签未通过，或者解密失败
		return nil, "", fmt.Errorf("parse notify request fail: %w", err)
	}
	// 处理通知内容
	// fmt.Println(notifyReq.Summary)
	// fmt.Println(transaction.TransactionId)
	return &Trade{Transaction: transaction}, *transaction.OutTradeNo, nil
}
