// Package redis
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 09:44
package redis

import (
	"context"
	"strings"
	"time"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/pkg/errors"
	"github.com/LLiuHuan/gin-template/pkg/timeutil"
	"github.com/LLiuHuan/gin-template/pkg/trace"

	"github.com/redis/go-redis/v9"
)

type Option func(*option)

type Trace = trace.T

type option struct {
	Trace *trace.Trace
	Redis *trace.Redis
}

const (
	// ClusterMode using clusterClient
	ClusterMode string = "cluster"
	// SimpleMode using Client
	SimpleMode string = "simple"
	// FailoverMode using Failover sentinel client
	FailoverMode string = "failover"
)

func newOption() *option {
	return &option{}
}

var _ Repo = (*cacheRepo)(nil)

type Repo interface {
	i()
	Set(key, value string, ttl time.Duration, options ...Option) error
	Get(key string, options ...Option) (string, error)
	TTL(key string) (time.Duration, error)
	Expire(key string, ttl time.Duration) bool
	ExpireAt(key string, ttl time.Time) bool
	Del(key string, options ...Option) bool
	Exists(keys ...string) bool
	Incr(key string, options ...Option) int64
	Close() error
	Version() string
}

type cacheRepo struct {
	client redis.UniversalClient
}

func New() (Repo, error) {
	cfg := configs.Get().Redis

	var client redis.UniversalClient

	switch cfg.Mode {
	case SimpleMode:
		client = NewClient(cfg)
	case FailoverMode:
		client = NewFailoverClient(cfg)
	case ClusterMode:
		client = NewClusterClient(cfg)
	default:
		panic("invalid redis mode")
	}

	// TODO: context
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {

		return nil, errors.Wrap(err, "ping redis err")
	}

	return &cacheRepo{
		client: client,
	}, nil
}

func (c *cacheRepo) i() {}

// NewClient 默认使用Simple模式
func NewClient(cfg configs.Redis) redis.UniversalClient {
	if cfg.Addr == "" {
		//elog.ErrorCtx(emptyCtx, eredis_config.PkgName, elog.FieldName(c.Name), elog.FieldError(fmt.Errorf(`invalid "addr" config, "addr" is empty but with stub mode"`)))
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		//IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Millisecond,
	})
	return client
}

// NewFailoverClient 创建哨兵模式
func NewFailoverClient(cfg configs.Redis) redis.UniversalClient {
	if len(cfg.Addrs) == 0 {
		//elog.ErrorCtx(emptyCtx, `invalid "addrs" config, "addrs" has none addresses but with failover mode`, elog.FieldName(c.Name))
		return nil
	}
	if cfg.MasterName == "" {
		//elog.ErrorCtx(emptyCtx, `invalid "masterName" config, "masterName" is empty but with sentinel mode"`, elog.FieldName(c.Name))
		return nil
	}
	failoverClient := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    cfg.MasterName,
		SentinelAddrs: cfg.Addrs,
		Password:      cfg.Pass,
		DB:            cfg.DB,
		MaxRetries:    cfg.MaxRetries,
		DialTimeout:   time.Duration(cfg.DialTimeout) * time.Millisecond,
		ReadTimeout:   time.Duration(cfg.ReadTimeout) * time.Millisecond,
		WriteTimeout:  time.Duration(cfg.WriteTimeout) * time.Millisecond,
		PoolSize:      cfg.PoolSize,
		MinIdleConns:  cfg.MinIdleConns,
		//IdleTimeout:   time.Duration(cfg.IdleTimeout) * time.Millisecond,
	})
	return failoverClient
}

// NewClusterClient 集群模式
func NewClusterClient(cfg configs.Redis) redis.UniversalClient {
	if len(cfg.Addrs) == 0 {
		//elog.ErrorCtx(emptyCtx, "invalid addrs config, addrs has none addresses but with cluster mode", elog.FieldName(c.Name))
		return nil
	}
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.Addrs,
		MaxRedirects: cfg.MaxRetries,
		Password:     cfg.Pass,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Millisecond,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Millisecond,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		//IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Millisecond,
	})
	return clusterClient
}

// Set set some <key,value> into redis
func (c *cacheRepo) Set(key, value string, ttl time.Duration, options ...Option) error {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "set"
			opt.Redis.Key = key
			opt.Redis.Value = value
			opt.Redis.TTL = ttl.Minutes()
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	// TODO: context
	ctx := context.Background()
	if err := c.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return errors.Wrapf(err, "redis set key: %s err", key)
	}

	return nil
}

// Get some key from redis
func (c *cacheRepo) Get(key string, options ...Option) (string, error) {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "get"
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	// TODO: context
	ctx := context.Background()
	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", errors.Wrapf(err, "redis get key: %s err", key)
	}

	return value, nil
}

// TTL get some key from redis
func (c *cacheRepo) TTL(key string) (time.Duration, error) {
	// TODO: context
	ctx := context.Background()
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return -1, errors.Wrapf(err, "redis get key: %s err", key)
	}

	return ttl, nil
}

// Expire expire some key
func (c *cacheRepo) Expire(key string, ttl time.Duration) bool {
	// TODO: context
	ctx := context.Background()
	ok, _ := c.client.Expire(ctx, key, ttl).Result()
	return ok
}

// ExpireAt expire some key at some time
func (c *cacheRepo) ExpireAt(key string, ttl time.Time) bool {
	// TODO: context
	ctx := context.Background()
	ok, _ := c.client.ExpireAt(ctx, key, ttl).Result()
	return ok
}

func (c *cacheRepo) Exists(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	// TODO: context
	ctx := context.Background()
	value, _ := c.client.Exists(ctx, keys...).Result()
	return value > 0
}

func (c *cacheRepo) Del(key string, options ...Option) bool {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "del"
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if key == "" {
		return true
	}

	// TODO: context
	ctx := context.Background()
	value, _ := c.client.Del(ctx, key).Result()
	return value > 0
}

func (c *cacheRepo) Incr(key string, options ...Option) int64 {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "incr"
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}
	// TODO: context
	ctx := context.Background()
	value, _ := c.client.Incr(ctx, key).Result()
	return value
}

// Close redis client
func (c *cacheRepo) Close() error {
	return c.client.Close()
}

// WithTrace 设置trace信息
func WithTrace(t Trace) Option {
	return func(opt *option) {
		if t != nil {
			opt.Trace = t.(*trace.Trace)
			opt.Redis = new(trace.Redis)
		}
	}
}

// Version redis server version
func (c *cacheRepo) Version() string {
	// TODO: context
	ctx := context.Background()
	server := c.client.Info(ctx, "server").Val()
	spl1 := strings.Split(server, "# Server")
	spl2 := strings.Split(spl1[1], "redis_version:")
	spl3 := strings.Split(spl2[1], "redis_git_sha1:")
	return spl3[0]
}
