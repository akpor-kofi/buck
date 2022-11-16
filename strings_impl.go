package buckis

func (d *dict) incrBy(flag int, key string, incr int) (int, error) {
	de, err := d.stringsLookup(key)

	if err != nil {
		return 0, err
	}

	switch de.Values.(type) {
	case int:
		val := de.Values.(int)
		val += incr
		de.Values = val

		// send command
		if flag == SAVE {
			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(INCRBY, key, incr)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}
		return val, nil

	// TODO: for some reason after unmarshalling an int it turns into a float
	case float64:
		// convert to int
		val := int(de.Values.(float64))
		val += incr
		de.Values = val

		return val, nil
	default:
		return 0, ErrNotAnInteger

	}
}

func (d *dict) set(flag int, key string, value any) error {
	i := d.hash(key)

	// if no set value before
	if d.Ht[Strings][i] == nil {
		d.Ht[Strings][i] = &DictEntry{
			Key:    key,
			Values: value,
			Next:   nil,
		}

		// send command
		if flag == SAVE {
			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(SET, key, value)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}

		return nil

	}

	de := &DictEntry{
		Key:    key,
		Values: value,
		Next:   d.Ht[Strings][i],
	}

	d.Ht[Strings][i] = de

	// send command
	if flag == SAVE {

		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(SET, key, value)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}
