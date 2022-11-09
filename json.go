package buckis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// arr append, arr len, arr insert, arr pop, arr trim
// get
//

func (d *dict) jsonSet(flag int, key, path, value string) error {
	jde, err := d.jsonLookup(key)
	i := d.hash(key)

	if err != nil {
		// create the json, thinking of storing it as map[string]any type
		var o map[string]any
		err := json.Unmarshal([]byte(value), &o)
		if err != nil {
			return err
		}

		de := &dictEntry{
			key:    key,
			values: o,
			next:   d.ht[Json][i],
		}

		d.ht[Json][i] = de

		if flag == SAVE {
			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(JSONSET, key, path, value)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}

		return nil
	}

	jsonObj := jde.values.(map[string]any)

	jsonSetInPath(jsonObj, nil, path, value)

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(JSONSET, key, path, value)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}

func (d *dict) JSONSet(key string, path string, value string) error {
	return d.jsonSet(SAVE, key, path, value)
}

func (d *dict) JSONGet(key string, path string) (any, error) {
	jde, err := d.jsonLookup(key)
	if err != nil {
		return "", err
	}

	value := jsonGetFromPath(jde.values.(map[string]any), path)

	fmt.Println(value, "here")

	b, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		return nil, err
	}

	return string(b), err
}

func (d *dict) jsonLookup(key string) (*dictEntry, error) {
	i := d.hash(key)

	currentEntry := d.ht[Json][i]

	if currentEntry == nil {
		return &dictEntry{}, ErrSetNotFound
	}

	for {
		if currentEntry.key == key {
			return currentEntry, nil
		}

		if currentEntry.next == nil {
			return &dictEntry{}, ErrSetNotFound
		}

		currentEntry = currentEntry.next
	}
}

func jsonSetInPath(jsonObj, parentObj map[string]any, path, value string) {
	paths := strings.Split(path, ".")
	fmt.Println(paths)

	if len(paths) <= 1 || path == "." {
		// do the setting here
		var pathJsonValue map[string]any

		err := json.Unmarshal([]byte(value), &pathJsonValue)
		if err != nil {
			// if value possible scenario is it's just a value
			fmt.Println("kkk")
			parentObj[paths[0]] = value
			return
		}

		for k, v := range pathJsonValue {
			jsonObj[k] = v
		}

		return
	}

	childPath := paths[1]

	if objVal, ok := jsonObj[childPath]; !ok {
		// no value was found so create an object
		jsonObj[childPath] = map[string]any{}
	} else {
		switch objVal.(type) {
		case map[string]any:
		default:
			jsonObj[childPath] = map[string]any{}
		}
	}

	jsonSetInPath(jsonObj[childPath].(map[string]any), jsonObj, strings.Join(paths[1:], "."), value)
}

func jsonGetFromPath(jsonObj map[string]any, path string) any {
	paths := strings.Split(path, ".")

	if path == "." {
		return jsonObj
	}

	currentObj := jsonObj

	for i := 1; i < len(paths); i++ {
		value := jsonObj[paths[i]]

		switch value.(type) {
		case map[string]any:
			currentObj = value.(map[string]any)

			if i == len(paths)-1 {
				return currentObj
			}

			continue
		default:
			return value
		}

	}

	return nil
}
