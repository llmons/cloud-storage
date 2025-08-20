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
	shareBasicCachePrefixKey = "shareBasic:"
	// ShareBasicExpireTime expire time
	ShareBasicExpireTime = 5 * time.Minute
)

var _ ShareBasicCache = (*shareBasicCache)(nil)

// ShareBasicCache cache interface
type ShareBasicCache interface {
	Set(ctx context.Context, id uint64, data *model.ShareBasic, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.ShareBasic, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.ShareBasic, error)
	MultiSet(ctx context.Context, data []*model.ShareBasic, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// shareBasicCache define a cache struct
type shareBasicCache struct {
	cache cache.Cache
}

// NewShareBasicCache new a cache
func NewShareBasicCache(cacheType *database.CacheType) ShareBasicCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.ShareBasic{}
		})
		return &shareBasicCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.ShareBasic{}
		})
		return &shareBasicCache{cache: c}
	}

	return nil // no cache
}

// GetShareBasicCacheKey cache key
func (c *shareBasicCache) GetShareBasicCacheKey(id uint64) string {
	return shareBasicCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *shareBasicCache) Set(ctx context.Context, id uint64, data *model.ShareBasic, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetShareBasicCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *shareBasicCache) Get(ctx context.Context, id uint64) (*model.ShareBasic, error) {
	var data *model.ShareBasic
	cacheKey := c.GetShareBasicCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *shareBasicCache) MultiSet(ctx context.Context, data []*model.ShareBasic, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetShareBasicCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *shareBasicCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.ShareBasic, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetShareBasicCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.ShareBasic)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.ShareBasic)
	for _, id := range ids {
		val, ok := itemMap[c.GetShareBasicCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *shareBasicCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetShareBasicCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *shareBasicCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetShareBasicCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *shareBasicCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
