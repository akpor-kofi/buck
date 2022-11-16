package buckis

import (
	"errors"
	"github.com/aviddiviner/go-murmur"
	"github.com/huandu/skiplist"
)

var ErrSortedSetNotFound = errors.New("sorted set not found")
var ErrRangeOutOfBounds = errors.New("range out of bounds")
var ErrScoreNotFound = errors.New("score not found")
var ErrMemberNotFound = errors.New("member not found")
var FlagRemoveElement = errors.New("remove element")

type ZEntry struct {
	score  int
	member string
	next   *ZEntry
}

type zdict struct {
	skiplist *skiplist.SkipList
	zht      [50]*ZEntry
}

type Z struct {
	member string
	score  int
}

func newZdict() *zdict {
	sl := skiplist.New(skiplist.Int)

	return &zdict{
		skiplist: sl,
		zht:      [50]*ZEntry{},
	}
}

func NewZ(member string, score int) *Z {
	return &Z{member: member, score: score}
}

func (d *dict) ZAdd(key string, member string, score int) error {
	return d.zadd(SAVE, key, member, score)
}

func (d *dict) ZScore(key string, member string) int {
	ssde, err := d.sortedSetLookup(key)

	if err == ErrSortedSetNotFound {
		return 0
	}

	zd := ssde.Values.(*zdict)

	currentZEntry := zd.zht[d.zhash(member)]

	for {
		if currentZEntry.member == member {
			return currentZEntry.score
		}

		if currentZEntry.next == nil {
			return 0 //meaning could not find member
		}

		currentZEntry = currentZEntry.next
	}
}

// ZRange lb is lower bound, ub is upper bound
func (d *dict) ZRange(key string, lb, ub int) (members []string, err error) {

	ssde, err := d.sortedSetLookup(key)

	if err == ErrSortedSetNotFound {
		return
	}

	zd := ssde.Values.(*zdict)
	firstElementIndex := 0
	numIterations := lb - firstElementIndex
	lowerBound := zd.skiplist.Front()

	// Fix bug: if you put a very large number it returns the last element
	for i := 0; i < numIterations; i++ {
		if lowerBound.Next() != nil {
			lowerBound = lowerBound.Next()
		}
	}

	for i := lb; i <= ub; i++ {
		// loop through znode linked list
		list := lowerBound.Value.(*lexTree)

		members = append(members, list.toSlice()...)

		if lowerBound.Next() == nil {
			break
		}
		lowerBound = lowerBound.Next()
	}

	return members, nil
}

func (d *dict) ZRangeByScore(key string, lb, ub int) (members []string, err error) {
	if lb > ub {
		err = ErrRangeOutOfBounds
		return
	}

	ssde, err := d.sortedSetLookup(key)

	if err == ErrSortedSetNotFound {
		return
	}

	zd := ssde.Values.(*zdict)

	for i := lb; i < ub; i++ {
		if val, ok := zd.skiplist.GetValue(i); !ok {
			continue
		} else {
			list := val.(*lexTree)

			members = append(members, list.toSlice()...)
		}
	}

	return
}

// ZRangeStore TODO: dk - new store key, sk -source key
func (d *dict) ZRangeStore(dk, sk string, lb, ub int) error {
	// TODO: implement this

	return nil
}

func (d *dict) ZRangeByLex(key string, score int, lb, ub string) (results []string, err error) {
	ssde, err := d.sortedSetLookup(key)

	if err == ErrSortedSetNotFound {
		return
	}

	zd := ssde.Values.(*zdict)
	val, exists := zd.skiplist.GetValue(score)

	if !exists {
		return
	}

	tree := val.(*lexTree)

	results = tree.findByLex(lb, ub)

	return
}

func (d *dict) ZIncrBy(key, member string, incr int) error {
	return d.zincrby(SAVE, key, member, incr)
}

func (d *dict) ZUnion(dk, sk string) (members []string, err error) {
	// TODO: implement this
	return nil, nil
}

func (d *dict) ZRank(key string, member string) int {
	ssde, err := d.sortedSetLookup(key)

	if err == ErrSortedSetNotFound {
		return 0
	}

	// Algorithm
	// 1) get member score from hashtable
	score := d.ZScore(key, member)
	// 2) get score rank from hashtable
	zd := ssde.Values.(*zdict)

	rank, err := zd.getRank(score)
	if err != nil {
		return 0
	}

	return rank
}

func (d *dict) ZRem(key string, member string) error {
	return d.zrem(SAVE, key, member)
}

func (d *dict) ZRandMember(key string) (member string, err error) {
	// TODO: implement this

	return "", nil
}

func (zd *zdict) getRank(score int) (int, error) {
	count := 0

	zEl := zd.skiplist.Front()

	for {
		if zEl.Key() == score {
			return count, nil
		}

		if zEl.Next() == nil {
			return 0, ErrScoreNotFound
		}

		zEl = zEl.Next()
		count++
	}
}

func (d *dict) zhash(member string) int {

	b := []byte(member)

	h := murmur.MurmurHash2(b, 0)

	return int(h % 50)
}

func (d *dict) sortedSetLookup(key string) (*DictEntry, error) {
	i := d.hash(key)

	currentEntry := d.Ht[SortedSet][i]

	if currentEntry == nil {
		return &DictEntry{}, ErrSortedSetNotFound
	}

	for {
		if currentEntry.Key == key {
			return currentEntry, nil
		}

		if currentEntry.Next == nil {
			return &DictEntry{}, ErrSortedSetNotFound
		}

		currentEntry = currentEntry.Next
	}
}

func findMembersScore(entry *ZEntry, member string) (int, error) {
	currNode := entry

	for {
		if currNode.member == member {
			return currNode.score, nil
		}

		if currNode.next == nil {
			return 0, ErrMemberNotFound
		}

		currNode = currNode.next
	}
}
