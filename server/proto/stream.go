package proto

type SetModeReq struct {
	Mode string `form:"mode" validate:"required"`
}

type SetQualityReq struct {
	Quality int `form:"quality" validate:"required"`
}

type SetGopReq struct {
	Gop int `form:"gop" validate:"required"`
}
