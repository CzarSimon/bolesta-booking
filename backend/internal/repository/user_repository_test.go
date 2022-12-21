package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/httputil/id"
	"github.com/CzarSimon/httputil/testutil"
	"github.com/CzarSimon/httputil/timeutil"
	"github.com/stretchr/testify/assert"
)

func Test_userRepo_Save(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	u1 := models.User{
		ID:        id.New(),
		Name:      "Some Name",
		Email:     "mail@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	var rowCount int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_account").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(0, rowCount)

	err = repo.Save(ctx, u1)
	assert.NoError(err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_account").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(1, rowCount)

	err = repo.Save(ctx, u1)
	assert.Error(err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_account").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(1, rowCount)

	u2 := models.User{
		ID:        id.New(),
		Name:      "Some Other",
		Email:     "other@mail.com",
		Password:  u1.Password,
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	err = repo.Save(ctx, u2)
	assert.Error(err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_account").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(1, rowCount)

	u2.Password = id.New()
	u2.Salt = u1.Salt

	err = repo.Save(ctx, u2)
	assert.Error(err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_account").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(1, rowCount)

	u2.Salt = id.New()

	err = repo.Save(ctx, u2)
	assert.NoError(err)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM user_account").Scan(&rowCount)
	assert.NoError(err)
	assert.Equal(2, rowCount)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	u3 := models.User{
		ID:        id.New(),
		Name:      "Cancel Name",
		Email:     "cancel@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	err = repo.Save(ctx, u3)
	assert.Error(err)
	assert.True(errors.Is(err, context.Canceled))
}

func Test_userRepo_Find(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	u1 := models.User{
		ID:        id.New(),
		Name:      "Some Name",
		Email:     "mail@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	user, found, err := repo.Find(ctx, u1.ID)
	assert.NoError(err)
	assert.False(found)
	assert.Empty(user)

	err = repo.Save(ctx, u1)
	assert.NoError(err)

	user, found, err = repo.Find(ctx, u1.ID)
	assert.NoError(err)
	assert.True(found)
	assert.Equal(u1, user)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, _, err = repo.Find(ctx, u1.ID)
	assert.Error(err)
	assert.True(errors.Is(err, context.Canceled))
}

func Test_userRepo_FindAll(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	repo := repository.NewUserRepository(db)
	ctx := context.Background()

	u1 := models.User{
		ID:        id.New(),
		Name:      "Some Name",
		Email:     "mail@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}
	u2 := models.User{
		ID:        id.New(),
		Name:      "Some Other Name",
		Email:     "other@mail.com",
		Password:  id.New(),
		Salt:      id.New(),
		CreatedAt: timeutil.Now(),
		UpdatedAt: timeutil.Now(),
	}

	users, err := repo.FindAll(ctx)
	assert.NoError(err)
	assert.Len(users, 0)

	err = repo.Save(ctx, u1)
	assert.NoError(err)
	err = repo.Save(ctx, u2)
	assert.NoError(err)

	users, err = repo.FindAll(ctx)
	assert.NoError(err)
	assert.Len(users, 2)

	um := make(map[string]models.User)
	for _, u := range users {
		um[u.ID] = u
	}

	u, ok := um[u1.ID]
	assert.True(ok)
	assert.Equal(u1, u)

	u, ok = um[u2.ID]
	assert.True(ok)
	assert.Equal(u2, u)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, err = repo.FindAll(ctx)
	assert.Error(err)
	assert.True(errors.Is(err, context.Canceled))
}
