package main

import (
	"flag"
	"fmt"
	"github.com/rhodeon/expression-parser/pkg/solver"
	"log"
	"os"
)

func main() {
	expr := flag.String("expr", "", "expression to solve")
	flag.Parse()

	if *expr == "" {
		println("Exparse: Input an expression with the 'expr' flag")
		os.Exit(0)
	}

	_, err := solver.Validate(*expr)
	if err != nil {
		log.Fatalln(err)
	}

	result := solver.Solve(*expr)
	fmt.Printf("Exparse: %s = %s\n", *expr, result)
}
