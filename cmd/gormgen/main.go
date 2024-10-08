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
	flagStructs := flag.String("structs", "", "[必需] 要为其生成结构的模式结构的名称，以逗号分隔\n")
	flagInput := flag.String("input", "", "[必填] 输入文件的名称 dir\n")
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
	//fmt.Println(gen, p, structs)
	if err := gen.ParserAST(p, structs).Generate().Format().Flush(); err != nil {
		log.Fatalln(err)
	}
}
