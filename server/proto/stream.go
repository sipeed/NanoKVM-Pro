package proto

type SetRateControlModeReq struct {
	Mode string `form:"mode" validate:"required"` // cbr / vbr
}

type SetModeReq struct {
	Mode string `form:"mode" validate:"required"`
}

type GetModeRsp struct {
	Mode   string `json:"mode"`
	Active string `json:"active"`
}

type SetQualityReq struct {
	Quality int `form:"quality" validate:"required"`
}

type SetGopReq struct {
	Gop int `form:"gop" validate:"required"`
}

type SetFpsReq struct {
	Fps int `form:"fps" validate:"omitempty"`
}
