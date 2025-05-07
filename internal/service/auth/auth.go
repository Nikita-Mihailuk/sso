package auth

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"sso/internal/domain/models"
	"sso/internal/repository"
	"sso/pkg/jwt"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	log          *zap.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

// NewAuthService returns new AuthService
func NewAuthService(log *zap.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passwordHash []byte) (userID int64, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID int) (models.App, error)
}

func (a *AuthService) Login(ctx context.Context, email, password string, appID int) (string, error) {
	log := a.log.With(zap.String("email", email))

	log.Info("login user")

	user, err := a.userProvider.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			log.Info("user not found", zap.Error(err))
			return "", ErrInvalidCredentials
		}
		log.Error("failed to get user", zap.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		log.Error("invalid credentials", zap.Error(err))
		return "", ErrInvalidCredentials
	}

	app, err := a.appProvider.GetApp(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
			log.Info("app not found", zap.Error(err))
			return "", ErrInvalidAppId
		}
		log.Error("failed to get app", zap.Error(err))
		return "", err
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to create token", zap.Error(err))
		return "", err
	}

	return token, nil
}

func (a *AuthService) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	log := a.log.With(zap.String("email", email))

	log.Info("registering new user")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate hash password", zap.Error(err))
		return 0, err
	}

	id, err := a.userSaver.SaveUser(ctx, email, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			log.Error("user already exists", zap.Error(err))
			return 0, ErrUserExists
		}
		log.Error("failed to save user", zap.Error(err))
		return 0, err
	}

	return id, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	log := a.log.With(zap.Int64("user_id", userID))
	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, err
	}

	log.Info("checked if user is admin", zap.Bool("isAdmin", isAdmin))
	return isAdmin, nil
}
