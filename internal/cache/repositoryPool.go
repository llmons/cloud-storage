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
	repositoryPoolCachePrefixKey = "repositoryPool:"
	// RepositoryPoolExpireTime expire time
	RepositoryPoolExpireTime = 5 * time.Minute
)

var _ RepositoryPoolCache = (*repositoryPoolCache)(nil)

// RepositoryPoolCache cache interface
type RepositoryPoolCache interface {
	Set(ctx context.Context, id uint64, data *model.RepositoryPool, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.RepositoryPool, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.RepositoryPool, error)
	MultiSet(ctx context.Context, data []*model.RepositoryPool, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// repositoryPoolCache define a cache struct
type repositoryPoolCache struct {
	cache cache.Cache
}

// NewRepositoryPoolCache new a cache
func NewRepositoryPoolCache(cacheType *database.CacheType) RepositoryPoolCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.RepositoryPool{}
		})
		return &repositoryPoolCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.RepositoryPool{}
		})
		return &repositoryPoolCache{cache: c}
	}

	return nil // no cache
}

// GetRepositoryPoolCacheKey cache key
func (c *repositoryPoolCache) GetRepositoryPoolCacheKey(id uint64) string {
	return repositoryPoolCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *repositoryPoolCache) Set(ctx context.Context, id uint64, data *model.RepositoryPool, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetRepositoryPoolCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *repositoryPoolCache) Get(ctx context.Context, id uint64) (*model.RepositoryPool, error) {
	var data *model.RepositoryPool
	cacheKey := c.GetRepositoryPoolCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *repositoryPoolCache) MultiSet(ctx context.Context, data []*model.RepositoryPool, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetRepositoryPoolCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *repositoryPoolCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.RepositoryPool, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetRepositoryPoolCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.RepositoryPool)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.RepositoryPool)
	for _, id := range ids {
		val, ok := itemMap[c.GetRepositoryPoolCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *repositoryPoolCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetRepositoryPoolCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *repositoryPoolCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetRepositoryPoolCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *repositoryPoolCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
