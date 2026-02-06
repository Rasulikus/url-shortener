package url

import (
	"context"
	"errors"
	"testing"

	"github.com/Rasulikus/url-shortener/internal/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/Rasulikus/url-shortener/internal/service"
	"github.com/Rasulikus/url-shortener/internal/service/url/mocks"
	"github.com/Rasulikus/url-shortener/internal/utils/generator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newService(t *testing.T, repo URLRepository) *Service {
	t.Helper()

	gen, err := generator.NewRandom(generator.DefaultLength)
	require.NoError(t, err)

	s, err := NewService("http://localhost:8080", gen, repo)
	require.NoError(t, err)
	return s
}

func TestService_GetLongURLByAlias(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		repoRet   string
		repoErr   error
		wantURL   string
		wantErrIs error
	}{
		{
			name:      "success",
			alias:     "aa",
			repoRet:   "http://example.com",
			repoErr:   nil,
			wantURL:   "http://example.com",
			wantErrIs: nil,
		},
		{
			name:      "not found error",
			alias:     "bb",
			repoRet:   "",
			repoErr:   repository.ErrNotFound,
			wantURL:   "",
			wantErrIs: service.ErrNotFound,
		},
		{
			name:      "unexpected error",
			alias:     "cc",
			repoRet:   "",
			repoErr:   errors.New("some error"),
			wantURL:   "",
			wantErrIs: service.ErrInternalError,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := new(mocks.MockURLRepository)

			repo.On("GetLongURLByAlias", mock.Anything, tc.alias).
				Return(tc.repoRet, tc.repoErr).
				Once()

			s := newService(t, repo)

			got, err := s.GetLongURLByAlias(context.Background(), tc.alias)

			require.Equal(t, tc.wantURL, got)

			if tc.wantErrIs != nil {
				require.ErrorIs(t, err, tc.wantErrIs)
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestService_CreateOrGet_Success(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := new(mocks.MockURLRepository)

		repo.On("CreateOrGet", mock.Anything, mock.MatchedBy(func(u *model.URL) bool {
			return u != nil &&
				u.ID == 0 &&
				u.LongURL == "http://example.com" &&
				u.Alias != "" &&
				u.CreatedAt.IsZero()
		})).
			Return(&model.URL{
				Alias: "aa",
			}, nil).
			Once()

		s := newService(t, repo)

		got, err := s.CreateOrGet(context.Background(), "http://example.com")

		require.NoError(t, err)
		require.Equal(t, "http://localhost:8080/aa", got)

		repo.AssertExpectations(t)
	})
}

func TestService_CreateOrGet_InvalidInput(t *testing.T) {
	t.Run("invalid input error", func(t *testing.T) {
		repo := new(mocks.MockURLRepository)

		s := newService(t, repo)
		got, err := s.CreateOrGet(context.Background(), "noturl")
		require.Error(t, err)
		require.Zero(t, got)

		repo.AssertNotCalled(t, "CreateOrGet", mock.Anything, mock.Anything)
	})
}

func TestService_CreateOrGet_Conflict(t *testing.T) {
	repo := new(mocks.MockURLRepository)

	repo.On("CreateOrGet", mock.Anything, mock.MatchedBy(func(u *model.URL) bool {
		return u != nil &&
			u.ID == 0 &&
			u.LongURL == "http://example.com" &&
			u.Alias != "" &&
			u.CreatedAt.IsZero()
	})).
		Return(nil, repository.ErrConflict).
		Once()

	s := newService(t, repo)

	gotAlias, err := s.CreateOrGet(context.Background(), "http://example.com")

	require.ErrorIs(t, err, service.ErrConflict)
	require.Zero(t, gotAlias)

	repo.AssertExpectations(t)
}

func TestService_CreateOrGet_UnexpectedRepoError(t *testing.T) {
	repo := new(mocks.MockURLRepository)

	repo.On("CreateOrGet", mock.Anything, mock.MatchedBy(func(u *model.URL) bool {
		return u != nil &&
			u.ID == 0 &&
			u.LongURL == "http://example.com" &&
			u.Alias != "" &&
			u.CreatedAt.IsZero()
	})).
		Return((*model.URL)(nil), errors.New("db down")).
		Once()

	s := newService(t, repo)

	gotAlias, err := s.CreateOrGet(context.Background(), "http://example.com")
	require.ErrorIs(t, err, service.ErrInternalError)
	require.Zero(t, gotAlias)

	repo.AssertExpectations(t)
}
