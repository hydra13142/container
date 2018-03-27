package sbt

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

// SBT树的节点
type Node struct {
	mrk uint8
	cnt uint
	ptA *Node
	ptB *Node
	ptO *Node
	item
}

// SBT树
type SBT struct {
	root *Node
}

// 简化代码用的，用来代替空节点的节点
var null = &Node{}

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

// 维护SBT树
func maintain(r, p *Node) *Node {
	var anchor = &Node{mrk: 2, ptA: r}
	r.ptO = anchor
	for p != anchor {
		l, r, o := p.Lson(), p.Rson(), p.ptO
		sp := (o.ptA == p)
		switch {
		case l.cnt > r.cnt:
			b, d := l.Lson(), l.Rson()
			if b.cnt >= d.cnt {
				if b.cnt > r.cnt {
					l.ptB, l.mrk = p, l.mrk|1
					if d == null {
						p.ptA, p.mrk = l, p.mrk&1
					} else {
						p.ptA, p.mrk = d, p.mrk|2
					}
					p.cnt = d.cnt + r.cnt + 1
					l.cnt = b.cnt + p.cnt + 1
					d.ptO, p.ptO, l.ptO = p, l, o
					if sp {
						o.ptA = l
					} else {
						o.ptB = l
					}
					p = o
					continue
				}
			} else {
				if d.cnt > r.cnt {
					x, y := d.Lson(), d.Rson()
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
					l.cnt = b.cnt + x.cnt + 1
					p.cnt = y.cnt + r.cnt + 1
					d.cnt = l.cnt + p.cnt + 1
					l.ptO, p.ptO, d.ptO = d, d, o
					x.ptO, y.ptO = l, p
					if sp {
						o.ptA = d
					} else {
						o.ptB = d
					}
					p = o
					continue
				}
			}
		case l.cnt < r.cnt:
			b, d := r.Lson(), r.Rson()
			if b.cnt <= d.cnt {
				if d.cnt > l.cnt {
					r.ptA, r.mrk = p, r.mrk|2
					if b == null {
						p.ptB, p.mrk = r, p.mrk&2
					} else {
						p.ptB, p.mrk = b, p.mrk|1
					}
					p.cnt = l.cnt + b.cnt + 1
					r.cnt = p.cnt + d.cnt + 1
					b.ptO, p.ptO, r.ptO = p, r, o
					if sp {
						o.ptA = r
					} else {
						o.ptB = r
					}
					p = o
					continue
				}
			} else {
				if b.cnt > l.cnt {
					x, y := b.Lson(), b.Rson()
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
					p.cnt = l.cnt + x.cnt + 1
					r.cnt = y.cnt + d.cnt + 1
					b.cnt = p.cnt + r.cnt + 1
					p.ptO, r.ptO, b.ptO = b, b, o
					x.ptO, y.ptO = p, r
					if sp {
						o.ptA = b
					} else {
						o.ptB = b
					}
					p = o
					continue
				}
			}
		}
		p.cnt = l.cnt + r.cnt + 1
		p = o
	}
	r = anchor.ptA
	r.ptO = nil
	return r
}

// 通过旋转将节点转变成叶节点
func toleaf(r, p *Node) *Node {
	var anchor = &Node{mrk: 2, ptA: r}
	r.ptO = anchor
	l, r, o := p.Lson(), p.Rson(), p.ptO
	sp := (o.ptA == p)
	for {
		switch {
		case l.cnt > r.cnt:
			t := l.Rson()
			l.cnt = 0
			l.ptB, l.mrk = p, l.mrk|1
			if t == null {
				p.ptA, p.mrk = l, p.mrk&1
			} else {
				p.ptA, p.mrk = t, p.mrk|2
			}
			t.ptO, p.ptO, l.ptO = p, l, o
			if sp {
				o.ptA = l
			} else {
				o.ptB = l
			}
			l, o, sp = t, l, false
		case r == null:
			r = anchor.ptA
			r.ptO = nil
			return r
		default:
			t := r.Lson()
			r.cnt = 0
			r.ptA, r.mrk = p, r.mrk|2
			if t == null {
				p.ptB, p.mrk = r, p.mrk&2
			} else {
				p.ptB, p.mrk = t, p.mrk|1
			}
			t.ptO, p.ptO, r.ptO = p, r, o
			if sp {
				o.ptA = r
			} else {
				o.ptB = r
			}
			r, o, sp = t, r, true
		}
	}
}

// 创建一个SBT线索树
func New() *SBT {
	p := new(SBT)
	p.root = null
	return p
}

