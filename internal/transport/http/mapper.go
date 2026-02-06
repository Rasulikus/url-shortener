package http

import (
	"errors"
	"net/http"

	"github.com/Rasulikus/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
}

func ErrorToHttp(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: "invalid input"})
	case errors.Is(err, service.ErrInvalidInput):
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: "invalid input"})
	case errors.Is(err, service.ErrNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, errorResponse{Error: "not found"})
	case errors.Is(err, service.ErrAliasCollision):
		c.AbortWithStatusJSON(http.StatusConflict, errorResponse{Error: "alias collision"})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
