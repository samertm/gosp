package parse

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var _ = fmt.Printf // debugging; delete when done (which will be never, basically)

type Atom struct {
	Value interface{}
	// possible types:
	// int, float, string, function, nil, t
	Type string
}

func Nil() *Atom {
	return &Atom{Value: "nil", Type: "nil"}
}

func T() *Atom {
	return &Atom{Value: "t", Type: "t"}
}

func tokenize(s string) []string {
	parseString := strings.Replace(s, "(", " ( ", -1)
	parseString = strings.Replace(parseString, ")", " ) ", -1)
	parsed := strings.Split(parseString, " ")
	// remove empty strings (split does not remove them)
	newParsed := make([]string, 0)
	for _, s := range parsed {
		if s != "" {
			newParsed = append(newParsed, s)
		}
	}
	return newParsed
}

func pop(strs []string) (string, []string) {
	if strs == nil || len(strs) == 0 {
		return "", nil
	} else if len(strs) == 1 {
		return strs[0], nil
	} else {
		return strs[0], strs[1:len(strs)]
	}
}

func atomize(t string) (a *Atom, e error) {
	a = new(Atom)
	i, err := strconv.Atoi(t)
	if err == nil {
		a.Value = i
		a.Type = "int"
		return
	}
	f, err := strconv.ParseFloat(t, 64)
	if err == nil {
		a.Value = f
		a.Type = "float"
		return
	}
	// keywords
	switch t {
	case "nil":
		a = Nil()
	case "t":
		a = T()
	default:		
	a.Value = t
	a.Type = "symbol"
	}
	return
}

func genAst(s []string) (*list.List, []string, error) {
	// s is the input stream, without the leading '(' token
	ast := list.New()
	for t, s := pop(s); t != ""; t, s = pop(s) {
		if t == ")" {
			return ast, s, nil
		} else if t == "(" {
			l, strs, err := genAst(s)
			s = strs
			if err != nil {
				return nil, s, err
			}
			ast.PushBack(l)
		} else {
			a, err := atomize(t)
			if err != nil {
				return nil, s, err
			}
			ast.PushBack(a)
		}
	}
	return nil, s, errors.New("Unbalanced parentheses")
}

// returns *Atom or *list.List
func Parse(s string) (interface{}, error) {
	strs := tokenize(s)
	t, strs := pop(strs)
	if t != "(" {
		// handle bare symbols
		// only one token (the symbol) allowed
		if len(strs) != 0 {
			return nil, errors.New("Too many tokens")
		}
		return atomize(t)
	}
	var ast *list.List
	ast, _, err := genAst(strs)
	if err != nil {
		return nil, err
	}
	return ast, nil
}
