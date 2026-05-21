package controller

import (
	"net/http"
	"strconv"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/response"
	"github.com/gin-gonic/gin"
)

// @Summary      List Players by Team
// @Description  Get all players in a team with pagination
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        teamId    path  int  true  "Team ID"
// @Param        page      query  int  false  "Page number"
// @Param        page_size query  int  false  "Page size"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.PlayersResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /players/team/{teamId} [get]
func (app *Application) listPlayersByTeam(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid team id",
		})
		return
	}

	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: err.Error(),
		})
		return
	}

	result, err := app.Service.IPlayer.GetByTeam(c, uint(teamID), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.PlayersResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success"},
		Players:      result.Data,
		Pagination:   result.Pagination,
	})
}

// @Summary      Get Player
// @Description  Get player by ID
// @Tags         players
// @Produce      json
// @Param        id   path  int  true  "Player ID"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.PlayerResponse
// @Failure      404  {object}  response.BaseResponse
// @Router       /players/{id} [get]
func (app *Application) getPlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	player, err := app.Service.IPlayer.GetByID(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, response.BaseResponse{
			Status: http.StatusNotFound, Message: "Player not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.PlayerResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success"},
		Player:       player,
	})
}

// @Summary      Create Player
// @Description  Add a new player to a team (admin only)
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        teamId   path  int                         true  "Team ID"
// @Param        request  body  request.AddPlayerRequest     true  "Player data"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.PlayerResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /players/team/{teamId} [post]
func (app *Application) createPlayer(c *gin.Context) {
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid team id",
		})
		return
	}

	var req request.AddPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: err.Error(),
		})
		return
	}

	player, err := app.Service.IPlayer.Create(c, uint(teamID), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.PlayerResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success add player"},
		Player:       player,
	})
}

// @Summary      Update Player
// @Description  Update player by ID (admin only)
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id       path  int                         true  "Player ID"
// @Param        request  body  request.AddPlayerRequest     true  "Player data"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.PlayerResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /players/{id} [put]
func (app *Application) updatePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	var req request.AddPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: err.Error(),
		})
		return
	}

	player, err := app.Service.IPlayer.Update(c, uint(id), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.PlayerResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success update player"},
		Player:       player,
	})
}

// @Summary      Delete Player
// @Description  Delete player by ID (admin only, soft delete)
// @Tags         players
// @Produce      json
// @Param        id   path  int  true  "Player ID"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /players/{id} [delete]
func (app *Application) deletePlayer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	err = app.Service.IPlayer.Delete(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status: http.StatusOK, Message: "Success delete player",
	})
}
