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

		de := &DictEntry{
			Key:    key,
			Values: bitArray,
			Next:   d.Ht[Bloom][idx],
		}

		d.Ht[Bloom][idx] = de

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

	bitArray := bfde.Values.([100]bool)

	for k := 1; k <= numberOfHashFuncs; k++ {
		bitArray[bloomHashFunc(value, k)] = true
	}

	bfde.Values = bitArray

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
