package parser

import (
	"fmt"
	"testing"
)

func TestParseNodes(t *testing.T) {
	var test = []string{
		"1",
		"var z",
		"var a = 2",
		"var b = (+3)",
		"var c = ((+4) + (-5))",
		"var d = (((-a) + b) + (+c))",
		"var e = (e + ((a + (b + (-c))) + d))",
		"var f = ((((8 + a) + b) / 7) * c)",
		"\"hello\"",
		"var g = \"hello world\"",
		"true",
		"var h = (!false)",
		"var i = (a + ((-b) / 7))",
		"{\n\t\"hello\";\n}",
		"if ((!true)) {\n\t\"hello\";\n}",
		"if ((!a)) {\n\t\"hello\";\n} else {\n\t\"world\";\n}",
		"fun j(a, b) {\n\t\"hello\";\n}",
		"b = a",
		"loop {\n\t\"hello\";\n}",
		"(a == b)",
		"j(5, (1 + (-1)))",
		"(5 + (-j(a, b, c)))",
		"return j()",
		"break",
		"{\n\t(5 + (a * 2));\n}",
	}
	var code = `
1; 
var z;
var a = 2;
var b = +3;
var c = +4 + -5;
var d = -a + b + +c;
var e = e + ( ( a + ( b + -c) ) + ( d ) );
var f = 8 + a + b / 7 * c;
"hello";
var g = "hello world";
true;
var h = !false;
var i = a + (-b / 7);
{ "hello"; }
if (!true) { "hello"; }
if (!a) { "hello"; } else { "world"; }
fun j(a, b) { "hello"; }
b = a;
loop { "hello"; }
a == b;
j(5, 1 + -1);
5 + -j(a, b, c);
return j();
break;
{ 5 + (a * 2); }
`
	p := NewParser()
	program := p.ParsePackage(code, "test")
	fmt.Println(program.Errors.String())
	for i, line := range program.Nodes {
		if line.String() != test[i] {
			t.Fatalf("failed at line %d expected \"%s\" got \"%s\"",
				i, test[i], line.String())
		}
	}
	fmt.Print(program.String())
}
