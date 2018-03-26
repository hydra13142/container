use 'godoc cmd/github.com/hydra13142/container/avl' for documentation on the github.com/hydra13142/container/avl command 

PACKAGE DOCUMENTATION

package avl
    import "github.com/hydra13142/container/avl"


TYPES

	typeA、typeB、typeC分别是int64、string和interface{}的别名。主要目的方便复用。

type AVL struct {
    // contains filtered or unexported fields
}
    AVL树

func New() *AVL
    创建一个AVL线索树

func (this *AVL) Delete(n typeA, s typeB)
    根据键删除键值对所对应的节点

func (this *AVL) Insert(n typeA, s typeB, v typeC)
    不管键已存在或不存在，都插入新的键值对

func (this *AVL) Max() *Node
    返回最大键的节点

func (this *AVL) Min() *Node
    返回最小键的节点

func (this *AVL) Search(n typeA, s typeB) *Node
    根据键查找键值对所对应的节点

func (this *AVL) Show(f func(*Node) string) string
    用来以文本格式显示二叉树（AVL包装版本）

func (this *AVL) Update(n typeA, s typeB, v typeC)
    如果键已存在，更新值；如果不存在，插入新的键值对

type Key struct {
    N typeA
    S typeB
}
    由int64和string复合构成的键

type Node struct {
    // contains filtered or unexported fields
}
    AVL树的节点

func (this *Node) Key() (typeA, typeB)
    获得节点的键，采用函数避免误修改

func (this *Node) Lson() *Node
    获得节点的左子节点

func (this *Node) Next() *Node
    获得后继节点，用于简单迭代

func (this *Node) Prev() *Node
    获得前驱节点，用于简单迭代

func (this *Node) Rson() *Node
    获得节点的右子节点

func (this *Node) Set(v typeC)
    设置节点的值

func (this *Node) Show(f func(*Node) string) (int, []string)
    用来以文本格式显示二叉树

func (this *Node) Val() typeC
    获得节点的值，采用函数避免误修改


