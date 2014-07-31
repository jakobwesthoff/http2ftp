package http2ftp

import "github.com/yob/graval"

// HTTPDriverFactory creates a new httpDriver for each graval connection
type HTTPDriverFactory struct {
	Configurations map[string]Configuration
}

// NewDriver creates a new correctly configured HTTPDriver
func (factory *HTTPDriverFactory) NewDriver() (graval.FTPDriver, error) {
	return &httpDriver{
			Configurations:    factory.Configurations,
			AuthenticatedUser: nil,
			UserConfiguration: nil,
		},
		nil
}
