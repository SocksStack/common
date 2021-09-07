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
	cfg		HttpConfig
}

func Default(mode string) *server {
	gin.SetMode(mode)
	HttpServer = new(server)
	HttpServer.Engine = gin.Default()
	return HttpServer
}

func (s *server) Run()  {
	// 启动前操作
	if s.before != nil {
		s.before()
	}

	// gin 配置
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port),
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

func (s *server) Route(routes ...constract.Route) constract.IServer {
	for _, item := range routes {
		handler, ok := item.(contract.Handler)
		if !ok {
			log.Fatalln("handler 必须实现 [github.com/SockStack/common/contract#Handler]")
		}
		handler.Route(s.Engine)
	}
	return s
}

func (s *server) Use(middlewares ...constract.Middleware) constract.IServer {
	for _, middleware := range middlewares {
		m, ok := middleware.(gin.HandlerFunc)
		if !ok {
			log.Fatalln("handler 必须实现 [github.com/SockStack/common/contract#Handler]")
		}
		s.Engine.Use(m)
	}
	return s
}

func (s *server) SetConfig(cfg constract.IHttpConfig) constract.IServer {
	config, ok := cfg.(HttpConfig)
	if !ok {
		panic("gin 配置错误")
	}
	s.cfg = config
	return s
}