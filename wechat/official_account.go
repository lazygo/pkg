package wechat

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
)

type OfficialAccountConfig struct {
	AppID          string `json:"appid" toml:"appid"`
	AppSecret      string `json:"app_secret" toml:"app_secret"`
	Token          string `json:"token" toml:"token"`
	EncodingAESKey string `json:"encoding_aes_key" toml:"encoding_aes_key"`
}

var officialAccountApp *officialAccount.OfficialAccount

// InitOfficialAccount 初始化默认公众号实例
// cache 传 nil 时使用进程内缓存，多实例部署可传入 redis 等共享缓存驱动（kernel.NewRedisClient）
func InitOfficialAccount(config OfficialAccountConfig, cache kernel.CacheInterface) error {
	app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  config.AppID,
		Secret: config.AppSecret,
		Token:  config.Token,
		AESKey: config.EncodingAESKey,
		Cache:  cache,
	})
	if err != nil {
		return err
	}
	officialAccountApp = app
	return nil
}

// OfficialAccount 返回 InitOfficialAccount 初始化的实例
func OfficialAccount() *officialAccount.OfficialAccount {
	return officialAccountApp
}
