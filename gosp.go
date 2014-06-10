package main

import (
	"github.com/samertm/gosp/env"
	"github.com/samertm/gosp/parse"

	"bufio"
	"container/list"
	"errors"
	"fmt"
	"os"
	"strings"
)

func printList(t *list.Element) {
	fmt.Println("(")
	for ; t != nil; t = t.Next() {
		// my first go type switch!
		switch ty := t.Value.(type) {
		case *list.List:
			l := t.Value.(*list.List)
			printList(l.Front())
		case *parse.Atom:
			a := t.Value.(*parse.Atom)
			fmt.Print(a.Value)
		default:
			fmt.Println("error", ty)
		}
	}
	fmt.Println(")")
}

func eval(i interface{}, s *env.Scope) (*parse.Atom, error) {
	switch i.(type) {
	case *list.List:
		e := i.(*list.List).Front()
		t := e.Value.(*parse.Atom)
		if t.Type != "symbol" {
			return nil, errors.New("Expected symbol")
		}
		// built ins
		if builtin, ok := builtins[t.Value.(string)]; ok {
			return builtin(e, s)
		}
		fun, ok := env.Find(s, t.Value.(string))
		if ok == false {
			return nil, errors.New("Symbol not in function table")
		}
		args := make([]*parse.Atom, 0)
		for e = e.Next(); e != nil; e = e.Next() {
			// eval step
			val, err := eval(e.Value, s)
			if err != nil {
				return nil, err
			}
			args = append(args, val)
		}
		// fun returns two values
		return fun(args)
	case *parse.Atom:
		a := i.(*parse.Atom)
		switch a.Type {
		case "int":
			return a, nil
		case "symbol":
			val, ok := env.Find(s, a.Value.(string))
			if ok == false {
				return nil, errors.New("Symbol '" + a.Value.(string) + "' not found")
			}
			return val([]*parse.Atom{})
		case "nil":
			return a, nil
		case "t":
			return a, nil
		default:
			return nil, errors.New("Type (" + a.Type + ") not recognized")
		}
	}
	return nil, errors.New("End of eval error")
}

func main() {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("gosp> ")
		input, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}
		input = strings.TrimSpace(input)
		astOrSymbol, err := parse.Parse(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		a, err := eval(astOrSymbol, env.GlobalScope)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(a.Value)
	}
}
