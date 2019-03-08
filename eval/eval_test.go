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

func testEval(input string) object.Object {
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()

    return Eval(program)
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