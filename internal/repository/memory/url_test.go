package memory

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Rasulikus/url-shortener/internal/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var repo *Repo

func TestMain(m *testing.M) {
	var err error
	repo, err = NewRepository(New())
	if err != nil {
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func TestRepo_GetByAlias(t *testing.T) {
	ctx := context.Background()

	u := &model.URL{
		LongURL: "https://rkrkrkrk.com",
		Alias:   "aa",
	}

	newU, err := repo.CreateOrGet(ctx, u)
	require.NoError(t, err)

	cases := []struct {
		name  string
		alias string
		check func(t *testing.T, get *model.URL, err error)
	}{
		{
			name:  "found",
			alias: "aa",
			check: func(t *testing.T, get *model.URL, err error) {
				require.NoError(t, err)
				assert.Equal(t, newU.ID, get.ID)
				assert.Equal(t, newU.LongURL, get.LongURL)
				assert.Equal(t, newU.Alias, get.Alias)
				assert.True(t, newU.CreatedAt.Equal(get.CreatedAt))
			},
		},
		{
			name:  "not found",
			alias: "bb",
			check: func(t *testing.T, get *model.URL, err error) {
				require.Nil(t, get)
				require.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			get, err := repo.GetByAlias(ctx, tc.alias)
			tc.check(t, get, err)
		})
	}
}

func TestRepo_GetLongURLByAlias(t *testing.T) {
	ctx := context.Background()

	u := &model.URL{
		LongURL: "https://rkrkrkrk.com",
		Alias:   "aa",
	}

	newU, err := repo.CreateOrGet(ctx, u)
	require.NoError(t, err)

	cases := []struct {
		name  string
		alias string
		check func(t *testing.T, get string, err error)
	}{
		{
			name:  "found",
			alias: "aa",
			check: func(t *testing.T, get string, err error) {
				require.NoError(t, err)
				assert.Equal(t, newU.LongURL, get)
			},
		},
		{
			name:  "not found",
			alias: "bb",
			check: func(t *testing.T, get string, err error) {
				require.Zero(t, get)
				require.ErrorIs(t, err, repository.ErrNotFound)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			get, err := repo.GetLongURLByAlias(ctx, tc.alias)
			tc.check(t, get, err)
		})
	}
}

func TestRepo_CreateOrGet(t *testing.T) {
	ctx := context.Background()

	u := &model.URL{
		LongURL: "https://rkrkrkrk.com",
		Alias:   "aa",
	}
	uWithExistingLongURL := &model.URL{
		LongURL: u.LongURL,
		Alias:   "bb",
	}
	uWithExistingAlias := &model.URL{
		LongURL: "http://aaa.com",
		Alias:   "aa",
	}
	cases := []struct {
		name  string
		url   *model.URL
		check func(t *testing.T, get *model.URL, err error)
	}{
		{
			name: "created",
			url:  u,
			check: func(t *testing.T, get *model.URL, err error) {
				require.NoError(t, err)
				assert.NotZero(t, get.ID)
				assert.Equal(t, u.LongURL, get.LongURL)
				assert.Equal(t, u.Alias, get.Alias)
				assert.Equal(t, u.CreatedAt, get.CreatedAt)
				assert.WithinDuration(t, time.Now(), u.CreatedAt, 2*time.Second)
			},
		},
		{
			name: "get",
			url:  uWithExistingLongURL,
			check: func(t *testing.T, get *model.URL, err error) {
				require.NoError(t, err)
				assert.Equal(t, u, uWithExistingLongURL)
			},
		},
		{
			name: "err conflict",
			url:  uWithExistingAlias,
			check: func(t *testing.T, get *model.URL, err error) {
				require.Nil(t, get)
				require.ErrorIs(t, err, repository.ErrConflict)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			get, err := repo.CreateOrGet(ctx, tc.url)
			tc.check(t, get, err)
		})
	}
}
