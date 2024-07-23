// Package captcha
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-23 23:12
package captcha

import (
	"context"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"time"
)

type RedisStore struct {
	expiration time.Duration
	preKey     string
	context    context.Context
	cache      redis.Repo
	logger     *zap.Logger
}

func NewDefaultRedisStore(cache redis.Repo, logger *zap.Logger) *RedisStore {
	return &RedisStore{
		expiration: time.Second * 180,
		preKey:     "CAPTCHA_",
		context:    context.TODO(),
		cache:      cache,
		logger:     logger,
	}
}

func (rs *RedisStore) UseWithCtx(ctx context.Context) base64Captcha.Store {
	rs.context = ctx
	return rs
}

func (rs *RedisStore) Set(id string, value string) error {
	err := rs.cache.Set(rs.preKey+id, value, rs.expiration)
	if err != nil {
		rs.logger.Error("RedisStoreSetError!", zap.Error(err))
		return err
	}
	return nil
}

func (rs *RedisStore) Get(key string, clear bool) string {
	val, err := rs.cache.Get(key)
	if err != nil {
		rs.logger.Error("RedisStoreGetError!", zap.Error(err))
		return ""
	}
	if clear {
		if b := rs.cache.Del(key); !b {
			rs.logger.Error("RedisStoreClearError!", zap.Error(err))
			return ""
		}
	}
	return val
}

func (rs *RedisStore) Verify(id, answer string, clear bool) bool {
	key := rs.preKey + id
	v := rs.Get(key, clear)
	return v == answer
}
