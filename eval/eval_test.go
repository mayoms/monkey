package eval

import (
	"monkey/lexer"
	"monkey/parser"
	"testing"
)

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`{"foo"->5}["foo"]`, 5},
		{`{"foo"->5}["bar"]`, nil},
		{`let key = "foo";{"foo"->5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5->5}[5]`, 5},
		{`{true->5}[true]`, 5},
		{`{false->5}[false]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
	let two = "two";
	{
		"one"        -> 10 - 9,
		two          -> 1 + 1,
		"thr" + "ee" -> 6 /2,
		4            -> 4,
		true         -> 5,
		false        -> 6
	}`

	evaluated := testEval(input)
	hash, ok := evaluated.(*Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T, (%+v)", evaluated, evaluated)
	}
	expected := map[HashKey]int64{
		(&String{Value: "one"}).HashKey():   1,
		(&String{Value: "two"}).HashKey():   2,
		(&String{Value: "three"}).HashKey(): 3,
		(&Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                      5,
		FALSE.HashKey():                     6,
	}
	if len(hash.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. expected=%d, got=%d", len(expected), len(hash.Pairs))
	}
}

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}
	if diff1.HashKey() == hello1.HashKey() {
		t.Errorf("strings with different content have same hash key")
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			3,
		},
		{
			"let myArray = [1, 2, 3, 4, 5]; let i = myArray[0:]; let mySlice = myArray[1:]; mySlice[0]",
			2,
		},
		{
			"let myArray = [1, 2, 3, 4, 5]; let i = myArray[0]; let mySlice = myArray[:1]; mySlice[0]",
			1,
		},
		{
			"let myArray = [1, 2, 3, 4, 5]; let mySlice = myArray[:]; mySlice[0]",
			1,
		},
		{
			"let myArray = [1, 2, 3, 4, 5];let mySlice = myArray[:]; mySlice[-1]",
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			if err, ok := evaluated.(*Error); ok {
				if err.Message != "index error: '3' out of range" {
					t.Errorf("wrong error message. got=%s", err.Message)
				}
			} else {
				t.Errorf("evaluated not array or error. got=%T", evaluated)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	results, ok := evaluated.(*Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T", evaluated)
	}
	if len(results.Members) != 3 {
		t.Fatalf("array has wrong number of elements. got=%d", len(results.Members))
	}
	testIntegerObject(t, results.Members[0], 1)
	testIntegerObject(t, results.Members[1], 4)
	testIntegerObject(t, results.Members[2], 6)
}

func TestFunctionalMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`let a = [1,2].map(fn(x) {x + 1}); fn(x) { if (x[0] == 2) { if (x[1] == 3) { return true; }} else { return false }}(a)`, true},
		{`let a = [1,2].filter(fn(x) {x == 1}); fn(x) { if (x.len() == 1) { if (x[0] == 1) { return true; }} else { return false }}(a)`, true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestArrayMethods(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1,2,3].pop()`, 3},
		{`[1,2,3].pop(0)`, 1},
		{`[1,2,3].pop(2)`, 3},
		{`let a = [1,2,3].push(4);a.pop()`, 4},
		{`let a = [1,2,3].pop(1)`, 2},
		{`let a = [1,2,3]; a.pop(1); len(a)`, 2},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*Error)
			if !ok {
				t.Errorf("object is not error. got=%T (%+v)", evaluated, evaluated)
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", errObj.Message, expected)
			}
		}
	}
}

func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len([1, 3, 5])`, 3},
		{`len([1,2,3])`, 3},
		{`"string".plus()`, "undefined method 'plus' for object STRING"},
		{`"string".plus`, "undefined method 'string.plus' for object STRING"},
		{`len("one", "two")`, "wrong number of arguments. expected=1, got=2"},
		{`len(1)`, "undefined method 'len' for object INTEGER"},
		{`int("1")`, 1},
		{`int("100")`, 100},
		{`int(1)`, 1},
		{`int("one")`, `unsupported input type 'STRING: one' for function or method int`},
		{`int([])`, `unsupported input type 'ARRAY' for function or method int`},
		{`int({})`, `unsupported input type 'HASH' for function or method int`},
		{`str(1)`, "1"},
		{`str(true)`, `true`},
		{`str(false)`, `false`},
		{`str("string")`, `string`},
		{`str("one")`, `one`},
		{`str([])`, `[]`},
		{`str({})`, `{}`},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			switch s := evaluated.(type) {
			case *Boolean:

			case *String:
				testStringObject(t, evaluated, expected)
			case *Error:
				if s.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q", expected, s.Message)
				}
			default:
				t.Errorf("object is not error. got=%T (%+v)", evaluated, evaluated)
			}
		}
	}
}
func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello, World!"`, "Hello, World!"},
		{`"Hello," + " " + "World!"`, "Hello, World!"},
	}

	for _, tt := range tests {
		testStringObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
let first = 10;
let second = 10;
let third = 10;

let ourFunction = fn(first) {
  let second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{"let fact = fn(n) { if(n==1) { return n } else { return n * fact(n-1) } }; fact(5);", 120},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2 };"

	evaluated := testEval(input)

	fn, ok := evaluated.(*Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Literal.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Literal.Parameters)
	}
	if fn.Literal.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Literal.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Literal.Body.String() != expectedBody {
		t.Fatalf("body is not '(x + 2). got=%q", fn.Literal.Body)
	}
}

