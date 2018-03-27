use 'godoc cmd/github.com/hydra13142/container/treap' for documentation on the github.com/hydra13142/container/treap command 

PACKAGE DOCUMENTATION

package treap
    import "github.com/hydra13142/container/treap"


TYPES

type BST struct {
    Treap
}
    以树堆为底层结构的二叉搜索树

func NewBST() *BST
    创建一个以树堆为底层结构的二叉搜索树

func (this *BST) Delete(n typeA, s typeB)
    删除键值对

func (this *BST) Insert(n typeA, s typeB, v typeC)
    添加键值对，即使键已存在仍然添加

func (this *BST) Search(n typeA, s typeB) typeC
    根据键查找值

func (this *BST) Update(n typeA, s typeB, v typeC)
    添加键值对或者更新已存在的键对应的值

type Key struct {
    N typeA
    S typeB
}
    由typeA和typeB复合构成的键

type Node struct {
    // contains filtered or unexported fields
}
    树堆的节点

func (this *Node) Key() (typeA, typeB)
    获得节点的键，采用函数避免误修改

func (this *Node) Set(v typeC)
    设置节点的值

func (this *Node) Val() typeC
    获得节点的值，采用函数避免误修改

func (this *Node) Weight() int64
    获得节点的优先级，越小越优先

type PQ struct {
    Treap
}
    使用树堆为底层结构的优先级队列

func NewPQ() *PQ
    使用树堆为底层结构的优先级队列

func (this *PQ) Insert(w int64, v typeC)
    不管是否存在同一优先级的任务都添加任务

func (this *PQ) Peek() *Node
    释放最高优先级的任务

func (this *PQ) Pop() *Node
    释放最高优先级的任务

func (this *PQ) Update(w int64, v typeC)
    添加任务或者更新同一优先级的任务

type Treap struct {
    // contains filtered or unexported fields
}
    树堆

func NewTreap() *Treap
    创建一个树堆

func (this *Treap) Insert(w int64, n typeA, s typeB, v typeC)
    插入键值对，不管键存不存在，都插入新的键值对。w为优先级；n、s构成键；v为值。

func (this *Treap) Update(w int64, n typeA, s typeB, v typeC)
    插入键值对，如果键已存在，则更新值。w为优先级；n、s构成键；v为值。


