package main

import (
	"encoding/hex"
	"fmt"
	"github.com/ontio/group_tool/config"
	"github.com/ontio/ontology/account"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/native/utils"
)

type Signer struct {
	id    []byte
	index uint32
}

func SerializeSigners(s []Signer) []byte {
	sink := common.NewZeroCopySink(nil)
	utils.EncodeVarUint(sink, uint64(len(s)))
	for _, v := range s {
		sink.WriteVarBytes(v.id)
		utils.EncodeVarUint(sink, uint64(v.index))
	}
	return sink.Bytes()
}

type Group struct {
	Members   []interface{} `json:"members"`
	Threshold uint          `json:"threshold"`
}

func (g *Group) Serialize() []byte {
	sink := common.NewZeroCopySink(nil)
	utils.EncodeVarUint(sink, uint64(len(g.Members)))
	for _, m := range g.Members {
		switch t := m.(type) {
		case []byte:
			sink.WriteVarBytes(t)
		case *Group:
			sink.WriteVarBytes(t.Serialize())
		default:
			fmt.Println(t)
			panic("invlid member type")
		}
	}
	utils.EncodeVarUint(sink, uint64(g.Threshold))
	return sink.Bytes()
}

func main() {
	err := config.DefConfig.Init("config.json")
	if err != nil {
		fmt.Printf("DefConfig.Init error:%s", err)
		return
	}

	id1, _ := account.GenerateID()
	id2, _ := account.GenerateID()
	id3, _ := account.GenerateID()
	subGroup := new(Group)
	subGroup.Threshold = 2
	subGroup.Members = []interface{}{[]byte(id2), []byte(id3)}
	group := new(Group)
	group.Threshold = 1
	group.Members = []interface{}{[]byte(id1), subGroup}

	signers1 := []Signer{{[]byte(id1), 1}}
	signers2 := []Signer{{[]byte(id2), 1}, {[]byte(id3), 1}}

	g := group.Serialize()
	s1 := SerializeSigners(signers1)
	s2 := SerializeSigners(signers2)
	fmt.Println("group: ", hex.EncodeToString(g))
	fmt.Println("signers1: ", hex.EncodeToString(s1))
	fmt.Println("signers2: ", hex.EncodeToString(s2))
}
