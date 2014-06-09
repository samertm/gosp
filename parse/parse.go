package parse

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
	"strconv"
)

var _ = fmt.Printf // debugging; delete when done (which will be never, basically)

type Atom struct {
	Value interface{}
	// possible types:
	// int, float, string
	Type  string
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

// TODO find errors
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
	a.Value = t
	a.Type = "symbol"
	return
}

// TODO: MAKE THESE FUNCTION NAMES NOT TERRIBLE
func recurse(s []string) (*list.List, []string, error) {
	// s is the input stream, without the leading '(' token
	ast := list.New()
	for t, s := pop(s); t != ""; t, s = pop(s) {
		if t == ")" {
			return ast, s, nil
		} else if t == "(" {
			l, strs, err := recurse(s)
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

func Parse(s string) (*list.List, error) {
	strs := tokenize(s)
	if len(strs) < 3 {
		return nil, errors.New("Expected > 2 tokens")
	}
	t, strs := pop(strs)
	if t != "(" {
		return nil, errors.New("Expected '('")
	}
	ast, _, err := recurse(strs)
	if err != nil {
		return nil, err
	}
	return ast, nil
}
