package util

import (
	"bytes"
	"fmt"
	"github.com/cloudflare/cfssl/scan/crypto/sha256"
	"testing"
)

func TestGetNewMerkelTree(t *testing.T) {
	tss := [][]byte{}
	ts1 := []byte("第一条交易")
	ts2 := []byte("第二条交易")
	ts3 := []byte("第三条交易")
	ts4 := []byte("第四条交易")
	ts5 := []byte("第五条交易")
	ts6 := []byte("第六条交易")
	tss = append(tss, ts1, ts2, ts3, ts4, ts5, ts6)
	nt := NewMerkelTree(tss)
	t.Logf("根hash:%x", nt.MerkelRootNode.Data)
}

func TestFindMerkelNode(t *testing.T) {
	tss := [][]byte{}
	ts1 := []byte("第一条交易")
	ts2 := []byte("第二条交易")
	ts3 := []byte("第三条交易")
	ts4 := []byte("第四条交易")
	ts5 := []byte("第五条交易")
	ts6 := []byte("第六条交易")
	tss = append(tss, ts1, ts2, ts3, ts4, ts5, ts6)
	nt := NewMerkelTree(tss)

	findTs := []byte("第一条交易")
	findTsHash := sha256.Sum256(findTs)
	fmt.Printf("findTs的hash：%x\n", findTsHash)
	mn := findMK(nt.MerkelRootNode, findTsHash[:])
	if mn == nil {
		fmt.Printf("没有找到\n")
	} else {
		fmt.Printf("找到了，hash为:%x\n", mn.Data)
	}
}

func findMK(mn *MerkelNode, findTsHash []byte) *MerkelNode {
	if mn.Left == nil && mn.Right == nil {
		return nil
	}
	if bytes.Equal(mn.Left.Data, findTsHash) {
		return mn.Left
	} else if bytes.Equal(mn.Right.Data, findTsHash) {
		return mn.Right
	}

	nmn := &MerkelNode{}
	lmn := findMK(mn.Left, findTsHash)
	if lmn != nil {
		nmn = lmn
		return nmn
	}
	rmn := findMK(mn.Right, findTsHash)
	if rmn != nil {
		nmn = rmn
		return nmn
	}
	return nil
}
