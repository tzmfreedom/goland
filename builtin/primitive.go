package builtin

import (
	"strings"

	"github.com/tzmfreedom/goland/ast"
)

var PrimitiveClasses []*ClassType

func NewClassMapWithPrimivie(classTypes []*ClassType) *ClassMap {
	classMap := PrimitiveClassMap()
	for _, classType := range classTypes {
		classMap.Set(classType.Name, classType)
	}
	return classMap
}

var primitiveClassMap = &ClassMap{
	Data: map[string]*ClassType{
		"integer": IntegerType,
		"string":  StringType,
		"double":  DoubleType,
		"boolean": BooleanType,
	},
}

func PrimitiveClassMap() *ClassMap {
	return primitiveClassMap
}

var nameSpaceStore = NewNameSpaceStore()

func GetNameSpaceStore() *NameSpaceStore {
	return nameSpaceStore
}

type Object struct {
	ClassType      *ClassType
	InstanceFields *ObjectMap
	GenericType    []*ClassType
	Extra          map[string]interface{}
}

func CreateObject(t *ClassType) *Object {
	return &Object{
		ClassType:      t,
		InstanceFields: NewObjectMap(),
		GenericType:    []*ClassType{},
		Extra:          map[string]interface{}{},
	}
}

func (o *Object) Value() interface{} {
	return o.Extra["value"]
}

func (o *Object) IntegerValue() int {
	return o.Value().(int)
}

func (o *Object) DoubleValue() float64 {
	return o.Value().(float64)
}

func (o *Object) BoolValue() bool {
	return o.Value().(bool)
}

func (o *Object) StringValue() string {
	return o.Value().(string)
}

/**
 * ObjectMap
 */
type ObjectMap struct {
	Data map[string]*Object
}

func NewObjectMap() *ObjectMap {
	return &ObjectMap{
		Data: map[string]*Object{},
	}
}

func (m *ObjectMap) Set(k string, n *Object) {
	m.Data[strings.ToLower(k)] = n
}

func (m *ObjectMap) Get(k string) (*Object, bool) {
	n, ok := m.Data[strings.ToLower(k)]
	return n, ok
}

func (m *ObjectMap) All() map[string]*Object {
	return m.Data
}

var publicModifier = &ast.Modifier{Name: "public"}
var privateModifier = &ast.Modifier{Name: "private"}
var protectedModifier = &ast.Modifier{Name: "protected"}
var globalModifier = &ast.Modifier{Name: "global"}

func PublicModifier() *ast.Modifier {
	return publicModifier
}

func CreateClass(
	name string,
	constructors []*Method,
	instanceMethods *MethodMap,
	staticMethods *MethodMap,
) *ClassType {
	return &ClassType{
		Name:            name,
		Modifiers:       []ast.Node{PublicModifier()},
		Constructors:    constructors,
		InstanceFields:  NewFieldMap(),
		StaticFields:    NewFieldMap(),
		InstanceMethods: instanceMethods,
		StaticMethods:   staticMethods,
		InnerClasses:    NewClassMap(),
	}
}

func CreateMethod(
	name string,
	returnType ast.Node,
	parameters []ast.Node,
	nativeFunction func(*Object, []*Object, map[string]interface{}) interface{},
) *Method {
	return &Method{
		Name:           name,
		Modifiers:      []ast.Node{PublicModifier()},
		ReturnType:     returnType,
		Parameters:     parameters,
		NativeFunction: nativeFunction,
	}
}

func CreateField(
	name string,
	fieldType *ast.TypeRef,
) *Field {
	return &Field{
		Name:      name,
		Modifiers: []ast.Node{PublicModifier()},
		Type:      fieldType,
	}
}
