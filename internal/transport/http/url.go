package http

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type URLService interface {
	CreateOrGet(ctx context.Context, longURL string) (string, error)
	GetLongURLByAlias(ctx context.Context, alias string) (string, error)
}

var ErrInvalidInput = errors.New("invalid input")

type URLHandler struct {
	s URLService
}

func NewURLHandler(s URLService) *URLHandler {
	return &URLHandler{
		s: s,
	}
}

type CreateUrlRequest struct {
	LongURL string `json:"long_url" binding:"required"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

func (h *URLHandler) Create(c *gin.Context) {
	var req CreateUrlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorToHttp(c, ErrInvalidInput)
		return
	}

	sUrl, err := h.s.CreateOrGet(c.Request.Context(), req.LongURL)
	if err != nil {
		ErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, CreateUrlResponse{ShortUrl: sUrl})
}

type getURLResponse struct {
	LongURL string `json:"long_url"`
}

func (h *URLHandler) GetLongURLByAlias(c *gin.Context) {
	alias := strings.TrimSpace(c.Param("alias"))
	if alias == "" {
		ErrorToHttp(c, ErrInvalidInput)
		return
	}

	url, err := h.s.GetLongURLByAlias(c.Request.Context(), alias)
	if err != nil {
		ErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, getURLResponse{
		LongURL: url,
	})
}

func (h *URLHandler) Redirect(c *gin.Context) {
	alias := strings.TrimSpace(c.Param("alias"))
	if alias == "" {
		ErrorToHttp(c, ErrInvalidInput)
		return
	}

	longURL, err := h.s.GetLongURLByAlias(c.Request.Context(), alias)
	if err != nil {
		ErrorToHttp(c, err)
		return
	}

	c.Redirect(http.StatusFound, longURL)
}
