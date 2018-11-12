package main

import (
	"fmt"
	"os"

	"flag"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/tzmfreedom/goland/ast"
	"github.com/tzmfreedom/goland/compiler"
	"github.com/tzmfreedom/goland/parser"
	"github.com/tzmfreedom/goland/visitor"
)

func main() {
	f := flag.String("f", "", "file")
	_ = flag.String("d", "", "directory")
	cmd := os.Args[1]
	os.Args = os.Args[1:]

	flag.Parse()

	t := parseFile(*f)
	switch cmd {
	case "format":
		tos(t)
	case "run":
		root, err := convert(t)
		if err != nil {
			handleError(err)
		}
		err = run(root)
		if err != nil {
			handleError(err)
		}
	case "check":
		root, err := convert(t)
		if err != nil {
			handleError(err)
		}
		err = check(root)
		if err != nil {
			handleError(err)
		}
	}
}

func parseFile(f string) ast.Node {
	input, err := antlr.NewFileStream(f)
	if err != nil {
		handleError(err)
	}
	return parse(input, f)
}

func parseString(c string) ast.Node {
	input := antlr.NewInputStream(c)
	return parse(input, "")
}

func parse(input antlr.CharStream, f string) ast.Node {
	lexer := parser.NewapexLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := parser.NewapexParser(stream)
	p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))
	p.BuildParseTrees = true
	tree := p.CompilationUnit()
	t := tree.Accept(&AstBuilder{
		CurrentFile: f,
	})
	return t.(ast.Node)
}

func convert(n ast.Node) (ast.Node, error) {
	return n, nil
}

func check(n ast.Node) error {
	checker := &visitor.SoqlChecker{}
	_, err := n.Accept(checker)
	return err
}

func semantic_analysis(n ast.Node) error {
	register := &compiler.ClassRegisterVisitor{}
	t, err := n.Accept(register)
	if err != nil {
		return err
	}
	classTypes := make([]compiler.ClassType, 1)
	if tp, ok := t.(compiler.ClassType); ok {
		classTypes[1] = tp
	}
	typeChecker := &compiler.TypeChecker{
		ClassTypes: classTypes,
	}
	_, err = n.Accept(typeChecker)
	if err != nil {
		return err
	}
	return nil
}

func run(n ast.Node) error {
	interpreter := &compiler.Interpreter{}
	_, err := n.Accept(interpreter)
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
