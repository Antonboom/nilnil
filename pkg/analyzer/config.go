package analyzer

import (
	"fmt"
	"sort"
	"strings"
)

func newDefaultCheckedTypes() checkedTypes {
	return checkedTypes{
		ptrType:   struct{}{},
		funcType:  struct{}{},
		ifaceType: struct{}{},
		mapType:   struct{}{},
		chanType:  struct{}{},
	}
}

const separator = ','

const (
	ptrType   = "ptr"
	funcType  = "func"
	ifaceType = "iface"
	mapType   = "map"
	chanType  = "chan"
)

var knownTypes = []string{ptrType, funcType, ifaceType, mapType, chanType}

type checkedTypes map[string]struct{}

func (c checkedTypes) Contains(t string) bool {
	_, ok := c[t]
	return ok
}

func (c checkedTypes) String() string {
	result := make([]string, 0, len(c))
	for t := range c {
		result = append(result, t)
	}

	sort.Strings(result)
	return strings.Join(result, string(separator))
}

func (c checkedTypes) Set(s string) error {
	types := strings.FieldsFunc(s, func(c rune) bool { return c == separator })
	if len(types) == 0 {
		return nil
	}

	c.disableAll()
	for _, t := range types {
		switch t {
		case ptrType, funcType, ifaceType, mapType, chanType:
			c[t] = struct{}{}
		default:
			return fmt.Errorf("unknown checked type name %q (see help)", t)
		}
	}

	return nil
}

func (c checkedTypes) disableAll() {
	for k := range c {
		delete(c, k)
	}
}
