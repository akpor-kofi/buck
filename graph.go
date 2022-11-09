package buckis

import (
	"sort"
	"strings"
)

func (d *dict) gAdd(flag int, subject, predicate, object string) {
	spo := "spo:" + subject + ":" + predicate + ":" + object
	sop := "sop:" + subject + ":" + object + ":" + predicate
	ops := "ops:" + object + ":" + predicate + ":" + subject
	osp := "osp:" + object + ":" + subject + ":" + predicate
	pso := "pso:" + predicate + ":" + subject + ":" + object
	pos := "pos:" + predicate + ":" + object + ":" + subject

	d.hexastore = append(d.hexastore, spo, sop, ops, osp, pso, pos)

	// sort the store
	sort.Strings(d.hexastore)

	if flag == SAVE {
		d.waiter.Add(1)
		go func(ch chan command) {
			defer d.waiter.Done()
			cmd := newCommand(GADD, subject, predicate, object)
			ch <- *cmd
		}(d.commandChan)
		d.waiter.Wait()
	}
}

func (d *dict) GAdd(subject, predicate, object string) {
	d.gAdd(SAVE, subject, predicate, object)
}

// how many nodes are connected to a predicate
// do s p 0

func (d *dict) GRel(subject, object string) (result []string, err error) {
	// checking sop
	triple := "sop:" + subject + ":" + object

	index := sort.Search(len(d.hexastore), func(i int) bool {
		return d.hexastore[i] > triple
	})

	for i := index; strings.HasPrefix(d.hexastore[i], triple); i++ {
		result = append(result, strings.Split(d.hexastore[i], ":")[3])
	}

	return
}

func (d *dict) GMatch(subject, predicate string) (result []string, err error) {
	triple := "spo:" + subject + ":" + predicate

	index := sort.Search(len(d.hexastore), func(i int) bool {
		return d.hexastore[i] > triple
	})

	for i := index; strings.HasPrefix(d.hexastore[i], triple); i++ {
		result = append(result, strings.Split(d.hexastore[i], ":")[3])
	}

	return
}
