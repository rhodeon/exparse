package main

import (
	"flag"
	"fmt"
	"github.com/rhodeon/expression-parser/pkg/solver"
	"log"
)

func main() {
	expr := flag.String("expr", "1 + 2", "expression to solve")
	flag.Parse()

	_, err := solver.Validate(*expr)
	if err != nil {
		log.Fatalln(err)
	}

	result := solver.Solve(*expr)
	fmt.Printf("%s = %g\n", *expr, result)
}
