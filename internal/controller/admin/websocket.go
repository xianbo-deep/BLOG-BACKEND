package admin

import (
	"Blog-Backend/consts"
	"Blog-Backend/dto/common"
	"Blog-Backend/internal/ws"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketController struct {
	hub      *ws.Hub
	upgrader *websocket.Upgrader
}

func NewWebSocketController(hub *ws.Hub) *WebSocketController {
	return &WebSocketController{
		hub: hub,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  consts.DefaultBufferSize,
			WriteBufferSize: consts.DefaultBufferSize,
			// 跨域配置
			CheckOrigin: func(r *http.Request) bool {
				allow := map[string]bool{
					os.Getenv(consts.EnvBaseURL):  true,
					os.Getenv(consts.EnvAdminURL): true,
				}
				return allow[r.Header.Get("Origin")]
			},
		},
	}
}

func (wc *WebSocketController) Handle(c *gin.Context) {
	conn, err := wc.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, consts.CodeInternal, err.Error())
		return
	}
	client := ws.NewWsClient(conn, wc.hub)
	wc.hub.Register(client)
	// 开启读和写
	go client.Write()
	go client.Read()
}
