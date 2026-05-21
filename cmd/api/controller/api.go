package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ariefzainuri96/ayo-test/cmd/api/middleware"
	"github.com/ariefzainuri96/ayo-test/internal/service"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Config struct {
	HTTPPort    int
	ShutdownTTL time.Duration
}

type Application struct {
	Config  Config
	Service service.Service
}

func (app *Application) RunServer(ctx context.Context, cfg Config, logger *zap.Logger) error {
	router := gin.New()

	router.Use(middleware.Logging(logger))
	router.Use(middleware.Recoverer(logger))

	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", app.login)
			auth.POST("/register", app.register)
		}

		teams := v1.Group("/teams")
		teams.Use(middleware.Authentication())
		{
			teams.GET("", app.listTeams)
			teams.GET("/:id", app.getTeam)
			teams.POST("", app.createTeam)
			teams.PUT("/:id", app.updateTeam)
			teams.DELETE("/:id", app.deleteTeam)
		}

		players := v1.Group("/players")
		players.Use(middleware.Authentication())
		{
			players.GET("/team/:teamId", app.listPlayersByTeam)
			players.GET("/:id", app.getPlayer)
			players.POST("/team/:teamId", app.createPlayer)
			players.PUT("/:id", app.updatePlayer)
			players.DELETE("/:id", app.deletePlayer)
		}

		matches := v1.Group("/matches")
		matches.Use(middleware.Authentication())
		{
			matches.GET("", app.listMatches)
			matches.GET("/:id", app.getMatch)
			matches.POST("", app.scheduleMatch)
			matches.PUT("/:id", app.updateMatch)
			matches.DELETE("/:id", app.deleteMatch)
			matches.POST("/:id/report", app.reportMatch)
			matches.GET("/:id/report", app.getMatchReport)
		}
	}

	router.GET("/v1/swagger/*any", gin.WrapH(
		httpSwagger.Handler(httpSwagger.URL("doc.json")),
	))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.HTTPPort),
		Handler:      router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  1 * time.Minute,
	}

	serverErrCh := make(chan error, 1)
	go func() {
		logger.Info("starting http server", zap.Int("port", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
			return
		}
		serverErrCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		logger.Info("http server shutting down")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("http server shutdown error", zap.Error(err))
			return err
		}
		logger.Info("http server stopped")
		return nil
	case err := <-serverErrCh:
		return err
	}
}
