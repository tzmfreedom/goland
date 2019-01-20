package builtin

import "github.com/tzmfreedom/goland/ast"

var httpRequestType = &ast.ClassType{Name: "HttpRequest"}
var httpRequestTypeParameter = &ast.Parameter{
	Type: httpRequestType,
	Name: "_",
}

func init() {
	instanceMethods := ast.NewMethodMap()
	staticMethods := ast.NewMethodMap()
	httpRequestType.Constructors = []*ast.Method{
		ast.CreateMethod(
			"HttpRequest",
			nil,
			[]*ast.Parameter{},
			func(this *ast.Object, params []*ast.Object, extra map[string]interface{}) interface{} {
				this.Extra["headers"] = map[string]*ast.Object{}
				return nil
			},
		),
	}
	httpRequestType.InstanceMethods = instanceMethods
	httpRequestType.StaticMethods = staticMethods

	instanceMethods.Set(
		"setHeader",
		[]*ast.Method{
			ast.CreateMethod(
				"setHeader",
				nil,
				[]*ast.Parameter{
					stringTypeParameter,
					stringTypeParameter,
				},
				func(this *ast.Object, params []*ast.Object, extra map[string]interface{}) interface{} {
					key := params[0].StringValue()
					headers := this.Extra["headers"].(map[string]*ast.Object)
					headers[key] = params[1]
					return nil
				},
			),
		},
	)
	instanceMethods.Set(
		"setMethod",
		[]*ast.Method{
			ast.CreateMethod(
				"setMethod",
				nil,
				[]*ast.Parameter{
					stringTypeParameter,
				},
				func(this *ast.Object, params []*ast.Object, extra map[string]interface{}) interface{} {
					method := params[0].StringValue()
					this.Extra["method"] = method
					return nil
				},
			),
		},
	)
	instanceMethods.Set(
		"setEndpoint",
		[]*ast.Method{
			ast.CreateMethod(
				"setEndpoint",
				nil,
				[]*ast.Parameter{
					stringTypeParameter,
				},
				func(this *ast.Object, params []*ast.Object, extra map[string]interface{}) interface{} {
					endpoint := params[0].StringValue()
					this.Extra["endpoint"] = endpoint
					return nil
				},
			),
		},
	)

	primitiveClassMap.Set("HttpRequest", httpRequestType)
}
