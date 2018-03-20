package avl

const depth = 80

// 由int64和string复合构成的键
type key struct {
	N int64
	S string
}

// 二叉树的指针
type treePointerEX struct {
	Lsn *NodeEX
	Rsn *NodeEX
}

// 链表的指针
type linkPointerEX struct {
	Lst *NodeEX
	Nxt *NodeEX
}

// 与双向链表复合AVL树的节点
type NodeEX struct {
	hgt int8
	key
	val interface{}
	treePointerEX
	linkPointerEX
}

// 插入/删除时查询的缓存记录
type traceEX struct {
	st [depth]**NodeEX
	sp int
}

// AVL树
type ListAVL struct {
	root *NodeEX
}

// 简化代码用的，用来代替空节点的节点
var empty = &NodeEX{}

// 方便计算深度
func max(x, y int8) int8 {
	if x > y {
		return x
	} else {
		return y
	}
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
func (this *NodeEX) Prev() *NodeEX {
	return this.Lst
}

// 获得后继节点，用于简单迭代
func (this *NodeEX) Next() *NodeEX {
	return this.Nxt
}

// 获得节点的键，采用函数避免误修改
func (this *NodeEX) Key() (int64, string) {
	return this.key.N, this.key.S
}

// 获得节点的值，采用函数避免误修改
func (this *NodeEX) Val() interface{} {
	return this.val
}

// 设置节点的值
func (this *NodeEX) Set(v interface{}) {
	this.val = v
}

// 从枝叶到根节点维护AVL树
func (this *traceEX) Maintain() {
	for i := this.sp - 1; i >= 0; i-- {
		p := *this.st[i]
		s := p.hgt
		l, r := p.Lsn, p.Rsn
		switch t := l.hgt - r.hgt; {
		case t > +1:
			b, d := l.Lsn, l.Rsn
			if b.hgt >= d.hgt {
				*this.st[i] = l
				l.Rsn, p.Lsn = p, d
				p.hgt = max(d.hgt, r.hgt) + 1
				l.hgt = max(b.hgt, p.hgt) + 1
				p = l
			} else {
				x, y := d.Lsn, d.Rsn
				*this.st[i] = d
				d.Lsn, d.Rsn = l, p
				l.Rsn, p.Lsn = x, y
				l.hgt = max(b.hgt, x.hgt) + 1
				p.hgt = max(y.hgt, r.hgt) + 1
				d.hgt = max(l.hgt, p.hgt) + 1
				p = d
			}
		case t < -1:
			b, d := r.Lsn, r.Rsn
			if b.hgt <= d.hgt {
				*this.st[i] = r
				r.Lsn, p.Rsn = p, b
				p.hgt = max(l.hgt, b.hgt) + 1
				r.hgt = max(p.hgt, d.hgt) + 1
				p = r
			} else {
				x, y := b.Lsn, b.Rsn
				*this.st[i] = b
				b.Lsn, b.Rsn = p, r
				p.Rsn, r.Lsn = x, y
				p.hgt = max(l.hgt, x.hgt) + 1
				r.hgt = max(y.hgt, d.hgt) + 1
				b.hgt = max(p.hgt, r.hgt) + 1
				p = b
			}
		default:
			p.hgt = max(l.hgt, r.hgt) + 1
		}
		if s == p.hgt {
			break
		}
	}
}

// 将根节点或者分支节点旋转到成为叶节点
func (this *traceEX) ToLeaf() {
	i := this.sp
	p := *this.st[i-1]
	l, r := p.Lsn, p.Rsn
loop:
	for {
		switch {
		case l.hgt > r.hgt:
			t := l.Rsn
			*this.st[i-1] = l
			this.st[i] = &l.Rsn
			l.hgt = 0 // 目的在于maintain时，通知maintain函数该节点是变化过的节点
			l.Rsn = p
			p.Lsn = t
			l = t
		case r == empty:
			break loop
		default:
			t := r.Lsn
			*this.st[i-1] = r
			this.st[i] = &r.Lsn
			r.hgt = 0 // 目的在于maintain时，通知maintain函数该节点是变化过的节点
			r.Lsn = p
			p.Rsn = t
			r = t
		}
		i++
	}
	this.sp = i
}

// 搜索键对应的节点，记录的是保存节点位置的变量的地址。
// 如果存在，最后一个数据为该节点；
// 否则，最后一个数据为如果插入新的节点，保存该节点地址的变量的地址。
func (this *traceEX) Search(x **NodeEX, k *key) bool {
	i, p := 0, *x
	for p != empty {
		this.st[i] = x
		i++
		switch compare(k, &p.key) {
		case -1:
			x, p = &p.Lsn, p.Lsn
		case +1:
			x, p = &p.Rsn, p.Rsn
		default:
			this.sp = i
			return true
		}
	}
	this.st[i] = x
	this.sp = i + 1
	return false
}

// 加入键值对，如果已有键则更新值。
func (this *traceEX) Update(x **NodeEX, k *key, v interface{}) {
	ok := this.Search(x, k)
	if ok {
		(*this.st[this.sp-1]).val = v
		return
	}
	p := new(NodeEX)
	*p = NodeEX{1, *k, v, treePointerEX{empty, empty}, linkPointerEX{}}
	*this.st[this.sp-1] = p
	this.sp--
	if this.sp > 0 {
		x := this.st[this.sp-1]
		q := *x
		if compare(k, &q.key) >= 0 {
			r := q.Nxt
			q.Nxt = p
			p.Lst = q
			p.Nxt = r
			if r != nil {
				r.Lst = p
			}
		} else {
			r := q.Lst
			q.Lst = p
			p.Nxt = q
			p.Lst = r
			if r != nil {
				r.Nxt = p
			}
		}
	}
	this.Maintain()
}

// 加入键值对，即使已有键仍加入新的键值对。
func (this *traceEX) Insert(x **NodeEX, k *key, v interface{}) {
	ok := this.Search(x, k)
	if ok {
		i := this.sp
		p := *this.st[i-1]
		if p.Lsn.hgt < p.Rsn.hgt {
			this.st[i] = &p.Lsn
			i++
			for p := p.Lsn; p != empty; p = p.Rsn {
				this.st[i] = &p.Rsn
				i++
			}
		} else {
			this.st[i] = &p.Rsn
			i++
			for p := p.Rsn; p != empty; p = p.Lsn {
				this.st[i] = &p.Lsn
				i++
			}
		}
		this.sp = i
	}
	p := new(NodeEX)
	*p = NodeEX{1, *k, v, treePointerEX{empty, empty}, linkPointerEX{}}
	*this.st[this.sp-1] = p
	this.sp--
	if this.sp > 0 {
		x := this.st[this.sp-1]
		q := *x
		if compare(k, &q.key) >= 0 {
			r := q.Nxt
			q.Nxt = p
			p.Lst = q
			p.Nxt = r
			if r != nil {
				r.Lst = p
			}
		} else {
			r := q.Lst
			q.Lst = p
			p.Nxt = q
			p.Lst = r
			if r != nil {
				r.Nxt = p
			}
		}
	}
	this.Maintain()
}

// 删除键值对
func (this *traceEX) Delete(x **NodeEX, k *key) {
	if !this.Search(x, k) {
		return
	}
	this.ToLeaf()
	x = this.st[this.sp-1]
	p := *x
	*x = empty
	this.sp--
	l, r := p.Lst, p.Nxt
	if l != nil {
		l.Nxt = r
	}
	if r != nil {
		r.Lst = l
	}
	this.Maintain()
}

// 创建一个与链表复合的AVL树
func NewListAVL() *ListAVL {
	p := new(ListAVL)
	p.root = empty
	return p
}

// 如果键已存在，更新值；如果不存在，插入新的键值对
func (this *ListAVL) Update(n int64, s string, v interface{}) {
	var tr traceEX
	tr.Update(&this.root, &key{n, s}, v)
}

// 不管键已存在或不存在，都插入新的键值对
func (this *ListAVL) Insert(n int64, s string, v interface{}) {
	var tr traceEX
	tr.Insert(&this.root, &key{n, s}, v)
}

// 根据键删除键值对所对应的节点
func (this *ListAVL) Delete(n int64, s string) {
	var tr traceEX
	tr.Delete(&this.root, &key{n, s})
}

// 根据键查找键值对所对应的节点
func (this *ListAVL) Search(n int64, s string) *NodeEX {
	k := &key{n, s}
	for p := this.root; p != empty; {
		switch compare(k, &p.key) {
		case -1:
			p = p.Lsn
		case +1:
			p = p.Rsn
		default:
			return p
		}
	}
	return nil
}

// 返回最小键的节点
func (this *ListAVL) Min() *NodeEX {
	p := this.root
	if p == empty {
		return nil
	}
	for p.Lsn != empty {
		p = p.Lsn
	}
	return p
}

// 返回最大键的节点
func (this *ListAVL) Max() *NodeEX {
	p := this.root
	if p == empty {
		return nil
	}
	for p.Rsn != empty {
		p = p.Rsn
	}
	return p
}
