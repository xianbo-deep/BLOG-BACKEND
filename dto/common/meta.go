package common

// 存放到ctx的结构体
type RequestMeta struct {
	IP        string
	Referer   string
	UserAgent string
	Origin    string

	Device  string
	OS      string
	Browser string

	Medium string
	Source string
}
