package controller

import (
	"net/http"
	"strconv"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/response"
	"github.com/ariefzainuri96/ayo-test/cmd/api/middleware"
	"github.com/gin-gonic/gin"
)

// @Summary      List Teams
// @Description  Get all teams with pagination
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        page       query  int  false  "Page number"
// @Param        page_size  query  int  false  "Page size"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.TeamsResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /teams [get]
func (app *Application) listTeams(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	result, err := app.Service.ITeam.GetAll(c, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.TeamsResponse{
		BaseResponse: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Success",
		},
		Teams:      result.Data,
		Pagination: result.Pagination,
	})
}

// @Summary      Get Team
// @Description  Get team by ID
// @Tags         teams
// @Produce      json
// @Param        id   path  int  true  "Team ID"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.TeamResponse
// @Failure      404  {object}  response.BaseResponse
// @Router       /teams/{id} [get]
func (app *Application) getTeam(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid id",
		})
		return
	}

	team, err := app.Service.ITeam.GetByID(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, response.BaseResponse{
			Status:  http.StatusNotFound,
			Message: "Team not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.TeamResponse{
		BaseResponse: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Success",
		},
		Team: team,
	})
}

// @Summary      Create Team
// @Description  Add new team (admin only)
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        request  body  request.AddTeamRequest  true  "Team data"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.TeamResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /teams [post]
func (app *Application) createTeam(c *gin.Context) {
	user, ok := middleware.GetUserFromGin(c)
	if !ok || !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, response.BaseResponse{
			Status:  http.StatusForbidden,
			Message: "Admin access required",
		})
		return
	}

	var req request.AddTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	team, err := app.Service.ITeam.Create(c, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.TeamResponse{
		BaseResponse: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Success create team",
		},
		Team: team,
	})
}

// @Summary      Update Team
// @Description  Update team by ID (admin only)
// @Tags         teams
// @Accept       json
// @Produce      json
// @Param        id       path  int                   true  "Team ID"
// @Param        request  body  request.AddTeamRequest  true  "Team data"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.TeamResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /teams/{id} [put]
func (app *Application) updateTeam(c *gin.Context) {
	user, ok := middleware.GetUserFromGin(c)
	if !ok || !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, response.BaseResponse{
			Status:  http.StatusForbidden,
			Message: "Admin access required",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid id",
		})
		return
	}

	var req request.AddTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	team, err := app.Service.ITeam.Update(c, uint(id), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.TeamResponse{
		BaseResponse: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Success update team",
		},
		Team: team,
	})
}

// @Summary      Delete Team
// @Description  Delete team by ID (admin only, soft delete)
// @Tags         teams
// @Produce      json
// @Param        id   path  int  true  "Team ID"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /teams/{id} [delete]
func (app *Application) deleteTeam(c *gin.Context) {
	user, ok := middleware.GetUserFromGin(c)
	if !ok || !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, response.BaseResponse{
			Status:  http.StatusForbidden,
			Message: "Admin access required",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid id",
		})
		return
	}

	err = app.Service.ITeam.Delete(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Success delete team",
	})
}
