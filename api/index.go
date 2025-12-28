package api

import (
	"Blog-Backend/core"
	"Blog-Backend/router"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	rEngine *gin.Engine
	once    sync.Once
)

func Handler(w http.ResponseWriter, r *http.Request) {
	/* 初始化数据库 */
	core.Init()

	/* 初始化路由 */
	once.Do(func() {
		rEngine = router.SetupRouter()
	})

	/* 请求转给gin */
	rEngine.ServeHTTP(w, r)
}
