package skiplist

import (
	"errors"
	"math/rand"
	"time"
)

var (
	rd = rand.New(rand.NewSource(time.Now().Unix()))
	er = errors.New("skiplist: out of range")
)

type key struct {
	N int64
	S string
}

type element struct {
	C       int
	K       key
	V       *interface{}
	L, R, D *element
}

// 跳表类型
type SkipList struct {
	root *element
	line []*element
	path []int
	long int
}

// 代表一个元素
type Element struct {
	p *element
	n int
}

func (x *key) cmp(y *key) int8 {
	switch {
	case x.N < y.N:
		return -1
	case x.N > y.N:
		return +1
	}
	switch {
	case x.S < y.S:
		return -1
	case x.S > y.S:
		return +1
	}
	return 0
}

// 创建一个跳表
func New() *SkipList {
	return &SkipList{new(element), make([]*element, 1), make([]int, 1), 0}
}

// 返回跳表长度
func (this *SkipList) Len() int {
	return this.long
}

// 清空重置跳表
func (this *SkipList) Reset() {
	var p, q, a, b *element
	for p = this.root; p != nil; p = q {
		q = p.D
		for a = p; a != nil; a = b {
			b = a.R
			a.L, a.R, a.D, a.V = nil, nil, nil, nil
		}
	}
	this.root = new(element)
	this.line = make([]*element, 1)
	this.path = make([]int, 1)
	this.long = 0
}

func (this *SkipList) search(x *key) (i int, c int) {
	for i, p := 0, this.root; p != nil; {
		t, q := int8(1), p.R
		if q != nil {
			t = q.K.cmp(x)
		}
		switch t {
		case -1:
			p, c = q, c+q.C
		case +1:
			this.line[i] = p
			this.path[i] = c
			i, p = i+1, p.D
		default:
			c += q.C
			this.line[i] = q
			this.path[i] = c
			return i, c
		}
	}
	return -1, c
}

func (this *SkipList) member(n int) (i int, c int) {
	n++
	c = n
	for i, p := 0, this.root; p != nil; {
		if n == 0 {
			this.line[i] = p
			this.path[i] = c
			return i, c
		}
		if q := p.R; q != nil && q.C <= n {
			p, n = q, n-q.C
			continue
		}
		this.line[i] = p
		this.path[i] = c - n
		i, p = i+1, p.D
	}
	return -1, c
}

func (this *SkipList) insert(i int, c int, k *key, v interface{}) {
	var p, q *element
	t := int(rd.ExpFloat64()) + 1
	l := len(this.line)
	if t > l {
		u := make([]*element, t)
		v := make([]int, t)
		copy(u[t-l:], this.line)
		copy(v[t-l:], this.path)
		for i := t - l - 1; i >= 0; i-- {
			p = new(element)
			p.D = this.root
			this.root = p
			u[i] = p
		}
		this.line = u
		this.path = v
		l = t
	}
	i, c = 0, c+1
	for i = 0; i < l-t; i++ {
		if p = this.line[i].R; p != nil {
			p.C++
		}
	}
	q = &element{}
	x := *k
	y := new(interface{})
	for ; i < l; i++ {
		p = this.line[i]
		q.D = new(element)
		q = q.D
		q.C = c - this.path[i]
		q.K, q.V = x, y
		q.L, q.R = p, p.R
		p.R = q
		if q.R != nil {
			q.R.L = q
			q.R.C -= q.C - 1
		}
	}
	*y = v
	this.long++
}

func (this *SkipList) remove(t int) {
	i, l := 0, len(this.line)
	for ; i < t; i++ {
		p := this.line[i]
		if p.R != nil {
			p.R.C--
		}
	}
	for p := this.line[i]; i < l; i++ {
		q := p.D
		if p.L != nil {
			p.L.R = p.R
		}
		if p.R != nil {
			p.R.C += p.C - 1
			p.R.L = p.L
		}
		p.L, p.R, p.D = nil, nil, nil
		p = q
	}
	this.long--
}

// 插入跳表新的值，如已存在相同键的值，不作修改并返回假
func (this *SkipList) Insert(n int64, s string, v interface{}) bool {
	k := &key{n, s}
	i, c := this.search(k)
	if i >= 0 {
		return false
	} else {
		this.insert(i, c, k, v)
	}
	return true
}

// 如该键不存在值则插入新值，如已存在则更新旧值
func (this *SkipList) InsertUpdate(n int64, s string, v interface{}) {
	k := &key{n, s}
	i, c := this.search(k)
	if i < 0 {
		this.insert(i, c, k, v)
	} else {
		*this.line[i].V = v
	}
}

// 删除跳表已有的值，如该键不存在值，返回假
func (this *SkipList) Delete(n int64, s string) bool {
	t, _ := this.search(&key{n, s})
	if t < 0 {
		return false
	} else {
		this.remove(t)
	}
	return true
}

// 删除跳表第n个值
func (this *SkipList) DeleteByIndex(n int) bool {
	if n < 0 || n >= this.long {
		return false
	}
	t, _ := this.member(n)
	if t < 0 {
		return false
	} else {
		this.remove(t)
	}
	return true
}

// 从跳表中搜索int64和string的参数为键
func (this *SkipList) Search(n int64, s string) *Element {
	c, k := 0, &key{n, s}
	for i, p := 0, this.root; p != nil; {
		t, q := int8(1), p.R
		if q != nil {
			t = q.K.cmp(k)
		}
		switch t {
		case -1:
			p, c = q, c+q.C
		case +1:
			i, p = i+1, p.D
		default:
			return &Element{q.R, c + q.C - 1}
		}
	}
	return nil
}

// 返回跳表中第n个值（下标从0开始）
func (this *SkipList) SearchByIndex(n int) *Element {
	if n < 0 || n >= this.long {
		return nil
	}
	c, n := n, n+1
	for i, p := 0, this.root; p != nil; {
		if n == 0 {
			return &Element{p, c}
		}
		if q := p.R; q != nil && q.C <= n {
			p, n = q, n-q.C
		} else {
			i, p = i+1, p.D
		}
	}
	return nil
}

// 向前遍历
func (this *Element) Prev() bool {
	this.p = this.p.L
	this.n--
	return this.p != nil
}

// 向后遍历
func (this *Element) Next() bool {
	this.p = this.p.R
	this.n++
	return this.p != nil
}

// 返回当前元素的键
func (this *Element) Key() (int64, string) {
	return this.p.K.N, this.p.K.S
}

// 返回当前元素的值
func (this *Element) Val() interface{} {
	return *this.p.V
}

// 设置当前元素的值
func (this *Element) Set(v interface{}) {
	*this.p.V = v
}
