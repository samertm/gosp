package main

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/samertm/gosp/parse"
	"log"
	"os"
	"strings"
)

func printList(t *list.Element) {
	fmt.Println("startlist")
	for ; t != nil; t = t.Next() {
		// my first go type switch!
		switch ty := t.Value.(type) {
		case *list.List:
			l := t.Value.(*list.List)
			printList(l.Front())
		case *parse.Atom:
			a := t.Value.(*parse.Atom)
			fmt.Println(a.Value, a.Type)
		default:
			fmt.Println("error", ty)
		}
	}
	fmt.Println("endlist")
}

var _ = parse.Parse // debugging

func main() {
	fmt.Println("write some lisssssp yo")
	r := bufio.NewReader(os.Stdin)
	input, err := r.ReadString('\n')
	if err != nil {
		log.Fatal("main", err)
	}
	input = strings.TrimSpace(input)
	ast, err := parse.Parse(input)
	if err != nil {
		log.Fatal("main", err)
	}
	printList(ast.Front())
}
