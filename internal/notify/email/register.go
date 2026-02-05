package email

import (
	"Blog-Backend/consts"
	"os"
)

func RegisterEmail() {
	cfg := EmailConfig{
		Host: os.Getenv(consts.EnvBaseURL),
		Port: os.Getenv(consts.EnvPort), // TODO 记得修改
		User: "",
		Pass: "",
	}

	emailClient := NewEmailClient(cfg)

}
