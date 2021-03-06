package eval

import (
    "monkey_interpreter/lexer"
    "monkey_interpreter/parser"
    "monkey_interpreter/object"
    "testing"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"5", 5},
        {"46", 46},
        {"-46", -46},
        {"5 + 5 + 5 + 5 - 10", 10},
        {"2 + 3 * 4", 14},
        {"2 * 3 + 4", 10},
        {"2 + 3 * 4 + 5", 19},
        {"(2 + 3) * (4 + 5)", 45},
        {"2 *  -10", -20},
        {"(2 + 8)/ 5 - 10", -8},
    }

    for _, test := range tests {
        evaledObj := testEval(test.input)
        testIntegerObject(t, evaledObj, test.expected)
    }
}

func TestEvalBooleanExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"true", true},
        {"false", false},
        {"1 < 2", true},
        {"1 > 2", false},
        {"1 == 2", false},
        {"1 != 2", true},
        {`"tom" == "tom"`, true},
        {`"bom" == "rom"`, false},
        {`"bom" != "rom"`, true},
    }

    for _, test := range tests {
        evaledObj := testEval(test.input)
        testBooleanObject(t, evaledObj, test.expected)
    }
}

func TestBangOperator(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    }{
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!!false", false},
        {"!!5", true},
    }

    for _, test := range tests {
        evaledObj := testEval(test.input)
        testBooleanObject(t, evaledObj, test.expected)
    }
}

func TestIfElseExpression(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {"if (true) { 10 }", 10},
        {"if (false) { 10 }", nil},
        {"if (1) { 10 }", 10},
        {"if (1<2) { 10 }", 10},
        {"if (1>2) { 10 }", nil},
        {"if (1<2) { 10 } else { 20 }", 10},
        {"if (1>2) { 10 } else { 20 }", 20},
        {"if (1+2+3 == 1*2*3) { 10 } else { 20 }", 10},
    }

    for _, test := range tests {
        evaled := testEval(test.input)
        int, ok := test.expected.(int)
        if ok {
            testIntegerObject(t, evaled, int64(int))
        } else {
            testNullObject(t, evaled)
        }
    }
}

func TestReturnStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"return 10;", 10},
        {"return 10; 9;", 10},
        {"return 2 * 5; 9;", 10},
        {"11; return 2 * 5; 9;", 10},
        {"10;11; 12;", 12},
        {`
        if (10 > 1) {
            if (10 > 1) {
                return 10;
            }
            return 1;
        }`, 10},
    }

    for _, test := range tests {
        evaled := testEval(test.input)
        testIntegerObject(t, evaled, test.expected)
    }
}

func TestErrorHandling(t *testing.T) {
    tests := []struct {
        input string
        expectedMsg string
    }{
        {
            "5 + true;",
            "type mismatch: INTEGER + BOOLEAN",
        },
        {
            "5 + true; 5;",
            "type mismatch: INTEGER + BOOLEAN",
        },
        {
            "-true;",
            "unknown operator: -BOOLEAN",
        },
        {
            "true + false;",
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            "5; true + false; 5",
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            "if (10 > 1) { true + false; }",
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            `
            if (10 > 1) {
                if (10 > 1) {
                    return true + false;
                }
                return 1;
            }`,
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            "foo",
            "identifier not found: foo",
        },
        {
            `"Hello" - "World"`,
            "unknown operator: STRING - STRING",
        },
        {
            `{"name": "Monkey"}[fn(x) { x }];`,
            "unusable as hash key: FUNCTION",
        },
    }

    for _, test := range tests {
        evaled := testEval(test.input)

        errObj, ok := evaled.(*object.Error)
        if !ok {
            t.Errorf("no error object returned")
            continue
        }

        if errObj.Msg != test.expectedMsg {
            t.Errorf("wrong error message\n\"%s\" expected, but got \"%s\"", test.expectedMsg, errObj.Msg)
        }
    }
}

func TestLetStatement(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"let a = 5; a;", 5},
        {"let a = 5 * 5; a;", 25},
        {"let a = 5; let b = a; b;", 5},
        {"let a = 5; let b = a; let c = a + b + 5; c;", 15},
    }

    for _, test := range tests {
        testIntegerObject(t, testEval(test.input), test.expected)
    }
}

func TestFunctionObject(t *testing.T) {
    input := "fn(x) { x + 2; };"

    evaled := testEval(input)
    fn, ok := evaled.(*object.Function)
    if !ok {
        t.Fatalf("object is not Function")
    }

    if len(fn.Params) != 1 {
        t.Fatalf("length of Params is incorrect")
    }

    if fn.Params[0].Value != "x" {
        t.Fatalf("incorrect params value")
    }

    expectedBody := "(x + 2)"
    if fn.Body.String() != expectedBody {
        t.Fatalf("incorrect body value")
    }
}

func TestFunctionApplication(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    }{
        {"let id = fn(x) { x; }; id(5);", 5},
        {"let id = fn(x) { return x; }; id(5);", 5},
        {"let double = fn(x) { x * 2; }; double(5);", 10},
        {"let id = fn(x) { x; }; id(5);", 5},
        {"let add = fn(x, y) { x + y; }; add(5, 10);", 15},
        {"let add = fn(x, y) { x + y; }; add(5+5, add(5,5));", 20},
        {"fn(x) { x; }(5)", 5},
        {"let i = 5; let p = fn(i) { i; }; p(10); i;", 5},
        {"let i = 5; let p = fn(i) { i; }; i; p(10);", 10},
    }

    for _, test := range tests {
        testIntegerObject(t, testEval(test.input), test.expected)
    }
}

