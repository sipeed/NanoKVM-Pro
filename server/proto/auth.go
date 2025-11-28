package proto

type LoginReq struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type LoginRsp struct {
	Token string `json:"token"`
	Count int    `json:"count"`
}

type GetAccountRsp struct {
	Username string `json:"username"`
}

type ChangePasswordReq struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type IsPasswordUpdatedRsp struct {
	IsUpdated bool `json:"isUpdated"`
}

type ConnectWifiReq struct {
	Ssid     string `validate:"required" json:"ssid"`
	Password string `validate:"omitempty" json:"password"`
}
