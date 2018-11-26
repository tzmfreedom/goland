package main

import (
	"errors"
	"fmt"
	"os"

	"flag"
	"strings"

	"io/ioutil"

	"github.com/tzmfreedom/goland/ast"
	"github.com/tzmfreedom/goland/builtin"
	"github.com/tzmfreedom/goland/compiler"
	"github.com/tzmfreedom/goland/interpreter"
	"github.com/tzmfreedom/goland/visitor"
)

var classMap = builtin.NewClassMap()
var preprocessors = []ast.PreProcessor{
	func(src string) string {
		return strings.Replace(src, "// #debugger", "Debugger.debug();", -1)
	},
}

type option struct {
	SubCommand string
	Action     string
	Files      []string
}

func main() {
	option, err := parseOption(os.Args)
	if err != nil {
		handleError(err)
		return
	}

	trees := make([]ast.Node, len(option.Files))
	for i, file := range option.Files {
		trees[i], err = ast.ParseFile(file, preprocessors...)
		if err != nil {
			handleError(err)
		}
	}
	switch option.SubCommand {
	case "format":
		for _, t := range trees {
			tos(t)
		}
	case "run":
		classTypes := make([]*builtin.ClassType, len(trees))
		for i, t := range trees {
			root, err := convert(t)
			if err != nil {
				handleError(err)
			}
			classTypes[i], err = register(root)
		}
		for _, t := range classTypes {
			if err = semanticAnalysis(t); err != nil {
				handleError(err)
			}
		}
		err = run(option.Action, classTypes)
		if err != nil {
			handleError(err)
		}
	case "check":
		newTrees := make([]*builtin.ClassType, len(trees))
		for i, t := range trees {
			root, err := convert(t)
			if err != nil {
				handleError(err)
			}
			newTrees[i], err = register(root)
		}
		for _, t := range newTrees {
			err = semanticAnalysis(t)
			if err != nil {
				handleError(err)
			}
		}
	}
}

func parseOption(args []string) (*option, error) {
	flg := flag.NewFlagSet(args[0], flag.ExitOnError)
	fileName := flg.String("f", "", "file")
	directory := flg.String("d", "", "directory")
	action := flg.String("a", "", "action")

	err := flg.Parse(args[2:])
	if err != nil {
		return nil, err
	}

	if fileName == nil && directory == nil {
		return nil, errors.New("-f FILE or -d DIRECTORY is required")
	}

	if *action == "" {
		return nil, errors.New("-a CLASS#METHOD is required")
	}
	var files []string
	if fileName != nil {
		files = []string{*fileName}
	} else {
		filesInDirectory, err := ioutil.ReadDir(*directory)
		if err != nil {
			handleError(err)
		}
		files = []string{}
		for _, f := range filesInDirectory {
			if f.IsDir() {
				continue
			}
			files = append(files, f.Name())
		}
	}
	return &option{
		SubCommand: os.Args[1],
		Action:     *action,
		Files:      files,
	}, nil
}

func convert(n ast.Node) (ast.Node, error) {
	return n, nil
}

func check(n ast.Node) error {
	checker := &visitor.SoqlChecker{}
	_, err := n.Accept(checker)
	return err
}

func register(n ast.Node) (*builtin.ClassType, error) {
	register := &compiler.ClassRegisterVisitor{}
	t, err := n.Accept(register)
	if err != nil {
		return nil, err
	}
	classType := t.(*builtin.ClassType)
	classMap.Set(classType.Name, classType)
	return classType, nil
}

func semanticAnalysis(t *builtin.ClassType) error {
	typeChecker := compiler.NewTypeChecker()
	typeChecker.Context.ClassTypes = builtin.PrimitiveClassMap()
	typeChecker.Context.ClassTypes.Set(t.Name, t)
	_, err := typeChecker.VisitClassType(t)
	if len(typeChecker.Errors) != 0 {
		for _, e := range typeChecker.Errors {
			fmt.Fprintf(os.Stderr, "%s\n", e.Message)
		}
	}
	return err
}

func run(action string, classTypes []*builtin.ClassType) error {
	method := "action"
	args := strings.Split(action, "#")
	if len(args) > 1 {
		method = args[1]
	}
	interpreter := interpreter.NewInterpreter(builtin.PrimitiveClassMap())
	for _, classType := range classTypes {
		interpreter.Context.ClassTypes.Set(classType.Name, classType)
	}
	invoke := &ast.MethodInvocation{
		NameOrExpression: &ast.Name{
			Value: []string{args[0], method},
		},
	}
	_, err := invoke.Accept(interpreter)
	return err
}

func tos(n ast.Node) {
	visitor := &ast.TosVisitor{}
	r, _ := n.Accept(visitor)
	fmt.Println(r)
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}

func validate() {
	return
}
