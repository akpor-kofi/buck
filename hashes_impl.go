package buckis

func (d *dict) hset(flag int, key string, hashes ...any) error {
	i := d.hash(key)

	// need to lookup and check already existing hash
	hashMap := make(map[string]any)

	for k, j := 0, 1; k < len(hashes); k += 2 {

		hashMap[hashes[k].(string)] = hashes[j]

		j += 2
	}

	if d.ht[Hashes][i] == nil {
		d.ht[Hashes][i] = &dictEntry{
			key:    key,
			values: hashMap,
			next:   nil,
		}

		if flag == SAVE {

			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(HSET, key, hashes)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}

		return nil
	}

	de := &dictEntry{
		key:    key,
		values: hashMap,
		next:   d.ht[Hashes][i],
	}

	d.ht[Hashes][i] = de

	if flag == SAVE {

		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(HSET, key, hashes)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}
