package endpoints

import (
	"informatikaas_crud/config"
	"informatikaas_crud/database"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	HEAD   Method = "HEAD"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
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
	// Initialize JWT
	InitJWT(ed.Cfg)

	// Public routes
	for _, endpoint := range endpoints {
		if strings.HasPrefix(endpoint.Path, "/api/auth") {
			router.Handle(string(endpoint.Method), endpoint.Path, func(ctx *gin.Context) {
				endpoint.Func(ctx, ed)
			})
		} else {
			// Protected routes
			router.Handle(string(endpoint.Method), endpoint.Path, AuthMiddleware(), func(ctx *gin.Context) {
				endpoint.Func(ctx, ed)
			})
		}
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

func init() {
	endpoints = append(endpoints, Endpoint{
		Path:   "/",
		Method: HEAD,
		Func: func(ctx *gin.Context, ed *EndpointData) {
			ctx.Status(http.StatusOK)
		},
	})

	// Authentication endpoints (public)
	endpoints = append(endpoints, Endpoint{
		Path:   "/api/auth/login",
		Method: POST,
		Func:   Login,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/auth/register",
		Method: POST,
		Func:   Register,
	})

	// User endpoints (protected)
	endpoints = append(endpoints, Endpoint{
		Path:   "/api/users",
		Method: GET,
		Func:   ListUsers,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/users/:user_id",
		Method: GET,
		Func:   GetUser,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/users/:user_id",
		Method: PUT,
		Func:   UpdateUser,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/users/:user_id",
		Method: DELETE,
		Func:   DeleteUser,
	})

	// Game endpoints (protected)
	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games",
		Method: GET,
		Func:   ListGames,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games",
		Method: POST,
		Func:   CreateGame,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games/:game_id",
		Method: GET,
		Func:   GetGame,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games/:game_id/user",
		Method: GET,
		Func:   GetGameWithUser,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games/:game_id",
		Method: PUT,
		Func:   UpdateGame,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games/:game_id",
		Method: DELETE,
		Func:   DeleteGame,
	})

	// Comment endpoints (protected)
	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games/:game_id/comments",
		Method: GET,
		Func:   ListCommentsByGame,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/comments",
		Method: POST,
		Func:   CreateComment,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/comments/:comment_id",
		Method: GET,
		Func:   GetComment,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/comments/:comment_id/user/game",
		Method: GET,
		Func:   GetCommentWithUserAndGame,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/comments/:comment_id",
		Method: PUT,
		Func:   UpdateComment,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/comments/:comment_id",
		Method: DELETE,
		Func:   DeleteComment,
	})

	// Favorite endpoints (protected)
	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games/:game_id/favorite",
		Method: POST,
		Func:   ToggleFavorite,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/users/favorites",
		Method: GET,
		Func:   GetUserFavorites,
	})

	endpoints = append(endpoints, Endpoint{
		Path:   "/api/games/:game_id/favorites",
		Method: GET,
		Func:   GetGameFavorites,
	})
}
