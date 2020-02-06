package conf

import "errors"

// Configuration is the root object defined in the configuration file.
type Configuration struct {
	API   *API
	Flags []*Flag
}

// Check if the configuration state seems credible.
func (c *Configuration) Check() error {
	// The API is a required field
	if c.API == nil {
		return errors.New("missing API for configuration")
	} else if err := c.API.Check(); err != nil {
		return err
	}
	// The Flags are required...
	if len(c.Flags) == 0 {
		return errors.New("missing flags for configuration")
	}
	// ...and can't be empty
	for _, f := range c.Flags {
		if f == nil {
			return errors.New("missing flag in flags")
		} else if err := c.API.Check(); err != nil {
			return err
		}
	}
	return nil
}
