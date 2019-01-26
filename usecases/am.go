package usecases

import (
	"fmt"
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
	if !strings.Contains(prefix, "::") {
		services, err := am.listServices(prefix)
		if err != nil {
			return nil, err
		}
		ams := make([]models.AM, len(services))
		for i := range services {
			ams[i] = models.AM{
				Prefix:   services[i],
				Complete: false,
			}
		}
		return ams, nil
	}
	parts := strings.Split(prefix, "::")
	if len(parts) == 2 {
		actions, err := am.listActions(parts[0], parts[1])
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
	ams, err := am.listResourceHierarchies(parts[0], parts[1])
	if err != nil {
		return nil, err
	}
	return append(
		[]models.AM{models.AM{Prefix: "*", Complete: true}}, ams...,
	), nil
}

func (am *am) listServices(prefix string) ([]string, error) {
	svcs := []string{constants.AppInfo.Name}
	filtered := []string{}
	for i := range svcs {
		if strings.HasPrefix(svcs[i], prefix) {
			filtered = append(filtered, svcs[i])
		}
	}
	return filtered, nil
}

func (am *am) listActions(service, prefix string) ([]string, error) {
	// TODO: check if service exists and user can list its actions
	all := append(constants.RolesActions, constants.ServiceAccountsActions...)
	keep := []string{}
	for i := range all {
		if ok := strings.HasPrefix(all[i], prefix); ok {
			keep = append(keep, fmt.Sprintf("%s::%s", service, all[i]))
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
