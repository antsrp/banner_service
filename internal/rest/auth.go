package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/antsrp/banner_service/internal/domain/models"
	"github.com/antsrp/banner_service/internal/domain/models/requests"
	"github.com/antsrp/banner_service/internal/service"
	"github.com/antsrp/banner_service/pkg/logger"
	"github.com/gin-gonic/gin"
)

const authusertag = "auth-user-tag-data"

type authHandler struct {
	storage service.UserStorager
	logger  logger.Logger
}

func newAuthHandler(us service.UserStorager, logger logger.Logger) authHandler {
	return authHandler{
		storage: us,
		logger:  logger,
	}
}

func (h authHandler) parseToken(ctx *gin.Context) (models.User, error) {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return models.User{}, fmt.Errorf("no token provided")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return models.User{}, fmt.Errorf("bad token provided")
	}

	if parts[0] != "Bearer" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return models.User{}, fmt.Errorf("not a bearer token")
	}

	user, err := h.storage.UserByToken(parts[1])
	if err != nil {
		//h.logger.Error("error while parsing token: %v", err.Cause().Error())
		if err.IsInternal() {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, service.ErrDefaultInternalError.Error())
		} else {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		return models.User{}, err.Cause()
	}

	return user, nil
}

func (h authHandler) authRequired(ctx *gin.Context) {
	user, err := h.parseToken(ctx)
	if err != nil {
		return
	}
	ctx.Set(authusertag, user)
}

func (h authHandler) adminAuthRequired(ctx *gin.Context) {
	user, err := h.parseToken(ctx)
	if err != nil {
		return
	}
	if !user.IsAdmin {
		h.logger.Info("auth failed: not an admin")
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	ctx.Set(authusertag, user)
}

func (h authHandler) signIn(c *gin.Context) {
	var input requests.SignInRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.storage.SignIn(input.Name)
	if err != nil {
		status := http.StatusBadRequest
		if err.IsInternal() {
			status = http.StatusInternalServerError
		}
		c.AbortWithStatusJSON(status, gin.H{"error": err.Cause().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
