package application

const (
	StableURL  = "https://cdn.sipeed.com/nanokvm"
	PreviewURL = "https://cdn.sipeed.com/nanokvm/preview"

	AppDir  = "/kvmapp"
	TempDir = "/root/.kvmcache"
)

var appNames = []string{"nanokvmpro", "pikvm", "kvmcomm"}

type Service struct{}

func NewService() *Service {
	return &Service{}
}
