package redis

import (
	"container/list"
	"github.com/pkg/errors"
)

const ListType = uint64(1)
const ListTypeFancy = "list"

var _ Item = (*List)(nil)

type List struct {
	goList *list.List
}

func NewList() *List {
	return &List{goList: list.New()}
}

func (l *List) Value() interface{} {
	return l.goList
}

func (l *List) Type() uint64 {
	return ListType
}

func (l *List) TypeFancy() string {
	return ListTypeFancy
}

func (l *List) OnDelete(key *string, db *RedisDb) {
	panic("implement me")
}

// LLen returns number of elements.
func (l *List) LLen() int {
	return l.goList.Len()
}

// LPush returns the length of the list after the push operation.
func (l *List) LPush(values ...*string) int {
	for _, v := range values {
		l.goList.PushFront(*v)
	}
	return l.LLen()
}

// RPush returns the length of the list after the push operation.
func (l *List) RPush(values ...*string) int {
	for _, v := range values {
		l.goList.PushBack(*v)
	}
	return l.LLen()
}

// LInsert see redis doc
func (l *List) LInsert(isBefore bool, pivot, value *string) int {
	for e := l.goList.Front(); e.Next() != nil; e = e.Next() {
		if *vts(e) == *pivot {
			if isBefore {
				l.goList.InsertBefore(*value, e)
			} else {
				l.goList.InsertAfter(*value, e)
			}
			return l.LLen()
		}
	}
	return -1
}

// LPop returns popped value and false -
// returns true if list is now emptied so the key can be deleted.
func (l *List) LPop() (*string, bool) {
	if e := l.goList.Front(); e == nil {
		return nil, true
	} else {
		l.goList.Remove(e)
		return vts(e), false
	}
}

// RPop returns popped value and false -
// returns true if list is now emptied so the key can be deleted.
func (l *List) RPop() (*string, bool) {
	if e := l.goList.Back(); e == nil {
		return nil, true
	} else {
		l.goList.Remove(e)
		return vts(e), false
	}
}

// LRem see redis doc
func (l *List) LRem(count int, value *string) int {
	// count > 0: Remove elements equal to value moving from head to tail.
	// count < 0: Remove elements equal to value moving from tail to head.
	// count = 0: Remove all elements equal to value.
	var rem int
	if count >= 0 {
		for e := l.goList.Front(); e.Next() != nil; {
			if *vts(e) == *value {
				r := e
				e = e.Next()
				l.goList.Remove(r)
				rem++
				if count != 0 && rem == count {
					break
				}
			} else {
				e = e.Next()
			}
		}
	} else if count < 0 {
		count = abs(count)
		for e := l.goList.Back(); e.Prev() != nil; {
			if *vts(e) == *value {
				r := e
				e = e.Prev()
				l.goList.Remove(r)
				rem++
				if count != 0 && rem == count {
					break
				}
			} else {
				e = e.Prev()
			}
		}
	}
	return rem
}

// LSet see redis doc
func (l *List) LSet(index int, value *string) error {
	e := atIndex(index, l.goList)
	if e == nil {
		return errors.New("index out of range")
	}
	e.Value = *value
	return nil
}

// LIndex see redis doc
func (l *List) LIndex(index int) (*string, error) {
	e := atIndex(index, l.goList)
	if e == nil {
		return nil, errors.New("index out of range")
	}
	return vts(e), nil
}

// LRange see redis doc
func (l *List) LRange(start int, end int) []string {
	values := make([]string, 0)
	// from index to index
	from, to := startEndIndexes(start, end, l.LLen())
	if from > to {
		return values
	}
	// get start element
	e := atIndex(from, l.goList)
	if e == nil { // shouldn't happen
		return values
	}
	// fill with values
	values = append(values, *vts(e))
	for i := 0; i < to; i++ {
		e = e.Next()
		values = append(values, *vts(e))
	}
	return values
}

// LTrim see redis docs - returns true if list is now emptied so the key can be deleted.
func (l *List) LTrim(start int, end int) bool {
	// from index to index
	from, to := startEndIndexes(start, end, l.LLen())
	if from > to {
		l.goList.Init()
		return true
	}
	// trim before
	if from > 0 {
		i := 0
		e := l.goList.Front()
		for e != nil && i < from {
			del := e
			e = e.Next()
			l.goList.Remove(del)
			i++
		}
	}
	// trim after
	if to < l.LLen() {
		i := l.LLen()
		e := l.goList.Back()
		for e != nil && i > to {
			del := e
			e = e.Prev()
			l.goList.Remove(del)
			i--
		}
	}
	return false
}

func startEndIndexes(start, end int, listLen int) (int, int) {
	if end > listLen-1 {
		end = listLen - 1
	}
	return toIndex(start, listLen), toIndex(end, listLen)
}

// atIndex finds element at given index or nil.
func atIndex(index int, list *list.List) *list.Element {
	index = toIndex(index, list.Len())
	e, i := list.Front(), 0
	for ; e.Next() != nil && i < index; i++ {
		if e.Next() == nil {
			return nil
		}
		e = e.Next()
	}
	return e
}

// Converts to real index.
//
// E.g. i=5, len=10 -> returns 5
//
// E.g. i=-1, len=10 -> returns 10
//
// E.g. i=-10, len=10 -> returns 0
//
// E.g. i=-3, len=10 -> returns 7
func toIndex(i int, len int) int {
	if i < 0 {
		if len+i > 0 {
			return len + i
		} else {
			return 0
		}
	}
	return i
}

// Value of a list element to string.
func vts(e *list.Element) *string {
	v := e.Value.(string)
	return &v
}

// Return positive x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
