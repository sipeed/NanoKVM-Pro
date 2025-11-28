package proto

type WakeOnLANReq struct {
	Mac string `form:"mac" validate:"required"`
}

type GetMacRsp struct {
	Macs []string `json:"macs"`
}

type DeleteMacReq struct {
	Mac string `form:"mac" validate:"required"`
}

type SetMacNameReq struct {
	Mac  string `form:"mac" validate:"required"`
	Name string `form:"name" validate:"required"`
}

type GetStaticIPRsp struct {
	Enabled bool   `json:"enabled"`
	IP      string `json:"ip"`
}

type SetStaticIPReq struct {
	Enabled bool   `form:"enabled" validate:"omitempty"`
	IP      string `form:"ip" validate:"omitempty"`
}

type WiFi struct {
	SSID      string `json:"ssid"`
	BSSID     string `json:"bssid"`
	Signal    int    `json:"signal"`
	Frequency int    `json:"frequency"`
	Security  string `json:"security"`
}

type GetWifiRsp struct {
	Supported bool  `json:"supported"`
	APMode    bool  `json:"apMode"`
	Connected bool  `json:"connected"`
	Wifi      *WiFi `json:"wifi"`
}

type ScanWifiRsp struct {
	WifiList []*WiFi `json:"wifiList"`
}
