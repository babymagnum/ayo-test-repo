package controller

import (
	"net/http"

	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/request"
	"github.com/ariefzainuri96/ayo-test/cmd/api/dto/response"
	"github.com/gin-gonic/gin"
)

// @Summary      Login
// @Description  Perform login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request	body	  request.LoginRequest	true "Login request"
// @Success      200  		{object}  response.LoginResponse
// @Failure      400  		{object}  response.BaseResponse
// @Failure      404  		{object}  response.BaseResponse
// @Router       /auth/login	[post]
func (app *Application) login(c *gin.Context) {
	var data request.LoginRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	user, token, err := app.Service.IAuth.Login(c, data)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid email/password!",
		})
		return
	}

	c.JSON(http.StatusOK, response.LoginResponse{
		BaseResponse: response.BaseResponse{
			Status:  http.StatusOK,
			Message: "Success",
		},
		Data: response.LoginData{
			ID:    int(user.ID),
			Token: token,
			Email: user.Email,
		},
	})
}

// @Summary      Register
// @Description  Perform register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request		body	  request.RegisterRequest	true "Register request"
// @Success      200  			{object}  response.BaseResponse
// @Failure      400  			{object}  response.BaseResponse
// @Failure      404  			{object}  response.BaseResponse
// @Router       /auth/register	[post]
func (app *Application) register(c *gin.Context) {
	var data request.RegisterRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.BaseResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	_, err := app.Service.IAuth.Register(c, data)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.BaseResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Status:  http.StatusOK,
		Message: "Success register account",
	})
}
