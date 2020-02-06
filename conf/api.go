package conf

import (
	"errors"
	"net/url"
)

// API contains the required data to contact the CTFd API endpoints.
// TODO: CTFd should support API tokens as the current approach is a hack.
type API struct {
	Session string `yaml:"session"`
	CSRF    string `yaml:"csrf"`
	URL     string `yaml:"url"`
}

// Check if the API state seems credible.
func (a *API) Check() error {
	// The API requires at least a session token...
	if len(a.Session) == 0 {
		return errors.New("missing API session")
	}
	// ...a CSRF token...
	if len(a.CSRF) == 0 {
		return errors.New("missing API CSRF token")
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
