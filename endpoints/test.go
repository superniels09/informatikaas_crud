package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	endpoints = append(endpoints, Endpoint{
		Path:   "ping",
		Method: GET,
		Func: func(ctx *gin.Context, ed *EndpointData) {
			ctx.String(http.StatusOK, "pong!")
		},
	})
	endpoints = append(endpoints, Endpoint{
		Path:   "reverse",
		Method: GET,
		Func: func(ctx *gin.Context, ed *EndpointData) {
			str := ctx.Query("str")
			if str == "" {
				ctx.String(http.StatusBadRequest, "missing str query parameter")
				return
			}

			// Convert to runes, reverse and convert back
			runes := []rune(str)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			str = string(runes)

			ctx.String(http.StatusOK, str)
		},
	})
}
