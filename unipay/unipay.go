package unipay

import (
	"net/http"

	"github.com/smartwalle/alipay/v3"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
)

type unipay interface {
	TradePreCreate(outTradeNo string, subject string, amount uint) (any, string, error)
	TradeQuery(outTradeNo string) (*Trade, error)
	AckNotification(w http.ResponseWriter) error
	Notify(req *http.Request) (*Trade, string, error)
}

type Trade struct {
	TradeQueryRsp *alipay.TradeQueryRsp `json:"TradeQueryRsp,omitempty"`
	Notification  *alipay.Notification  `json:"Notification,omitempty"`
	Transaction   *payments.Transaction `json:"Transaction,omitempty"`
}
