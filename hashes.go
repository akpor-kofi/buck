package buckis

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

type structable struct {
	err    error
	hashes map[string]any
}

func (s *structable) Scan(o any) {
	if s.err != nil {
		panic(s.err)
	}

	mapBytes, err := json.Marshal(s.hashes)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(mapBytes, o)
	if err != nil {
		panic(err)
	}

}

var ErrHashNotFound = errors.New("hash not found")

// HSet func (d *dict) HSet(key string, hashes ...string) error {
func (d *dict) HSet(key string, entity any) error {
	var hashes []any

	t := reflect.TypeOf(entity)
	v := reflect.ValueOf(entity)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		k := field.Tag.Get("buckis")

		if k == "" {
			continue
		}

		value := v.Field(i)

		switch v.Field(i).Kind() {
		case reflect.String:
			hashes = append(hashes, k, value.String())
		case reflect.Bool:
			hashes = append(hashes, k, value.Bool())
		case reflect.Int:
			hashes = append(hashes, k, value.Int())
		case reflect.Float32, reflect.Float64:
			hashes = append(hashes, k, value.Float())
		default:
			return fmt.Errorf("value type not supported: %s", k)
		}
	}

	return d.hset(SAVE, key, hashes...)
}

func (d *dict) HGetAll(key string) *structable {
	de, err := d.hashesLookup(key)

	if err != nil {
		return &structable{
			err: err,
		}
	}

	hashesStruct := &structable{
		hashes: de.values.(map[string]any),
	}

	return hashesStruct
}

func (d *dict) hashesLookup(key string) (*dictEntry, error) {
	i := d.hash(key)
	currentEntry := d.ht[Hashes][i]

	if currentEntry == nil {
		return &dictEntry{}, ErrHashNotFound
	}

	for {
		if currentEntry.key == key {
			return currentEntry, nil
		}

		if currentEntry.next == nil {
			return &dictEntry{}, ErrHashNotFound
		}

		currentEntry = currentEntry.next
	}
}

func (d *dict) HGet(key, hashKey string) (any, error) {
	de, err := d.hashesLookup(key)

	if err != nil {
		return nil, err
	}

	return de.values.(map[string]any)[hashKey], nil
}

// HIncrBy TODO: implement this
func (d *dict) HIncrBy(key, hashKey string, incr int) (int, error) {
	de, err := d.hashesLookup(key)
	if err != nil {
		return 0, nil
	}

	value, ok := de.values.(map[string]any)[hashKey]

	if !ok {
		return 0, fmt.Errorf("attribute does not exists")
	}

	switch value.(type) {
	case int, int64:
		a := int(value.(int64))
		a += incr
		de.values.(map[string]any)[hashKey] = a
		return a, nil
	}

	return 0, fmt.Errorf("not an integer")
}
