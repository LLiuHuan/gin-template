package tool

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/pkg/currency"
	"github.com/LLiuHuan/gin-template/pkg/env"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cast"
)

type projectInfoRequest struct{}

type projectInfoResponse struct {
	MemTotal       string  // 内存总量
	MemUsed        string  // 内存使用量
	MemUsedPercent float64 // 内存使用率

	DiskTotal       string  // 磁盘总量
	DiskUsed        string  // 磁盘使用量
	DiskUsedPercent float64 // 磁盘使用率

	HostOS   string // 操作系统
	HostName string // 主机名

	CpuName        string  // CPU 名称
	CpuCores       int32   // CPU 核数
	CpuUsedPercent float64 // CPU 使用率

	GoPath      string // GoPath
	GoVersion   string // Go 版本
	Goroutine   int    // Goroutine 数量
	ProjectPath string // 项目路径
	Env         string // 运行环境
	Host        string // 主机地址
	GoOS        string // GoOS
	GoArch      string // GoArch

	ProjectVersion  string // 项目版本
	DatabaseVersion string // 数据库版本
	RedisVersion    string // Redis 版本
}

// ProjectInfo 项目基础信息
//
//	@Summary		项目基础信息
//	@Description	项目基础信息
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		projectInfoRequest	true	"请求信息"
//	@Success		200		{object}	projectInfoResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/project/info [get]
func (h *handler) ProjectInfo() core.HandlerFunc {
	type mysqlVersion struct {
		Ver string
	}

	databaseVer := new(mysqlVersion)
	if h.db != nil {
		h.db.GetDB().Raw("SELECT version() as ver").Scan(databaseVer)
	}

	redisVer := ""
	if h.cache != nil {
		redisVer = h.cache.Version()
	}
	return func(ctx core.Context) {
		memInfo, _ := mem.VirtualMemory()
		diskInfo, _ := disk.Usage("/")
		hostInfo, _ := host.Info()
		cpuInfo, _ := cpu.Info()
		cpuPercent, _ := cpu.Percent(time.Microsecond, false)

		obj := new(projectInfoResponse)
		obj.MemTotal = currency.FormatFileSize(cast.ToInt64(memInfo.Total))
		obj.MemUsed = currency.FormatFileSize(cast.ToInt64(memInfo.Used))
		obj.MemUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", memInfo.UsedPercent), 64)

		obj.DiskTotal = currency.FormatFileSize(cast.ToInt64(diskInfo.Total))
		obj.DiskUsed = currency.FormatFileSize(cast.ToInt64(diskInfo.Used))
		obj.DiskUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", diskInfo.UsedPercent), 64)

		obj.HostOS = fmt.Sprintf("%s(%s) %s", hostInfo.Platform, hostInfo.PlatformFamily, hostInfo.PlatformVersion)
		obj.HostName = hostInfo.Hostname

		if len(cpuInfo) > 0 {
			obj.CpuName = cpuInfo[0].ModelName
			obj.CpuCores = cpuInfo[0].Cores
		}

		if len(cpuPercent) > 0 {
			obj.CpuUsedPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", cpuPercent[0]), 64)
		}

		obj.GoPath = runtime.GOROOT()
		obj.GoVersion = runtime.Version()
		obj.Goroutine = runtime.NumGoroutine()
		dir, _ := os.Getwd()
		obj.ProjectPath = strings.Replace(dir, "\\", "/", -1)
		obj.Host = ctx.Host()
		obj.Env = env.Active().Value()
		obj.GoOS = runtime.GOOS
		obj.GoArch = runtime.GOARCH
		obj.ProjectVersion = configs.ProjectVersion
		obj.DatabaseVersion = databaseVer.Ver
		obj.RedisVersion = redisVer

		ctx.Payload(obj)
	}
}
