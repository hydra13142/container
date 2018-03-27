package avl

// 这个深度，已足以保存最少4万亿个数据，最多115亿亿个数据。
// 列出一个深度和最小装载数、最大装载数的表：
// 19 => 1,0946					=>	5.242e05
// 38 => 1,0233,4155			=>	2.748e11
// 58 => 1,5480,0875,5920		=>	2.882e17
// 77 => 1,4472,3340,2467,6260 	=>	1.511e23
const depth = 60

type typeA = int64

type typeB = string

type typeC = interface{}

// 由typeA和typeB复合构成的键
type Key struct {
	N typeA
	S typeB
}

// 二叉树的指针
type item struct {
	Key
	Val typeC
}

// AVL树的节点
type Node struct {
	mrk uint8
	hgt int8
	ptA *Node
	ptB *Node
	item
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

// 方便计算深度
func max(x, y int8) int8 {
	if x > y {
		return x
	} else {
		return y
	}
}

// 键值的比较函数
func compare(x, y *Key) int8 {
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
	p := this.ptA
	if this.mrk&2 != 0 {
		for ; p.mrk&1 != 0; p = p.ptB {
		}
	}
	return p
}

// 获得后继节点，用于简单迭代
func (this *Node) Next() *Node {
	p := this.ptB
	if this.mrk&1 != 0 {
		for ; p.mrk&2 != 0; p = p.ptA {
		}
	}
	return p
}

// 获得节点的左子节点
func (this *Node) Lson() *Node {
	if this.mrk&2 == 0 {
		return null
	} else {
		return this.ptA
	}
}

// 获得节点的右子节点
func (this *Node) Rson() *Node {
	if this.mrk&1 == 0 {
		return null
	} else {
		return this.ptB
	}
}

// 获得节点的键，采用函数避免误修改
func (this *Node) Key() (typeA, typeB) {
	return this.item.Key.N, this.item.Key.S
}

// 获得节点的值，采用函数避免误修改
func (this *Node) Val() typeC {
	return this.item.Val
}

// 设置节点的值
func (this *Node) Set(v typeC) {
	this.item.Val = v
}

// 用来以文本格式显示二叉树
func (this *Node) Show(f func(*Node) string) (int, []string) {
	var (
		s string
		v []string
	)
	if this == null {
		return -1, nil
	}
	v = []string{f(this)} // 需要生成一个表示节点的字符串
	i, x := this.Lson().Show(f)
	j, y := this.Rson().Show(f)
	if i < 0 && j < 0 {
		return 0, v
	}
	for t := 0; t < len(x); t++ {
		switch {
		case t < i:
			s = "     "
		case t > i:
			s = "|    "
		default:
			s = "+----"
		}
		x[t] = s + x[t]
	}
	for t := 0; t < len(y); t++ {
		switch {
		case t < j:
			s = "|    "
		case t > j:
			s = "     "
		default:
			s = "+----"
		}
		y[t] = s + y[t]
	}
	if len(x) > 0 {
		x = append(x, "|")
	}
	if len(y) > 0 {
		v = append(v, "|")
	}
	return len(x), append(append(x, v...), y...)
}

// 从枝叶到根节点维护AVL树
func (this *trace) Maintain() {
	for i := this.sp - 1; i >= 0; i-- {
		p := *this.st[i]
		s := p.hgt
		l, r := p.Lson(), p.Rson()
		switch t := l.hgt - r.hgt; {
		case t > +1:
			b, d := l.Lson(), l.Rson()
			if b.hgt >= d.hgt {
				*this.st[i] = l
				l.ptB, l.mrk = p, l.mrk|1
				if d == null {
					p.ptA, p.mrk = l, p.mrk&1
				} else {
					p.ptA, p.mrk = d, p.mrk|2
				}
				p.hgt = max(d.hgt, r.hgt) + 1
				l.hgt = max(b.hgt, p.hgt) + 1
				p = l
			} else {
				x, y := d.Lson(), d.Rson()
				*this.st[i] = d
				d.ptA, d.ptB, d.mrk = l, p, 3
				if x == null {
					l.ptB, l.mrk = d, l.mrk&2
				} else {
					l.ptB, l.mrk = x, l.mrk|1
				}
				if y == null {
					p.ptA, p.mrk = d, p.mrk&1
				} else {
					p.ptA, p.mrk = y, p.mrk|2
				}
				l.hgt = max(b.hgt, x.hgt) + 1
				p.hgt = max(y.hgt, r.hgt) + 1
				d.hgt = max(l.hgt, p.hgt) + 1
				p = d
			}
		case t < -1:
			b, d := r.Lson(), r.Rson()
			if b.hgt <= d.hgt {
				*this.st[i] = r
				r.ptA, r.mrk = p, r.mrk|2
				if b == null {
					p.ptB, p.mrk = r, p.mrk&2
				} else {
					p.ptB, p.mrk = b, p.mrk|1
				}
				p.hgt = max(l.hgt, b.hgt) + 1
				r.hgt = max(p.hgt, d.hgt) + 1
				p = r
			} else {
				x, y := b.Lson(), b.Rson()
				*this.st[i] = b
				b.ptA, b.ptB, b.mrk = p, r, 3
				if x == null {
					p.ptB, p.mrk = b, p.mrk&2
				} else {
					p.ptB, p.mrk = x, p.mrk|1
				}
				if y == null {
					r.ptA, r.mrk = b, r.mrk&1
				} else {
					r.ptA, r.mrk = y, r.mrk|2
				}
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
	l, r := p.Lson(), p.Rson()
loop:
	for {
		switch {
		case l.hgt > r.hgt:
			t := l.Rson()
			*this.st[i-1] = l
			this.st[i] = &l.ptB
			l.hgt = 0 // 目的在于maintain时，通知maintain函数该节点是变化过的节点
			l.ptB, l.mrk = p, l.mrk|1
			if t == null {
				p.ptA, p.mrk = l, p.mrk&1
			} else {
				p.ptA, p.mrk = t, p.mrk|2
			}
			l = t
		case r == null:
			break loop
		default:
			t := r.Lson()
			*this.st[i-1] = r
			this.st[i] = &r.ptA
			r.hgt = 0 // 目的在于maintain时，通知maintain函数该节点是变化过的节点
			r.ptA, r.mrk = p, r.mrk|2
			if t == null {
				p.ptB, p.mrk = r, p.mrk&2
			} else {
				p.ptB, p.mrk = t, p.mrk|1
			}
			r = t
		}
		i++
	}
	this.sp = i
}

// 搜索键对应的节点，记录的是保存节点位置的变量的地址。
// 如果存在，最后一个数据为该节点；
// 否则，最后一个数据为如果插入新的节点，保存该节点地址的变量的地址。
func (this *trace) Search(x **Node, k *Key) bool {
	i, p := 0, *x
	for p != null {
		this.st[i] = x
		i++
		switch compare(k, &p.item.Key) {
		case -1:
			x, p = &p.ptA, p.Lson()
		case +1:
			x, p = &p.ptB, p.Rson()
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
func (this *trace) Update(x **Node, k *Key, v typeC) {
	ok := this.Search(x, k)
	if ok {
		(*this.st[this.sp-1]).item.Val = v
		return
	}
	p := new(Node)
	*p = Node{0, 1, nil, nil, item{*k, v}}
	x = this.st[this.sp-1]
	t := *x
	*x = p
	this.sp--
	if this.sp > 0 {
		q := *this.st[this.sp-1]
		if &q.ptB == x {
			q.mrk |= 1
			p.ptA, p.ptB = q, t
		} else {
			q.mrk |= 2
			p.ptA, p.ptB = t, q
		}
	}
	this.Maintain()
}

// 加入键值对，即使已有键仍加入新的键值对。
func (this *trace) Insert(x **Node, k *Key, v typeC) {
	ok := this.Search(x, k)
	if ok {
		i := this.sp
		p := *this.st[i-1]
		l, r := p.Lson(), p.Rson()
		if l.hgt < r.hgt {
			this.st[i] = &p.ptA
			i++
			for ; l != null; l = l.Rson() {
				this.st[i] = &l.ptB
				i++
			}
		} else {
			this.st[i] = &p.ptB
			i++
			for ; r != null; r = r.Lson() {
				this.st[i] = &r.ptA
				i++
			}
		}
		this.sp = i
	}
	p := new(Node)
	*p = Node{0, 1, nil, nil, item{*k, v}}
	x = this.st[this.sp-1]
	t := *x
	*x = p
	this.sp--
	if this.sp > 0 {
		q := *this.st[this.sp-1]
		if &q.ptB == x {
			q.mrk |= 1
			p.ptA, p.ptB = q, t
		} else {
			q.mrk |= 2
			p.ptA, p.ptB = t, q
		}
	}
	this.Maintain()
}

// 删除键值对
func (this *trace) Delete(x **Node, k *Key) {
	if !this.Search(x, k) {
		return
	}
	this.ToLeaf()
	p := *this.st[this.sp-1]
	this.sp--
	if this.sp > 0 {
		q := *this.st[this.sp-1]
		l, r := p.ptA, p.ptB
		if l == q {
			q.ptB, q.mrk = r, q.mrk&2
		} else {
			q.ptA, q.mrk = l, q.mrk&1
		}
	} else {
		*this.st[0] = null
	}
	this.Maintain()
}

// 创建一个AVL线索树
func New() *AVL {
	p := new(AVL)
	p.root = null
	return p
}

// 如果键已存在，更新值；如果不存在，插入新的键值对
func (this *AVL) Update(n typeA, s typeB, v typeC) {
	var tr trace
	tr.Update(&this.root, &Key{n, s}, v)
}

// 不管键已存在或不存在，都插入新的键值对
func (this *AVL) Insert(n typeA, s typeB, v typeC) {
	var tr trace
	tr.Insert(&this.root, &Key{n, s}, v)
}

// 根据键删除键值对所对应的节点
func (this *AVL) Delete(n typeA, s typeB) {
	var tr trace
	tr.Delete(&this.root, &Key{n, s})
}

// 根据键查找键值对所对应的节点
func (this *AVL) Search(n typeA, s typeB) *Node {
	k := &Key{n, s}
	for p := this.root; p != null; {
		switch compare(k, &p.item.Key) {
		case -1:
			p = p.Lson()
		case +1:
			p = p.Rson()
		default:
			return p
		}
	}
	return nil
}

// 返回最小键的节点
func (this *AVL) Min() *Node {
	p := this.root
	if p == null {
		return nil
	}
	for {
		if q := p.Lson(); q == null {
			return p
		} else {
			p = q
		}
	}
}

// 返回最大键的节点
func (this *AVL) Max() *Node {
	p := this.root
	if p == null {
		return nil
	}
	for {
		if q := p.Rson(); q == null {
			return p
		} else {
			p = q
		}
	}
}

// 用来以文本格式显示二叉树（AVL包装版本）
func (this *AVL) Show(f func(*Node) string) string {
	n, str := this.root.Show(f)
	block := ""
	for i, line := range str {
		if i == n {
			block += "----" + line + "\r\n"
		} else {
			block += "    " + line + "\r\n"
		}
	}
	return block
}