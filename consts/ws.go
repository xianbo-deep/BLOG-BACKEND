package consts

const (
	WriteWait  = 10 * TimeRangeSecond // 写入数据操作的超时时间
	PongWait   = 60 * TimeRangeSecond // 等待客户端发送ping的超时时间
	PingPeriod = (PongWait * 9) / 10  // 发送ping的时间间隔，给客户端6s响应

)

const (
	MaxMessageSize    = 512  // 单大信息最大大小，字节
	SendBufferSize    = 256  // 通道的缓冲大小
	DefaultBufferSize = 1024 // upgrader的写和读缓冲区大小
)
