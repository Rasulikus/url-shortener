package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Rasulikus/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

var InvalidInputError = errors.New("invalid input")

type URLHandler struct {
	s service.URLService
}

func NewURLHandler(s service.URLService) *URLHandler {
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
		ErrorToHttp(c, InvalidInputError)
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
	URL string `json:"url"`
}

func (h *URLHandler) GetByAlias(c *gin.Context) {
	alias := strings.TrimSpace(c.Param("alias"))
	if alias == "" {
		ErrorToHttp(c, InvalidInputError)
		return
	}

	url, err := h.s.GetByAlias(c.Request.Context(), alias)
	if err != nil {
		ErrorToHttp(c, err)
		return
	}
	c.JSON(http.StatusOK, getURLResponse{
		URL: url.LongURL,
	})
}
