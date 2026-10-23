package endpoints

import (
	"database/sql"
	"errors"
	"informatikaas_crud/sql_queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserResponse represents a user response.
type UserResponse struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UserUpdateRequest represents a request to update a user.
type UserUpdateRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

// GetUser retrieves a user by ID.
func GetUser(ctx *gin.Context, ed *EndpointData) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := ed.DB.Queries.GetUserByID(ed.DB.Ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendError(ctx, http.StatusNotFound, "User not found")
			return
		}
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	createdAtString := ""
	if user.CreatedAt.Valid {
		createdAtString = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	updatedAtString := ""
	if user.UpdatedAt.Valid {
		updatedAtString = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}

	ctx.JSON(http.StatusOK, UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName.String,
		Bio:         user.Bio.String,
		CreatedAt:   createdAtString,
		UpdatedAt:   updatedAtString,
	})
}

// ListUsers lists all users with pagination.
func ListUsers(ctx *gin.Context, ed *EndpointData) {
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

	users, err := ed.DB.Queries.ListUsers(ed.DB.Ctx, sql_queries.ListUsersParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		createdAtString := ""
		if user.CreatedAt.Valid {
			createdAtString = user.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		updatedAtString := ""
		if user.UpdatedAt.Valid {
			updatedAtString = user.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}

		userResponses = append(userResponses, UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName.String,
			Bio:         user.Bio.String,
			CreatedAt:   createdAtString,
			UpdatedAt:   updatedAtString,
		})
	}

	ctx.JSON(http.StatusOK, userResponses)
}

// UpdateUser updates a user's information.
func UpdateUser(ctx *gin.Context, ed *EndpointData) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var updateReq UserUpdateRequest
	if _, ok := BindJSON[UserUpdateRequest](ctx); !ok {
		return
	}

	// Check if the user exists
	_, err = ed.DB.Queries.GetUserByID(ed.DB.Ctx, userID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// Update the user
	_, err = ed.DB.Queries.UpdateUser(ed.DB.Ctx, sql_queries.UpdateUserParams{
		Username:    updateReq.Username,
		Email:       updateReq.Email,
		DisplayName: sql.NullString{String: updateReq.DisplayName, Valid: true},
		Bio:         sql.NullString{String: updateReq.Bio, Valid: true},
		ID:          userID,
	})

	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to update user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
	})
}

// DeleteUser deletes a user by ID.
func DeleteUser(ctx *gin.Context, ed *EndpointData) {
	userID, err := strconv.ParseInt(ctx.Param("user_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	// Check if the user exists
	_, err = ed.DB.Queries.GetUserByID(ed.DB.Ctx, userID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	// Delete the user
	err = ed.DB.Queries.DeleteUser(ed.DB.Ctx, userID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
