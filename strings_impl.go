package buckis

func (d *dict) incrBy(flag int, key string, incr int) (int, error) {
	de, err := d.stringsLookup(key)

	if err != nil {
		return 0, err
	}

	switch de.values.(type) {
	case int:
		val := de.values.(int)
		val += incr
		de.values = val

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
		val := int(de.values.(float64))
		val += incr
		de.values = val

		return val, nil
	default:
		return 0, ErrNotAnInteger

	}
}

func (d *dict) set(flag int, key string, value any) error {
	i := d.hash(key)

	// if no set value before
	if d.ht[Strings][i] == nil {
		d.ht[Strings][i] = &dictEntry{
			key:    key,
			values: value,
			next:   nil,
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

	de := &dictEntry{
		key:    key,
		values: value,
		next:   d.ht[Strings][i],
	}

	d.ht[Strings][i] = de

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
