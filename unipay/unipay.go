package unipay

import (
	"net/http"

	wxresponse "github.com/ArtisanCloud/PowerWeChat/v3/src/payment/order/response"
	"github.com/smartwalle/alipay/v3"
)

type unipay interface {
	TradePreCreate(outTradeNo string, subject string, amount uint) (any, string, error)
	TradeQuery(outTradeNo string) (*Trade, error)
	AckNotification(w http.ResponseWriter) error
	Notify(req *http.Request) (*Trade, string, error)
}

type Trade struct {
	TradeQueryRsp *alipay.TradeQueryRsp     `json:"TradeQueryRsp,omitempty"`
	Notification  *alipay.Notification      `json:"Notification,omitempty"`
	Transaction   *wxresponse.ResponseOrder `json:"Transaction,omitempty"`
}
