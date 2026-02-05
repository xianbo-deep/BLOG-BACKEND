package email

type EmailConfig struct {
	Host string // 执行发送的主机
	Port string // 执行发送的端口
	User string // 目标用户邮箱
	Pass string // SMTP授权码
}

type EmailClient struct {
	cfg EmailConfig
}

func NewEmailClient(cfg EmailConfig) *EmailClient {
	return &EmailClient{cfg: cfg}
}
