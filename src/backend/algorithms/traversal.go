package algorithms

import (
	"strings"
	"tubes2/src/backend/parser"
)

func parseSelectors(sel string) []string {
	sel = strings.ReplaceAll(sel, " > ", ">")
	sel = strings.ReplaceAll(sel, "> ", ">")
	sel = strings.ReplaceAll(sel, " >", ">")

	sel = strings.ReplaceAll(sel, " + ", "+")
	sel = strings.ReplaceAll(sel, "+ ", "+")
	sel = strings.ReplaceAll(sel, " +", "+")

	sel = strings.ReplaceAll(sel, " ~ ", "~")
	sel = strings.ReplaceAll(sel, "~ ", "~")
	sel = strings.ReplaceAll(sel, " ~", "~")

	sel = strings.TrimSpace(sel)

	var tokens []string
	curr := ""
	for i := 0; i < len(sel); i++ {
		c := sel[i]
		if c == ' ' || c == '>' || c == '+' || c == '~' {
			if curr != "" {
				tokens = append(tokens, curr)
				curr = ""
			}
			if c == ' ' {
				if len(tokens) > 0 && tokens[len(tokens)-1] != " " && tokens[len(tokens)-1] != ">" && tokens[len(tokens)-1] != "+" && tokens[len(tokens)-1] != "~" {
					tokens = append(tokens, " ")
				}
			} else {
				tokens = append(tokens, string(c))
			}
		} else {
			curr += string(c)
		}
	}
	if curr != "" {
		tokens = append(tokens, curr)
	}
	return tokens
}

func matchSimple(n *parser.Node, sel string) bool {
	if n.IsText {
		return false
	}
	if sel == "" || sel == "*" {
		return true
	}

	tag := ""
	id := ""
	var classes []string

	curr := ""
	mode := 0
	for i := 0; i < len(sel); i++ {
		c := sel[i]
		switch c {
		case '#':
			switch mode {
			case 0:
				tag = curr
			case 1:
				id = curr
			case 2:
				classes = append(classes, curr)
			}
			curr = ""
			mode = 1
		case '.':
			switch mode {
			case 0:
				tag = curr
			case 1:
				id = curr
			case 2:
				classes = append(classes, curr)
			}
			curr = ""
			mode = 2
		default:
			curr += string(c)
		}
	}
	if curr != "" {
		switch mode {
		case 0:
			tag = curr
		case 1:
			id = curr
		case 2:
			classes = append(classes, curr)
		}
	}

	if tag != "" && n.Tag != tag {
		return false
	}
	if id != "" && n.Attributes["id"] != id {
		return false
	}
	if len(classes) > 0 {
		nodeClasses := strings.Fields(n.Attributes["class"])
		for _, reqClass := range classes {
			found := false
			for _, nc := range nodeClasses {
				if nc == reqClass {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

func evaluateChain(n *parser.Node, tokens []string, tokenIdx int) bool {
	if tokenIdx < 0 {
		return true
	}
	if n == nil {
		return false
	}

	targetSel := tokens[tokenIdx]
	if !matchSimple(n, targetSel) {
		return false
	}

	if tokenIdx == 0 {
		return true
	}

	comb := tokens[tokenIdx-1]
	prevSelIdx := tokenIdx - 2

	if comb == ">" {
		return evaluateChain(n.Parent, tokens, prevSelIdx)
	} else if comb == " " {
		p := n.Parent
		for p != nil {
			if evaluateChain(p, tokens, prevSelIdx) {
				return true
			}
			p = p.Parent
		}
		return false
	} else if comb == "+" {
		if n.Parent == nil {
			return false
		}
		var prevElement *parser.Node
		for _, child := range n.Parent.Children {
			if child == n {
				break
			}
			if !child.IsText {
				prevElement = child
			}
		}
		if prevElement != nil {
			return evaluateChain(prevElement, tokens, prevSelIdx)
		}
		return false
	} else if comb == "~" {
		if n.Parent == nil {
			return false
		}
		for _, child := range n.Parent.Children {
			if child == n {
				break
			}
			if !child.IsText {
				if evaluateChain(child, tokens, prevSelIdx) {
					return true
				}
			}
		}
		return false
	}

	return false
}

func matchComplex(n *parser.Node, sel string) bool {
	if n.IsText {
		return false
	}
	tokens := parseSelectors(sel)
	if len(tokens) == 0 {
		return true
	}
	return evaluateChain(n, tokens, len(tokens)-1)
}

func formatLog(n *parser.Node) string {
	res := n.Tag
	if id, ok := n.Attributes["id"]; ok && id != "" {
		res += "#" + id
	}
	if cls, ok := n.Attributes["class"]; ok && cls != "" {
		res += "." + strings.ReplaceAll(strings.TrimSpace(cls), " ", ".")
	}
	if res == "" {
		res = "text"
	}
	return res
}

func BFS(root *parser.Node, selector string, limit int) ([]*parser.Node, []string, int) {
	var results []*parser.Node
	var log []string
	count := 0

	if root == nil {
		return results, log, count
	}

	q := []*parser.Node{root}
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if !curr.IsText {
			count++
			log = append(log, formatLog(curr))
		}

		if matchComplex(curr, selector) {
			results = append(results, curr)
			if limit > 0 && len(results) >= limit {
				break
			}
		}

		for _, child := range curr.Children {
			if child != nil {
				q = append(q, child)
			}
		}
	}

	return results, log, count
}

func DFS(root *parser.Node, selector string, limit int) ([]*parser.Node, []string, int) {
	var results []*parser.Node
	var log []string
	count := 0

	if root == nil {
		return results, log, count
	}

	stack := []*parser.Node{root}
	for len(stack) > 0 {
		idx := len(stack) - 1
		curr := stack[idx]
		stack = stack[:idx]

		if !curr.IsText {
			count++
			log = append(log, formatLog(curr))
		}

		if matchComplex(curr, selector) {
			results = append(results, curr)
			if limit > 0 && len(results) >= limit {
				break
			}
		}

		for i := len(curr.Children) - 1; i >= 0; i-- {
			if curr.Children[i] != nil {
				stack = append(stack, curr.Children[i])
			}
		}
	}

	return results, log, count
}
