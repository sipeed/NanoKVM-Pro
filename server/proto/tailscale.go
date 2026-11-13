package proto

type TailscaleState string

const (
	TailscaleNotInstall TailscaleState = "notInstall"
	TailscaleNotRunning TailscaleState = "notRunning"
	TailscaleNotLogin   TailscaleState = "notLogin"
	TailscaleStopped    TailscaleState = "stopped"
	TailscaleRunning    TailscaleState = "running"
)

type GetTailscaleStatusRsp struct {
	State    TailscaleState `json:"state"`
	Name     string         `json:"name"`
	IP       string         `json:"ip"`
	Account  string         `json:"account"`
	Serve    bool           `json:"serve"`
	ServeUrl string         `json:"serveUrl"`
}

type ServeTailscaleReq struct {
	Enable bool `json:"enable"`
}

// ServeTailscaleRsp carries the consent/authorization URL (if any) and the resulting HTTPS serve URL.
type ServeTailscaleRsp struct {
	AuthUrl  string `json:"authUrl"`
	ServeUrl string `json:"serveUrl"`
}

type LoginTailscaleRsp struct {
	Url string `json:"url"`
}