// 如果键已存在，更新值；如果不存在，插入新的键值对
func (this *SBT) Update(n typeA, s typeB, v typeC) {
	var (
		p, q *Node
		sp   int8
	)
	k := &Key{n, s}
loop:
	for q, p = nil, this.root; p != null; {
		sp = compare(k, &p.item.Key)
		switch sp {
		case -1:
			q, p = p, p.Lson()
		case +1:
			q, p = p, p.Rson()
		default:
			break loop
		}
	}
	if p != null {
		p.item.Val = v
		return
	}
	p = new(Node)
	*p = Node{0, 1, nil, nil, q, item{*k, v}}
	if q == nil {
		this.root = p
		return
	}
	if sp < 0 {
		t := q.ptA
		q.ptA, q.mrk = p, q.mrk|2
		p.ptA, p.ptB = t, q
	} else {
		t := q.ptB
		q.ptB, q.mrk = p, q.mrk|1
		p.ptA, p.ptB = q, t
	}
	this.root = maintain(this.root, p)
}

// 不管键已存在或不存在，都插入新的键值对
func (this *SBT) Insert(n typeA, s typeB, v typeC) {
	var (
		p, q *Node
		sp   int8
	)
	k := &Key{n, s}
loop:
	for q, p = nil, this.root; p != null; {
		sp = compare(k, &p.item.Key)
		switch sp {
		case -1:
			q, p = p, p.Lson()
		case +1:
			q, p = p, p.Rson()
		default:
			break loop
		}
	}
	if p != null {
		l, r := p.Lson(), p.Rson()
		if l.cnt < r.cnt {
			if l == null {
				q, sp = p, -1
			} else {
				for q, p = p, l; p != null; q, p = p, p.Rson() {
				}
				sp = +1
			}
		} else {
			if r == null {
				q, sp = p, +1
			} else {
				for q, p = p, r; p != null; q, p = p, p.Lson() {
				}
				sp = -1
			}
		}
	}
	p = new(Node)
	*p = Node{0, 1, nil, nil, q, item{*k, v}}
	if q == nil {
		this.root = p
		return
	}
	if sp < 0 {
		t := q.ptA
		q.ptA, q.mrk = p, q.mrk|2
		p.ptA, p.ptB = t, q
	} else {
		t := q.ptB
		q.ptB, q.mrk = p, q.mrk|1
		p.ptA, p.ptB = q, t
	}
	this.root = maintain(this.root, p)
}

// 删除节点
func (this *SBT) Delete(p *Node) {
	if p == nil || p == null {
		return
	}
	this.root = toleaf(this.root, p)
	l, r, o := p.ptA, p.ptB, p.ptO
	if o == nil {
		this.root = null
		return
	}
	if l == o {
		o.ptB, o.mrk = r, o.mrk&2
	} else {
		o.ptA, o.mrk = l, o.mrk&1
	}
	this.root = maintain(this.root, o)
}

// 根据键查找键值对所对应的节点
func (this *SBT) Search(n typeA, s typeB) *Node {
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

// 根据索引查找值
func (this *SBT) Index(n uint) *Node {
	if p := this.root; p.cnt > n {
		for {
			L, R := p.Lson(), p.Rson()
			switch {
			case n == L.cnt:
				return p
			case n < L.cnt:
				p = L
			default:
				n -= L.cnt + 1
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
	for {
		if q := p.Lson(); q == null {
			return p
		} else {
			p = q
		}
	}
}

// 返回最大键的节点
func (this *SBT) Max() *Node {
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

// 用来以文本格式显示二叉树（SBT包装版本）
func (this *SBT) Show(f func(*Node) string, spin bool) string {
	n, str := this.root.Show(f)
	for i, line := range str {
		if i == n {
			str[i] = "----" + line
		} else {
			str[i] = "    " + line
		}
	}
	block := ""
	x, y := len(str), 0
	if spin {
		for i := 0; i < x; i++ {
			if j := len(str[i]); y < j {
				y = j
			}
		}
		out := make([][]byte, y)
		for j := 0; j < y; j++ {
			out[j] = make([]byte, x)
		}
		for i := 0; i < x; i++ {
			t := len(str[i])
			for j := 0; j < y; j++ {
				if j < t {
					switch c := str[i][j]; c {
					case '-':
						out[j][i] = '|'
					case '|':
						out[j][i] = '-'
					default:
						out[j][i] = c
					}
				} else {
					out[j][i] = ' '
				}
			}
		}
		for i := 0; i < y; i++ {
			block += string(out[i]) + "\r\n"
		}
	} else {
		for i := 0; i < x; i++ {
			block += str[i] + "\r\n"
		}
	}
	return block
}
