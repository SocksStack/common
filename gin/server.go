package gin

import (
	"context"
	"errors"
	"fmt"
	"github.com/SocksStack/common/constract"
	"github.com/SocksStack/common/gin/contract"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	HttpServer *server
)

type server struct {
	Engine 	*gin.Engine
	before 	constract.Callback
	after 	constract.Callback
}

func Default() *server {
	HttpServer = new(server)
	HttpServer.Engine = gin.Default()
	return HttpServer
}

func (s *server) Run(config constract.IHttpConfig)  {
	cfg, ok := config.(HttpConfig)
	if !ok {
		panic("请使用gin的配置")
		return
	}
	// 启动前操作
	if s.before != nil {
		s.before()
	}

	// gin 配置
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: s.Engine,
	}

	// 启动
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if s.after != nil {
		s.after()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	} else {
		fmt.Println("gin http server have been stopped!")
	}
}

func (s *server) Route(routes ...constract.Route) {
	for _, item := range routes {
		handler, ok := item.(contract.Handler)
		if !ok {
			log.Fatalln("handler 必须实现 [github.com/SockStack/common/contract#Handler]")
		}
		handler.Route(s.Engine)
	}
}

func (s *server) Middleware(middlewares ...constract.Middleware) {
	for _, middleware := range middlewares {
		m, ok := middleware.(gin.HandlerFunc)
		if !ok {
			log.Fatalln("handler 必须实现 [github.com/SockStack/common/contract#Handler]")
		}
		s.Engine.Use(m)
	}
}

func (s *server) SetMode(mode string) {
	gin.SetMode(mode)
}
