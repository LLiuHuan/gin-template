package install

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/proposal/tablesqls"
	"github.com/LLiuHuan/gin-template/internal/repository/socket"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type installRequest struct {
	Language  string `form:"language" json:"language"`     // 语言包
	RedisAddr string `form:"redis_addr" json:"redis_addr"` // 连接地址，例如：127.0.0.1:6379
	RedisPass string `form:"redis_pass" json:"redis_pass"` // 连接密码
	RedisDb   int    `form:"redis_db" json:"redis_db"`     // 连接 db
	RedisMode string `form:"redis_mode" json:"redis_mode"` // 连接模式，例如：simple、cluster、sentinel

	IsOpenReadDB int    `form:"isOpenReadDB" json:"isOpenReadDB"` // 是否开启读写分离
	DatabaseHost string `form:"database_host" json:"database_host"`
	DatabasePort int    `form:"database_port" json:"database_port"`
	DatabaseUser string `form:"database_user" json:"database_user"`
	DatabasePass string `form:"database_pass" json:"database_pass"`
	DatabaseName string `form:"database_name" json:"database_name"`
	DatabaseMode string `form:"database_mode" json:"database_mode"`
}

type installResponse struct{}

// Install 安装
//
//	@Summary		安装
//	@Description	安装
//	@Tags			API.install
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		installRequest	true	"请求信息"
//	@Success		200		{object}	installResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/v1/api/install [post]
func (h *handler) Install() core.HandlerFunc {
	installTableList := map[string]map[string]string{
		"authorized": {
			"table_sql":      tablesqls.CreateAuthorizedTableSql(),
			"table_data_sql": tablesqls.CreateAuthorizedTableDataSql(),
		},
		"authorized_api": {
			"table_sql":      tablesqls.CreateAuthorizedAPITableSql(),
			"table_data_sql": tablesqls.CreateAuthorizedAPITableDataSql(),
		},
		"admin": {
			"table_sql":      tablesqls.CreateAdminTableSql(),
			"table_data_sql": tablesqls.CreateAdminTableDataSql(),
		},
		"admin_menu": {
			"table_sql":      tablesqls.CreateAdminMenuTableSql(),
			"table_data_sql": tablesqls.CreateAdminMenuTableDataSql(),
		},
		"menu": {
			"table_sql":      tablesqls.CreateMenuTableSql(),
			"table_data_sql": tablesqls.CreateMenuTableDataSql(),
		},
		"menu_action": {
			"table_sql":      tablesqls.CreateMenuActionTableSql(),
			"table_data_sql": tablesqls.CreateMenuActionTableDataSql(),
		},
		"cron_task": {
			"table_sql":      tablesqls.CreateCronTaskTableSql(),
			"table_data_sql": "",
		},
	}

	return func(ctx core.Context) {
		server, err := socket.New(h.logger, h.db, h.cache, ctx.ResponseWriter(), ctx.Request(), nil)
		req := new(installRequest)

		// 获取传过来的参数
		//接收消息
		_, message, err := server.ReadMessage()
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			server.OnClose()
			return
		}

		if err = json.Unmarshal(message, &req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			server.OnClose()
			return
		}

		// region 验证 version
		versionStr := runtime.Version()
		version := cast.ToFloat32(versionStr[2:6])
		if version < configs.MinGoVersion {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.GoVersionError,
				code.Text(code.GoVersionError)),
			)
			return
		}
		// endregion

		// region 验证 Redis 配置
		cfg := configs.Get()
		redisClient := redis.NewClient(&redis.Options{
			Addr:         req.RedisAddr,
			Password:     req.RedisPass,
			DB:           cast.ToInt(req.RedisDb),
			MaxRetries:   cfg.Redis.MaxRetries,
			PoolSize:     cfg.Redis.PoolSize,
			MinIdleConns: cfg.Redis.MinIdleConns,
		})

		if err := redisClient.Ping(ctx.GetCtx()).Err(); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.RedisConnectError,
				code.Text(code.RedisConnectError)).WithError(err),
			)
			return
		}

		defer redisClient.Close()

		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.SocketConnectError,
				code.Text(code.SocketConnectError)).WithError(err),
			)
			return
		}
		err = server.OnSend([]byte("已检测 Redis 配置可用。"))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.SocketSendError,
				code.Text(code.SocketSendError)).WithError(err),
			)
		}
		// endregion

		// region 验证 MySQL 配置
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
			req.DatabaseUser,
			req.DatabasePass,
			req.DatabaseHost,
			req.DatabasePort,
			req.DatabaseName,
			true,
			"Local")

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			//Logger: logger.Default.LogMode(logger.Info), // 日志配置
		})

		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.MySQLConnectError,
				code.Text(code.MySQLConnectError)).WithError(err),
			)
			return
		}

		db.Set("gorm:table_options", "CHARSET=utf8mb4")

		dbClient, _ := db.DB()
		defer dbClient.Close()

		err = server.OnSend([]byte("已检测 MySQL 配置可用。"))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.SocketSendError,
				code.Text(code.SocketSendError)).WithError(err),
			)
		}
		// endregion

		// region 写入配置文件
		config := configs.GetContainer()
		config.Set("project.local", req.Language)

		config.Set("redis.addr", req.RedisAddr)
		config.Set("redis.pass", req.RedisPass)
		config.Set("redis.db", req.RedisDb)
		config.Set("redis.mode", req.RedisMode)

		config.Set("database.mode", req.DatabaseMode)
		dbMode := fmt.Sprintf("database.%s", req.DatabaseMode)
		dbRead := dbMode + ".read"
		dbWrite := dbMode + ".write"
		dbBase := dbMode + ".base"

		config.Set(dbMode+"isOpenReadDB", req.IsOpenReadDB)

		config.Set(dbRead+".host", req.DatabaseHost)
		config.Set(dbRead+".port", req.DatabasePort)
		config.Set(dbRead+".user", req.DatabaseUser)
		config.Set(dbRead+".pass", req.DatabasePass)
		config.Set(dbRead+".name", req.DatabaseName)

		config.Set(dbWrite+".host", req.DatabaseHost)
		config.Set(dbWrite+".port", req.DatabasePort)
		config.Set(dbWrite+".user", req.DatabaseUser)
		config.Set(dbWrite+".pass", req.DatabasePass)
		config.Set(dbWrite+".name", req.DatabaseName)

		config.Set(dbBase+".host", req.DatabaseHost)
		config.Set(dbBase+".port", req.DatabasePort)
		config.Set(dbBase+".user", req.DatabaseUser)
		config.Set(dbBase+".pass", req.DatabasePass)
		config.Set(dbBase+".name", req.DatabaseName)
		err = config.WriteConfig()
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.WriteConfigError,
				code.Text(code.WriteConfigError)).WithError(err),
			)
			return
		}

		err = server.OnSend([]byte("语言包 " + req.Language + " 配置成功。"))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.WriteConfigError,
				code.Text(code.WriteConfigError)).WithError(err),
			)
		}
		err = server.OnSend([]byte("配置项 Redis、MySQL 配置成功。"))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.WriteConfigError,
				code.Text(code.WriteConfigError)).WithError(err),
			)
		}
		// endregion

		// region 初始化表结构 + 默认数据
		for k, v := range installTableList {
			if v["table_sql"] != "" {
				// region 初始化表结构
				if err = db.Exec(v["table_sql"]).Error; err != nil {
					ctx.AbortWithError(core.Error(
						http.StatusBadRequest,
						code.MySQLExecError,
						code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
					)
					return
				}

				err = server.OnSend([]byte("初始化 MySQL 数据表：" + k + " 成功。"))
				if err != nil {
					ctx.AbortWithError(core.Error(
						http.StatusInternalServerError,
						code.MySQLExecError,
						code.Text(code.MySQLExecError)).WithError(err),
					)
				}

				// 有个延时状态
				time.Sleep(300 * time.Microsecond)
				// endregion

				// region 初始化默认数据
				if v["table_data_sql"] != "" {
					if err = db.Exec(v["table_data_sql"]).Error; err != nil {
						ctx.AbortWithError(core.Error(
							http.StatusBadRequest,
							code.MySQLExecError,
							code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
						)
						return
					}

					err = server.OnSend([]byte("初始化 MySQL 数据表：" + k + " 默认数据成功。"))

					// 有个延时状态
					time.Sleep(300 * time.Microsecond)
					if err != nil {
						ctx.AbortWithError(core.Error(
							http.StatusInternalServerError,
							code.MySQLExecError,
							code.Text(code.MySQLExecError)).WithError(err),
						)
					}
				}
				// endregion
			}
		}
		// endregion

		// region 生成 install 完成标识
		f, err := os.Create(configs.ProjectInstallMark)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.MySQLExecError,
				code.Text(code.MySQLExecError)+" "+err.Error()).WithError(err),
			)
			return
		}
		defer f.Close()
		// endregion

		// 有个延时状态
		time.Sleep(time.Second)

		err = server.OnSend([]byte("初始化成功。"))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.MySQLExecError,
				code.Text(code.MySQLExecError)).WithError(err),
			)
		}

		server.OnClose()
	}
}
