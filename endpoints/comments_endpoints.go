package endpoints

import (
	"database/sql"
	"errors"
	"informatikaas_crud/sql_queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentResponse represents a comment response.
type CommentResponse struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	UserID    int64  `json:"user_id"`
	GameID    int64  `json:"game_id"`
	CreatedAt string `json:"created_at"`
}

// CommentWithUserAndGameResponse represents a comment with user and game information.
type CommentWithUserAndGameResponse struct {
	ID          int64  `json:"id"`
	Content     string `json:"content"`
	UserID      int64  `json:"user_id"`
	GameID      int64  `json:"game_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	GameTitle   string `json:"game_title"`
	CreatedAt   string `json:"created_at"`
}

// CreateComment creates a new comment on a game.
func CreateComment(ctx *gin.Context, ed *EndpointData) {
	var commentReq struct {
		Content string `json:"content" binding:"required"`
		GameID  int64  `json:"game_id" binding:"required"`
	}
	if _, ok := BindJSON[struct {
		Content string `json:"content"`
		GameID  int64  `json:"game_id"`
	}](ctx); !ok {
		return
	}

	userID, exists := ctx.Get("userID")
	if !exists {
		SendError(ctx, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userIDInt, err := strconv.ParseInt(userID.(string), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to parse user ID")
		return
	}

	commentID, err := ed.DB.Queries.CreateComment(ed.DB.Ctx, sql_queries.CreateCommentParams{
		Content: commentReq.Content,
		UserID:  userIDInt,
		GameID:  commentReq.GameID,
	})

	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to create comment")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Comment created successfully",
		"comment_id": commentID,
	})
}

// GetComment retrieves a comment by ID.
func GetComment(ctx *gin.Context, ed *EndpointData) {
	commentID, err := strconv.ParseInt(ctx.Param("comment_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	comment, err := ed.DB.Queries.GetCommentByID(ed.DB.Ctx, commentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendError(ctx, http.StatusNotFound, "Comment not found")
			return
		}
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch comment")
		return
	}

	createdAtString := ""
	if comment.CreatedAt.Valid {
		createdAtString = comment.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}

	ctx.JSON(http.StatusOK, CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		GameID:    comment.GameID,
		CreatedAt: createdAtString,
	})
}

// GetCommentWithUserAndGame retrieves a comment with user and game information.
func GetCommentWithUserAndGame(ctx *gin.Context, ed *EndpointData) {
	commentID, err := strconv.ParseInt(ctx.Param("comment_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	commentWithUserAndGame, err := ed.DB.Queries.GetCommentWithUserAndGame(ed.DB.Ctx, commentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendError(ctx, http.StatusNotFound, "Comment not found")
			return
		}
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch comment")
		return
	}

	createdAtString := ""
	if commentWithUserAndGame.CreatedAt.Valid {
		createdAtString = commentWithUserAndGame.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}

	ctx.JSON(http.StatusOK, CommentWithUserAndGameResponse{
		ID:          commentWithUserAndGame.ID,
		Content:     commentWithUserAndGame.Content,
		UserID:      commentWithUserAndGame.UserID,
		GameID:      commentWithUserAndGame.GameID,
		Username:    commentWithUserAndGame.UserUsername,
		DisplayName: commentWithUserAndGame.UserDisplayName.String,
		GameTitle:   commentWithUserAndGame.GameTitle,
		CreatedAt:   createdAtString,
	})
}

// ListCommentsByGame lists all comments for a specific game.
func ListCommentsByGame(ctx *gin.Context, ed *EndpointData) {
	gameID, err := strconv.ParseInt(ctx.Param("game_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid game ID")
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid page size")
		return
	}

	offset := (page - 1) * pageSize

	comments, err := ed.DB.Queries.ListGameCommentsWithPagination(ed.DB.Ctx, sql_queries.ListGameCommentsWithPaginationParams{
		GameID: gameID,
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch comments")
		return
	}

	var commentResponses []CommentResponse
	for _, comment := range comments {
		createdAtString := ""
		if comment.CreatedAt.Valid {
			createdAtString = comment.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		commentResponses = append(commentResponses, CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			GameID:    comment.GameID,
			CreatedAt: createdAtString,
		})
	}

	ctx.JSON(http.StatusOK, commentResponses)
}

// UpdateComment updates a comment's content.
func UpdateComment(ctx *gin.Context, ed *EndpointData) {
	commentID, err := strconv.ParseInt(ctx.Param("comment_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	var updateReq struct {
		Content string `json:"content" binding:"required"`
	}
	if _, ok := BindJSON[struct {
		Content string `json:"content"`
	}](ctx); !ok {
		return
	}

	// Check if the comment exists
	_, err = ed.DB.Queries.GetCommentByID(ed.DB.Ctx, commentID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch comment")
		return
	}

	// Update the comment
	_, err = ed.DB.Queries.UpdateComment(ed.DB.Ctx, sql_queries.UpdateCommentParams{
		Content: updateReq.Content,
		ID:      commentID,
	})
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to update comment")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comment updated successfully",
	})
}

// DeleteComment deletes a comment by ID.
func DeleteComment(ctx *gin.Context, ed *EndpointData) {
	commentID, err := strconv.ParseInt(ctx.Param("comment_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid comment ID")
		return
	}

	// Check if the comment exists
	_, err = ed.DB.Queries.GetCommentByID(ed.DB.Ctx, commentID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch comment")
		return
	}

	// Delete the comment
	err = ed.DB.Queries.DeleteComment(ed.DB.Ctx, commentID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to delete comment")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
	})
}
