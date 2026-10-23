package endpoints

import (
	"informatikaas_crud/sql_queries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FavoriteResponse represents a favorite response.
type FavoriteResponse struct {
	UserID    int64  `json:"user_id"`
	GameID    int64  `json:"game_id"`
	CreatedAt string `json:"created_at"`
}

// ToggleFavorite toggles a favorite (add or remove).
func ToggleFavorite(ctx *gin.Context, ed *EndpointData) {
	gameID, err := strconv.ParseInt(ctx.Param("game_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid game ID")
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

	// Check if the game exists
	_, err = ed.DB.Queries.GetGameByID(ed.DB.Ctx, gameID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch game")
		return
	}

	// Check if the favorite already exists
	_, err = ed.DB.Queries.IsGameFavorited(ed.DB.Ctx, sql_queries.IsGameFavoritedParams{
		UserID: userIDInt,
		GameID: gameID,
	})

	// If no error, it means the favorite exists, so remove it
	if err == nil {
		err = ed.DB.Queries.RemoveFavorite(ed.DB.Ctx, sql_queries.RemoveFavoriteParams{
			UserID: userIDInt,
			GameID: gameID,
		})
		if err != nil {
			SendError(ctx, http.StatusInternalServerError, "Failed to remove favorite")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Game unfavorited successfully",
		})
	} else {
		// Otherwise, add the favorite
		err = ed.DB.Queries.AddFavorite(ed.DB.Ctx, sql_queries.AddFavoriteParams{
			UserID: userIDInt,
			GameID: gameID,
		})
		if err != nil {
			SendError(ctx, http.StatusInternalServerError, "Failed to add favorite")
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Game favorited successfully",
		})
	}
}

// GetUserFavorites retrieves all games favorited by a user.
func GetUserFavorites(ctx *gin.Context, ed *EndpointData) {
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

	favorites, err := ed.DB.Queries.GetUserFavoritesWithGameInfo(ed.DB.Ctx, userIDInt)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch favorites")
		return
	}

	var favoriteResponses []FavoriteResponse
	for _, favorite := range favorites {
		createdAtString := ""
		if favorite.CreatedAt.Valid {
			createdAtString = favorite.FavoritedAt.Time.Format("2006-01-02 15:04:05")
		}
		favoriteResponses = append(favoriteResponses, FavoriteResponse{
			UserID:    favorite.UserID,
			GameID:    favorite.ID,
			CreatedAt: createdAtString,
		})
	}

	ctx.JSON(http.StatusOK, favoriteResponses)
}

// GetGameFavorites retrieves all users who favorited a game.
func GetGameFavorites(ctx *gin.Context, ed *EndpointData) {
	gameID, err := strconv.ParseInt(ctx.Param("game_id"), 10, 64)
	if err != nil {
		SendError(ctx, http.StatusBadRequest, "Invalid game ID")
		return
	}

	favorites, err := ed.DB.Queries.GetGameFavoritedByUsers(ed.DB.Ctx, gameID)
	if err != nil {
		SendError(ctx, http.StatusInternalServerError, "Failed to fetch favorites")
		return
	}

	var favoriteResponses []FavoriteResponse
	for _, favorite := range favorites {
		createdAtString := ""
		if favorite.CreatedAt.Valid {
			createdAtString = favorite.FavoritedAt.Time.Format("2006-01-02 15:04:05")
		}
		favoriteResponses = append(favoriteResponses, FavoriteResponse{
			UserID:    favorite.ID,
			GameID:    gameID,
			CreatedAt: createdAtString,
		})
	}

	ctx.JSON(http.StatusOK, favoriteResponses)
}
