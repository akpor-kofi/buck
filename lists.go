package buckis

import "errors"

func (d *dict) RPush(key, item string) (index int, err error) {
	return d.rPush(SAVE, key, item)
}

func (d *dict) LPush(key, item string) (index int, err error) {
	return d.lPush(SAVE, key, item)
}

func (d *dict) LPop(key string) (index int, err error) {
	return d.lPop(SAVE, key)
}

func (d *dict) RPop(key string) (index int, err error) {
	return d.rPop(SAVE, key)
}

// LRange ub - upper bound, lb - lower bound
func (d *dict) LRange(key string, ub int, lb int) (arr []string, err error) {

	lde, err := d.listLookup(key)
	if err != nil {
		return
	}

	lastPos := len(lde.values.([]string)) - 1
	if lb > lastPos {
		lb = lastPos
	}

	for i := ub; i <= lb; i++ {
		arr = append(arr, lde.values.([]string)[i])
	}

	return
}

func (d *dict) LTop(key string) (top string, err error) {
	lde, err := d.listLookup(key)
	if err != nil {
		return
	}

	top = lde.values.([]string)[len(lde.values.([]string))-1]

	return
}

func (d *dict) LLen(key string) (length int, err error) {
	lde, err := d.listLookup(key)
	if err != nil {
		return
	}

	length = len(lde.values.([]string))

	return
}

func (d *dict) LRevRange(key string, opts ...int) (arr []string, err error) {
	count := opts[0]

	lde, err := d.listLookup(key)
	if err != nil {
		return
	}

	if count == 0 {
		count = len(lde.values.([]string))
	}

	lastPos := len(lde.values.([]string)) - 1

	for i := lastPos; i >= 0; i-- {
		if count == 0 {
			break
		}
		arr = append(arr, lde.values.([]string)[i])
		count--
	}

	return
}

func (d *dict) LIndex(key, item string) (i int, err error) {
	lde, err := d.listLookup(key)
	if err != nil {
		return
	}

	lastPos := len(lde.values.([]string)) - 1

	for ; i <= lastPos; i++ {
		if lde.values.([]string)[i] == item {
			return
		}
	}

	err = errors.New("item not found")

	return
}

func (d *dict) listLookup(key string) (*dictEntry, error) {
	i := d.hash(key)

	currentEntry := d.ht[List][i]

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
