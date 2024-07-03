// Package gormgen
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:45
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/LLiuHuan/gin-template/cmd/gormgen/pkg"
)

var (
	input   string
	structs []string
)

func init() {
	flagStructs := flag.String("structs", "", "[Required] The name of schema structs to generate structs for, comma seperated\n")
	flagInput := flag.String("input", "", "[Required] The name of the input file dir\n")
	flag.Parse()

	if *flagStructs == "" || *flagInput == "" {
		flag.Usage()
		os.Exit(1)
	}

	structs = strings.Split(*flagStructs, ",")
	input = *flagInput
}

func main() {
	gen := pkg.NewGenerator(input)
	p := pkg.NewParser(input)
	if err := gen.ParserAST(p, structs).Generate().Format().Flush(); err != nil {
		log.Fatalln(err)
	}
}
