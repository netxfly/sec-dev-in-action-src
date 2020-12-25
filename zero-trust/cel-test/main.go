package main

import (
	"fmt"
	"log"

	"github.com/Knetic/govaluate"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func main() {
	//err := testCel1()
	//fmt.Printf("call testCel1(), return value: %v\n", err)

	//err = testCel2()
	//fmt.Printf("call testCel2(), return value: %v\n", err)

	err := testCelFunc()
	fmt.Printf("call testCelFunc(), return value: %v\n", err)

	//err := testGoValuate()
	//fmt.Printf("call testGoValuate, return value: %v\n", err)
}

func testCel1() error {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("name", decls.String),
			decls.NewVar("group", decls.String),
			decls.NewVar("site", decls.String),
		),
	)
	if err != nil {
		return err
	}

	ast, issues := env.Compile(`name.startsWith("/group/" + group)`)
	if issues != nil && issues.Err() != nil {
		log.Panicf("ast: %v, issues: %v\n", ast, issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return err
	}

	value := make(map[string]interface{})
	value["name"] = "/group/xsec.io/type/sec"
	value["group"] = "xsec.io"
	value["site"] = "sec.lu"

	out, detail, err := prg.Eval(value)
	fmt.Printf("out: %v, detail: %v, err: %v\n", out, detail, err)
	return err
}

func testCel2() error {
	env, err := cel.NewEnv(
		cel.Declarations(
			decls.NewVar("name", decls.String),
			decls.NewVar("group", decls.String),
			decls.NewVar("site", decls.String),
		),
	)
	if err != nil {
		return err
	}
	exp := `name.startsWith("/group/" + group)`
	parsed, _ := env.Parse(exp)
	ast, issues := env.Check(parsed)
	if issues != nil && issues.Err() != nil {
		log.Panicf("ast: %v, issues: %v\n", ast, issues.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return err
	}

	value := make(map[string]interface{})
	value["name"] = "/group/xsec.io/type/sec"
	value["group"] = "xsec.io"
	value["site"] = "sec.lu"

	out, detail, err := prg.Eval(value)
	fmt.Printf("out: %v, detail: %v, err: %v\n", out, detail, err)
	return err
}

func testCelFunc() error {
	dec := cel.Declarations(
		decls.NewVar("i", decls.String),
		decls.NewVar("you", decls.String),
		decls.NewFunction("func_test", decls.NewOverload("func_test_string_string",
			[]*exprpb.Type{decls.String, decls.String},
			decls.String,
		),
		),
	)

	testFunc := &functions.Overload{
		Operator:     "func_test_string_string",
		OperandTrait: 0,
		Unary:        nil,
		Function:     nil,
		Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
			return types.String(fmt.Sprintf("%v请%v吃饭啊", lhs, rhs))
		},
	}

	env, err := cel.NewEnv(dec)
	if err != nil {
		return err
	}

	ast, iss := env.Compile(`func_test(i, you)`)
	if iss.Err() != nil {
		return err
	}

	prg, err := env.Program(ast, cel.Functions(testFunc))
	if err != nil {
		return err
	}

	out, _, err := prg.Eval(map[string]interface{}{
		"i":   "我",
		"you": "你",
	})

	fmt.Printf("out: %v, err: %v\n", out, err)
	return err
}

func testGoValuate() error {
	expression, err := govaluate.NewEvaluableExpression("最近14天未出京 && (体温<=36.1) && (AreYouOk == 'ok')")
	if err != nil {
		return err
	}

	personA := make(map[string]interface{})
	personA["姓名"] = "张飞"
	personA["最近14天未出京"] = true
	personA["体温"] = 35.5
	personA["AreYouOk"] = "ok"

	personB := make(map[string]interface{})
	personB["姓名"] = "李逵"
	personB["最近14天未出京"] = false
	personB["体温"] = 35.5
	personA["AreYouOk"] = "ok"

	result, err := expression.Evaluate(personA)
	fmt.Printf("check %v, result: %v, err: %v\n", personA["姓名"], result, err)

	result, err = expression.Evaluate(personB)
	fmt.Printf("check %v, result: %v, err: %v\n", personB["姓名"], result, err)

	// 以下为函数表过式的测试
	functions := map[string]govaluate.ExpressionFunction{
		"装备": func(arguments ...interface{}) (interface{}, error) {
			arg := arguments[0].(string)
			result := ""
			switch arg {
			case "张飞":
				// result = "丈八蛇矛"
				result = "码农双肩包"
			case "李逵":
				result = "2把大板斧"
			case "吕布":
				result = "方天画戟"
			}
			return result, nil
		},
	}

	expressionFunc, _ := govaluate.NewEvaluableExpressionWithFunctions("最近14天未出京 && "+
		"(体温<=36.1) && (AreYouOk == 'ok') && (装备(姓名) == '码农双肩包')",
		functions)

	result, err = expressionFunc.Evaluate(personA)
	fmt.Printf("check %v, result: %v, err: %v\n", personA["姓名"], result, err)

	result, err = expressionFunc.Evaluate(personB)
	fmt.Printf("check %v, result: %v, err: %v\n", personB["姓名"], result, err)

	return err
}
