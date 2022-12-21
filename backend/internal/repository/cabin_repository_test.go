package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/CzarSimon/bolesta-booking/backend/internal/models"
	"github.com/CzarSimon/bolesta-booking/backend/internal/repository"
	"github.com/CzarSimon/httputil/testutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_cabinRepo_Find(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	repo := repository.NewCabinRepository(db)
	ctx := context.Background()

	cabin, found, err := repo.Find(ctx, "a4b4f496-767e-423e-9816-83b71e1cfa89")
	assert.NoError(err)
	assert.True(found)
	assert.Equal("Bölestastugan", cabin.Name)

	cabin, found, err = repo.Find(ctx, "this-id-does-not-exist")
	assert.NoError(err)
	assert.False(found)
	assert.Empty(cabin)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, _, err = repo.Find(ctx, "a4b4f496-767e-423e-9816-83b71e1cfa89")
	assert.Error(err)
	assert.True(errors.Is(err, context.Canceled))
}

func Test_cabinRepo_FindAll(t *testing.T) {
	assert := assert.New(t)
	db := testutil.InMemoryDB(true, "../../resources/db/sqlite")
	repo := repository.NewCabinRepository(db)
	ctx := context.Background()

	cabins, err := repo.FindAll(ctx)
	assert.NoError(err)
	assert.Len(cabins, 3)

	cm := make(map[string]models.Cabin)
	for _, c := range cabins {
		cm[c.ID] = c
	}

	c, ok := cm["a4b4f496-767e-423e-9816-83b71e1cfa89"]
	assert.True(ok)
	assert.Equal("Bölestastugan", c.Name)
	assert.NotEmpty(c.CreatedAt)
	assert.NotEmpty(c.UpdatedAt)

	c, ok = cm["63e71fef-0037-451f-b731-27249c0164d9"]
	assert.True(ok)
	assert.Equal("Gulhuset", c.Name)
	assert.NotEmpty(c.CreatedAt)
	assert.NotEmpty(c.UpdatedAt)

	c, ok = cm["2aa15162-2443-48f1-9b8f-6314f90faf9a"]
	assert.True(ok)
	assert.Equal("Bergebo", c.Name)
	assert.NotEmpty(c.CreatedAt)
	assert.NotEmpty(c.UpdatedAt)

	ctx, cancel := context.WithCancel(ctx)
	cancel()

	_, err = repo.FindAll(ctx)
	assert.Error(err)
	assert.True(errors.Is(err, context.Canceled))
}
