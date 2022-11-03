package buckis

import "errors"

var ErrHashNotFound = errors.New("hash not found")

func (d *dict) HSet(key string, hashes ...string) error {
	return d.hset(SAVE, key, hashes...)
}

func (d *dict) HGetAll(key string) (map[string]any, error) {
	de, err := d.hashesLookup(key)

	if err != nil {
		return nil, err
	}

	return de.values.(map[string]any), nil
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

	return 0, nil
}
