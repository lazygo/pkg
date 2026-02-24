package wechat

import (
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
)

var officialAccountConfig *offConfig.Config

type OfficialAccountConfig struct {
	AppID          string `json:"appid" toml:"appid"`
	AppSecret      string `json:"app_secret" toml:"app_secret"`
	Token          string `json:"token" toml:"token"`
	EncodingAESKey string `json:"encoding_aes_key" toml:"encoding_aes_key"`
}

func InitOfficialAccount(config OfficialAccountConfig, cache cache.Cache) error {
	// 使用memcache保存access_token，也可选择redis或自定义cache
	officialAccountConfig = &offConfig.Config{
		AppID:          config.AppID,
		AppSecret:      config.AppSecret,
		Token:          config.Token,
		EncodingAESKey: config.EncodingAESKey,
		Cache:          cache,
	}
	return nil
}

func OfficialAccount() *officialaccount.OfficialAccount {
	return wechat.NewWechat().GetOfficialAccount(officialAccountConfig)
}
