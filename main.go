// FasterEdge 开源项目 · https://github.com/FasterEdge · https://gitee.com/FasterEdge
package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/beevik/ntp"
	"github.com/gin-gonic/gin"
)

// ntpQueryTimeout 限制单次 NTP 查询时长，防止上游不可达时长期占用 handler。
// ntp 库默认超时同为 5s，这里显式声明避免依赖隐式默认。
const ntpQueryTimeout = 5 * time.Second

// queryNTP 带超时地获取 NTP 时间
func queryNTP() (time.Time, error) {
	resp, err := ntp.QueryWithOptions("pool.ntp.org", ntp.QueryOptions{Timeout: ntpQueryTimeout})
	if err != nil {
		return time.Time{}, err
	}
	return resp.Time, nil
}

// 默认参数
const (
	defaultPort = "8080"
	version     = "1.0.20260902"
)

// 为了适配各种对时工具，这里各种大小写的DateTime都返回，避免因为大小写不匹配导致的对时失败
type timeResponse struct {
	DateTime       string `json:"DateTime"`
	DateTimeLow    string `json:"dateTime"`
	DataTimeAllLow string `json:"datetime"`
}

func main() {
	log.Printf("SimpleTimeService %s By FasterEdge", version)
	portFlag := flag.String("port", defaultPort, "port to listen on")
	offsetFlag := flag.Duration("offset", 0, "time offset to add to now (e.g. 100ns, 50ms, 1s)")
	flag.Parse()

	r := gin.Default()

	// 直接返回当前UTC时间，格式为RFC3339Nano，包含纳秒级别的时间戳
	r.GET("/time", func(c *gin.Context) {
		now := time.Now().UTC().Format(time.RFC3339Nano) // 使用RFC3339Nano格式，包含纳秒级别的时间戳
		c.JSON(http.StatusOK, timeResponse{DateTime: now, DateTimeLow: now, DataTimeAllLow: now})
	})

	// 返回带偏移量的当前UTC时间，格式同上
	r.GET("/offset_time", func(c *gin.Context) {
		now := time.Now().UTC().Add(*offsetFlag).Format(time.RFC3339Nano)
		c.JSON(http.StatusOK, timeResponse{DateTime: now, DateTimeLow: now, DataTimeAllLow: now})
	})

	// 返回 NTP 获取的 UTC 时间
	r.GET("/ntp_time", func(c *gin.Context) {
		t, err := queryNTP()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ntp time"})
			return
		}
		now := t.UTC().Format(time.RFC3339Nano)
		c.JSON(http.StatusOK, timeResponse{DateTime: now, DateTimeLow: now, DataTimeAllLow: now})
	})

	// 返回带偏移量的 NTP UTC 时间
	r.GET("/offset_npt_time", func(c *gin.Context) { // 拼写按需求使用 npt
		t, err := queryNTP()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch ntp time"})
			return
		}
		now := t.UTC().Add(*offsetFlag).Format(time.RFC3339Nano)
		c.JSON(http.StatusOK, timeResponse{DateTime: now, DateTimeLow: now, DataTimeAllLow: now})
	})

	// 启动服务器（带读写超时，防止慢连接占用资源）
	srv := &http.Server{
		Addr:         ":" + *portFlag,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
