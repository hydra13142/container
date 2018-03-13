// 平衡二叉查找树: 宽度平衡二叉树 —— Size Balanced Tree —— SBT树
package sbt

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
	right int64
	Value interface{}
	treePointer
	linkPointer
}

// 出于简化代码考虑，隐藏的末端节点
var null = &Node{}

// 宽度平衡二叉树，因为维护二叉树会导致根节点改变，另设一类型维护
type SBT struct {
	root      *Node
	duplicate bool
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

// 创建一个SBT树，重复键只会更新对应的值
func NewSBT() *SBT {
	p := new(SBT)
	p.root = null
	p.duplicate = false
	return p
}

// 创建一个可以存入重复键值对的SBT树
func NewDuplicateSBT() *SBT {
	p := new(SBT)
	p.root = null
	p.duplicate = true
	return p
}

// 插入新值
func (this *SBT) Insert(k int64, v interface{}) {
	var p, q, r *Node
	for q, p = null, this.root; p != null; {
		switch {
		case k < p.right:
			q, p = p, p.Lson
		case k > p.right:
			q, p = p, p.Rson
		default:
			if !this.duplicate {
				q.Value = v
				return
			}
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
	if k > q.right {
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

// 删除值
// 内部删除的节点会直接进行清零处理后释放
func (this *SBT) Delete(k int64) {
	if p := this.Search(k); p != nil {
		this.DeleteNode(p)
		*p = Node{}
	}
}

// 删除树的某个节点，使用者应确保该节点确为当前树的节点，否则结果不可预知。
// 删除的节点会保持原本的指针和键值对，对齐清理/利用由使用者负责
func (this *SBT) DeleteNode(p *Node) {
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
func (this *SBT) Search(k int64) *Node {
	for p := this.root; p != null; {
		switch {
		case p.right > k:
			p = p.Lson
		case p.right < k:
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

// 获得节点的键，采用函数避免误修改
func (this *Node) Key() int64 {
	return this.right
}

// 获得后继节点，用于简单迭代
func (this *Node) Next() *Node {
	return this.Nxt
}

// 获得前驱节点，用于简单迭代
func (this *Node) Last() *Node {
	return this.Pre
}
