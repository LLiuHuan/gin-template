// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:32
// @description: 数据库配置
package configs

type DataBase struct {
	Mode string `toml:"mode"` // Mode 数据库模式 database|sqlite3|postgres|sqlserver
	// MySql
	MySql DataBaseConf `toml:"mysql"`
	// SqlServer
	SqlServer DataBaseConf `toml:"sqlserver"`
	// PostgreSql
	PostgreSql DataBaseConf `toml:"postgresql"`
}

type DataBaseConf struct {
	IsOpenReadDB int                `toml:"isOpenReadDB"` // 是否开启读写分离，如果不开启的话就默认忽略Read库的配置
	Read         DataBaseConfDetail `toml:"read"`
	Write        DataBaseConfDetail `toml:"write"`
	Base         DataBaseConfBase   `toml:"base"`
}

type DataBaseConfDetail struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
	DataBase string `toml:"dataBase"`
	Charset  string `toml:"charset"`
	Prefix   string `toml:"prefix"` // 表前缀
	// 下面是PGSql的
	SSLMode     string `toml:"sslMode"`     // SSLMode 是否开启SSL
	SSLRootCert string `toml:"sslRootCert"` // SSLRootCert 根证书路径
	SSLCert     string `toml:"sslCert"`     // SSLCert 证书路径
	TimeZone    string `toml:"timeZone"`    // TimeZone 时区
	//下面是sqlite3的
	Path string `toml:"path"` // Path 数据库路径
}

type DataBaseConfBase struct {
	MaxIdleConn     int `toml:"maxIdleConn"`
	MaxOpenConn     int `toml:"maxOpenConn"`
	ConnMaxLifeTime int `toml:"connMaxLifeTime"`
}

func (d *DataBase) GetDataBaseConfig() DataBaseConf {
	switch d.Mode {
	case "mysql":
		return d.MySql
	case "sqlserver":
		return d.SqlServer
	case "postgresql":
		return d.PostgreSql
	default:
		return d.MySql
	}
}
