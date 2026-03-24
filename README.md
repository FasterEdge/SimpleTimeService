<div align="center">
  <img src="./Logo.png" alt="logo" width="96" />
  <h3>SimpleTimeService · 本地/NTP时间服务</h3>
</div>

### 项目简介
- 提供本地UTC时间与可偏移时间的JSON接口（多种DateTime大小写，兼容不同客户端）
- 支持基于NTP（pool.ntp.org）获取更准确的UTC时间，亦可叠加自定义偏移
- 时间格式统一为RFC3339Nano，便于纳秒级解析（兼容示例fetchNetworkTime）

### 快速开始
```bash
# 启动服务，监听 8080，偏移 0（可用 ns/us/ms/s 等 Go duration）
go run main.go --port 8080 --offset 0

# 示例：偏移 100 纳秒
go run main.go --port 8080 --offset 100ns

# 访问接口示例
curl http://localhost:8080/time
curl http://localhost:8080/offset_time
curl http://localhost:8080/ntp_time
curl http://localhost:8080/offset_npt_time
```

### 启动参数
| 参数      | 默认值 | 说明                               |
|-----------|--------|------------------------------------|
| --port    | 8080   | 监听端口                           |
| --offset  | 0      | 给时间叠加的偏移量（支持 ns/us/ms/s 等） |

### API 一览
| 接口              | 方法 | 说明                                               |
|-------------------|------|----------------------------------------------------|
| `/time`           | GET  | 返回当前本地 UTC 时间（DateTime/dateTime/datetime） |
| `/offset_time`    | GET  | 返回当前本地 UTC 时间叠加 `--offset`               |
| `/ntp_time`       | GET  | 从 NTP 获取 UTC 时间                               |
| `/offset_npt_time`| GET  | 从 NTP 获取 UTC 时间并叠加 `--offset`              |

> 提示：所有时间字段同时提供 `DateTime`、`dateTime`、`datetime` 三种大小写，格式均为 RFC3339Nano。
