package evaluator

import (
	"fmt"
	"github.com/Onelio/Eldrlang/parser"
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
		"-1",
		"",
		"hello",
		"world",
		"",
		"",
		"hello world",
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
{ var b = -1; b; }
b;
if (true) { "hello"; } else { "world"; }
if (!true) { "hello"; } else { "world"; }
if (!true) { "hello"; }
fun hello(second) { return "hello " + second; }
hello("world");
`
	p := parser.NewParser()
	eval := NewEvaluator()
	parsed := p.ParsePackage(code, "main")
	if parsed.Errors.Len() != 0 {
		fmt.Println("Parse-time:")
		fmt.Print(parsed.Errors.String())
		t.Fatalf("PARSE-FAIL")
	}
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
	if eval.errors.Len() != 0 {
		fmt.Println("Run-time:")
		fmt.Print(eval.errors.String())
		t.Log("RUN-FAIL")
	}
}
