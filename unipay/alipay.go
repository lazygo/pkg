package unipay

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/smartwalle/alipay/v3"
)

var _ unipay = (*Alipay)(nil)

type AlipayConfig struct {
	AppId        string `json:"appid" toml:"appid"`
	AliPublicKey string `json:"ali_public_key" toml:"ali_public_key"`
	PrivateKey   string `json:"mch_private_key" toml:"mch_private_key"`
	NotifyURL    string `json:"notify_url" toml:"notify_url"`
}

type Alipay struct {
	config *AlipayConfig
	client *alipay.Client
}

var AlipayClient *Alipay

func InitAlipay(conf *AlipayConfig) error {
	var err error
	AlipayClient, err = NewAlipay(conf)
	return err
}

func NewAlipay(conf *AlipayConfig) (*Alipay, error) {
	client, err := alipay.New(conf.AppId, conf.PrivateKey, true, func(c *alipay.Client) {
		c.LoadAliPayPublicKey(conf.AliPublicKey)
	})
	if err != nil {
		return nil, err
	}
	ali := &Alipay{
		config: conf,
		client: client,
	}

	return ali, nil
}

// TradePreCreate
// amount 单位为分
// https://opendocs.alipay.com/apis/api_1/alipay.trade.precreate
func (a *Alipay) TradePreCreate(outTradeNo string, subject string, amount uint) (any, string, error) {

	trade := alipay.TradePreCreate{}
	// 支付宝回调地址（需要在支付宝后台配置）
	// 支付成功后，支付宝会发送一个POST消息到该地址
	trade.NotifyURL = a.config.NotifyURL
	// 支付成功之后，浏览器将会重定向到该 URL
	// trade.ReturnURL = "http://localhost:8088/return"
	//支付标题
	trade.Subject = subject
	//订单号，一个订单号只能支付一次
	trade.OutTradeNo = outTradeNo
	//销售产品码，与支付宝签约的产品码名称,目前仅支持FACE_TO_FACE_PAYMENT
	trade.ProductCode = "FACE_TO_FACE_PAYMENT"
	//金额
	trade.TotalAmount = fmt.Sprintf("%.2f", float64(amount)/100)

	//pay.GoodsDetail = []*alipay.GoodsDetailItem{}
	res, err := a.client.TradePreCreate(context.Background(), trade)
	if err != nil {
		return res, "", err
	}
	if !res.IsSuccess() {
		return res, "", errors.New(res.Msg)
	}

	return res, res.QRCode, nil
}

// https://opendocs.alipay.com/apis/api_1/alipay.trade.query
func (a *Alipay) TradeQuery(outTradeNo string) (*Trade, error) {
	trade := alipay.TradeQuery{}
	trade.OutTradeNo = outTradeNo
	res, err := a.client.TradeQuery(context.Background(), trade)
	if err != nil {
		return nil, err
	}
	if !res.IsSuccess() {
		return nil, errors.New(res.Msg)
	}
	return &Trade{TradeQueryRsp: res}, nil
}

func (a *Alipay) AckNotification(w http.ResponseWriter) error {
	a.client.AckNotification(w)
	return nil
}

func (a *Alipay) Client() *alipay.Client {
	return a.client
}

func (a *Alipay) Notify(req *http.Request) (*Trade, string, error) {
	noti, err := a.client.GetTradeNotification(req)
	if err != nil {
		return nil, "", err
	}
	if noti == nil {
		return nil, "", errors.New("notify is empty")
	}

	return &Trade{Notification: noti}, noti.OutTradeNo, nil
}
