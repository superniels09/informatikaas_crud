package endpoints

import (
	"informatikaas_crud/config"
	"informatikaas_crud/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

type EndpointData struct {
	Cfg *config.Config
	DB  *database.DB
}

type EndpointFunc func(ctx *gin.Context, ed *EndpointData)

type Endpoint struct {
	Path   string
	Method Method
	Func   EndpointFunc
}

var endpoints = []Endpoint{}

func Attach(router *gin.Engine, ed *EndpointData) {
	for _, endpoint := range endpoints {
		router.Handle(string(endpoint.Method), endpoint.Path, func(ctx *gin.Context) {
			endpoint.Func(ctx, ed)
		})
	}
}

func InternalError(ctx *gin.Context, message string, err error) {
	log.Printf("%s: %v", message, err)
	SendError(ctx, http.StatusInternalServerError, message)
}

func SendError(ctx *gin.Context, code int, message string) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": message})
}

func BindJSON[T any](ctx *gin.Context) (T, bool) {
	var req T
	if err := ctx.BindJSON(&req); err != nil {
		SendError(ctx, http.StatusBadRequest, "invalid json")
		return req, false
	}
	return req, true
}

func Message(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{"message": message})
}