func TestStringLiteral(t *testing.T) {
    test := `"howdy? toasa."`
    evaled := testEval(test)

    s, ok := evaled.(*object.String)
    if !ok {
        t.Errorf("type assertion error")
    }

    if s.Value != "howdy? toasa." {
        t.Errorf("expected howdy? toasa. but got %s", s.Value)
    }
}

func TestStringConcatenation(t *testing.T) {
    input := `"Howdy?" + " " + "toasa"`
    evaled := testEval(input)

    s, ok := evaled.(*object.String)
    if !ok {
        t.Errorf("type assertion error")
    }

    if s.Value != "Howdy? toasa" {
        t.Errorf("expected Hoydy? toasa but got %s", s.Value)
    }
}

func TestBuiltinFunctions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {`len("")`, 0},
        {`len("five")`, 4},
        {`len("hello toasa")`, 11},
        {`len(1)`, "argument to `len` not supported, got INTEGER"},
        {`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
    }

    for _, test := range tests {
        evaled := testEval(test.input)

        switch expected :=  test.expected.(type){
        case int:
            testIntegerObject(t, evaled, int64(expected))
        case string:
            errObj, ok := evaled.(*object.Error)
            if !ok {
                t.Errorf("object is not error")
                continue
            }
            if errObj.Msg != expected {
                t.Errorf("wrong err message")
            }
        }
    }
}

func TestArrayLiterals(t *testing.T) {
    input := "[1, 2 * 2, 3 + 3]"
    evaled := testEval(input)
    a, ok := evaled.(*object.Array)
    if !ok {
        t.Errorf("type assertion error")
        return
    }
    testIntegerObject(t, a.Elems[0], int64(1))
    testIntegerObject(t, a.Elems[1], int64(4))
    testIntegerObject(t, a.Elems[2], int64(6))
}

func TestArrayIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
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
            "let i = 0; [1][i]",
            1,
        },
        {
            "[1, 2, 3][1 + 1]",
            3,
        },
        {
            "let arr = [1, 2, 3]; arr[0] + arr[1] + arr[2];",
            6,
        },
        {
            "let arr = [1, 2, 3]; let i = arr[0]; arr[i]",
            2,
        },
        {
            "[1, 2, 3][3]",
            nil,
        },
        {
            "[1, 2, 3][-1]",
            nil,
        },
    }

    for _, test := range tests {
        evaled := testEval(test.input)
        i, ok := test.expected.(int)
        if ok {
            testIntegerObject(t, evaled, int64(i))
        } else {
            testNullObject(t, evaled)
        }
    }
}

func TestHashLiteral(t *testing.T) {
    input := `let two = "two";
    {
        "one": 10 - 9,
        two: 1 + 1,
        "thr" + "ee": 6/ 2,
        4: 4,
        true: 5,
        false: 6,
    }`

    evaled := testEval(input)
    h, ok := evaled.(*object.Hash)
    if !ok {
        t.Fatalf("Eval didn't return Hash")
    }

    expected := map[object.HashKey]int64 {
        (&object.String{Value: "one"}).HashKey(): 1,
        (&object.String{Value: "two"}).HashKey(): 2,
        (&object.String{Value: "three"}).HashKey(): 3,
        (&object.Integer{Value: 4}).HashKey(): 4,
        (&object.Boolean{Value: true}).HashKey(): 5,
        (&object.Boolean{Value: false}).HashKey(): 6,
    }

    if len(h.Pairs) != len(expected) {
        t.Fatalf("Hash has wrong num of pairs")
    }

    for expectedKey, expectedVal := range expected {
        pair, ok := h.Pairs[expectedKey]
        if !ok {
            t.Errorf("no pair for given key")
        }

        testIntegerObject(t, pair.Value, expectedVal)
    }
}

func TestHashIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {
            `{"foo": 5}["foo"]`,
            5,
        },
        {
            `{"foo": 5}["bar"]`,
            nil,
        },
        {
            `let key = "foo"; {"foo": 5}[key]`,
            5,
        },
        {
            `{}["foo"]`,
            nil,
        },
        {
            `{5: 5}[5]`,
            5,
        },
        {
            `{true: 5}[true]`,
            5,
        },
        {
            `{false: 5}[false]`,
            5,
        },
    }

    for _, test := range tests {
        evaled := testEval(test.input)
        i, ok := test.expected.(int)
        if ok {
            testIntegerObject(t, evaled, int64(i))
        } else {
            testNullObject(t, evaled)
        }
    }
}

func testEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    env := object.NewEnv()

    return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
    i, ok := obj.(*object.Integer)
    if !ok {
        t.Errorf("object is not Integer, got = %T", obj)
        return false
    }

    if i.Value != expected {
        t.Errorf("%d expected, but got %d", expected, i.Value)
        return false
    }

    return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
    b, ok := obj.(*object.Boolean)
    if !ok {
        t.Errorf("object is not Boolean, got = %T", obj)
        return false
    }

    if b.Value != expected {
        t.Errorf("object has incorrect value")
        return false
    }

    return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
    if obj == NULL {
        return true
    }
    t.Errorf("obj is not Null object, but got %+v (%T)", obj, obj)
    return false
}
