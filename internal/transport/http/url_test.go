package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Rasulikus/url-shortener/internal/service"
	"github.com/Rasulikus/url-shortener/internal/transport/http/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRouter(h *URLHandler) *gin.Engine {
	r := gin.New()
	r.POST("/api", h.Create)
	r.GET("/api/:alias", h.GetLongURLByAlias)
	r.GET("/:alias", h.Redirect)
	return r
}
func TestURLHandler_Create_OK(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("CreateOrGet", mock.Anything, "http://example.com").
			Return("http://localhost:8080/aa", nil).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(`{"long_url":"http://example.com"}`))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, `{"short_url":"http://localhost:8080/aa"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_Create_InvalidJSON(t *testing.T) {
	s := mocks.NewMockURLService(t)

	h := NewURLHandler(s)
	r := setupRouter(h)

	cases := []struct {
		name string
		body string
	}{
		{
			"invalid json",
			`{"lon":""}`,
		},
		{
			"invalid json value",
			`{"long_url":""}`,
		},
		{
			"invalid json syntax",
			"{",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(tc.body))
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code)
			require.JSONEq(t, `{"error":"invalid input"}`, w.Body.String())

			s.AssertNotCalled(t, "CreateOrGet")
		})
	}
}

func TestURLHandler_Create_ServiceInvalidInput(t *testing.T) {
	t.Run("invalid input", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("CreateOrGet", mock.Anything, "http://example.com").
			Return("", service.ErrInvalidInput).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(`{"long_url":"http://example.com"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.JSONEq(t, `{"error":"invalid input"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_Create_AliasCollision(t *testing.T) {
	t.Run("alias collision", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("CreateOrGet", mock.Anything, "http://example.com").
			Return("", service.ErrConflict).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(`{"long_url":"http://example.com"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusConflict, w.Code)
		require.JSONEq(t, `{"error":"conflict"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_Create_DefaultError(t *testing.T) {
	t.Run("default error", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("CreateOrGet", mock.Anything, "http://example.com").
			Return("", errors.New("some err")).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodPost, "/api", strings.NewReader(`{"long_url":"http://example.com"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.JSONEq(t, `{"error":"internal server error"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_GetLongURLByAlias_OK(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("GetLongURLByAlias", mock.Anything, "aa").
			Return("http://example.com", nil).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/api/aa", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, `{"long_url":"http://example.com"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_GetByAlias_ServiceNotFound(t *testing.T) {
	t.Run("service not found", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("GetLongURLByAlias", mock.Anything, "aa").
			Return("", service.ErrNotFound).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/api/aa", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.JSONEq(t, `{"error":"not found"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_GetByAlias_DefaultError(t *testing.T) {
	t.Run("default error", func(t *testing.T) {
		s := mocks.NewMockURLService(t)
		s.On("GetLongURLByAlias", mock.Anything, "aa").
			Return("", errors.New("some err")).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/api/aa", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.JSONEq(t, `{"error":"internal server error"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_Redirect_OK(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("GetLongURLByAlias", mock.Anything, "aa").
			Return("http://example.com", nil).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/aa", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusFound, w.Code)
		require.Equal(t, "http://example.com", w.Header().Get("Location"))

		s.AssertExpectations(t)
	})
}

func TestURLHandler_Redirect_ServiceNotFound(t *testing.T) {
	t.Run("service not found", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("GetLongURLByAlias", mock.Anything, "aa").
			Return("", service.ErrNotFound).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/aa", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
		require.JSONEq(t, `{"error":"not found"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}

func TestURLHandler_Redirect_DefaultError(t *testing.T) {
	t.Run("default error", func(t *testing.T) {
		s := mocks.NewMockURLService(t)

		s.On("GetLongURLByAlias", mock.Anything, "aa").
			Return("", errors.New("some err")).
			Once()

		h := NewURLHandler(s)
		r := setupRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/aa", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.JSONEq(t, `{"error":"internal server error"}`, w.Body.String())

		s.AssertExpectations(t)
	})
}
