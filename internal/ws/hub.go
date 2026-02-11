package ws

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/response"
	"encoding/json"
	"log"
)

type Hub struct {
	clients    map[*WsClient]struct{} // 客户端map
	register   chan *WsClient         // 注册新客户端的通道
	unregister chan *WsClient         // 注销老客户端的通道
	broadcast  chan []byte            // 广播消息的通道
}

// 用通道是为了防止多个协程同时操作map导致线程不安全
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*WsClient]struct{}),
		register:   make(chan *WsClient), // 无缓冲通道
		unregister: make(chan *WsClient), // 无缓冲通道
		broadcast:  make(chan []byte, consts.SendBufferSize),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client) // 从map删除
				close(client.send)        // 关闭client的信息发送通道
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				// 通道已满 放不进去
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (h *Hub) Register(ws *WsClient) {
	h.register <- ws
}

func (h *Hub) Unregister(ws *WsClient) {
	h.unregister <- ws
}

func (h *Hub) Broadcast(message []byte) {
	select {
	case h.broadcast <- message:
	default:
		log.Printf("broadcast channel is full")
	}
}

func (h *Hub) BroadcastJSON(event response.Event) {
	b, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling event: %s", err)
		return
	}
	h.Broadcast(b)
}
