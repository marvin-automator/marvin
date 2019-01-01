package actions

import (
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/internal/db"
	"reflect"
)

const configsStoreName = "requirements_configs"

func loadRequirementConfig(providerName string, req actions.Requirement) error {
	s := db.GetStore(configsStoreName)

	c := req.Config()
	isPointer := true
	t := reflect.TypeOf(c)

	if t.Kind() != reflect.Ptr {
		isPointer = false
		t = t.Elem()
	}

	i := reflect.New(t).Interface()

	key := makeKey(providerName, req)
	err := s.Get(key, i)

	if _, ok := err.(db.KeyNotFoundError); ok {
		return nil // There's no stored configuration for this requirement, not a problem.
	}

	if err != nil {
		return err
	}

	if !isPointer {
		i = reflect.Indirect(reflect.ValueOf(i)).Interface()
	}

	return req.SetConfig(i)
}

func storeRequirementConfig(providerName string, req actions.Requirement) error {
	s := db.GetStore(configsStoreName)

	key := makeKey(providerName, req)
	return s.Set(key, req.Config())
}

func makeKey(providerName string, req actions.Requirement) string {
	return providerName + "|" + req.Name()
}
