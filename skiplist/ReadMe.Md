use 'godoc cmd/github.com/hydra13142/container/skiplist' for documentation on the github.com/hydra13142/container/skiplist command 

PACKAGE DOCUMENTATION

package skiplist
    import "github.com/hydra13142/container/skiplist"


TYPES

type Node struct {
    // contains filtered or unexported fields
}
    跳表的节点

func (this *Node) Key() (int64, string)
    返回当前元素的键

func (this *Node) Next() *Node
    获得后继节点，用于简单迭代

func (this *Node) Prev() *Node
    获得前驱节点，用于简单迭代

func (this *Node) Set(v interface{})
    设置当前元素的值

func (this *Node) Val() interface{}
    返回当前元素的值

type Skiplist struct {
    // contains filtered or unexported fields
}
    跳表类型

func New() *Skiplist
    创建一个跳表

func (this *Skiplist) Delete(n int64, s string)
    删除键值对

func (this *Skiplist) Insert(n int64, s string, v interface{})
    插入跳表新的键值对，即使已存在该键，仍进行插入

func (this *Skiplist) Max() *Node
    返回最大键的（位于最底层的）节点

func (this *Skiplist) Min() *Node
    返回最小键的（位于最底层的）节点

func (this *Skiplist) Search(n int64, s string) *Node
    根据键来查找节点

func (this *Skiplist) Update(n int64, s string, v interface{})
    如该键不存在值则插入新键值对，如已存在则更新旧值


