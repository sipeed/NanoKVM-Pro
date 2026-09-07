package proto

type GetVersionRsp struct {
	Current string `json:"current"`
	Latest  string `json:"latest"`
}

type GetPreviewRsp struct {
	Enabled bool `json:"enabled"`
}

type SetPreviewReq struct {
	Enable bool `validate:"omitempty"`
}

type GetUpdateSourceRsp struct {
	// URL is the root used to derive application and firmware update paths.
	URL        string `json:"url"`
	IsOfficial bool   `json:"isOfficial"`
	Enabled    bool   `json:"enabled"`
}

type SetUpdateSourceReq struct {
	// URL may be HTTP or HTTPS; the service validates and normalizes it before saving.
	URL string `json:"url"`
}

type ConfirmManualUpdateReq struct {
	// ID identifies the staged package that the user has reviewed.
	ID string `json:"id" validate:"required"`
}

type ManualUpdateRsp struct {
	// State is staged, installing, succeeded, failed, or reboot_scheduled.
	ID      string `json:"id"`
	Version string `json:"version"`
	Size    int64  `json:"size"`
	State   string `json:"state"`
	Error   string `json:"error,omitempty"`
}
