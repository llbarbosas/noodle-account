package infra

import (
	"errors"
	"reflect"
)

var ErrNotFound = errors.New("Cannot found data")

type QueryFunc func(entity reflect.Value) bool
type GenericUpdateFunc func(entity reflect.Value) reflect.Value

type GenericRepositoryCreaters interface {
	Create(entity interface{}) error
}

type GenericRepositoryGetters interface {
	Get(out interface{}, query QueryFunc) error
	GetOne(out interface{}, query QueryFunc) error
	GetAll(out interface{}) error
	Count() int
}

type GenericRepositoryUpdaters interface {
	Update(query QueryFunc, updateFunc GenericUpdateFunc) error
}

type GenericRepositoryDeleters interface {
	Delete(query QueryFunc) error
}

type GenericRepository interface {
	GenericRepositoryCreaters
	GenericRepositoryGetters
	GenericRepositoryUpdaters
	GenericRepositoryDeleters
}

func QueryByStringField(field string, value string) QueryFunc {
	return func(v reflect.Value) bool {
		return v.FieldByName(field).String() == value
	}
}
