package object

import "fmt"

type ObjectType string

const (
    INTEGER_OBJ = "INTEGER"
    BOOLEAN_OBJ = "BOOLEAN"
    NULL_OBJ = "NULL"
    RETURN_VALUE_OBJ = "RETURN_VALUE"
    ERROR_OBJ = "ERROR"
)

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Integer struct {
    Value int64
}

func (i *Integer) Inspect() string {
    return fmt.Sprintf("%d", i.Value)
}
func (i *Integer) Type() ObjectType {
    return INTEGER_OBJ
}

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType {
    return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
    return fmt.Sprintf("%t", b.Value)
}

type ReturnValue struct {
    Value Object
}

func (rv *ReturnValue) Type() ObjectType {
    return RETURN_VALUE_OBJ
}
func (rc *ReturnValue) Inspect() string {
    return rc.Value.Inspect()
}

type Null struct {}

func (n *Null) Type() ObjectType {
    return NULL_OBJ
}
func (n *Null) Inspect() string {
    return "null"
}

type Error struct {
    Msg string
}

func (e *Error) Type() ObjectType {
    return ERROR_OBJ
}
func (e *Error) Inspect() string {
    return "ERROR: " + e.Msg
}
