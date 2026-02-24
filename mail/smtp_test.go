package mail

import "testing"

func TestMainTo(t *testing.T) {
	Init(SmtpConfig{
		Host:     "smtp.qcloudmail.com",
		Port:     465,
		From:     "admin@p2link.cn",
		Password: "as2r0OSH89KxM74Q",
	})
	err := MailTo("lzp9421@qq.com", "测试", "测试内容")
	if err != nil {
		t.Errorf("send mail fial %v", err)
	}
}
