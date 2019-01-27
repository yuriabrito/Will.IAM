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
	service := parts[0]
	if service == constants.AppInfo.Name {
		return am.listWillIAMPermissions(prefix)
	}
	// TODO: find service by name and use it's /am (AMURL)
	return []models.AM{}, nil
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

func (am *am) listWillIAMPermissions(prefix string) ([]models.AM, error) {
	parts := strings.Split(prefix, "::")
	if len(parts) == 2 {
		actions, err := am.listWillIAMActions(parts[1])
		if err != nil {
			return nil, err
		}
		ams := make([]models.AM, len(actions))
		for i := range actions {
			ams[i] = models.AM{
				Prefix:   fmt.Sprintf("%s::%s", parts[0], actions[i]),
				Complete: false,
			}
		}
		return ams, nil
	}
	ams, err := am.listWillIAMResourceHierarchies(parts[1], parts[2])
	if err != nil {
		return nil, err
	}
	for i := range ams {
		ams[i].Prefix = fmt.Sprintf("%s::%s::%s", parts[0], parts[1], ams[i].Prefix)
	}
	return append(
		[]models.AM{models.AM{
			Prefix: fmt.Sprintf("%s::%s::*", parts[0], parts[1]), Complete: true,
		}}, ams...,
	), nil
}

func (am *am) listWillIAMActions(prefix string) ([]string, error) {
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

func (am *am) listWillIAMResourceHierarchies(
	action, prefix string,
) ([]models.AM, error) {
	if actionsContains(constants.RolesActions, action) {
		return am.listRolesActionsRH(action, prefix)
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
	action, prefix string,
) ([]models.AM, error) {
	if action == "CreateRole" || action == "ListRoles" {
		return []models.AM{}, nil
	}
	rs, err := am.rsUC.WithNamePrefix(prefix, 10)
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
