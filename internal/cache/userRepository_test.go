package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-dev-frame/sponge/pkg/gotest"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"cloud-storage/internal/database"
	"cloud-storage/internal/model"
)

func newUserRepositoryCache() *gotest.Cache {
	record1 := &model.UserRepository{}
	record1.ID = 1
	record2 := &model.UserRepository{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewUserRepositoryCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_userRepositoryCache_Set(t *testing.T) {
	c := newUserRepositoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserRepository)
	err := c.ICache.(UserRepositoryCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(UserRepositoryCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_userRepositoryCache_Get(t *testing.T) {
	c := newUserRepositoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserRepository)
	err := c.ICache.(UserRepositoryCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserRepositoryCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(UserRepositoryCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_userRepositoryCache_MultiGet(t *testing.T) {
	c := newUserRepositoryCache()
	defer c.Close()

	var testData []*model.UserRepository
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserRepository))
	}

	err := c.ICache.(UserRepositoryCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserRepositoryCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.UserRepository))
	}
}

func Test_userRepositoryCache_MultiSet(t *testing.T) {
	c := newUserRepositoryCache()
	defer c.Close()

	var testData []*model.UserRepository
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserRepository))
	}

	err := c.ICache.(UserRepositoryCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userRepositoryCache_Del(t *testing.T) {
	c := newUserRepositoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserRepository)
	err := c.ICache.(UserRepositoryCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userRepositoryCache_SetCacheWithNotFound(t *testing.T) {
	c := newUserRepositoryCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserRepository)
	err := c.ICache.(UserRepositoryCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(UserRepositoryCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewUserRepositoryCache(t *testing.T) {
	c := NewUserRepositoryCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewUserRepositoryCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewUserRepositoryCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
