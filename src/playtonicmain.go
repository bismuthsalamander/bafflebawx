package main

import (
	"encoding/json"
	"fmt"

	"github.com/bismuthsalamander/bafflebawx/playtonic"
)

type response2 struct {
	Page   int      `json:"aaaa"`
	Fruits []string `json:"xxxx"`
}

type WHA struct {
	difficulty int `json:"difficulty"`
	modifier   int `json:"modifier"`
}

func main() {
	x := response2{2, []string{"apple", "peach", "pear"}}
	a, b := json.Marshal(x)
	fmt.Printf("%v %v %v\n", x, string(a), b)
	y := WHA{9, 2}
	a, b = json.Marshal(y)
	fmt.Printf("%v %v %v\n", y, string(a), b)
	/*
		var e error
		var r playtonic.RollType
		r, e = playtonic.ParseDescriptor("4d6*^")
		fmt.Printf("%v, %v\n", r, e)
		for i := 0; i < 30; i++ {
			a, b, e := playtonic.SkillCheck(0, 5)
			fmt.Printf("Skill check: %v, %v %v\n", a, b, e)
		}

		r, e = playtonic.ParseDescriptor("d6-3")
		fmt.Printf("%v, %v\n", r, e)
		r, e = playtonic.ParseDescriptor("d99+2")
		fmt.Printf("%v, %v\n", r, e)
		r, e = playtonic.ParseDescriptor("d7+-3")
		fmt.Printf("%v, %v\n", r, e)
		r, e = playtonic.ParseDescriptor("1d")
		fmt.Printf("%v, %v\n", r, e)
		r, e = playtonic.ParseDescriptor("d6")
		fmt.Printf("%v, %v\n", r, e)
	*/
	s, er := playtonic.Server(playtonic.PlaytonicConfig{"http://localhost:8080"})
	if er == nil {
		s.Run(":8081")
	}
}
