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

func newShareBasicCache() *gotest.Cache {
	record1 := &model.ShareBasic{}
	record1.ID = 1
	record2 := &model.ShareBasic{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewShareBasicCache(&database.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_shareBasicCache_Set(t *testing.T) {
	c := newShareBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ShareBasic)
	err := c.ICache.(ShareBasicCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(ShareBasicCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_shareBasicCache_Get(t *testing.T) {
	c := newShareBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ShareBasic)
	err := c.ICache.(ShareBasicCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ShareBasicCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(ShareBasicCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_shareBasicCache_MultiGet(t *testing.T) {
	c := newShareBasicCache()
	defer c.Close()

	var testData []*model.ShareBasic
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.ShareBasic))
	}

	err := c.ICache.(ShareBasicCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(ShareBasicCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[utils.StrToUint64(k)], v.(*model.ShareBasic))
	}
}

func Test_shareBasicCache_MultiSet(t *testing.T) {
	c := newShareBasicCache()
	defer c.Close()

	var testData []*model.ShareBasic
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.ShareBasic))
	}

	err := c.ICache.(ShareBasicCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_shareBasicCache_Del(t *testing.T) {
	c := newShareBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ShareBasic)
	err := c.ICache.(ShareBasicCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_shareBasicCache_SetCacheWithNotFound(t *testing.T) {
	c := newShareBasicCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.ShareBasic)
	err := c.ICache.(ShareBasicCache).SetPlaceholder(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	b := c.ICache.(ShareBasicCache).IsPlaceholderErr(err)
	t.Log(b)
}

func TestNewShareBasicCache(t *testing.T) {
	c := NewShareBasicCache(&database.CacheType{
		CType: "",
	})
	assert.Nil(t, c)
	c = NewShareBasicCache(&database.CacheType{
		CType: "memory",
	})
	assert.NotNil(t, c)
	c = NewShareBasicCache(&database.CacheType{
		CType: "redis",
	})
	assert.NotNil(t, c)
}
