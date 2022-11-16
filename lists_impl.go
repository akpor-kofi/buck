package buckis

func (d *dict) rPush(flag int, key, item string) (index int, err error) {
	// 1) check if the list exists first
	// lde => list dict meta

	lde, err := d.listLookup(key)
	if err != nil {
		// make a list for it and add it to

		de := &DictEntry{
			Key:    key,
			Values: []string{item},
		}

		i := d.hash(key)
		if d.Ht[List][i] == nil {
			de.Next = nil
		} else {
			de.Next = d.Ht[List][i]
		}

		d.Ht[List][i] = de

		if flag == SAVE {
			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(RPUSH, key, item)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}

		return 0, nil
	}

	// queue implementation
	list := lde.Values.([]string)
	list = append(list, item)

	lde.Values = list

	// set index
	index = len(list) - 1

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(RPUSH, key, item)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return
}

func (d *dict) lPush(flag int, key, item string) (index int, err error) {
	lde, err := d.listLookup(key)
	if err != nil {
		return 0, err
	}

	// can only push if there space that is rear is greater than 0
	list := lde.Values.([]string)

	list = append([]string{item}, list...)
	lde.Values = list

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(LPUSH, key, item)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return
}

func (d *dict) lPop(flag int, key string) (index int, err error) {
	// 1) check if the list exists first
	// lde => list dict meta

	lde, err := d.listLookup(key)
	if err != nil {
		return 0, err
	}

	// queue implementation
	list := lde.Values.([]string)

	lde.Values = list[1:]

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(LPOP, key)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return
}

func (d *dict) rPop(flag int, key string) (index int, err error) {
	// 1) check if the list exists first
	// lde => list dict meta

	lde, err := d.listLookup(key)
	if err != nil {
		return
	}

	list := lde.Values.([]string)
	index = len(list) - 2
	lde.Values = list[:index]

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(RPOP, key)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return
}
