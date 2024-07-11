// Package configs
//
//	@program:		gin-template
//	@author:		[lliuhuan](https://github.com/lliuhuan)
//	@create:		2024-07-02 21:33
//	@description:	redis配置
package configs

type Redis struct {
	Mode         string   `toml:"mode"`         // Mode Redis模式 cluster|simple|failover
	MasterName   string   `toml:"masterName"`   // MasterName 哨兵主节点名称，sentinel模式下需要配置此项
	Addr         string   `toml:"addr"`         // Addr Simple 实例配置地址
	Addrs        []string `toml:"addrs"`        // Addrs Cluster,Failover实例配置地址
	Pass         string   `toml:"pass"`         // Pass 密码
	DB           int      `toml:"db"`           // DB，默认为0, 一般应用不推荐使用DB分片
	MaxRetries   int      `toml:"maxRetries"`   // MaxRetries 网络相关的错误最大重试次数 默认8次
	PoolSize     int      `toml:"poolSize"`     // PoolSize 集群内每个节点的最大连接池限制 默认每个CPU10个连接
	MinIdleConns int      `toml:"minIdleConns"` // MinIdleConns 最小空闲连接数
	DialTimeout  int      `toml:"dialTimeout"`  // DialTimeout 拨超时时间
	ReadTimeout  int      `toml:"readTimeout"`  // ReadTimeout 读超时 默认3s
	WriteTimeout int      `toml:"writeTimeout"` // WriteTimeout 读超时 默认3s
	//IdleTimeout  int      `toml:"writeTimeout"` // IdleTimeout 连接最大空闲时间，默认60s, 超过该时间，连接会被主动关闭
}
