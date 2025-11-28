package speedtest

func NewServer() *SpeedTestServer {
	server := &SpeedTestServer{
		dataPool: make(map[string][]byte),
	}
	return server
}
