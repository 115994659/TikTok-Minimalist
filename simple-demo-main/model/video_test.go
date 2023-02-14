package model

import (
	"fmt"
	"testing"
)

func TestVideoDAO_QueryVideoListByUserId(t *testing.T) {
	InitDB()
	s := make([]*Video, 8)
	err := NewVideoDAO().QueryVideoListByUserId(1, &s)
	if err != nil {
		panic(err)
	}
	for _, v := range s {
		fmt.Printf("%#v\n", *v)
	}
}
