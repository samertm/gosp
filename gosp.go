package main

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/samertm/gosp/env"
	"github.com/samertm/gosp/parse"
	"log"
	"os"
	"strings"
	"errors"
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
var _ = env.Keys    // debugging

func eval(e *list.Element) (*parse.Atom, error) {
	t := e.Value.(*parse.Atom)
	if t.Type != "symbol" {
		return nil, errors.New("Expected symbol")
	}
	fun, ok := env.Keys[t.Value.(string)]
	if ok == false {
		return nil, errors.New("Symbol not in function table")
	}
	args := make([]*parse.Atom, 0)
	for e = e.Next(); e != nil; e = e.Next() {
		// eval step for ints
		args = append(args, e.Value.(*parse.Atom))
	}
	return fun(args), nil
}

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
	//printList(ast.Front())
	a, err := eval(ast.Front())
	fmt.Println(a.Value)
}
