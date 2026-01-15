package common

import "strings"

type RequestHeader struct {
	Authorization string `header:"Authorization" binding:"omitempty"`
	Referer       string `header:"Referer" binding:"omitempty"`
	UserAgent     string `header:"User-Agent" binding:"omitempty"`
	ForwardIP     string `header:"X-Forwarded-For" binding:"omitempty"`
	RealIP        string `header:"X-Real-IP" binding:"omitempty,ip"`
	Origin        string `header:"Origin" binding:"omitempty"`
}

func (r RequestHeader) GetFirstFowardIP() string {
	if r.ForwardIP == "" {
		return ""
	}
	ip := strings.Split(r.ForwardIP, ",")[0]
	return strings.TrimSpace(ip)
}
