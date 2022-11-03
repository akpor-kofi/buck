package buckis

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var ErrMemberNotInTree = errors.New("member not found")

type lexTree []string

func (lt *lexTree) Add(s string) *lexTree {
	*lt = append(*lt, s)
	lt.sort()

	return lt
}

func (lt *lexTree) sort() {
	sort.Strings(*lt)
}

func (lt *lexTree) len() int {
	return len(*lt)
}

func (lt *lexTree) toSlice() []string {
	return *lt
}

func (lt *lexTree) remove(s string) (*lexTree, error) {
	i := sort.SearchStrings(*lt, s)

	if len(*lt) <= 1 {
		return nil, FlagRemoveElement
	}

	if i == len(*lt) {
		return nil, ErrMemberNotInTree
	}

	temp := *lt

	temp = append(temp[:i], temp[i+1:]...)

	*lt = temp

	return lt, nil
}

func (lt *lexTree) findByLex(args ...string) (results []string) {
	// [banana + count 2
	lb := args[0] //lower bound
	ub := args[1] //upper bound

	arr := *lt

	fmt.Println(arr)

	if strings.HasPrefix(lb, "[") {
		query := strings.TrimLeft(lb, "[")

		startingIndex := sort.Search(lt.len(), func(i int) bool {
			fmt.Println(arr[i])
			return arr[i] > query
		})

		for i := startingIndex; i < len(arr); i++ {
			// should be a place to check if we want strict result
			if strings.HasPrefix(arr[i], query) {
				results = append(results, arr[i])
			}

			if strings.HasPrefix(arr[i], strings.TrimLeft(ub, "(")) {
				break
			}
		}
	}

	return
}
