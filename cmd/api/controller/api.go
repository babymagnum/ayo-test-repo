package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/middleware"
	"github.com/ariefzainuri96/go-logstream/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Config struct {
	HTTPPort    int
	ShutdownTTL time.Duration
}

type Application struct {
	Config    Config
	Service   service.Service
	Validator *validator.Validate
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
	}

	// swagger — Gin butuh ANY wildcard
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

	// Start server in goroutine so we can watch ctx
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
		// Shutdown with timeout
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
