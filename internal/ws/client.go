package ws

import (
	"Blog-Backend/consts"
	"time"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	conn *websocket.Conn // 连接对象
	hub  *Hub
	send chan []byte // 客户端自己的消息队列
}

func NewWsClient(conn *websocket.Conn, hub *Hub) *WsClient {
	return &WsClient{
		conn: conn,
		hub:  hub,
		send: make(chan []byte, consts.SendBufferSize),
	}
}

func (ws *WsClient) Read() {
	// 客户端断连
	defer func() {
		ws.hub.Unregister(ws)
		_ = ws.conn.Close()
	}()
	// 设定传输的最大字节
	ws.conn.SetReadLimit(consts.MaxMessageSize)
	// 设定超时时间
	_ = ws.conn.SetReadDeadline(time.Now().Add(consts.PongWait))
	// 心跳检测
	ws.conn.SetPongHandler(func(string) error {
		_ = ws.conn.SetReadDeadline(time.Now().Add(consts.PongWait))
		return nil
	})
	for {
		if _, _, err := ws.conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (ws *WsClient) Write() {
	// 创建心跳定时器
	ticker := time.NewTicker(consts.PingPeriod)

	// 清理资源
	defer func() {
		ticker.Stop()
		ws.conn.Close()
	}()

	for {
		select {
		case message, ok := <-ws.send:
			if !ok {
				_ = ws.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := ws.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			_ = ws.conn.SetWriteDeadline(time.Now().Add(consts.WriteWait))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
