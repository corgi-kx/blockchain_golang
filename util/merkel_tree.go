/*
	默克尔树
*/
package util

import (
	"crypto/sha256"
)

type MerkelTree struct {
	MerkelRootNode *MerkelNode
}

type MerkelNode struct {
	Left  *MerkelNode
	Right *MerkelNode
	Data  []byte
}

func NewMerkelTree(data [][]byte) *MerkelTree {
	//首先需要知道data的数量，如果是奇数需要在结尾拷贝一份凑成偶数
	if len(data)%2 != 0 {
		data = append(data, data[len(data)-1])
	}
	//将普通交易计算成默克尔树最远叶节点，保存到切片里
	nodes := []MerkelNode{}
	for i := 0; i < len(data); i++ {
		mn := BuildMerkelNode(nil, nil, data[i])
		nodes = append(nodes, mn)
	}
	//循环获得根节点
	for {
		if len(nodes) == 1 {
			break
		}
		newNotes := []MerkelNode{}
		for i := 0; i < len(nodes); i = i + 2 {
			mn := BuildMerkelNode(&nodes[i], &nodes[i+1], nil)
			newNotes = append(newNotes, mn)
		}
		nodes = newNotes
		//防止奇数叶节点
		if len(nodes) != 1 && len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}
	}
	return &MerkelTree{&nodes[0]}
}

func BuildMerkelNode(left, right *MerkelNode, data []byte) MerkelNode {
	if left == nil && right == nil {
		datum := sha256.Sum256(data)
		mn := MerkelNode{nil, nil, datum[:]}
		return mn
	}
	sumData := append(left.Data, right.Data...)
	finalData := sha256.Sum256(sumData)
	mn := MerkelNode{left, right, finalData[:]}
	return mn
}
