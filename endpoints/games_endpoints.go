package endpoints

import (
	"database/sql"
	"errors"
	"informatikaas_crud/sql_queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GameResponse represents a game response.
type GameResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int64  `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GameWithUserResponse represents a game with user information.
type GameWithUserResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// GameUpdateRequest represents a request to update a game.
type GameUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// CreateGame creates a new game.
func CreateGame(ctx *gin.Context, ed *EndpointData) {
	var gameReq struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Script      []byte `json:"script" binding:"required"`
	}
	if _, ok := BindJSON[struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Script      []byte `json:"script"`
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

	gameID, err := ed.DB.Queries.CreateGame(ed.DB.Ctx, sql_queries.CreateGameParams{
		Title:       gameReq.Title,
		Description: gameReq.Description,
		Script:      gameReq.Script,
		UserID:      userIDInt,
	})

	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to create game")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Game created successfully",
		"game_id": gameID,
	})
}

// GetGame retrieves a game by ID.
func GetGame(ctx *gin.Context, ed *EndpointData) {
	gameID, err := strconv.ParseInt(ctx.Param("game_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid game ID")
		return
	}

	game, err := ed.DB.Queries.GetGameByID(ed.DB.Ctx, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendError(ctx, http.StatusNotFound, "Game not found")
			return
		}
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch game")
		return
	}

	createdAtString := ""
	if game.CreatedAt.Valid {
		createdAtString = game.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	updatedAtString := ""
	if game.UpdatedAt.Valid {
		updatedAtString = game.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}

	ctx.JSON(http.StatusOK, GameResponse{
		ID:          game.ID,
		Title:       game.Title,
		Description: game.Description,
		UserID:      game.UserID,
		CreatedAt:   createdAtString,
		UpdatedAt:   updatedAtString,
	})
}

// GetGameWithUser retrieves a game with user information.
func GetGameWithUser(ctx *gin.Context, ed *EndpointData) {
	gameID, err := strconv.ParseInt(ctx.Param("game_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid game ID")
		return
	}

	gameWithUser, err := ed.DB.Queries.GetGameWithUser(ed.DB.Ctx, gameID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			SendError(ctx, http.StatusNotFound, "Game not found")
			return
		}
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch game")
		return
	}

	createdAtString := ""
	if gameWithUser.CreatedAt.Valid {
		createdAtString = gameWithUser.CreatedAt.Time.Format("2006-01-02 15:04:05")
	}
	updatedAtString := ""
	if gameWithUser.UpdatedAt.Valid {
		updatedAtString = gameWithUser.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	}

	ctx.JSON(http.StatusOK, GameWithUserResponse{
		ID:          gameWithUser.ID,
		Title:       gameWithUser.Title,
		Description: gameWithUser.Description,
		UserID:      gameWithUser.UserID,
		Username:    gameWithUser.UserUsername,
		DisplayName: gameWithUser.UserDisplayName.String,
		CreatedAt:   createdAtString,
		UpdatedAt:   updatedAtString,
	})
}

// ListGames lists all games with pagination.
func ListGames(ctx *gin.Context, ed *EndpointData) {
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

	games, err := ed.DB.Queries.ListGames(ed.DB.Ctx, sql_queries.ListGamesParams{
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch games")
		return
	}

	var gameResponses []GameResponse
	for _, game := range games {
		createdAtString := ""
		if game.CreatedAt.Valid {
			createdAtString = game.CreatedAt.Time.Format("2006-01-02 15:04:05")
		}
		updatedAtString := ""
		if game.UpdatedAt.Valid {
			updatedAtString = game.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		}

		gameResponses = append(gameResponses, GameResponse{
			ID:          game.ID,
			Title:       game.Title,
			Description: game.Description,
			UserID:      game.UserID,
			CreatedAt:   createdAtString,
			UpdatedAt:   updatedAtString,
		})
	}

	ctx.JSON(http.StatusOK, gameResponses)
}

// UpdateGame updates a game's information.
func UpdateGame(ctx *gin.Context, ed *EndpointData) {
	gameID, err := strconv.ParseInt(ctx.Param("game_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid game ID")
		return
	}

	var updateReq GameUpdateRequest
	if _, ok := BindJSON[GameUpdateRequest](ctx); !ok {
		return
	}

	// Check if the game exists
	_, err = ed.DB.Queries.GetGameByID(ed.DB.Ctx, gameID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch game")
		return
	}

	// Update the game
	_, err = ed.DB.Queries.UpdateGame(ed.DB.Ctx, sql_queries.UpdateGameParams{
		Title:       updateReq.Title,
		Description: updateReq.Description,
		ID:          gameID,
	})

	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to update game")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Game updated successfully",
	})
}

// DeleteGame deletes a game by ID.
func DeleteGame(ctx *gin.Context, ed *EndpointData) {
	gameID, err := strconv.ParseInt(ctx.Param("game_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid game ID")
		return
	}

	// Check if the game exists
	_, err = ed.DB.Queries.GetGameByID(ed.DB.Ctx, gameID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch game")
		return
	}

	// Delete the game
	err = ed.DB.Queries.DeleteGame(ed.DB.Ctx, gameID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to delete game")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Game deleted successfully",
	})
}
