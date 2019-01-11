package usecases

import (
	"strings"

	"github.com/ghostec/Will.IAM/constants"
	"github.com/ghostec/Will.IAM/models"
)

// AM define entrypoints for Access Management actions
type AM interface {
	List(prefix string) ([]models.AM, error)
}

type am struct {
	rsUC Roles
}

func (am *am) List(prefix string) ([]models.AM, error) {
	if strings.Contains(prefix, "::") {
		parts := strings.Split(prefix, "::")
		ams, err := am.listResourceHierarchies(parts[0], parts[1])
		if err != nil {
			return nil, err
		}
		return append(
			[]models.AM{models.AM{Prefix: "*", Complete: true}}, ams...,
		), nil
	}
	actions, err := am.listActions(prefix)
	if err != nil {
		return nil, err
	}
	ams := make([]models.AM, len(actions))
	for i := range actions {
		ams[i] = models.AM{
			Prefix:   actions[i],
			Complete: false,
		}
	}
	return ams, nil
}

func (am *am) listActions(prefix string) ([]string, error) {
	all := append(constants.RolesActions, constants.ServiceAccountsActions...)
	keep := []string{}
	for i := range all {
		if ok := strings.HasPrefix(all[i], prefix); ok {
			keep = append(keep, all[i])
		}
	}
	return keep, nil
}

func actionsContains(actions []string, action string) bool {
	for _, aa := range actions {
		if aa == action {
			return true
		}
	}
	return false
}

func (am *am) listResourceHierarchies(
	action, rhPrefix string,
) ([]models.AM, error) {
	if actionsContains(constants.RolesActions, action) {
		return am.listRolesActionsRH(rhPrefix)
	}
	if actionsContains(constants.ServiceAccountsActions, action) {
		return []models.AM{}, nil
	}
	if actionsContains(constants.ServicesActions, action) {
		return []models.AM{}, nil
	}
	return []models.AM{}, nil
}

func (am *am) listRolesActionsRH(
	rhPrefix string,
) ([]models.AM, error) {
	rs, err := am.rsUC.WithNamePrefix(rhPrefix, 10)
	if err != nil {
		return nil, err
	}
	ams := make([]models.AM, len(rs))
	for i := range rs {
		ams[i] = models.AM{
			Prefix:   rs[i].ID,
			Alias:    rs[i].Name,
			Complete: true,
		}
	}
	return ams, nil
}

// NewAM am ctor
func NewAM(rsUC Roles) AM {
	return &am{rsUC: rsUC}
}
