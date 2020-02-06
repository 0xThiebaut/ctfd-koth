package conf

import "errors"

// Award contains the fields needed to inform the CTFd API of awarded points.
type Award struct {
	User        int    `json:"user_id" yaml:"-"`
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description"`
	Value       int    `json:"value" yaml:"value"`
	Category    string `json:"category" yaml:"category"`
	Icon        string `json:"icon,omitempty" yaml:"icon"`
}

// Check if the award state seems credible.
func (a *Award) Check() error {
	// An award should at least have a name...
	if len(a.Name) == 0 {
		return errors.New("missing name for award")
	}
	// ...and a category
	if len(a.Category) == 0 {
		return errors.New("missing category for award")
	}
	return nil
}
