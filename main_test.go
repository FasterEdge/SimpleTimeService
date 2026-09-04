package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// fetch 发起请求并解析 timeResponse。
func fetch(t *testing.T, h http.Handler, path string, wantStatus int) timeResponse {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != wantStatus {
		t.Fatalf("GET %s = %d, want %d (body: %s)", path, rec.Code, wantStatus, rec.Body.String())
	}
	var resp timeResponse
	if wantStatus == http.StatusOK {
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("GET %s: bad JSON: %v", path, err)
		}
	}
	return resp
}

// assertTimeFields 校验三字段一致且为可解析的 RFC3339Nano 时间。
func assertTimeFields(t *testing.T, path string, resp timeResponse) {
	t.Helper()
	if resp.DateTime == "" || resp.DateTimeLow == "" || resp.DataTimeAllLow == "" {
		t.Fatalf("GET %s: empty time fields: %+v", path, resp)
	}
	if resp.DateTime != resp.DateTimeLow || resp.DateTime != resp.DataTimeAllLow {
		t.Fatalf("GET %s: fields differ: %+v", path, resp)
	}
	if _, err := time.Parse(time.RFC3339Nano, resp.DateTime); err != nil {
		t.Fatalf("GET %s: not RFC3339Nano: %q (%v)", path, resp.DateTime, err)
	}
}

func TestTimeEndpoint(t *testing.T) {
	r := newRouter(0)
	before := time.Now().UTC()
	resp := fetch(t, r, "/time", http.StatusOK)
	after := time.Now().UTC()
	assertTimeFields(t, "/time", resp)
	parsed, err := time.Parse(time.RFC3339Nano, resp.DateTime)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Before(before.Add(-2*time.Second)) || parsed.After(after.Add(2*time.Second)) {
		t.Fatalf("/time out of range: %v (before=%v after=%v)", parsed, before, after)
	}
}

func TestOffsetTimeEndpoint(t *testing.T) {
	offset := 90 * time.Second
	r := newRouter(offset)
	resp := fetch(t, r, "/offset_time", http.StatusOK)
	assertTimeFields(t, "/offset_time", resp)
	parsed, _ := time.Parse(time.RFC3339Nano, resp.DateTime)
	if parsed.Sub(time.Now().UTC()) < 80*time.Second || parsed.Sub(time.Now().UTC()) > 100*time.Second {
		t.Fatalf("/offset_time offset wrong: got %v from now", parsed.Sub(time.Now().UTC()))
	}
}

// TestNtpTimeEndpoint 验证 /ntp_time 响应契约：
// 网络可用时 200 + 三字段一致时间；不可用时 500 + error 字段。都不应 panic/挂死。
func TestNtpTimeEndpoint(t *testing.T) {
	r := newRouter(0)
	req := httptest.NewRequest(http.MethodGet, "/ntp_time", nil)
	rec := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		defer close(done)
		r.ServeHTTP(rec, req)
	}()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("/ntp_time hung (no timeout)")
	}
	switch rec.Code {
	case http.StatusOK:
		var resp timeResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("/ntp_time bad JSON: %v", err)
		}
		assertTimeFields(t, "/ntp_time", resp)
	case http.StatusInternalServerError:
		// NTP 上游不可达时的降级路径——合法
	default:
		t.Fatalf("/ntp_time = %d, want 200 or 500", rec.Code)
	}
}

func TestOffsetNptTimeEndpointContract(t *testing.T) {
	r := newRouter(0)
	req := httptest.NewRequest(http.MethodGet, "/offset_npt_time", nil)
	rec := httptest.NewRecorder()
	done := make(chan struct{})
	go func() {
		defer close(done)
		r.ServeHTTP(rec, req)
	}()
	select {
	case <-done:
	case <-time.After(15 * time.Second):
		t.Fatal("/offset_npt_time hung (no timeout)")
	}
	switch rec.Code {
	case http.StatusOK:
		var resp timeResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("/offset_npt_time bad JSON: %v", err)
		}
		assertTimeFields(t, "/offset_npt_time", resp)
	case http.StatusInternalServerError:
		// 降级路径——合法
	default:
		t.Fatalf("/offset_npt_time = %d, want 200 or 500", rec.Code)
	}
}
