package buckis

//
//import (
//	"fmt"
//	"github.com/aviddiviner/go-murmur"
//)
//
//type listDictG[T any] struct {
//	ht [100]*listDictEntryG[T]
//}
//
//type listDictEntryG[T any] struct {
//	key  string
//	meta listMeta[T]
//	next *listDictEntryG[T]
//}
//
//func NewListDictG[T any]() *listDictG[T] {
//	return &listDictG[T]{
//		ht: [100]*listDictEntryG[T]{},
//	}
//}
//
//type list[T any] []T
//
//type listMeta[T any] struct {
//	front int
//	rear  int
//	list  list[T]
//}
//
//func newListMeta[T any]() *listMeta[T] {
//	return &listMeta[T]{
//		front: 0,
//		rear:  0,
//		list:  make(list[T], 20),
//	}
//}
//
//func (ld *listDictG[T]) RPush(key string, item T) (index int, err error) {
//	// 1) check if the list exists first
//	// lde => list dict meta
//
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		// make a list for it and add it to
//
//		de := &listDictEntryG[T]{
//			key:  key,
//			meta: *newListMeta[T](),
//		}
//
//		front := de.meta.front
//		de.meta.list[front] = item
//		de.meta.front++
//
//		i := ld.hash(key)
//		if ld.ht[i] == nil {
//			de.next = nil
//		} else {
//			de.next = ld.ht[i]
//		}
//
//		ld.ht[i] = de
//
//		return 0, nil
//	}
//
//	// queue implementation
//	front := lde.meta.front
//	lde.meta.list[front] = item
//	index = front
//
//	// increment front in meta
//	lde.meta.front++
//
//	return
//}
//
//func (ld *listDictG[T]) LPush(key string, item T) (index int, err error) {
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		return 0, err
//	}
//
//	rear := &lde.meta.rear
//
//	// can only push if there space that is rear is greater than 0
//	if *rear < 1 {
//		// TODO: might implement shifting later
//		return 0, fmt.Errorf("No space bruh")
//	}
//	lde.meta.rear--
//
//	fmt.Println(*rear)
//
//	lde.meta.list[*rear] = item
//
//	index = *rear
//	return
//}
//
//func (ld *listDictG[T]) LPop(key string) (index int, err error) {
//	// 1) check if the list exists first
//	// lde => list dict meta
//
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		return 0, err
//	}
//
//	// queue implementation
//	lde.meta.rear++
//
//	index = lde.meta.rear
//	return
//}
//
//func (ld *listDictG[T]) RPop(key string) (index int, err error) {
//	// 1) check if the list exists first
//	// lde => list dict meta
//
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		return
//	}
//
//	lde.meta.front--
//
//	index = lde.meta.front
//
//	return
//}
//
//// LRange ub - upper bound, lb - lower bound
//func (ld *listDictG[T]) LRange(key string, ub int, lb int) (arr []T, err error) {
//
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		return
//	}
//
//	if lb > lde.meta.front {
//		lb = lde.meta.front
//	}
//
//	for i := ub; i <= lb; i++ {
//		arr = append(arr, lde.meta.list[i])
//	}
//
//	return
//}
//
//func (ld *listDictG[T]) LTop(key string) (top T, err error) {
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		return
//	}
//
//	top = lde.meta.list[lde.meta.front]
//
//	return
//}
//
//func (ld *listDictG[T]) LLen(key string) (length int, err error) {
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		return
//	}
//
//	front := lde.meta.front
//	rear := lde.meta.rear
//
//	queue := lde.meta.list[rear:front]
//	length = len(queue)
//
//	return
//}
//
//func (ld *listDictG[T]) LRevRange(key string, opts ...int) (arr []T, err error) {
//	count := opts[0]
//
//	lde, err := ld.listLookup(key)
//	if err != nil {
//		return
//	}
//
//	front := lde.meta.front
//	rear := lde.meta.rear
//
//	for i := front; i >= rear; i-- {
//		if count == 0 {
//			break
//		}
//		arr = append(arr, lde.meta.list[i])
//		count--
//	}
//
//	return
//}
//
//func (ld *listDictG[T]) listLookup(key string) (*listDictEntryG[T], error) {
//	i := ld.hash(key)
//
//	currentEntry := ld.ht[i]
//
//	if currentEntry == nil {
//		return &listDictEntryG[T]{}, ErrSetNotFound
//	}
//
//	for {
//		if currentEntry.key == key {
//			return currentEntry, nil
//		}
//
//		if currentEntry.next == nil {
//			return &listDictEntryG[T]{}, ErrSetNotFound
//		}
//
//		currentEntry = currentEntry.next
//	}
//
//}
//
//func (ld *listDictG[T]) hash(key string) uint32 {
//	b := []byte(key)
//
//	h := murmur.MurmurHash2(b, 0)
//
//	return h % uint32(len(ld.ht))
//}
