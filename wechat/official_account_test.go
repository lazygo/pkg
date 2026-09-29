package wechat

import "testing"

func TestInitOfficialAccount(t *testing.T) {
	err := InitOfficialAccount(OfficialAccountConfig{
		AppID:     "appid",
		AppSecret: "secret",
		Token:     "token",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if OfficialAccount() == nil {
		t.Fatal("expected non-nil official account app")
	}
}