func TestLetStatements(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5", 25},
		{"let a = 5; let b = a;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range test {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"unsupported operator for infix expression: '+' and types INTEGER and BOOLEAN",
		},
		{
			"5 + true; 5;",
			"unsupported operator for infix expression: '+' and types INTEGER and BOOLEAN",
		},
		{
			"-true",
			"unsupported operator for prefix expression:'-' and type: BOOLEAN",
		},
		{
			"true + false;",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unsupported operator for infix expression: '+' and types BOOLEAN and BOOLEAN",
		},
		{"foobar", "unknown identifier: 'foobar' is not defined"},
		{`"abc" + 2`, "unsupported operator for infix expression: '+' and types STRING and INTEGER"},
		{`"abc" - "abc"`, "unsupported operator for infix expression: '-' and types STRING and STRING"},
		{`"abc" * "abc"`, "unsupported operator for infix expression: '*' and types STRING and STRING"},
		{`"abc" / "abc"`, "unsupported operator for infix expression: '/' and types STRING and STRING"},
		{`"abc" > "abc"`, "unsupported operator for infix expression: '>' and types STRING and STRING"},
		{`"abc" < "abc"`, "unsupported operator for infix expression: '<' and types STRING and STRING"},
		{`{"name"->"Monkey"}[fn(x) {x}];`, "key error: type FUNCTION is not hashable"},
		{`"abc"[0]`, "index error: type STRING is not indexable"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
		{"let x = 5; return x;", 5},
		{"return;", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if integer, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) {10}", nil},
		{"if (1) {10}", 10},
		{"if (1 > 2) {10}", nil},
		{"if (1 > 2) {10} else {20}", 20},
		{"if (1 < 2) {10} else {20}", 10},
		{"let x = 5;if(x == 5) { return;}", nil},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if integer, ok := tt.expected.(int); ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestEvalBooleanLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(2 < 1) == true", false},
		{"let x = 5;x == 5", true},
		{"let x = 5; x != 5", false},
		{"let x = 5; x > 5", false},
		{"let x = 4; x < 5", true},
		{"let x = 4; (x + 5) > 5", true},
		{`"abc" == "abc"`, true},
		{`"abc" == "bc"`, false},
		{`"abc" != "abc"`, false},
		{`"abc" != "bc"`, true},
		{`let x = "abc"; x == "abc"`, true},
		{`let x = fn(){ "abc" }; x() == "abc"`, true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj Object, expected bool) bool {
	result, ok := obj.(*Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%v, want=%v", result.Value, expected)
		return false
	}
	return true
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) Object {
	l := lexer.New(input)
	p := parser.New(l)
	s := NewScope(nil)
	program := p.ParseProgram()

	return Eval(program, s)
}

func testIntegerObject(t *testing.T, obj Object, expected int64) bool {
	result, ok := obj.(*Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj Object, expected string) bool {
	result, ok := obj.(*String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}
