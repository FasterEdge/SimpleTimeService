<div align="center">
  <img src="./Logo.png" alt="logo" width="100" />
  <h2>SimpleTimeService</h2>
  <h3>Local and NTP Time Service</h3>
</div>

### 1. Introduction
- Provides JSON endpoints for local UTC time and offset time (multiple DateTime casing variants for compatibility with different clients)
- Supports obtaining more accurate UTC time via NTP (pool.ntp.org), with optional custom offset
- Time format is unified as RFC3339Nano for nanosecond-level parsing (compatible with the `fetchNetworkTime` example)

### 2. Quick Start
```bash
# Start the service, listen on 8080, offset 0 (supports ns/us/ms/s etc. Go durations)
go run main.go --port 8080 --offset 0

# Example: offset of 100 nanoseconds
go run main.go --port 8080 --offset 100ns

# Endpoint examples
curl http://localhost:8080/time
curl http://localhost:8080/offset_time
curl http://localhost:8080/ntp_time
curl http://localhost:8080/offset_npt_time
```

### 3. Startup Arguments
| Argument | Default | Description                         |
|----------|---------|-------------------------------------|
| --port   | 8080    | Listening port                      |
| --offset | 0       | Offset added to the time (supports ns/us/ms/s etc.) |

### 4. API Overview
| Endpoint             | Method | Description                                              |
|----------------------|--------|----------------------------------------------------------|
| `/time`              | GET    | Returns the current local UTC time (DateTime/dateTime/datetime) |
| `/offset_time`       | GET    | Returns the current local UTC time plus `--offset`       |
| `/ntp_time`          | GET    | Gets UTC time from NTP                                   |
| `/offset_npt_time`   | GET    | Gets UTC time from NTP plus `--offset`                   |

> Tip: All time fields are provided in `DateTime`, `dateTime`, and `datetime` casing variants, all in RFC3339Nano format.
