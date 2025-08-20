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

func newRepositoryPoolCache() *gotest.Cache {
	record1 := &model.RepositoryPool{}
	record1.ID = 1
	record2 := &model.RepositoryPool{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewRepositoryPoolCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_repositoryPoolCache_Set(t *testing.T) {
	c := newRepositoryPoolCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RepositoryPool)
	err := c.ICache.(RepositoryPoolCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(RepositoryPoolCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_repositoryPoolCache_Get(t *testing.T) {
	c := newRepositoryPoolCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RepositoryPool)
	err := c.ICache.(RepositoryPoolCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RepositoryPoolCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(RepositoryPoolCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_repositoryPoolCache_MultiGet(t *testing.T) {
	c := newRepositoryPoolCache()
	defer c.Close()

	var testData []*model.RepositoryPool
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.RepositoryPool))
	}

	err := c.ICache.(RepositoryPoolCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(RepositoryPoolCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.RepositoryPool))
	}
}

func Test_repositoryPoolCache_MultiSet(t *testing.T) {
	c := newRepositoryPoolCache()
	defer c.Close()

	var testData []*model.RepositoryPool
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.RepositoryPool))
	}

	err := c.ICache.(RepositoryPoolCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_repositoryPoolCache_Del(t *testing.T) {
	c := newRepositoryPoolCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RepositoryPool)
	err := c.ICache.(RepositoryPoolCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_repositoryPoolCache_SetCacheWithNotFound(t *testing.T) {
	c := newRepositoryPoolCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.RepositoryPool)
	err := c.ICache.(RepositoryPoolCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(RepositoryPoolCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewRepositoryPoolCache(t *testing.T) {
	c := NewRepositoryPoolCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewRepositoryPoolCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewRepositoryPoolCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
