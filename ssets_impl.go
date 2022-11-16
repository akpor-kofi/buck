package buckis

func (d *dict) zadd(flag int, key string, member string, score int) error {
	z := NewZ(member, score)
	ssde, err := d.sortedSetLookup(key)

	if err != nil {
		// there is no z in place for this key
		ssetHash := d.zhash(z.member)

		// insert z into hash table
		zd := newZdict()

		if zd.Zht[ssetHash] == nil {
			zd.Zht[ssetHash] = &ZEntry{z.score, z.member, nil}
		} else {
			zd.Zht[ssetHash] = &ZEntry{z.score, z.member, zd.Zht[ssetHash]}
		}

		// link score to member on skip list
		tree := lexTree{}
		zd.Skiplist.Set(z.score, tree.Add(z.member))

		de := &DictEntry{
			Key:    key,
			Values: zd,
		}

		// store the entry
		d.Ht[SortedSet][d.hash(key)] = de

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

	zd := ssde.Values.(*Zdict)

	i := d.zhash(z.member)

	if zd.Zht[i] == nil {
		zd.Zht[i] = &ZEntry{z.score, z.member, nil}
	} else {
		zd.Zht[i] = &ZEntry{z.score, z.member, zd.Zht[i]}
	}

	// first check if score is already in
	if val, exists := zd.Skiplist.GetValue(z.score); exists {
		//TODO: sort the linked list some how maybe using Skiplist
		zd.Skiplist.Set(z.score, val.(*lexTree).Add(z.member))
	} else {

		// create a new lex tree

		lt := &lexTree{}
		zd.Skiplist.Set(z.score, lt.Add(z.member))
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
	previousScore, err := d.ZScore(key, member)
	if err != nil {
		return err
	}

	// 2) delete the previous member entry
	err = d.ZRem(key, member)
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

	zd := ssde.Values.(*Zdict)

	// delete member from hashtable
	i := d.zhash(member)

	// delete member from Skiplist
	if zd.Zht[d.zhash(member)] == nil {
		return ErrEntryNotFound
	}

	score, err := findMembersScore(zd.Zht[d.zhash(member)], member)

	if err != nil {
		return err
	}

	el := zd.Skiplist.Get(score)

	switch {
	case el != nil:
		tree := el.Value.(*lexTree)
		tree, err = tree.remove(member)

		if err == FlagRemoveElement {
			zd.Skiplist.RemoveElement(el)
		} else {
			zd.Skiplist.Set(score, tree)
		}
	}

	// traverse linked list
	currentNode := zd.Zht[i]
	var tempPrevNode *ZEntry

	if currentNode == nil {
		return ErrHashNotFound
	}

	for {
		if currentNode.member == member {
			// delete node
			if tempPrevNode == nil {
				zd.Zht[i] = nil
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

	//elm := zd.Skiplist.Get(105)
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
