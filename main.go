package main

import (
	"Blog-Backend/bootstrap"
	"Blog-Backend/consts"
	"Blog-Backend/core"
	"Blog-Backend/internal/task"
	"Blog-Backend/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rEngine *gin.Engine
	once    sync.Once
)

func main() {
	/* 初始化数据库 */
	err := core.Init()

	if err != nil {
		panic(err)
	}

	/* 初始化组件 它们依赖于数据库 */
	comps := bootstrap.InitComponet()

	/* 启动定时任务*/
	cronStop := task.InitCron(comps)

	/* 初始化路由 */
	rEngine = router.SetupRouter(comps)

	/* 获取端口 */
	port := os.Getenv(consts.EnvPort)
	if port == "" {
		port = "8080"
	}

	/* 创建server */
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           rEngine,
		ReadHeaderTimeout: 5 * consts.TimeRangeSecond,
	}

	/* 接收退出信号*/
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	/* 监听端口，启动后端 */
	go func() {
		log.Printf("后端服务启动，监听端口 %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %s\n", err)
		}
	}()
	
	/* 等待停止信号 */
	<-quit
	log.Println("服务正在关闭")

	/* 关闭定时任务 */
	if cronStop != nil {
		cronStop()
	}

	/* 优雅关闭http */
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("服务关闭出错: %v\n", err)
		// 强制关闭
		_ = srv.Close()
	}
	log.Printf("后端服务退出")

}
