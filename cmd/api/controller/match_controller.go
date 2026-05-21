package controller

import (
	"net/http"
	"strconv"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/response"
	"github.com/ariefzainuri96/go-logstream/cmd/api/middleware"
	"github.com/gin-gonic/gin"
)

// @Summary      List Matches
// @Description  Get all matches with pagination
// @Tags         matches
// @Accept       json
// @Produce      json
// @Param        page       query  int  false  "Page number"
// @Param        page_size  query  int  false  "Page size"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.MatchesResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /matches [get]
func (app *Application) listMatches(c *gin.Context) {
	var req request.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: err.Error(),
		})
		return
	}

	result, err := app.Service.IMatch.GetAll(c, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: "Internal server error",
		})
		return
	}

	matches := make([]response.MatchData, len(result.Data))
	for i, m := range result.Data {
		matches[i] = response.MatchData{
			ID: m.ID,
			HomeTeam: response.TeamData{
				ID:   m.HomeTeam.ID,
				Name: m.HomeTeam.Name,
				Logo: m.HomeTeam.LogoURL,
			},
			AwayTeam: response.TeamData{
				ID:   m.AwayTeamID,
				Name: m.AwayTeam.Name,
				Logo: m.AwayTeam.LogoURL,
			},
			MatchDate: m.MatchDate,
			MatchTime: m.MatchTime,
			Status:    m.Status,
		}
	}

	c.JSON(http.StatusOK, response.MatchesResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success"},
		Matches:      matches,
		Pagination:   result.Pagination,
	})
}

// @Summary      Get Match
// @Description  Get match by ID
// @Tags         matches
// @Produce      json
// @Param        id   path  int  true  "Match ID"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.MatchResponse
// @Failure      404  {object}  response.BaseResponse
// @Router       /matches/{id} [get]
func (app *Application) getMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	match, err := app.Service.IMatch.GetByID(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, response.BaseResponse{
			Status: http.StatusNotFound, Message: "Match not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.MatchResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success"},
		Match:        match,
	})
}

// @Summary      Schedule Match
// @Description  Create a new match schedule (admin only)
// @Tags         matches
// @Accept       json
// @Produce      json
// @Param        request  body  request.ScheduleMatchRequest  true  "Match schedule data"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.MatchResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /matches [post]
func (app *Application) scheduleMatch(c *gin.Context) {
	user, ok := middleware.GetUserFromGin(c)
	if !ok || !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, response.BaseResponse{
			Status: http.StatusForbidden, Message: "Admin access required",
		})
		return
	}

	var req request.ScheduleMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: err.Error(),
		})
		return
	}

	match, err := app.Service.IMatch.Schedule(c, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.MatchResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success schedule match"},
		Match:        match,
	})
}

// @Summary      Update Match
// @Description  Update match schedule by ID (admin only)
// @Tags         matches
// @Accept       json
// @Produce      json
// @Param        id       path  int                         true  "Match ID"
// @Param        request  body  request.ScheduleMatchRequest  true  "Match schedule data"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.MatchResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /matches/{id} [put]
func (app *Application) updateMatch(c *gin.Context) {
	user, ok := middleware.GetUserFromGin(c)
	if !ok || !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, response.BaseResponse{
			Status: http.StatusForbidden, Message: "Admin access required",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	var req request.ScheduleMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: err.Error(),
		})
		return
	}

	match, err := app.Service.IMatch.Update(c, uint(id), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.MatchResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success update match"},
		Match:        match,
	})
}

// @Summary      Delete Match
// @Description  Delete match by ID (admin only, soft delete)
// @Tags         matches
// @Produce      json
// @Param        id   path  int  true  "Match ID"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /matches/{id} [delete]
func (app *Application) deleteMatch(c *gin.Context) {
	user, ok := middleware.GetUserFromGin(c)
	if !ok || !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, response.BaseResponse{
			Status: http.StatusForbidden, Message: "Admin access required",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	err = app.Service.IMatch.Delete(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status: http.StatusOK, Message: "Success delete match",
	})
}

// @Summary      Report Match Result
// @Description  Submit match result with goals (admin only)
// @Tags         matches
// @Accept       json
// @Produce      json
// @Param        id       path  int                          true  "Match ID"
// @Param        request  body  request.ReportMatchRequest   true  "Match result data"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.BaseResponse
// @Failure      400  {object}  response.BaseResponse
// @Router       /matches/{id}/report [post]
func (app *Application) reportMatch(c *gin.Context) {
	user, ok := middleware.GetUserFromGin(c)
	if !ok || !user.IsAdmin {
		c.AbortWithStatusJSON(http.StatusForbidden, response.BaseResponse{
			Status: http.StatusForbidden, Message: "Admin access required",
		})
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	var req request.ReportMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: err.Error(),
		})
		return
	}

	err = app.Service.IMatch.ReportResult(c, uint(id), req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status: http.StatusInternalServerError, Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status: http.StatusOK, Message: "Success report match",
	})
}

// @Summary      Get Match Report
// @Description  Get match result with goals, winner, top scorer, and team win accumulations
// @Tags         matches
// @Produce      json
// @Param        id   path  int  true  "Match ID"
// @Security 	 ApiKeyAuth
// @Success      200  {object}  response.MatchReportResponse
// @Failure      404  {object}  response.BaseResponse
// @Router       /matches/{id}/report [get]
func (app *Application) getMatchReport(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status: http.StatusBadRequest, Message: "Invalid id",
		})
		return
	}

	matchReport, err := app.Service.IMatch.GetReport(c, uint(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, response.BaseResponse{
			Status: http.StatusNotFound, Message: err.Error(),
		})
		return
	}

	playerMap := make(map[uint]entity.Player)
	for _, p := range matchReport.PlayerGoals {
		playerMap[p.ID] = p
	}

	var goalData []response.MatchGoalData
	for _, g := range matchReport.Goals {
		goalData = append(goalData, response.MatchGoalData{
			ID:         g.ID,
			PlayerID:   g.PlayerID,
			PlayerName: playerMap[g.PlayerID].Name,
			Minute:     g.Minute,
			IsOwnGoal:  g.IsOwnGoal,
		})
	}

	c.JSON(http.StatusOK, response.MatchReportResponse{
		BaseResponse: response.BaseResponse{Status: http.StatusOK, Message: "Success"},
		Data: response.MatchReportData{
			Match:              matchReport.Match,
			HomeScore:          matchReport.Result.HomeScore,
			AwayScore:          matchReport.Result.AwayScore,
			Winner:             matchReport.Result.Winner,
			Goals:              goalData,
			TopScorer:          matchReport.TopScorer,
			TopScorerGoalCount: matchReport.TopScorerGoalCount,
			HomeTeam:           matchReport.HomeTeam,
			AwayTeam:           matchReport.AwayTeam,
			HomeTotalWins:      matchReport.HomeTeam.TotalWins,
			AwayTotalWins:      matchReport.AwayTeam.TotalWins,
		},
	})
}
