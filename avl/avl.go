package avl

// 二叉树的指针
type treePointer struct {
	Lsn *Node
	Rsn *Node
}

// 链表的指针
type linkPointer struct {
	Last *Node
	Nxt *Node
}

// AVL树的节点
type Node struct {
	hgt int8
	key
	val interface{}
	treePointer
}

// 插入/删除时查询的缓存记录
type trace struct {
	st [depth]**Node
	sp int
}

// AVL树
type AVL struct {
	root *Node
}

// 简化代码用的，用来代替空节点的节点
var null = &Node{}

// 从枝叶到根节点维护AVL树
func (this *trace) Maintain() {
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
func (this *trace) ToLeaf() {
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
		case r == null:
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
func (this *trace) Search(x **Node, k *key) bool {
	i, p := 0, *x
	for p != null {
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
func (this *trace) Update(x **Node, k *key, v interface{}) {
	ok := this.Search(x, k)
	if ok {
		(*this.st[this.sp-1]).val = v
		return
	}
	p := new(Node)
	*p = Node{1, *k, v, treePointer{null, null}}
	*this.st[this.sp-1] = p
	this.sp--
	this.Maintain()
}

// 加入键值对，即使已有键仍加入新的键值对。
func (this *trace) Insert(x **Node, k *key, v interface{}) {
	ok := this.Search(x, k)
	if ok {
		i := this.sp
		p := *this.st[i-1]
		if p.Lsn.hgt < p.Rsn.hgt {
			this.st[i] = &p.Lsn
			i++
			for p := p.Lsn; p != null; p = p.Rsn {
				this.st[i] = &p.Rsn
				i++
			}
		} else {
			this.st[i] = &p.Rsn
			i++
			for p := p.Rsn; p != null; p = p.Lsn {
				this.st[i] = &p.Lsn
				i++
			}
		}
		this.sp = i
	}
	p := new(Node)
	*p = Node{1, *k, v, treePointer{null, null}}
	*this.st[this.sp-1] = p
	this.sp--
	this.Maintain()
}

// 删除键值对
func (this *trace) Delete(x **Node, k *key) {
	if !this.Search(x, k) {
		return
	}
	this.ToLeaf()
	*this.st[this.sp-1] = null
	this.sp--
	this.Maintain()
}

// 创建一个AVL树
func NewAVL() *AVL {
	p := new(AVL)
	p.root = null
	return p
}

// 如果键已存在，更新值；如果不存在，插入新的键值对
func (this *AVL) Update(n int64, s string, v interface{}) {
	var tr trace
	tr.Update(&this.root, &key{n, s}, v)
}

// 不管键已存在或不存在，都插入新的键值对
func (this *AVL) Insert(n int64, s string, v interface{}) {
	var tr trace
	tr.Insert(&this.root, &key{n, s}, v)
}

// 根据键删除键值对所对应的节点
func (this *AVL) Delete(n int64, s string) {
	var tr trace
	tr.Delete(&this.root, &key{n, s})
}

// 根据键查找键值对所对应的值
func (this *AVL) Search(n int64, s string) (interface{}, bool) {
	k := &key{n, s}
	for p := this.root; p != null; {
		switch compare(k, &p.key) {
		case -1:
			p = p.Lsn
		case +1:
			p = p.Rsn
		default:
			return p.val, true
		}
	}
	return nil, false
}
