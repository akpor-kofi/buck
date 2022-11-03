package buckis

import "strings"

// GraphCreate implement using the sorted sets following the hexastore principle
const graphKey = "graph:relations"

// graph relationship
type graphRel struct {
	subject   string
	predicate string
	object    string
	relations []string
}

func newGraphRel(s, p, o string) *graphRel {
	var relations []string

	spo := "spo:" + s + ":" + p + ":" + o
	sop := "sop:" + s + ":" + o + ":" + p
	ops := "ops:" + o + ":" + p + ":" + s
	osp := "osp:" + o + ":" + s + ":" + p
	pso := "pso:" + p + ":" + s + ":" + o
	pos := "pos:" + p + ":" + o + ":" + s

	relations = append(relations, spo, sop, ops, osp, pso, pos)

	return &graphRel{
		subject:   s,
		predicate: p,
		object:    o,
		relations: relations,
	}
}

func (d *dict) RAdd(s, p, o string) error {
	gr := newGraphRel(s, p, o)

	for _, r := range gr.relations {
		d.ZAdd(graphKey, r, 0)
	}

	return nil
}

func (d *dict) RMatch(s, p string) (result []string, err error) {
	res, err := d.ZRangeByLex(graphKey, 0, "[spo:"+s+":"+p, "+")
	if err != nil {
		return result, err
	}

	for _, r := range res {
		result = append(result, strings.Split(r, ":")[3])
	}

	return result, nil
}
