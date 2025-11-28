package ui

type Response struct {
	Ok     bool   `json:"ok"`
	Result Result `json:"result"`
}

type Result struct {
	Delay      int      `json:"delay"`
	Deprecated []string `json:"deprecated"`
	Features   Features `json:"features"`
	Limits     Limits   `json:"limits"`
	Params     Params   `json:"params"`
	Snapshot   Snapshot `json:"snapshot"`
	Streamer   Streamer `json:"streamer"`
}

type Features struct {
	H264       bool `json:"h264"`
	Quality    bool `json:"quality"`
	Resolution bool `json:"resolution"`
}

type Limits struct {
	AvailableResolutions []string  `json:"available_resolutions"`
	DesiredFps           MinMaxInt `json:"desired_fps"`
	H264Bitrate          MinMaxInt `json:"h264_bitrate"`
	H264Gop              MinMaxInt `json:"h264_gop"`
}

type MinMaxInt struct {
	Max int `json:"max"`
	Min int `json:"min"`
}

type Params struct {
	DesiredFps  int    `json:"desired_fps"`
	H264Bitrate int    `json:"h264_bitrate"`
	H264Gop     int    `json:"h264_gop"`
	Resolution  string `json:"resolution"`
}

type Snapshot struct {
	Saved interface{} `json:"saved"`
}

type Streamer struct {
	Encoder    Encoder `json:"encoder"`
	H264       H264    `json:"h264"`
	InstanceID string  `json:"instance_id"`
	Sinks      Sinks   `json:"sinks"`
	Source     Source  `json:"source"`
	Stream     Stream  `json:"stream"`
}

type Encoder struct {
	Quality int    `json:"quality"`
	Type    string `json:"type"`
}

type H264 struct {
	Bitrate     int  `json:"bitrate"`
	RealBitrate int  `json:"real_bitrate"`
	Fps         int  `json:"fps"`
	Gop         int  `json:"gop"`
	Online      bool `json:"online"`
}

type Sinks struct {
	H264 SinkStatus `json:"h264"`
	Jpeg SinkStatus `json:"jpeg"`
}

type SinkStatus struct {
	HasClients bool `json:"has_clients"`
}

type Source struct {
	CapturedFps int        `json:"captured_fps"`
	DesiredFps  int        `json:"desired_fps"`
	Online      bool       `json:"online"`
	Resolution  Resolution `json:"resolution"`
}

type Resolution struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Stream struct {
	Clients     int                    `json:"clients"` // 0-H264 1-MJPEG
	ClientsStat map[string]interface{} `json:"clients_stat"`
	QueuedFps   int                    `json:"queued_fps"`
}
