package buckis

func (d *dict) sadd(flag int, key, value string) error {
	de, err := d.setLookup(key)

	if err == nil {
		de.Values.(map[string]void)[value] = void{}

		if flag == SAVE {

			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(SADD, key, value)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}

		return nil
	}

	i := d.hash(key)

	de = &DictEntry{
		Key: key,
		Values: map[string]void{
			value: {},
		},
	}

	if d.Ht[Set][i] != nil {
		de.Next = d.Ht[Set][i]
	}

	d.Ht[Set][i] = de

	if flag == SAVE {

		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(SADD, key, value)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}

func (d *dict) srem(flag int, key, value string) error {
	de, err := d.setLookup(key)

	if err != nil {
		// key does not exist so just return
		return err
	}

	delete(de.Values.(map[string]void), value)

	if flag == SAVE {

		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(SREM, key, value)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}

func (d *dict) smove(flag int, sk, dk, value string) error {
	sde, err := d.setLookup(sk)
	if err != nil {
		return err
	}

	delete(sde.Values.(map[string]void), value)

	// add the value to the set
	dde, err := d.setLookup(dk)
	if err != nil {
		return err
	}

	dde.Values.(map[string]void)[value] = void{}

	if flag == SAVE {

		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(SMOVE, sk, dk, value)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}
