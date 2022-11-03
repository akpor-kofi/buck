package buckis

import (
	"errors"
)

var ErrEntryNotFound = errors.New("entry not found")
var ErrNotAnInteger = errors.New("value not an integer")

func (d *dict) Set(key string, value any) error {
	return d.set(SAVE, key, value)
}

func (d *dict) Get(key string) (any, error) {
	de, err := d.stringsLookup(key)

	if err != nil {
		return "", err
	}

	return de.values, nil
}

func (d *dict) IncrBy(key string, incr int) (int, error) {
	return d.incrBy(SAVE, key, incr)
}

func (d *dict) stringsLookup(key string) (*dictEntry, error) {
	i := d.hash(key)

	currentEntry := d.ht[Strings][i]

	if currentEntry == nil {
		return &dictEntry{}, ErrEntryNotFound
	}

	for {
		if currentEntry.key == key {
			return currentEntry, nil
		}

		if currentEntry.next == nil {
			return &dictEntry{}, ErrEntryNotFound
		}

		currentEntry = currentEntry.next
	}
}
