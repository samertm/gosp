package main

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"github.com/samertm/gosp/env"
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
var _ = env.Keys    // debugging

func eval(i interface{}) (*parse.Atom, error) {
	switch i.(type) {
	case *list.List:
		e := i.(*list.List).Front()
		t := e.Value.(*parse.Atom)
		if t.Type != "symbol" {
			return nil, errors.New("Expected symbol")
		}
		// built ins
		switch t.Value.(string) {
		case "def":
			name := e.Next().Value.(*parse.Atom).Value.(string)
			val, err := eval(e.Next().Next().Value)
			if err != nil {
				log.Fatal("nope2")
			}
			return env.Def(name, val), nil
		default:
			fun, ok := env.Keys[t.Value.(string)]
			if ok == false {
				return nil, errors.New("Symbol not in function table")
			}
			args := make([]*parse.Atom, 0)
			for e = e.Next(); e != nil; e = e.Next() {
				// eval step
				val, err := eval(e.Value)
				if err != nil {
					log.Fatal(err)
				}
				args = append(args, val)
			}
			return fun(args), nil
		}
	case *parse.Atom:
		a := i.(*parse.Atom)
		switch a.Type {
		case "int":
			return a, nil
		case "symbol":
			val, ok := env.Keys[a.Value.(string)]
			if ok == false {
				log.Fatal("nope1")
			}
			return val([]*parse.Atom{}), nil
		}
	}
	return nil, errors.New("nope")
}

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("lisppppp> ")
		input, err := r.ReadString('\n')
		if err != nil {
			log.Fatal("main", err)
		}
		input = strings.TrimSpace(input)
		ast, err := parse.Parse(input)
		if err != nil {
			log.Fatal("main", err)
		}
		a, err := eval(ast)
		fmt.Println(a.Value)
	}
}
