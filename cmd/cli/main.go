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

	result, err := solver.Solve(*expr)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s = %g\n", *expr, result)
}
