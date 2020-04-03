package conf

import (
	"errors"
	"net/url"
)

// API contains the required data to contact the CTFd API endpoints.
type API struct {
	Token  	string `yaml:"token"`
	URL     string `yaml:"url"`
}

// Check if the API state seems credible.
func (a *API) Check() error {
	// The API requires at least a session token...
	if len(a.Token) == 0 {
		return errors.New("missing API token")
	}
	// ...and a valid URL
	if len(a.URL) == 0 {
		return errors.New("missing API URL")
	}
	if _, err := url.Parse(a.URL); err != nil {
		return err
	}
	return nil
}
