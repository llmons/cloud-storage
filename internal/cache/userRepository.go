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
	userRepositoryCachePrefixKey = "userRepository:"
	// UserRepositoryExpireTime expire time
	UserRepositoryExpireTime = 5 * time.Minute
)

var _ UserRepositoryCache = (*userRepositoryCache)(nil)

// UserRepositoryCache cache interface
type UserRepositoryCache interface {
	Set(ctx context.Context, id uint64, data *model.UserRepository, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.UserRepository, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserRepository, error)
	MultiSet(ctx context.Context, data []*model.UserRepository, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// userRepositoryCache define a cache struct
type userRepositoryCache struct {
	cache cache.Cache
}

// NewUserRepositoryCache new a cache
func NewUserRepositoryCache(cacheType *database.CacheType) UserRepositoryCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserRepository{}
		})
		return &userRepositoryCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserRepository{}
		})
		return &userRepositoryCache{cache: c}
	}

	return nil // no cache
}

// GetUserRepositoryCacheKey cache key
func (c *userRepositoryCache) GetUserRepositoryCacheKey(id uint64) string {
	return userRepositoryCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *userRepositoryCache) Set(ctx context.Context, id uint64, data *model.UserRepository, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserRepositoryCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *userRepositoryCache) Get(ctx context.Context, id uint64) (*model.UserRepository, error) {
	var data *model.UserRepository
	cacheKey := c.GetUserRepositoryCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *userRepositoryCache) MultiSet(ctx context.Context, data []*model.UserRepository, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserRepositoryCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *userRepositoryCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserRepository, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserRepositoryCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UserRepository)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.UserRepository)
	for _, id := range ids {
		val, ok := itemMap[c.GetUserRepositoryCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *userRepositoryCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserRepositoryCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *userRepositoryCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserRepositoryCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *userRepositoryCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
