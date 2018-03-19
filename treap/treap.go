package treap

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var rd = rand.New(rand.NewSource(time.Now().Unix()))

type key struct {
	N int64
	S string
}

// 树堆的节点
type node struct {
	wgt int64
	key
	val           interface{}
	Lsn, Rsn, Dad *node
}

// 树堆
type Treap struct {
	*node
}

// 使用树堆为底层结构的优先级队列
type PQ struct {
	Treap
}

// 以树堆为底层结构的二叉搜索树
type BST struct {
	Treap
}

var null = &node{wgt: int64(^uint64(0) >> 1)}

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

// 创建一个树堆
func NewTreap() *Treap {
	p := new(Treap)
	p.node = null
	return p
}

// 针对新加入的叶节点，从底向上维护树堆，返回维护后的树堆根节点
func arrange(p *node) *node {
	for {
		D, L, R := p.Dad, p.Lsn, p.Rsn
		if L.wgt < R.wgt {
			if L.wgt < p.wgt {
				r := L.Rsn
				p.Lsn, p.Dad, L.Rsn, L.Dad = r, L, p, D
				if r != null {
					r.Dad = p
				}
				if D == null {
					return L
				}
				if D.Lsn == p {
					D.Lsn = L
				} else {
					D.Rsn = L
				}
			}
		} else {
			if R.wgt < p.wgt {
				l := R.Lsn
				p.Rsn, p.Dad, R.Lsn, R.Dad = l, R, p, D
				if l != null {
					l.Dad = p
				}
				if D == null {
					return R
				}
				if D.Lsn == p {
					D.Lsn = R
				} else {
					D.Rsn = R
				}
			}
		}
		if D == null {
			return p
		}
		p = D
	}
}

// 将当前节点视为根节点多次旋转到成为叶节点后删除，并返回新的根节点
func release(p *node) *node {
	q := &node{Lsn: p, Rsn: p, Dad: p.Dad}
	D, L, R := q, p.Lsn, p.Rsn
	for {
		if L.wgt < R.wgt {
			r := L.Rsn
			p.Lsn, p.Dad, L.Rsn, L.Dad = r, L, p, D
			if r != null {
				r.Dad = p
			}
			if D.Lsn == p {
				D.Lsn = L
			} else {
				D.Rsn = L
			}
			D, L = L, r
		} else if R == null {
			if D.Lsn == p {
				D.Lsn = null
			} else {
				D.Rsn = null
			}
			p.Dad = null
			break
		} else {
			l := R.Lsn
			p.Rsn, p.Dad, R.Lsn, R.Dad = l, R, p, D
			if l != null {
				l.Dad = p
			}
			if D.Lsn == p {
				D.Lsn = R
			} else {
				D.Rsn = R
			}
			D, R = R, l
		}
	}
	D, p, q = q.Dad, q.Lsn, q.Rsn
	if D != null {
		if D.Lsn == q {
			D.Lsn = p
		} else {
			D.Rsn = p
		}
	}
	p.Dad = D
	return p
}

// 插入键值对，如果键已存在，则更新值。w为优先级；n、s构成键；v为值。
func (this *Treap) Update(w, n int64, s string, v interface{}) {
	var p, q *node
	k := key{n, s}
	for q, p = null, this.node; p != null; {
		switch p.key.cmp(&k) {
		case -1:
			q, p = p, p.Rsn
		case +1:
			q, p = p, p.Lsn
		default:
			p.val = v
			return
		}
	}
	p = new(node)
	*p = node{w, k, v, null, null, null}
	if q == null {
		this.node = p
		return
	}
	switch q.key.cmp(&k) {
	case -1:
		q.Rsn = p
	case +1:
		q.Lsn = p
	}
	p.Dad = q
	this.node = arrange(p)
}

// 插入键值对，不管键存不存在，都插入新的键值对。w为优先级；n、s构成键；v为值。
func (this *Treap) Insert(w, n int64, s string, v interface{}) {
	var p, q *node
	k := key{n, s}
	for q, p = null, this.node; p != null; {
		switch p.key.cmp(&k) {
		case -1:
			q, p = p, p.Rsn
		case +1:
			q, p = p, p.Lsn
		default:
			for q, p = p, p.Rsn; p != null; q, p = p, p.Lsn {
			}
		}
	}
	p = new(node)
	*p = node{w, k, v, null, null, null}
	if q == null {
		this.node = p
		return
	}
	switch q.key.cmp(&k) {
	case -1:
		q.Rsn = p
	case +1:
		q.Lsn = p
	}
	p.Dad = q
	this.node = arrange(p)
}

// 根据键查找值
func (this *Treap) Search(n int64, s string) interface{} {
	p := this.node
	k := key{n, s}
	for p != null {
		switch p.key.cmp(&k) {
		case 0:
			return p.val
		case 1:
			p = p.Lsn
		default:
			p = p.Rsn
		}
	}
	return nil
}

// 删除键值对
func (this *Treap) Delete(n int64, s string) {
	var p, q *node
	k := key{n, s}
	for q, p = null, this.node; p != null; {
		switch p.key.cmp(&k) {
		case -1:
			q, p = p, p.Rsn
		case +1:
			q, p = p, p.Lsn
		default:
			p = release(p)
			if q == null {
				this.node = p
			}
			return
		}
	}
}

// 使用树堆为底层结构的优先级队列
func NewPQ() *PQ {
	p := new(PQ)
	p.node = null
	return p
}

// 添加任务或者更新同一优先级的任务
func (this *PQ) Update(w int64, v interface{}) {
	this.Treap.Update(w, rd.Int63(), "", v)
}

// 不管是否存在同一优先级的任务都添加任务
func (this *PQ) Insert(w int64, v interface{}) {
	this.Treap.Insert(w, rd.Int63(), "", v)
}

// 释放最高优先级的任务
func (this *PQ) Pop() interface{} {
	if p := this.node; p != null {
		this.node = release(p)
		return p.val
	}
	return nil
}

// 创建一个以树堆为底层结构的二叉搜索树
func NewBST() *BST {
	p := new(BST)
	p.node = null
	return p
}

// 添加键值对或者更新已存在的键对应的值
func (this *BST) Update(n int64, s string, v interface{}) {
	this.Treap.Update(rd.Int63(), n, s, v)
}

// 添加键值对，即使键已存在仍然添加
func (this *BST) Insert(n int64, s string, v interface{}) {
	this.Treap.Insert(rd.Int63(), n, s, v)
}
