package speedtest

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// Speed test data sizes
	SmallTestSize  = 1024             // 1KB - For latency tests
	MediumTestSize = 1024 * 1024      // 1MB - For speed test benchmark
	LargeTestSize  = 10 * 1024 * 1024 // 10MB - For large file tests

	// Time and content size headers
	TimeHeader         = "X-Time"
	ContentSizeHeader  = "Content-Size"
	ResponseTimeHeader = "X-Response-Time"
	TestTypeHeader     = "X-Test-Type"

	// Test type definitions
	TestTypeLatency  = "latency"
	TestTypeDownload = "download"
	TestTypeUpload   = "upload"

	// Result units
	UnitKbps = "Kbps"
	UnitMbps = "Mbps"
)

// SpeedTestServer Server-side speed test handler
type SpeedTestServer struct {
	dataPool map[string][]byte // Pre-generated test data cache
	mu       sync.RWMutex
}

// SpeedTestResult Speed test result structure
type SpeedTestResult struct {
	TestType      string `json:"test_type"`      // Test type: latency/download/upload
	ClientTime    int64  `json:"client_time"`    // Client request timestamp (nanoseconds)
	ServerTime    int64  `json:"server_time"`    // Server receive timestamp (nanoseconds)
	ResponseTime  int64  `json:"response_time"`  // Server response timestamp (nanoseconds)
	ClientReceive int64  `json:"client_receive"` // Client response receive timestamp (nanoseconds)

	Latency  float64 `json:"latency"`   // One-way latency (ms)
	RTT      float64 `json:"rtt"`       // Round-trip time (ms)
	Speed    float64 `json:"speed"`     // Transfer speed (Mbps)
	DataSize int64   `json:"data_size"` // Transfer data size (bytes)

	Error     string `json:"error,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

// Ping Handle latency test
func (s *SpeedTestServer) Ping(c *gin.Context) {
	clientTimeStr := c.GetHeader(TimeHeader)
	clientTime, err := strconv.ParseInt(clientTimeStr, 10, 64)
	if err != nil || clientTime <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid time header",
		})
		return
	}

	serverTime := time.Now().UnixNano()

	c.Header(TimeHeader, clientTimeStr)
	c.Header("X-Server-Time", strconv.FormatInt(serverTime, 10))
	c.Header(ContentSizeHeader, "0")

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"timestamp":   time.Now().Unix(),
		"server_time": serverTime,
		"client_time": clientTime,
	})
}

// DownloadSpeedTest Handle download speed test
func (s *SpeedTestServer) DownloadSpeedTest(c *gin.Context) {
	clientTimeStr := c.GetHeader(TimeHeader)
	clientTime, err := strconv.ParseInt(clientTimeStr, 10, 64)
	if err != nil || clientTime <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Time header required",
		})
		return
	}

	// Get test data size (default 1MB)
	size := c.DefaultQuery("size", "1048576")
	dataSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil || dataSize <= 0 || dataSize > 50*1024*1024 {
		dataSize = MediumTestSize
	}

	// Generate random test data
	data := s.generateTestData(dataSize)

	serverTime := time.Now().UnixNano()

	c.Header(TimeHeader, clientTimeStr)
	c.Header("X-Server-Time", strconv.FormatInt(serverTime, 10))
	c.Header(ContentSizeHeader, strconv.FormatInt(dataSize, 10))
	c.Header(TestTypeHeader, TestTypeDownload)

	// Send binary data directly
	c.Data(http.StatusOK, "application/octet-stream", data)
}

// UploadSpeedTest Handle upload speed test
func (s *SpeedTestServer) UploadSpeedTest(c *gin.Context) {
	clientTimeStr := c.GetHeader(TimeHeader)
	clientTime, err := strconv.ParseInt(clientTimeStr, 10, 64)
	if err != nil || clientTime <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Time header required",
		})
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Failed to read upload data",
		})
		return
	}

	serverTime := time.Now().UnixNano()
	clientReceive := time.Now().UnixNano()

	dataSize := int64(len(data))

	// 计算速度
	rtt := float64(clientReceive-clientTime) / 1e6                // 转成毫秒
	latency := rtt / 2                                            // 单边延迟
	uploadSpeed := float64(dataSize*8) / (rtt / 1000) / 1_000_000 // Mbps

	c.JSON(http.StatusOK, SpeedTestResult{
		TestType:      TestTypeUpload,
		ClientTime:    clientTime,
		ServerTime:    serverTime,
		ResponseTime:  clientReceive,
		ClientReceive: clientReceive,
		Latency:       latency,
		RTT:           rtt,
		Speed:         uploadSpeed,
		DataSize:      dataSize,
		Timestamp:     time.Now().Unix(),
	})
}

// generateTestData Generate random test data and cache it
func (s *SpeedTestServer) generateTestData(size int64) []byte {
	s.mu.RLock()
	if data, exists := s.dataPool[strconv.FormatInt(size, 10)]; exists {
		s.mu.RUnlock()
		return data
	}
	s.mu.RUnlock()

	// Generate random data
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}

	s.mu.Lock()
	s.dataPool[strconv.FormatInt(size, 10)] = data
	s.mu.Unlock()

	return data
}
