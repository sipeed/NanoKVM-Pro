package proto

type GetImagesRsp struct {
	Files []string `json:"files"`
}

type MountImageReq struct {
	File     string `form:"file" validate:"omitempty"`
	Cdrom    bool   `form:"cdrom" validate:"omitempty"`
	ReadOnly bool   `form:"readOnly" validate:"omitempty"`
}

type GetMountedImageRsp struct {
	File     string `json:"file"`
	Cdrom    bool   `json:"cdrom"`
	ReadOnly bool   `json:"readOnly"`
}

type DeleteImageReq struct {
	File string `json:"file" validate:"required"`
}

type ImageEnabledRsp struct {
	Enabled bool `json:"enabled"`
}

type StatusImageRsp struct {
	Status     string `json:"status"`
	File       string `json:"file"`
	Percentage string `json:"percentage"`
}
