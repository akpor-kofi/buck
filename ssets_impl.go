package buckis

func (d *dict) zadd(flag int, key string, member string, score int) error {
	z := NewZ(member, score)
	ssde, err := d.sortedSetLookup(key)

	if err != nil {
		// there is no z in place for this key
		ssetHash := d.zhash(z.member)

		// insert z into hash table
		zd := newZdict()

		if zd.zht[ssetHash] == nil {
			zd.zht[ssetHash] = &ZEntry{z.score, z.member, nil}
		} else {
			zd.zht[ssetHash] = &ZEntry{z.score, z.member, zd.zht[ssetHash]}
		}

		// link score to member on skip list
		tree := lexTree{}
		zd.skiplist.Set(z.score, tree.Add(z.member))

		de := &dictEntry{
			key:    key,
			values: zd,
		}

		// store the entry
		d.ht[SortedSet][d.hash(key)] = de

		if flag == SAVE {
			d.waiter.Add(1)
			go func(ch chan command) {
				defer d.waiter.Done()
				cmd := newCommand(ZADD, key, member, score)
				ch <- *cmd
			}(d.commandChan)
			d.waiter.Wait()
		}

		return nil
	}

	zd := ssde.values.(*zdict)

	i := d.zhash(z.member)

	if zd.zht[i] == nil {
		zd.zht[i] = &ZEntry{z.score, z.member, nil}
	} else {
		zd.zht[i] = &ZEntry{z.score, z.member, zd.zht[i]}
	}

	// first check if score is already in
	if val, exists := zd.skiplist.GetValue(z.score); exists {
		//TODO: sort the linked list some how maybe using skiplist
		zd.skiplist.Set(z.score, val.(*lexTree).Add(z.member))
	} else {

		// create a new lex tree

		lt := &lexTree{}
		zd.skiplist.Set(z.score, lt.Add(z.member))
	}

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(ZADD, key, member, score)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}

func (d *dict) zincrby(flag int, key, member string, incr int) error {
	// 1) get the current score of the member
	previousScore := d.ZScore(key, member)

	// 2) delete the previous member entry
	err := d.ZRem(key, member)
	if err != nil {
		return err
	}

	// 3) insert a new member entry with the "incr" added to the previous score
	currentScore := previousScore + incr

	err = d.ZAdd(key, member, currentScore)
	if err != nil {
		return err
	}

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(ZINCRBY, key, member, incr)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}

func (d *dict) zrem(flag int, key string, member string) error {
	ssde, err := d.sortedSetLookup(key)

	if err == ErrSortedSetNotFound {
		return err
	}

	zd := ssde.values.(*zdict)

	// delete member from hashtable
	i := d.zhash(member)

	// delete member from skiplist
	if zd.zht[d.zhash(member)] == nil {
		return ErrEntryNotFound
	}

	score, err := findMembersScore(zd.zht[d.zhash(member)], member)

	if err != nil {
		return err
	}

	el := zd.skiplist.Get(score)

	switch {
	case el != nil:
		tree := el.Value.(*lexTree)
		tree, err = tree.remove(member)

		if err == FlagRemoveElement {
			zd.skiplist.RemoveElement(el)
		} else {
			zd.skiplist.Set(score, tree)
		}
	}

	// traverse linked list
	currentNode := zd.zht[i]
	var tempPrevNode *ZEntry

	if currentNode == nil {
		return ErrHashNotFound
	}

	for {
		if currentNode.member == member {
			// delete node
			if tempPrevNode == nil {
				zd.zht[i] = nil
			} else {
				tempPrevNode.next = currentNode.next
			}

			break
		}

		if currentNode.next == nil {
			return ErrHashNotFound
		}

		tempPrevNode = currentNode
		currentNode = currentNode.next
	}

	//elm := zd.skiplist.Get(105)
	//fmt.Println(elm.Value.(*ZNode).member)
	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(ZREM, key, member)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}

	return nil
}
