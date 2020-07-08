package evaluator

import (
	"Eldrlang/object"
	"Eldrlang/parser"
	"fmt"
	"testing"
)

func TestEvaluatorNodes(t *testing.T) {
	var test = []string{
		"1",
		"hello world",
		"-1",
		"1",
		"false",
		"2",
		"5",
		"0",
		"false",
		"true",
		"false",
		"true",
		"false",
		"true",
		"",
		"1",
		"4",
		"",
		"",
		"hello world",
		"7",
		"hello",
		"world",
		"null",
	}
	var code = `
1;
"hello world";
-1;
+1;
false;
1+1;
10/2;
1+-1;
10<5;
10>5;
10!=10;
10==10;
true==false;
true==true;
var a = 1;
a;
5 + -a;
var abc = "hello";
abc = abc + " world";
abc;
{ 5 + (a * 2); }
if (true) { "hello"; } else { "world"; }
if (!true) { "hello"; } else { "world"; }
if (!true) { "hello"; }
`
	p := parser.NewParser()
	parsed := p.ParsePackage(code, "main")
	if parsed.Errors.Len() != 0 {
		fmt.Print(parsed.Errors.String())
	}
	eval := NewEvaluator(object.NewContext())
	for i, node := range parsed.Nodes {
		eval := eval.EvaluateNode(node)
		if eval == nil {
			if test[i] != "" {
				t.Fatal("no return when expected ", test[i])
			}
			fmt.Println("\"\"")
		} else {
			if eval.Inspect() != test[i] {
				t.Fatalf("failed at line %d expected \"%s\" got \"%s\"",
					i, test[i], eval.Inspect())
			}
			fmt.Println(eval.Inspect())
		}
	}
}
