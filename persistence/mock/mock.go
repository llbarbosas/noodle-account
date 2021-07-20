package mock

import (
	"encoding/json"
	"reflect"
	"sync"

	core "github.com/llbarbosas/noodle-account/core/infra"
)

type Repository struct {
	mux      sync.RWMutex
	data     []reflect.Value
	dataType reflect.Type
}

func NewRepository(t interface{}) Repository {
	return Repository{
		data:     make([]reflect.Value, 0),
		dataType: reflect.TypeOf(t),
	}
}

func (r *Repository) Count() int {
	return len(r.data)
}

// getValueCopy is a hacky way to copy struct data to avoid Fiber/fasthttp
// zero allocation memory sharing
func getValueCopy(value reflect.Value, valueType reflect.Type) (*reflect.Value, error) {
	dataCopy, err := json.Marshal(value.Interface())

	if err != nil {
		return nil, err
	}

	dataPointer := reflect.New(valueType.Elem())

	if err = json.Unmarshal(dataCopy, dataPointer.Interface()); err != nil {
		return nil, err
	}

	elem := dataPointer.Elem()

	return &elem, nil
}

func (r *Repository) Create(data interface{}) error {
	dataValue := GetPointerElem(data)
	dataCopy, err := getValueCopy(dataValue, reflect.TypeOf(data))

	if err != nil {
		return err
	}

	dataValue = *dataCopy

	r.mux.Lock()
	r.data = append(r.data, dataValue)
	r.mux.Unlock()

	return nil
}

func (r *Repository) Get(out interface{}, query core.QueryFunc) error {
	values := reflect.MakeSlice(reflect.SliceOf(r.dataType), 0, r.Count())
	outValue := GetPointerElem(out)

	r.mux.RLock()
	for _, value := range r.data {
		if query(value) == true {
			reflect.Append(values, value)
		}
	}
	r.mux.RUnlock()

	if values.Len() == 0 {
		return core.ErrNotFound
	}

	outValue.Set(values)

	return nil
}

func (r *Repository) GetOne(out interface{}, query core.QueryFunc) error {
	outValue := GetPointerElem(out)

	r.mux.RLock()
	for _, value := range r.data {
		if query(value) == true {
			outValue.Set(value)
			r.mux.RUnlock()
			return nil
		}
	}
	r.mux.RUnlock()

	return core.ErrNotFound
}

func (r *Repository) GetAll(out interface{}) error {
	outValue := GetPointerElem(out)
	values := reflect.MakeSlice(reflect.SliceOf(r.dataType), 0, r.Count())

	r.mux.RLock()
	for _, value := range r.data {
		values = reflect.Append(values, value)
	}
	r.mux.RUnlock()

	outValue.Set(values)

	return nil
}

func (r *Repository) Update(query core.QueryFunc, updateFunc core.GenericUpdateFunc) error {
	r.mux.Lock()
	for i, value := range r.data {
		if query(value) == true {
			r.data[i] = updateFunc(value)
		}
	}
	r.mux.Unlock()

	return nil
}

func (r *Repository) UpdateOne(query core.QueryFunc, updateFunc core.GenericUpdateFunc) error {
	r.mux.Lock()
	for i, value := range r.data {
		if query(value) == true {
			r.data[i] = updateFunc(value)

			r.mux.Unlock()
			return nil
		}
	}
	r.mux.Unlock()

	return nil
}

func (r *Repository) Delete(query core.QueryFunc) error {
	r.mux.Lock()
	for i, value := range r.data {
		if query(value) == true {
			r.data = append(r.data[:i], r.data[i+1:]...)
		}
	}
	r.mux.Unlock()

	return nil
}

func (r *Repository) DeleteOne(query core.QueryFunc) error {
	r.mux.Lock()
	for i, value := range r.data {
		if query(value) == true {
			r.data = append(r.data[:i], r.data[i+1:]...)

			r.mux.Unlock()
			return nil
		}
	}
	r.mux.Unlock()

	return nil
}

func GetPointerElem(p interface{}) reflect.Value {
	v, ok := p.(reflect.Value)

	if !ok {
		v = reflect.ValueOf(p)
	}

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v
}
