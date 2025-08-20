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

func newUserBasicCache() *gotest.Cache {
	record1 := &model.UserBasic{}
	record1.ID = 1
	record2 := &model.UserBasic{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewUserBasicCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_userBasicCache_Set(t *testing.T) {
	c := newUserBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserBasic)
	err := c.ICache.(UserBasicCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(UserBasicCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_userBasicCache_Get(t *testing.T) {
	c := newUserBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserBasic)
	err := c.ICache.(UserBasicCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserBasicCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(UserBasicCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_userBasicCache_MultiGet(t *testing.T) {
	c := newUserBasicCache()
	defer c.Close()

	var testData []*model.UserBasic
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserBasic))
	}

	err := c.ICache.(UserBasicCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(UserBasicCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.UserBasic))
	}
}

func Test_userBasicCache_MultiSet(t *testing.T) {
	c := newUserBasicCache()
	defer c.Close()

	var testData []*model.UserBasic
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.UserBasic))
	}

	err := c.ICache.(UserBasicCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userBasicCache_Del(t *testing.T) {
	c := newUserBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserBasic)
	err := c.ICache.(UserBasicCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_userBasicCache_SetCacheWithNotFound(t *testing.T) {
	c := newUserBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.UserBasic)
	err := c.ICache.(UserBasicCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(UserBasicCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewUserBasicCache(t *testing.T) {
	c := NewUserBasicCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewUserBasicCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewUserBasicCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
