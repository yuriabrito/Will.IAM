// +build unit

package models_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ghostec/Will.IAM/models"
)

func TestValidatePermission(t *testing.T) {
	type testCase struct {
		str   string
		valid bool
		err   error
	}
	tt := []testCase{
		testCase{
			str:   "RX::ListSchedulers::Maestro::*",
			valid: false,
			err:   fmt.Errorf("OwnershipLevel needs to be RO or RL"),
		},
		testCase{
			str:   "RO::ListSchedulers::Maestro::*",
			valid: true,
			err:   nil,
		},
		testCase{
			str:   "RL::ListSchedulers::Maestro::*",
			valid: true,
			err:   nil,
		},
		testCase{
			str:   "RO::ListSchedulers::Maestro::",
			valid: false,
			err:   fmt.Errorf("No parts can be empty"),
		},
		testCase{
			str:   "RL::ListSchedulers::Maestro",
			valid: false,
			err: fmt.Errorf(
				"Incomplete permission. Expected format: " +
					"OwnershipLevel::Action::Service::{ResourceHierarchy}",
			),
		},
	}

	for _, tt := range tt {
		valid, err := models.ValidatePermission(tt.str)
		if valid != tt.valid {
			t.Errorf("Expected valid to be %t. Got %t", valid, tt.valid)
			if err != nil {
				t.Errorf(err.Error())
			}
		}
		if err != nil && err.Error() != tt.err.Error() {
			t.Errorf("Expected error to be %s. Got %s", err.Error(), tt.err.Error())
		}
	}
}

func TestBuildResourceHierarchy(t *testing.T) {
	type testCase struct {
		str       string
		open      bool
		size      int
		hierarchy []string
	}
	tt := []testCase{
		testCase{
			str:       "*",
			open:      true,
			size:      1,
			hierarchy: []string{"*"},
		},
		testCase{
			str:       "SomeGame::*",
			open:      true,
			size:      2,
			hierarchy: []string{"SomeGame", "*"},
		},
		testCase{
			str:       "SomeGame::some-sub-resource",
			open:      false,
			size:      3,
			hierarchy: []string{"SomeGame", "some-sub-resource"},
		},
	}

	for _, tt := range tt {
		rh := models.BuildResourceHierarchy(tt.str)
		if rh.Open != tt.open {
			fmt.Errorf("Expected Open to be %t. Got %t", tt.open, rh.Open)
		}
		if rh.Size != tt.size {
			fmt.Errorf("Expected Size to be %d. Got %d", tt.size, rh.Size)
		}
		if len(rh.Hierarchy) != len(tt.hierarchy) {
			fmt.Errorf(
				"Expected Hierarchy to be %#v. Got %#v", tt.hierarchy, rh.Hierarchy,
			)
		}
		for i := range rh.Hierarchy {
			if rh.Hierarchy[i] != tt.hierarchy[i] {
				fmt.Errorf(
					"Expected Hierarchy[%d] to be %s. Got %s",
					i, tt.hierarchy[i], rh.Hierarchy[i],
				)
			}
		}
	}
}

func TestBuildPermission(t *testing.T) {
	type testCase struct {
		str        string
		permission models.Permission
		err        error
	}
	tt := []testCase{
		testCase{
			str: "RO::ListSchedulers::Maestro::*",
			permission: models.Permission{
				OwnershipLevel: models.OwnershipLevels.Owner,
				Action:         models.BuildAction("ListSchedulers"),
				Service:        models.BuildService("Maestro"),
				ResourceHierarchy: models.ResourceHierarchy{
					Open:      true,
					Size:      1,
					Hierarchy: []string{"*"},
				},
			},
			err: nil,
		},
	}

	for _, tt := range tt {
		permission, _ := models.BuildPermission(tt.str)
		if equal := reflect.DeepEqual(permission, tt.permission); !equal {
			fmt.Errorf(
				"Expected permission to be %#v. Got %#v", tt.permission, permission,
			)
		}
	}
}
