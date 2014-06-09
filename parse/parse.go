package parse

import (
	"container/list"
	"errors"
	"fmt"
	"strings"
)

var _ = fmt.Printf // debugging; delete when done (which will be never, basically)

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
		return strs[0], strs[1:]
	}
}

// TODO: MAKE THESE FUNCTION NAMES NOT TERRIBLE
func recurse(s []string) (*list.List, []string, error) {
	// s is the input stream, without the leading '(' token
	ast := list.New()
	for t, s := pop(s); t != ""; t, s = pop(s) {
		if t == ")" {
			return ast, s, nil
		} else if t == "(" {
			l, str, err := recurse(s)
			s = str
			if err != nil {
				return nil, s, err
			}
			ast.PushBack(l)
		} else {
			ast.PushBack(t)
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
