package contract

import "github.com/gin-gonic/gin"

type Handler interface {
	Route(engine *gin.Engine)
}
