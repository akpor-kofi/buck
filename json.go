package buckis

import (
	"encoding/json"
	"fmt"
	"strings"
)

// arr append, arr len, arr insert, arr pop, arr trim
// get
//

func (d *dict) jsonSet(flag int, key, path string, value any) (err error) {
	jde, err := d.jsonLookup(key)
	i := d.hash(key)

	if err != nil {
		de := &dictEntry{
			key:  key,
			next: d.ht[Json][i],
		}

		switch value.(type) {
		case map[string]any:
			de.values = value
		case string:
			// try to unmarshal object string
			s := value.(string)
			var o map[string]any
			err = json.Unmarshal([]byte(s), &o)
			if err != nil {
				return
			}
			fmt.Println(o)
			de.values = o
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

		return
	}

	if path == "." {
		var o map[string]any
		switch value.(type) {
		case string:
			s := value.(string)
			err = json.Unmarshal([]byte(s), &o)
			if err != nil {
				return
			}
			jde.values = o
			return
		}

	} else {
		jsonObj := jde.values.(map[string]any)

		jsonSetInPath(jsonObj, path, value)
	}

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

func (d *dict) JSONSet(key string, path string, value any) error {
	return d.jsonSet(SAVE, key, path, value)
}

type jsonMap struct {
	isRoot  bool
	pathKey string
	value   any
	err     error
}

func (s *jsonMap) Scan(o any) error {
	if s.err != nil {
		panic(s.err)
	}

	// handle root get
	if s.isRoot {
		// is definitely type map[string]interface{} so i need to marshall and unmarshal
		mapBytes, err := json.Marshal(s.value)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(mapBytes, o)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(mapBytes))

		return nil
	}

	mapBytes, err := json.Marshal(s.value)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(mapBytes, o)
	if err != nil {
		panic(err)
	}

	return nil
}

func (d *dict) JSONGet(key string, path string) *jsonMap {
	jde, err := d.jsonLookup(key)
	if err != nil {
		return &jsonMap{err: err}
	}

	a := &jsonMap{}

	value := jsonGetFromPath(jde.values.(map[string]any), path)

	if value == nil {
		a.err = fmt.Errorf("no path found")
	}

	if path == "." {
		a.isRoot = true
		a.value = value
		return a
	}

	paths := strings.Split(path, ".")
	a.pathKey = paths[len(paths)-1]
	a.value = value

	return a
}

func (d *dict) JSONNumIncrBy(key string, path string, incr int) (err error) {

	jde, err := d.jsonLookup(key)

	if err != nil {
		return
	}

	jsonObj := jde.values.(map[string]any)

	pathValue := jsonGetFromPath(jsonObj, path)

	fmt.Println(pathValue)

	switch pathValue.(type) {
	case int:
		b := pathValue.(int)
		b += incr
		jsonSetInPath(jsonObj, path, b)
	case float64:
		b := pathValue.(float64)
		b += float64(incr)
		jsonSetInPath(jsonObj, path, b)
	case nil:
		err = fmt.Errorf("no path found")
	default:
		err = fmt.Errorf("not a number")
	}

	return
}

func (d *dict) JSONToggle(key string, path string) (err error) {
	jde, err := d.jsonLookup(key)

	if err != nil {
		return
	}

	jsonObj := jde.values.(map[string]any)

	pathValue := jsonGetFromPath(jsonObj, path)

	switch pathValue.(type) {
	case bool:
		jsonSetInPath(jsonObj, path, !pathValue.(bool))
	case nil:
		err = fmt.Errorf("no path found")
	default:
		err = fmt.Errorf("not a boolean")
	}

	return
}

// JSONDel delete a key recursively
func (d *dict) JSONDel(key string, path string) (err error) {
	return
}

// JSONArrAppend currently only support strings and numbers and not structs
func (d *dict) JSONArrAppend(key, path string, element ...string) (err error) {
	jde, err := d.jsonLookup(key)

	if err != nil {
		return
	}

	jsonObj := jde.values.(map[string]any)

	pathValue := jsonGetFromPath(jsonObj, path)

	switch pathValue.(type) {
	case []string:
		arr := pathValue.([]string)
		arr = append(arr, element...)
		jsonSetInPath(jsonObj, path, arr)
	}

	return
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

func jsonSetInPath(jsonObj map[string]any, path string, value any) {
	paths := strings.Split(path, ".")

	if len(paths) <= 1 || path == "." {
		jsonObj[paths[0]] = value

		return
	}

	childPathPosition := 1
	childNextChild := childPathPosition + 1
	childPath := paths[childPathPosition]

	if childObj, ok := jsonObj[childPath]; !ok {
		// no value was found so create an object
		jsonObj[childPath] = map[string]any{}

		// NOTE: a way to know if the child path is the last path
		if len(paths) <= childNextChild {
			jsonObj[childPath] = value
		} else {
			jsonSetInPath(jsonObj[childPath].(map[string]any), strings.Join(paths[1:], "."), value)
		}

	} else {
		switch childObj.(type) {
		case map[string]any:
			jsonSetInPath(jsonObj[childPath].(map[string]any), strings.Join(paths[1:], "."), value)
		default:
			jsonObj[childPath] = value
		}
	}

}

func jsonGetFromPath(jsonObj map[string]any, path string) any {
	paths := strings.Split(path, ".")

	if path == "." {
		return jsonObj
	}

	currentObj := jsonObj

	for i := 1; i < len(paths); i++ {
		value := currentObj[paths[i]]

		switch value.(type) {
		case map[string]any:

			if i == len(paths)-1 {
				return value
			}

			currentObj = value.(map[string]any)
			continue
		default:
			return value
		}

	}

	return nil
}
