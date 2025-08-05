package helper

import "sort"

// TreeNode 泛型接口，T 是指针类型
type TreeNode[T any] interface {
	GetID() uint
	GetParentID() uint
	SetChildren(children []T)
	GetSort() int
	SetSelected(selected bool)
}

// BuildTree 构建无限树
func BuildTree[T TreeNode[T]](nodes []T, parentID uint) []T {
	var tree []T
	for _, node := range nodes {
		if node.GetParentID() == parentID {
			children := BuildTree(nodes, node.GetID())
			if len(children) > 0 {
				node.SetChildren(children)
			}
			tree = append(tree, node)
		}
	}
	// 排序
	sort.SliceStable(tree, func(i, j int) bool {
		return tree[i].GetSort() < tree[j].GetSort()
	})

	return tree
}
func BuildTreeWithSelected[T TreeNode[T]](nodes []T, parentID uint, selectedIds []uint) []T {
	// 将切片转换成 map 提高查找效率
	selectedMap := make(map[uint]struct{}, len(selectedIds))
	for _, id := range selectedIds {
		selectedMap[id] = struct{}{}
	}
	var build func(nodes []T, parentID uint) []T
	build = func(nodes []T, parentID uint) []T {
		var tree []T
		for _, node := range nodes {
			if node.GetParentID() == parentID {
				// 递归构建子节点
				children := build(nodes, node.GetID())
				if len(children) > 0 {
					node.SetChildren(children)
				}
				// 判断是否选中
				if _, ok := selectedMap[node.GetID()]; ok {
					node.SetSelected(true)
				}
				tree = append(tree, node)
			}
		}
		// 排序
		sort.SliceStable(tree, func(i, j int) bool {
			return tree[i].GetSort() < tree[j].GetSort()
		})
		return tree
	}
	return build(nodes, parentID)
}
