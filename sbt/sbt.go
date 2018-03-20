// 平衡二叉查找树: 宽度平衡二叉树 —— Size Balanced Tree —— SBT树
package sbt

// 由int64和string复合构成的键
type key struct {
	N int64
	S string
}

// 二叉树的指针
type treePointer struct {
	Dad  *Node // 不维护该指针则必须在maintain时维护一个节点路径数组
	Lson *Node
	Rson *Node
}

// 链表的指针
type linkPointer struct {
	Pre *Node
	Nxt *Node
}

// 宽度平衡二叉树的节点，同时也是双向链表的节点
type Node struct {
	count uint
	key
	val   interface{}
	treePointer
	linkPointer
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

// 宽度平衡二叉树，因为维护二叉树会导致根节点改变，另设一类型维护
type SBT struct {
	root *Node
}

// 出于简化代码考虑，隐藏的末端节点
var null = &Node{}

// 获得前驱节点，用于简单迭代
func (this *Node) Prev() *Node {
	return this.Pre
}

// 获得后继节点，用于简单迭代
func (this *Node) Next() *Node {
	return this.Nxt
}

// 获得节点的键，采用函数避免误修改
func (this *Node) Key() (int64, string) {
	return this.N, this.S
}

// 获得节点的值，采用函数避免误修改
func (this *Node) Val() interface{} {
	return this.val
}

// 设置节点的值
func (this *Node) Set(v interface{}) {
	this.val = v
}

// 包私有的树维护函数
func maintain(p *Node) *Node {
	if p != null {
		for {
			D, L, R := p.Dad, p.Lson, p.Rson
			sp := (D.Lson == p)
			switch {
			case L.count < R.count:
				x, y := R.Lson, R.Rson
				if x.count <= y.count {
					if L.count < y.count {
						R.Dad, R.Lson = D, p
						p.Dad, p.Rson = R, x
						x.Dad = p
						p.count = L.count + x.count + 1
						L, R, p = p, y, R
					}
				} else {
					if L.count < x.count {
						m, n := x.Lson, x.Rson
						x.Dad, x.Lson, x.Rson = D, p, R
						p.Dad, p.Rson, m.Dad = x, m, p
						R.Dad, R.Lson, n.Dad = x, n, R
						p.count = L.count + m.count + 1
						R.count = n.count + y.count + 1
						L, p = p, x
					}
				}
			case L.count > R.count:
				x, y := L.Lson, L.Rson
				if x.count >= y.count {
					if R.count < x.count {
						L.Dad, L.Rson = D, p
						p.Dad, p.Lson = L, y
						y.Dad = p
						p.count = y.count + R.count + 1
						L, R, p = x, p, L
					}
				} else {
					if R.count < y.count {
						m, n := y.Lson, y.Rson
						y.Dad, y.Lson, y.Rson = D, L, p
						p.Dad, p.Lson, n.Dad = y, n, p
						L.Dad, L.Rson, m.Dad = y, m, L
						L.count = x.count + m.count + 1
						p.count = n.count + R.count + 1
						R, p = p, y
					}
				}
			}
			p.count = L.count + R.count + 1
			switch {
			case D == null:
				return p
			case sp:
				D.Lson = p
			default:
				D.Rson = p
			}
			p = D
		}
	}
	return null
}

// 创建一个SBT树
func New() *SBT {
	p := new(SBT)
	p.root = null
	return p
}

// 如果键已存在，更新值；如果不存在，插入新的键值对
func (this *SBT) Update(n int64, s string, v interface{}) {
	var p, q, r *Node
	k := key{n,s}
	for q, p = null, this.root; p != null; {
		switch compare(&k, &p.key) {
		case -1:
			q, p = p, p.Lson
		case +1:
			q, p = p, p.Rson
		default:
			q.val = v
			return
		}
	}
	p = new(Node)
	*p = Node{1, k, v, treePointer{q, null, null}, linkPointer{nil, nil}}
	if q == null {
		this.root = p
		return
	}
	if compare(&k, &q.key) > 0 {
		q.Rson = p
		r = q.Nxt
		q.Nxt, p.Pre, p.Nxt = p, q, r
		if r != nil {
			r.Pre = p
		}
	} else {
		q.Lson = p
		r = q.Pre
		q.Pre, p.Nxt, p.Pre = p, q, r
		if r != nil {
			r.Nxt = p
		}
	}
	this.root = maintain(q)
}

// 不管键已存在或不存在，都插入新的键值对
func (this *SBT) Insert(n int64, s string, v interface{}) {
	var p, q, r *Node
	k := key{n,s}
	for q, p = null, this.root; p != null; {
		switch compare(&k, &p.key) {
		case -1:
			q, p = p, p.Lson
		case +1:
			q, p = p, p.Rson
		default:
			if L, R := p.Lson, p.Rson; L.count < R.count {
				for q, p = p, R; p != null; q, p = p, p.Lson {
				}
			} else {
				for q, p = p, L; p != null; q, p = p, p.Rson {
				}
			}
		}
	}
	p = new(Node)
	*p = Node{1, k, v, treePointer{q, null, null}, linkPointer{nil, nil}}
	if q == null {
		this.root = p
		return
	}
	if compare(&k, &q.key) > 0 {
		q.Rson = p
		r = q.Nxt
		q.Nxt, p.Pre, p.Nxt = p, q, r
		if r != nil {
			r.Pre = p
		}
	} else {
		q.Lson = p
		r = q.Pre
		q.Pre, p.Nxt, p.Pre = p, q, r
		if r != nil {
			r.Nxt = p
		}
	}
	this.root = maintain(q)
}

// 删除树的某个节点，使用者应确保该节点确为当前树的节点，否则结果不可预知。
// 删除的节点会保持原本的指针和键值对，对其清理/利用由使用者负责
func (this *SBT) Delete(p *Node) {
	var q, r, s *Node
	if p == nil {
		return
	}
	D, L, R := p.Dad, p.Lson, p.Rson
	// 删除链表中的节点
	if q = p.Pre; q != nil {
		q.Nxt = p.Nxt
	}
	if q = p.Nxt; q != nil {
		q.Pre = p.Pre
	}
	switch {
	case p.count == 1: // 叶节点
		r, q = D, null
	case R == null: // 只有左子
		L.Dad = D
		r, q = D, L
	case L == null: // 只有右子
		R.Dad = D
		r, q = D, R
	case L.Rson == null: // 左子无右孙
		L.Dad, R.Dad, L.Rson = D, L, R
		r, q = L, L
	case R.Lson == null: // 右子无左孙
		R.Dad, L.Dad, R.Lson = D, R, L
		r, q = R, R
	case L.count > R.count:
		for q = L; q.Rson != null; q = q.Rson { // 循环到只有左子
		}
		r, s = q.Dad, q.Lson
		r.Rson = s
		if s != null {
			s.Dad = r
		}
		L.Dad, R.Dad = q, q
		q.treePointer = p.treePointer
	default:
		for q = R; q.Lson != null; q = q.Lson { // 循环到只有右子
		}
		r, s = q.Dad, q.Rson
		r.Lson = s
		if s != null {
			s.Dad = r
		}
		L.Dad, R.Dad = q, q
		q.treePointer = p.treePointer
	}
	switch {
	case D == null:
	case D.Lson == p:
		D.Lson = q
	default:
		D.Rson = q
	}
	this.root = maintain(r)
}

// 根据键查找值
func (this *SBT) Search(n int64, s string) *Node {
	k := key{n,s}
	for p := this.root; p != null; {
		switch compare(&k, &p.key) {
		case -1:
			p = p.Lson
		case +1:
			p = p.Rson
		default:
			return p
		}
	}
	return nil
}

// 根据索引查找值
func (this *SBT) Index(n uint) *Node {
	if p := this.root; p.count > n {
		for {
			L, R := p.Lson, p.Rson
			switch {
			case n == L.count:
				return p
			case n < L.count:
				p = L
			default:
				n -= L.count + 1
				p = R
			}
		}
	}
	return nil
}

// 返回最小键的节点
func (this *SBT) Min() *Node {
	p := this.root
	if p == null {
		return nil
	}
	for p.Lson != null {
		p = p.Lson
	}
	return p
}

// 返回最大键的节点
func (this *SBT) Max() *Node {
	p := this.root
	if p == null {
		return nil
	}
	for p.Rson != null {
		p = p.Rson
	}
	return p
}