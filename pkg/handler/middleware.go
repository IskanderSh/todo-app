package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(userCtx, fmt.Sprintf("%d", userId), 3600*12, "", "", false, true)
}

func getUserId(c *gin.Context) (int, error) {
	id, err := c.Cookie(userCtx)
	if err != nil {
		return 0, errors.New("user id not found")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, errors.New("user id not found")
	}

	return idInt, nil
}
