package skiplist

import (
	"math/rand"
	"time"
)

// 本函数用于随机生成跳表高度，从1到10，每个数字出现的几率都是其左侧邻居的约20%。
var height func() int

// 由int64和string复合构成的键
type key struct {
	N int64
	S string
}

// 键值对，键值对独立出来的目的是，跳表同一键值对的节点是成列的。
type item struct {
	key
	val interface{}
}

// 跳表的节点
type Node struct {
	*item
	lft, rgt, dwn *Node
}

// 本类型目的在于记录查询键时经过的节点，这些节点在之后的插入、删除等操作中都是需要的
type trace [10]*Node

// 跳表类型
type Skiplist struct {
	root *Node // 链表左上角的节点
}

// 键值的比较函数
func compare(x, y *key) int8 {
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

// 获得前驱节点，用于简单迭代
func (this *Node) Prev() *Node {
	return this.lft
}

// 获得后继节点，用于简单迭代
func (this *Node) Next() *Node {
	return this.rgt
}

// 返回当前元素的键
func (this *Node) Key() (int64, string) {
	return this.N, this.S
}

// 返回当前元素的值
func (this *Node) Val() interface{} {
	return this.val
}

// 设置当前元素的值
func (this *Node) Set(v interface{}) {
	this.item.val = v
}

// p为跳表左上角的节点，返回值，int值表示查询经历的层数，bool表示是否查询到该键。
func (this *trace) Search(p *Node, k *key) (int, bool) {
	var (
		q *Node
		i = 0
	)
	if p == nil || compare(k, &p.item.key) < 0 {
		return 0, false
	}
outer:
	for {
	inner:
		for p != nil {
			switch compare(&p.item.key, k) {
			case -1:
				q, p = p, p.rgt
			case 0:
				break outer
			default:
				break inner
			}
		}
		this[i] = q
		i++
		if p = q.dwn; p == nil {
			return i, false
		}
	}
	this[i] = p
	i++
	return i, true // 注意，此时不一定遍历到跳表最底层！
}

// 向跳表中插入新的键值对，注意trace的数据应为执行Search方法记录了查找轨迹的
func (this *trace) Insert(root *Node, i int, t *item) *Node {
	var l, r, p, q *Node
	if i == 0 {
		if root == nil {
			root = new(Node)
			root.item = t
			return root
		}
		v := root.item
		for p := root; p != nil; p = p.dwn {
			p.item = t
			this[i] = p
			i++
		}
		t = v
	}
	n := height()
	for i, q = i-1, nil; n > 0; i, n = i-1, n-1 {
		if i >= 0 {
			l = this[i]
			r = l.rgt
		} else {
			l = new(Node)
			l.item = root.item
			l.dwn = root
			root = l
			r = nil
		}
		p = new(Node)
		*p = Node{t, l, r, q}
		l.rgt = p
		if r != nil {
			r.lft = p
		}
		q = p
	}
	return root
}

// 删除键值对，使用者应确保该键确实存在于跳表中，注意trace的数据应为执行Search方法记录了查找轨迹的
func (this *trace) Delete(root *Node, i int, k *key) *Node {
	if compare(&root.item.key, k) == 0 {
		for p := root.dwn; p != nil; p = p.dwn {
			this[i] = p
			i++
		}
		p := this[i-1].rgt
		if p == nil {
			return nil
		}
		x := p.item
		for j := i; j > 0; j-- {
			this[j-1].item = x
			p := this[j-1].rgt
			if p != nil && p.item == x {
				i = j
			}
		}
		this[i-1] = this[i-1].rgt
	}
	for p := this[i-1]; p != nil; p = p.dwn {
		l, r := p.lft, p.rgt
		l.rgt = r
		if r != nil {
			r.lft = l
		}
	}
	return root
}

// 创建一个跳表
func New() *Skiplist {
	return new(Skiplist)
}

// 如该键不存在值则插入新键值对，如已存在则更新旧值
func (this *Skiplist) Update(n int64, s string, v interface{}) {
	var tr trace
	k := key{n, s}
	i, ok := tr.Search(this.root, &k)
	if ok {
		tr[i-1].item.val = v
		return
	}
	t := new(item)
	*t = item{k, v}
	this.root = tr.Insert(this.root, i, t)
}

// 插入跳表新的键值对，即使已存在该键，仍进行插入
func (this *Skiplist) Insert(n int64, s string, v interface{}) {
	var tr trace
	k := key{n, s}
	i, ok := tr.Search(this.root, &k)
	if ok {
		for p := tr[i-1].dwn; p != nil; p = p.dwn {
			tr[i] = p
			i++
		}
	}
	t := new(item)
	*t = item{k, v}
	this.root = tr.Insert(this.root, i, t)
}

// 删除键值对
func (this *Skiplist) Delete(n int64, s string) {
	var tr trace
	k := key{n, s}
	i, ok := tr.Search(this.root, &k)
	if ok {
		this.root = tr.Delete(this.root, i, &k)
	}
}

// 根据键来查找节点
func (this *Skiplist) Search(n int64, s string) *Node {
	var p, q *Node = this.root, nil
	k := key{n, s}
	if p == nil || compare(&k, &p.item.key) < 0 {
		return nil
	}
outer:
	for {
	inner:
		for p != nil {
			switch compare(&p.item.key, &k) {
			case -1:
				q, p = p, p.rgt
			case 0:
				break outer
			default:
				break inner
			}
		}
		if p = q.dwn; p == nil {
			return nil
		}
	}
	for p.dwn != nil {
		p = p.dwn
	}
	return p
}

// 返回最小键的（位于最底层的）节点
func (this *Skiplist) Min() *Node {
	p := this.root
	if p == nil {
		return nil
	}
	for p.dwn != nil {
		p = p.dwn
	}
	return p
}

// 返回最大键的（位于最底层的）节点
func (this *Skiplist) Max() *Node {
	p := this.root
	if p == nil {
		return nil
	}
	for {
		for p.rgt != nil {
			p = p.rgt
		}
		if p.dwn == nil {
			break
		}
		p = p.dwn
	}
	return p
}

func init() {
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	height = func() int {
		if i := int(rd.ExpFloat64()*0.618) + 1; i < 10 {
			return i
		}
		return 10
	}
}
