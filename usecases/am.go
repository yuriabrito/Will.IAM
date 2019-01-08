package usecases

import (
	"strings"

	"github.com/ghostec/Will.IAM/constants"
)

// AM define entrypoints for Access Management actions
type AM interface {
	List(prefix string) ([]string, error)
}

type am struct{}

func (am *am) List(prefix string) ([]string, error) {
	if strings.Contains(prefix, "::") {
		// skip actions
		return nil, nil
	}
	return am.listActions(prefix)
}

func (am *am) listActions(prefix string) ([]string, error) {
	all := append(constants.Actions.Roles, constants.Actions.ServiceAccounts...)
	keep := []string{}
	for i := range all {
		if ok := strings.HasPrefix(all[i], prefix); ok {
			keep = append(keep, all[i])
		}
	}
	return keep, nil
}

// NewAM am ctor
func NewAM() AM {
	return &am{}
}
