package main

import (
	"fmt"
	"github.com/rhodeon/expression-parser/pkg/solver"
	"log"
)

func main() {
	expr := "2+5-2(6+4)+3"
	_, err := solver.Validate(expr)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := solver.Solve(expr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
}
