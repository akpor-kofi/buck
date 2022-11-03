package buckis

import (
	"errors"
	"math/rand"
	"time"
)

var ErrSetNotFound = errors.New("set not found")

type void struct{}

func (d *dict) SAdd(key, value string) error {
	return d.sadd(SAVE, key, value)
}

func (d *dict) SIsMember(key, value string) (bool, error) {
	de, err := d.setLookup(key)

	if err != nil {
		return false, err
	}

	if _, exists := de.values.(map[string]void)[value]; !exists {
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

	for v, _ := range de.values.(map[string]void) {
		members = append(members, v)
	}

	return members, nil
}

func (d *dict) SCard(key string) int {
	de, err := d.setLookup(key)
	if err != nil {
		return 0
	}

	return len(de.values.(map[string]void))
}

func (d *dict) SRandMember(key string) string {
	de, err := d.setLookup(key)
	if err != nil {
		return ""
	}

	rand.Seed(time.Now().UnixMilli())
	random := rand.Intn(len(de.values.(map[string]void))) + 1
	i := 0

	for v, _ := range de.values.(map[string]void) {
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

	unionMap := firstSet.values.(map[string]void)

	for v, _ := range secondSet.values.(map[string]void) {
		unionMap[v] = void{}
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

	for v, _ := range firstSet.values.(map[string]void) {
		if _, exists := secondSet.values.(map[string]void)[v]; exists {
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

	for v, _ := range firstSet.values.(map[string]void) {
		if _, exists := secondSet.values.(map[string]void)[v]; !exists {
			diff = append(diff, v)
		}
	}

	// naive implementation
	for v, _ := range secondSet.values.(map[string]void) {
		if _, exists := firstSet.values.(map[string]void)[v]; !exists {
			diff = append(diff, v)
		}
	}

	return diff, nil
}

func (d *dict) setLookup(key string) (*dictEntry, error) {
	i := d.hash(key)

	currentEntry := d.ht[Set][i]

	if currentEntry == nil {
		return &dictEntry{}, ErrSetNotFound
	}

	for {
		if currentEntry.key == key {
			return currentEntry, nil
		}

		if currentEntry.next == nil {
			return &dictEntry{}, ErrSetNotFound
		}

		currentEntry = currentEntry.next
	}
}
