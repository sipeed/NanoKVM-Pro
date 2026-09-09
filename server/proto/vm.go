package proto

type IP struct {
	Name    string `json:"name"`
	Addr    string `json:"addr"`
	Version string `json:"version"`
	Type    string `json:"type"`
}

type GetInfoRsp struct {
	IPs          []IP   `json:"ips"`
	Mdns         string `json:"mdns"`
	Image        string `json:"image"`
	Application  string `json:"application"`
	DeviceKey    string `json:"deviceKey"`
	DeviceNumber string `json:"pn"`
	Arch         string `json:"arch"`
}

type GetHardwareRsp struct {
	Version string `json:"version"`
}

type SetGpioReq struct {
	Type     string `validate:"required"`  // reset / power
	Duration uint   `validate:"omitempty"` // press time (unit: milliseconds)
}

type GetGpioRsp struct {
	PWR bool `json:"pwr"` // power led
	HDD bool `json:"hdd"` // hdd led
}

type GetScriptsRsp struct {
	Files []string `json:"files"`
}

type UploadScriptRsp struct {
	File string `json:"file"`
}

type RunScriptReq struct {
	Name string `validate:"required"`
	Type string `validate:"required"` // foreground | background
}

type RunScriptRsp struct {
	Log string `json:"log"`
}

type DeleteScriptReq struct {
	Name string `validate:"required"`
}

type GetVirtualDeviceRsp struct {
	IsNetworkEnabled bool   `json:"isNetworkEnabled"`
	IsMicEnabled     bool   `json:"isMicEnabled"`
	MountedDisk      string `json:"mountedDisk"`
	IsEmmcExist      bool   `json:"isEmmcExist"`
	IsSdCardExist    bool   `json:"isSdCardExist"`
}

type UpdateVirtualDeviceReq struct {
	Device string `json:"device" validate:"required"` // disk | network | mic
	Type   string `json:"type" validate:"omitempty"`  // sdcard | emmc (only available for disk)
}

type RefreshVirtualDeviceReq struct {
	Device string `json:"device" validate:"required"`
}

type SetMemoryLimitReq struct {
	Enabled bool  `validate:"omitempty"`
	Limit   int64 `validate:"omitempty"`
}

type GetMemoryLimitRsp struct {
	Enabled bool  `json:"enabled"`
	Limit   int64 `json:"limit"`
}

type SetOledReq struct {
	Sleep int `validate:"omitempty"`
}

type GetOLEDRsp struct {
	Exist bool   `json:"exist"`
	Type  string `json:"type"`
	Sleep int    `json:"sleep"`
}

type GetLcdTimeFormatRsp struct {
	Format string `json:"format"`
}

type SetLcdTimeFormatReq struct {
	Format string `json:"format" validate:"required"`
}

type GetLcdScreenOffRsp struct {
	Supported   bool `json:"supported"`
	Enabled     bool `json:"enabled"`
	StartMinute int  `json:"startMinute"`
	EndMinute   int  `json:"endMinute"`
}

type SetLcdScreenOffReq struct {
	Enabled     *bool `json:"enabled"`
	StartMinute *int  `json:"startMinute"`
	EndMinute   *int  `json:"endMinute"`
}

type LcdDisplaySchedule struct {
	Supported          bool `json:"supported"`
	Enabled            bool `json:"enabled"`
	StartMinute        int  `json:"startMinute"`
	EndMinute          int  `json:"endMinute"`
	WakeTimeoutSeconds int  `json:"wakeTimeoutSeconds,omitempty"`
}

type GetLcdDisplayPolicyRsp struct {
	ScreenType         string             `json:"screenType"`
	SupportedModes     []string           `json:"supportedModes"`
	Mode               string             `json:"mode,omitempty"`
	ModeTimeoutSeconds map[string]int     `json:"modeTimeoutSeconds,omitempty"`
	Schedule           LcdDisplaySchedule `json:"schedule"`
}

type SetLcdDisplayPolicyReq struct {
	Mode     *string                `json:"mode"`
	Schedule *SetLcdDisplaySchedule `json:"schedule"`
}

type SetLcdDisplaySchedule struct {
	Enabled     *bool `json:"enabled"`
	StartMinute *int  `json:"startMinute"`
	EndMinute   *int  `json:"endMinute"`
}

type GetSSHStateRsp struct {
	Enabled bool `json:"enabled"`
}

type GetMouseJigglerRsp struct {
	Enabled bool   `json:"enabled"`
	Mode    string `json:"mode"`
}

type SetMouseJigglerReq struct {
	Enabled bool   `validate:"omitempty"`
	Mode    string `validate:"omitempty"`
}

type GetMdnsStateRsp struct {
	Enabled bool `json:"enabled"`
}

type GetLowPowerRsp struct {
	Enabled bool `json:"enabled"`
}

type SetLowPowerReq struct {
	Enable bool `form:"enable" validate:"omitempty"`
}

type SetHostnameReq struct {
	Hostname string `validate:"required"`
}

type GetHostnameRsp struct {
	Hostname string `json:"hostname"`
}

type SetWebTitleReq struct {
	Title string `validate:"omitempty"`
}

type GetWebTitleRsp struct {
	Title string `json:"title"`
}

type SetLedStripReq struct {
	On         bool `validate:"omitempty"`
	Hor        int  `validate:"omitempty"`
	Ver        int  `validate:"omitempty"`
	Brightness int  `validate:"omitempty"`
}

type GetEdidRsp struct {
	Edid string `json:"edid"`
}

type SwitchEdidReq struct {
	Edid string `form:"edid" validate:"required"`
}

type GetCustomEdidListRsp struct {
	EdidList []string `json:"edidList"`
}

type DeleteEdidReq struct {
	Edid string `form:"edid" validate:"required"`
}

type UploadEdidRsp struct {
	File string `json:"file"`
}

type GetHdmiCaptureRsp struct {
	Enabled bool `json:"enabled"`
}

type SetHdmiCaptureReq struct {
	Enabled bool `form:"enabled" validate:"omitempty"`
}

type GetHdmiPassthroughRsp struct {
	Enabled bool `json:"enabled"`
}

type SetHdmiPassthroughReq struct {
	Enabled bool `form:"enabled" validate:"omitempty"`
}

type SetTimeZoneReq struct {
	Timezone string `validate:"required"`
}

type GetTimeZoneRsp struct {
	Timezone string `json:"timezone"`
}

type GetTimeStatusRsp struct {
	IsSynchronized bool  `json:"isSynchronized"`
	LastSyncTime   int64 `json:"lastSyncTime"`
}

type MenuBarConfig struct {
	DisabledItems []string `json:"disabledItems"`
}

type SetMenuBarConfigReq struct {
	DisabledItems []string `form:"disabledItems" json:"disabledItems" validate:"omitempty"`
}

type GetMenuBarConfigRsp struct {
	DisabledItems []string `json:"disabledItems"`
}
