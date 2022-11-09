package buckis

func (d *dict) bfAdd(flag int, key, value string) error {
	bfde, err := d.bloomFilterLookup(key)
	idx := d.hash(key)

	if err != nil {
		// add a new bloom filter
		bitArray := [100]bool{}

		for k := 1; k <= numberOfHashFuncs; k++ {
			bitArray[bloomHashFunc(value, k)] = true
		}

		de := &dictEntry{
			key:    key,
			values: bitArray,
			next:   d.ht[Bloom][idx],
		}

		d.ht[Bloom][idx] = de

		if flag == SAVE {
			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(BFADD, key, value)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}

		return nil
	}

	bitArray := bfde.values.([100]bool)

	for k := 1; k <= numberOfHashFuncs; k++ {
		bitArray[bloomHashFunc(value, k)] = true
	}

	bfde.values = bitArray

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(BFADD, key, value)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}
