package service

import (
	"context"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"go.uber.org/zap"
)

type AuthServiceImpl struct {
	logger *zap.Logger
	store  store.Storage
}

func NewAuthService(store store.Storage, logger *zap.Logger) *AuthServiceImpl {
	return &AuthServiceImpl{
		logger: logger,
		store:  store,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, req request.RegisterRequest) (uint, error) {
	if req.Role == "" {
		req.Role = "user"
	}

	return s.store.IAuth.Register(ctx, req)
}

func (s *AuthServiceImpl) Login(ctx context.Context, req request.LoginRequest) (entity.User, string, error) {
	return s.store.IAuth.Login(ctx, req)
}