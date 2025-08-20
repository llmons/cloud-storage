package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"cloud-storage/internal/database"
	"cloud-storage/internal/model"
)

const (
	// cache prefix key, must end with a colon
	userBasicCachePrefixKey = "userBasic:"
	// UserBasicExpireTime expire time
	UserBasicExpireTime = 5 * time.Minute
)

var _ UserBasicCache = (*userBasicCache)(nil)

// UserBasicCache cache interface
type UserBasicCache interface {
	Set(ctx context.Context, id uint64, data *model.UserBasic, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.UserBasic, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserBasic, error)
	MultiSet(ctx context.Context, data []*model.UserBasic, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// userBasicCache define a cache struct
type userBasicCache struct {
	cache cache.Cache
}

// NewUserBasicCache new a cache
func NewUserBasicCache(cacheType *database.CacheType) UserBasicCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserBasic{}
		})
		return &userBasicCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserBasic{}
		})
		return &userBasicCache{cache: c}
	}

	return nil // no cache
}

// GetUserBasicCacheKey cache key
func (c *userBasicCache) GetUserBasicCacheKey(id uint64) string {
	return userBasicCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *userBasicCache) Set(ctx context.Context, id uint64, data *model.UserBasic, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserBasicCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *userBasicCache) Get(ctx context.Context, id uint64) (*model.UserBasic, error) {
	var data *model.UserBasic
	cacheKey := c.GetUserBasicCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *userBasicCache) MultiSet(ctx context.Context, data []*model.UserBasic, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserBasicCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *userBasicCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserBasic, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserBasicCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UserBasic)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.UserBasic)
	for _, id := range ids {
		val, ok := itemMap[c.GetUserBasicCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *userBasicCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserBasicCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *userBasicCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserBasicCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *userBasicCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
