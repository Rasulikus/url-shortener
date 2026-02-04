package url

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Rasulikus/url-shortener/internal/domain/model"
	"github.com/Rasulikus/url-shortener/internal/repository"
	"github.com/Rasulikus/url-shortener/internal/repository/postgres/test_postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	testPool, err = test_postgres.NewTestPool()
	if err != nil {
		os.Exit(1)
	}
	code := m.Run()
	testPool.Close()
	os.Exit(code)
}

type testSuite struct {
	pool    *pgxpool.Pool
	urlRepo *Repo
	ctx     context.Context
}

func (s *testSuite) ctx2s() (context.Context, context.CancelFunc) {
	return context.WithTimeout(s.ctx, 2*time.Second)
}

func setupTestSuite(t *testing.T) *testSuite {
	t.Helper()

	s := &testSuite{
		pool: testPool,
		ctx:  context.Background(),
	}

	var err error
	s.urlRepo, err = NewRepository(s.pool)
	require.NoError(t, err)

	ctx, cancel := s.ctx2s()
	defer cancel()
	require.NoError(t, test_postgres.TruncateUrls(ctx, s.pool))

	return s
}

func insertURL(t *testing.T, ctx context.Context, pool *pgxpool.Pool, url *model.URL) *model.URL {
	t.Helper()

	u := new(model.URL)
	const q = `
	INSERT INTO urls (long_url, alias)
	VALUES ($1, $2)
	RETURNING id, long_url, alias, created_at;
`

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := pool.QueryRow(ctx, q, url.LongURL, url.Alias).Scan(&u.ID, &u.LongURL, &u.Alias, &u.CreatedAt)
	require.NoError(t, err)

	require.NotZero(t, u.ID)
	require.NotZero(t, u.LongURL)
	require.NotZero(t, u.Alias)
	require.NotZero(t, u.CreatedAt)
	require.WithinDuration(t, time.Now(), u.CreatedAt, 2*time.Second)

	return u
}

func TestRepo_GetByAlias(t *testing.T) {
	s := setupTestSuite(t)

	u := &model.URL{
		LongURL: "https://rkrkrkrk.com",
		Alias:   "aa",
	}
	newU := insertURL(t, s.ctx, s.pool, u)

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
			ctx, candel := s.ctx2s()
			defer candel()
			get, err := s.urlRepo.GetByAlias(ctx, tc.alias)
			tc.check(t, get, err)
		})
	}
}

func TestRepo_GetByAliasCtxCancelErr(t *testing.T) {
	s := setupTestSuite(t)

	u := &model.URL{
		LongURL: "https://rkrkrkrk.com",
		Alias:   "aa",
	}
	newU := insertURL(t, s.ctx, s.pool, u)

	ctx, cancel := s.ctx2s()
	cancel()
	t.Run("ctx err", func(t *testing.T) {
		get, err := s.urlRepo.GetByAlias(ctx, newU.LongURL)
		require.ErrorIs(t, err, context.Canceled)
		require.Nil(t, get)
	})
}

func TestRepo_CreateOrGet(t *testing.T) {
	s := setupTestSuite(t)

	u := &model.URL{
		LongURL: "https://rkrkrkrk.com",
		Alias:   "aa",
	}
	uWithExistingLongURL := &model.URL{
		LongURL: u.LongURL,
		Alias:   "bb",
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := s.ctx2s()
			defer cancel()
			get, err := s.urlRepo.CreateOrGet(ctx, tc.url)
			tc.check(t, get, err)
		})
	}
}
