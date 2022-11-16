package buckis

import (
	"errors"
	"math/rand"
	"time"
)

var ErrSetNotFound = errors.New("set not found")

type Void struct{}

func (d *dict) SAdd(key, value string) error {
	return d.sadd(SAVE, key, value)
}

func (d *dict) SIsMember(key, value string) (bool, error) {
	de, err := d.setLookup(key)

	if err != nil {
		return false, err
	}

	if _, exists := de.Values.(map[string]Void)[value]; !exists {
		return false, nil
	}

	return true, nil
}

func (d *dict) SRem(key, value string) error {
	return d.srem(SAVE, key, value)
}

func (d *dict) SMembers(key string) ([]string, error) {
	de, err := d.setLookup(key)
	if err != nil {
		return []string{}, err
	}

	var members []string

	for v, _ := range de.Values.(map[string]Void) {
		members = append(members, v)
	}

	return members, nil
}

func (d *dict) SCard(key string) int {
	de, err := d.setLookup(key)
	if err != nil {
		return 0
	}

	return len(de.Values.(map[string]Void))
}

func (d *dict) SRandMember(key string) string {
	de, err := d.setLookup(key)
	if err != nil {
		return ""
	}

	rand.Seed(time.Now().UnixMilli())
	random := rand.Intn(len(de.Values.(map[string]Void))) + 1
	i := 0

	for v, _ := range de.Values.(map[string]Void) {
		i++
		if i == random {
			return v
		}
	}

	return ""
}

// SMove sk - src set, dk - dest set
func (d *dict) SMove(sk, dk, value string) error {
	return d.smove(SAVE, sk, dk, value)
}

func (d *dict) SUnion(s1, s2 string) ([]string, error) {
	firstSet, err := d.setLookup(s1)

	secondSet, err := d.setLookup(s2)

	if err != nil {
		return nil, err
	}

	unionMap := firstSet.Values.(map[string]Void)

	for v, _ := range secondSet.Values.(map[string]Void) {
		unionMap[v] = Void{}
	}

	var union []string

	for v, _ := range unionMap {
		union = append(union, v)
	}

	return union, nil
}

func (d *dict) SInter(s1, s2 string) ([]string, error) {
	firstSet, err := d.setLookup(s1)

	secondSet, err := d.setLookup(s2)

	if err != nil {
		return nil, err
	}

	var inter []string

	for v, _ := range firstSet.Values.(map[string]Void) {
		if _, exists := secondSet.Values.(map[string]Void)[v]; exists {
			inter = append(inter, v)
		}
	}

	return inter, nil
}

func (d *dict) SDiff(s1, s2 string) ([]string, error) {
	firstSet, err := d.setLookup(s1)

	secondSet, err := d.setLookup(s2)

	if err != nil {
		return nil, err
	}

	var diff []string

	for v, _ := range firstSet.Values.(map[string]Void) {
		if _, exists := secondSet.Values.(map[string]Void)[v]; !exists {
			diff = append(diff, v)
		}
	}

	// naive implementation
	for v, _ := range secondSet.Values.(map[string]Void) {
		if _, exists := firstSet.Values.(map[string]Void)[v]; !exists {
			diff = append(diff, v)
		}
	}

	return diff, nil
}

func (d *dict) setLookup(key string) (*DictEntry, error) {
	i := d.hash(key)

	currentEntry := d.Ht[Set][i]

	if currentEntry == nil {
		return &DictEntry{}, ErrSetNotFound
	}

	for {
		if currentEntry.Key == key {
			return currentEntry, nil
		}

		if currentEntry.Next == nil {
			return &DictEntry{}, ErrSetNotFound
		}

		currentEntry = currentEntry.Next
	}
}
