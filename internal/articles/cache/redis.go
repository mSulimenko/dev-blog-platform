package cache

import (
	"context"
	"github.com/go-redis/cache/v9"
	"github.com/mSulimenko/dev-blog-platform/internal/articles/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
	cache  *cache.Cache
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	mycache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &RedisCache{
		client: client,
		cache:  mycache,
	}
}

func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) GetLatestArticles(ctx context.Context) ([]*models.Article, error) {
	var articles []*models.Article
	err := r.cache.Get(ctx, "latest_articles:10", &articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *RedisCache) SetLatestArticles(ctx context.Context, articles []*models.Article) error {
	return r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   "latest_articles:10",
		Value: articles,
		TTL:   time.Minute * 5, // 5 минут TTL
	})
}

func (r *RedisCache) InvalidateLatestArticles(ctx context.Context) error {
	return r.cache.Delete(ctx, "latest_articles:10")
}
