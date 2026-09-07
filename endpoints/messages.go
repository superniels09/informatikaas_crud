package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MessagesPostRequest struct {
	Message string `json:"message"`
}

type MessagesListResponse struct {
	Messages []ResponseMessage `json:"messages"`
}

type ResponseMessage struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func init() {
	endpoints = append(endpoints, Endpoint{
		Path:   "/messages/post",
		Method: POST,
		Func: func(ctx *gin.Context, ed *EndpointData) {
			req, valid := BindJSON[MessagesPostRequest](ctx)
			if !valid {
				return
			}
			_, err := ed.DB.Queries.MessagesCreate(ed.DB.Ctx, req.Message)
			if err != nil {
				InternalError(ctx, "failed to create message", err)
				return
			}
			Message(ctx, http.StatusOK, "message created successfully")
		},
	})
	endpoints = append(endpoints, Endpoint{
		Path:   "/messages/list",
		Method: GET,
		Func: func(ctx *gin.Context, ed *EndpointData) {
			messageRows, err := ed.DB.Queries.MessagesList(ed.DB.Ctx)
			if err != nil {
				InternalError(ctx, "failed to list messages", err)
				return
			}
			resp := MessagesListResponse{
				Messages: make([]ResponseMessage, 0, len(messageRows)),
			}
			for _, messageRow := range messageRows {
				resp.Messages = append(resp.Messages, ResponseMessage{
					Content:   messageRow.Content,
					CreatedAt: messageRow.CreatedAt,
				})
			}
			ctx.JSON(http.StatusOK, resp)
		},
	})
}
