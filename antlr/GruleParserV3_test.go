//  Copyright hyperjumptech/grule-rule-engine Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package antlr

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/hyperjumptech/grule-rule-engine/antlr/parser/grulev3"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
)

const (
	rules = `
rule SpeedUp "When testcar is speeding up we keep increase the speed." salience 10 {
    when
        TestCar.SpeedUp == true && TestCar.Speed < TestCar.MaxSpeed
    then
        TestCar.Speed = TestCar.Speed + TestCar.SpeedIncrement;
		DistanceRecord.TotalDistance = DistanceRecord.TotalDistance + TestCar.Speed;
}

rule StartSpeedDown "When testcar is speeding up and over max speed we change to speed down." salience 10  {
    when
        TestCar.SpeedUp == true && TestCar.Speed >= TestCar.MaxSpeed
    then
        TestCar.SpeedUp = false;
		Log("Now we slow down");
}

rule SlowDown "When testcar is slowing down we keep decreasing the speed." salience 10  {
    when
        TestCar.SpeedUp == false && TestCar.Speed > 0
    then
        TestCar.Speed = TestCar.Speed - TestCar.SpeedIncrement;
		DistanceRecord.TotalDistance = DistanceRecord.TotalDistance + TestCar.Speed;
}

rule SetTime "When Distance Recorder time not set, set it." {
	when
		IsZero(DistanceRecord.TestTime)
	then
		ABC.abc = "cde";
		Log("Set the test time");
		DistanceRecord.TestTime = Now();
		Log(TimeFormat(DistanceRecord.TestTime,"Mon Jan _2 15:04:05 2006"));
		ABC.abc = "cde";

}
`
	invalidEscapeRule = `rule SetTime "When Distance Recorder time not set, set it." {
	when
		IsZero(DistanceRecord.TestTime)
	then
		ABC.abc = "abc\cde";
		Log("Set the test time");
		DistanceRecord.TestTime = Now();
		Log(TimeFormat(DistanceRecord.TestTime,"Mon Jan _2 15:04:05 2006"));

}`
	validEscapeRule = `rule SetTime "When Distance Recorder time not set, set it." {
	when
		IsZero(DistanceRecord.TestTime)
	then
		ABC.abc = "abc\\cde";
		Log("Set the test time");
		DistanceRecord.TestTime = Now();
		Log(TimeFormat(DistanceRecord.TestTime,"Mon Jan _2 15:04:05 2006"));

}`
)

type Person struct {
	Name         string
	ParentString string
	Child        *Child
}

type Child struct {
	Name        string
	ChildString string
	GrandChild  *GrandChild
}

type GrandChild struct {
	GrandChildString string
	Name             string
}

func TestV3Lexer(t *testing.T) {
	data, err := os.ReadFile("./sample4.grl")
	if err != nil {
		t.Fatal(err)
	} else {
		is := antlr.NewInputStream(string(data))

		// Create the Lexer
		lexer := parser.Newgrulev3Lexer(is)
		//lexer := parser.NewLdifParserLexer(is)

		// Read all tokens
		for {
			nt := lexer.NextToken()
			fmt.Println(nt.GetText())
			if nt.GetTokenType() == antlr.TokenEOF {
				break
			}
		}
	}

}

func TestV3Parser(t *testing.T) {
	// logrus.SetLevel(logrus.TraceLevel)
	data, err := os.ReadFile("./sample4.grl")
	if err != nil {
		t.Fatal(err)
	} else {
		sdata := string(data)

		is := antlr.NewInputStream(sdata)
		lexer := parser.Newgrulev3Lexer(is)
		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

		errReporter := &pkg.GruleErrorReporter{
			Errors: make([]error, 0),
		}

		kb := ast.NewKnowledgeLibrary().GetKnowledgeBase("T", "1")

		listener := NewGruleV3ParserListener(kb, errReporter)

		psr := parser.Newgrulev3Parser(stream)
		psr.BuildParseTrees = true
		antlr.ParseTreeWalkerDefault.Walk(listener, psr.Grl())

		if errReporter.HasError() {
			t.Log(errReporter.Error())
			t.FailNow()
		}
	}
}

func TestV3Parser2(t *testing.T) {
	// logrus.SetLevel(logrus.InfoLevel)

	sdata := rules

	is := antlr.NewInputStream(sdata)
	lexer := parser.Newgrulev3Lexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	errReporter := &pkg.GruleErrorReporter{
		Errors: make([]error, 0),
	}
	kb := ast.NewKnowledgeLibrary().GetKnowledgeBase("T", "1")

	listener := NewGruleV3ParserListener(kb, errReporter)

	psr := parser.Newgrulev3Parser(stream)
	psr.BuildParseTrees = true
	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Grl())

	if errReporter.HasError() {
		t.Log(errReporter.Error())
		t.FailNow()
	}
}

func TestV3ParserEscapedStringInvalid(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)

	is := antlr.NewInputStream(invalidEscapeRule)
	lexer := parser.Newgrulev3Lexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	errReporter := &pkg.GruleErrorReporter{
		Errors: make([]error, 0),
	}
	kb := ast.NewKnowledgeLibrary().GetKnowledgeBase("T", "1")

	listener := NewGruleV3ParserListener(kb, errReporter)

	psr := parser.Newgrulev3Parser(stream)
	psr.BuildParseTrees = true
	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Grl())

	if errReporter.HasError() == false {
		t.Fatal("Successfully parsed invalid string literal, should have gotten an error")
	}
}

func TestV3ParserEscapedStringValid(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)

	is := antlr.NewInputStream(validEscapeRule)
	lexer := parser.Newgrulev3Lexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	errReporter := &pkg.GruleErrorReporter{
		Errors: make([]error, 0),
	}
	kb := ast.NewKnowledgeLibrary().GetKnowledgeBase("T", "1")

	listener := NewGruleV3ParserListener(kb, errReporter)

	psr := parser.Newgrulev3Parser(stream)
	psr.BuildParseTrees = true
	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Grl())

	if errReporter.HasError() {
		t.Fatal("Failed to parse rule with escaped string constant")
	}
}

func TestV3ParserSnapshotEyeBalling(t *testing.T) {
	// logrus.SetLevel(logrus.TraceLevel)

	data := `
rule SpeedUp "When testcar is speeding up we keep increase the speed." salience 10 {
    when
        TestCar.SpeedUp == true && TestCar.Speed < TestCar.MaxSpeed
    then
        TestCar.Speed = TestCar.Speed + TestCar.SpeedIncrement;
		DistanceRecord.TotalDistance = DistanceRecord.TotalDistance + TestCar.Speed;
}
`

	is := antlr.NewInputStream(data)
	lexer := parser.Newgrulev3Lexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	errReporter := &pkg.GruleErrorReporter{
		Errors: make([]error, 0),
	}
	kb := ast.NewKnowledgeLibrary().GetKnowledgeBase("T", "1")

	listener := NewGruleV3ParserListener(kb, errReporter)

	psr := parser.Newgrulev3Parser(stream)
	psr.BuildParseTrees = true
	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Grl())
	if errReporter.HasError() {
		t.Log(errReporter.Error())
		t.FailNow()
	}

	listener.KnowledgeBase.WorkingMemory.IndexVariables()
	if listener.Grl.RuleEntries["SpeedUp"] == nil {
		t.Fatalf("Rule entry not exist")
	}
	if listener.Grl.RuleEntries["SpeedUp"].WhenScope == nil {
		t.Fatalf("When scope not exist")
	}
	if listener.Grl.RuleEntries["SpeedUp"].ThenScope == nil {
		t.Fatalf("Then scope not exist")
	}
	t.Log(listener.Grl.RuleEntries["SpeedUp"].GetSnapshot())
	if listener.Grl.RuleEntries["SpeedUp"].ThenScope.ThenExpressionList == nil {
		t.Fatalf("Then expression list is not exist")
	}
	if listener.Grl.RuleEntries["SpeedUp"].ThenScope.ThenExpressionList.ThenExpressions == nil {
		t.Fatalf("Then expression list array is not exist")
	}
	if len(listener.Grl.RuleEntries["SpeedUp"].ThenScope.ThenExpressionList.ThenExpressions) != 2 {
		t.Fatalf("Then expression list array %s contains not 2 but %d", listener.Grl.RuleEntries["SpeedUp"].ThenScope.ThenExpressionList.GetAstID(), len(listener.Grl.RuleEntries["SpeedUp"].ThenScope.ThenExpressionList.ThenExpressions))
	}
}

func prepareV3TestKnowledgeBase(t *testing.T, grl string) (*ast.KnowledgeBase, *ast.WorkingMemory) {
	is := antlr.NewInputStream(grl)
	lexer := parser.Newgrulev3Lexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	errReporter := &pkg.GruleErrorReporter{
		Errors: make([]error, 0),
	}
	kb := ast.NewKnowledgeLibrary().GetKnowledgeBase("T", "1")

	listener := NewGruleV3ParserListener(kb, errReporter)

	psr := parser.Newgrulev3Parser(stream)
	psr.BuildParseTrees = true

	psr.RemoveErrorListeners()
	psr.AddErrorListener(errReporter)

	antlr.ParseTreeWalkerDefault.Walk(listener, psr.Grl())
	assert.False(t, errReporter.HasError())
	listener.KnowledgeBase.WorkingMemory.IndexVariables()
	return kb, kb.WorkingMemory
}

func TestV3ConstantFunctionAndConstantFunctionChain(t *testing.T) {
	// logrus.SetLevel(logrus.InfoLevel)
	dctx := ast.NewDataContext()

	data := `
rule RuleOne "RuleOneDesc" salience 123 {
    when
        "    true    ".Trim() != "true"
    then
        Log("      Success    ".Trim().ToUpper());
}
`
	kb, wm := prepareV3TestKnowledgeBase(t, data)
	err := dctx.Add("DEFUNC", &ast.BuiltInFunctions{
		Knowledge:     kb,
		WorkingMemory: wm,
		DataContext:   dctx,
	})

	assert.NoError(t, err)

	whenVal, err := kb.RuleEntries["RuleOne"].WhenScope.Evaluate(dctx, wm)
	assert.NoError(t, err)
	assert.True(t, whenVal.IsValid())
	assert.Equal(t, reflect.Bool, whenVal.Kind())
	assert.False(t, whenVal.Bool())
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions)
	assert.Equal(t, 1, len(kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions))
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions[0])
	err = kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions[0].Execute(dctx, wm)
	assert.NoError(t, err)
}

func TestV3RuleRetract(t *testing.T) {
	// logrus.SetLevel(logrus.InfoLevel)
	dctx := ast.NewDataContext()

	data := `
rule RuleOne "RuleOneDesc" salience 123 {
    when
        "    true    ".Trim() != "true"
    then
        Retract("RuleOne");
}
`

	kb, wm := prepareV3TestKnowledgeBase(t, data)
	err := dctx.Add("DEFUNC", &ast.BuiltInFunctions{
		Knowledge:     kb,
		WorkingMemory: wm,
		DataContext:   dctx,
	})

	assert.False(t, kb.RuleEntries["RuleOne"].Retracted)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions)
	assert.Equal(t, 1, len(kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions))
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions[0])

	err = kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions[0].Execute(dctx, wm)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	assert.True(t, kb.RuleEntries["RuleOne"].Retracted)
	kb.Reset()
	assert.False(t, kb.RuleEntries["RuleOne"].Retracted)
}

func TestV3RuleAssignment(t *testing.T) {
	// logrus.SetLevel(logrus.TraceLevel)
	dctx := ast.NewDataContext()

	data := `
rule RuleOne "RuleOneDesc" salience 123 {
    when
        Person.Name == "Rudolf"
    then
        Person.Name="Pearson".ToUpper();
}
`

	kb, wm := prepareV3TestKnowledgeBase(t, data)
	err := dctx.Add("DEFUNC", &ast.BuiltInFunctions{
		Knowledge:     kb,
		WorkingMemory: wm,
		DataContext:   dctx,
	})

	p := &Person{
		Name: "Rudolf",
	}
	err = dctx.Add("Person", p)

	assert.NoError(t, err)
	assert.False(t, kb.RuleEntries["RuleOne"].Retracted)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList)
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions)

	ret, err := kb.RuleEntries["RuleOne"].WhenScope.Evaluate(dctx, wm)
	assert.NoError(t, err)
	assert.True(t, ret.IsValid())
	assert.Equal(t, reflect.Bool, ret.Kind())
	assert.True(t, ret.Bool())

	assert.Equal(t, 1, len(kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions))
	assert.NotNil(t, kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions[0])
	err = kb.RuleEntries["RuleOne"].ThenScope.ThenExpressionList.ThenExpressions[0].Execute(dctx, wm)
	assert.NoError(t, err)
	assert.Equal(t, "PEARSON", p.Name)

	ret, err = kb.RuleEntries["RuleOne"].WhenScope.Evaluate(dctx, wm)
	assert.NoError(t, err)
	assert.True(t, ret.IsValid())
	assert.Equal(t, reflect.Bool, ret.Kind())
	assert.False(t, ret.Bool())

}

func TestV3ParserGarbageInput(t *testing.T) {
	testCases := []string{
		`rule ExampleRuleName "One line rule description" salience 0  {
  when
    <FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" when 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" salience 0  {
  salience
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" salience 0  {
  then
    FACT.SomeField == "true"
  when
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" then 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" salience 0 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" salience "0"  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`then ExampleRuleName "One line rule description" salience 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" rule 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" salience 0  {
  when
    FACT.SomeField == "true"
  rule
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" salience 0  {
  when
    FACT.SomeField == "true"
  then
    rule FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName "One line rule description" salience 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" ); salience 2
}`,
		`rule ExampleRuleName then "One line rule description" salience 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
}`,
		`rule ExampleRuleName  "One line rule description" salience 0  {
  when
    FACT.SomeField == "true"
  then
    FACT.WriteMessage( "message to write" );
  then
    FACT.WriteMessage( "message to write" );
}`,
	}

	for i, tc := range testCases {
		is := antlr.NewInputStream(tc)
		lexer := parser.Newgrulev3Lexer(is)
		stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

		errReporter := &pkg.GruleErrorReporter{
			Errors: make([]error, 0),
		}
		kb := ast.NewKnowledgeLibrary().GetKnowledgeBase("T", "1")

		listener := NewGruleV3ParserListener(kb, errReporter)

		psr := parser.Newgrulev3Parser(stream)
		psr.BuildParseTrees = true

		psr.RemoveErrorListeners()
		psr.AddErrorListener(errReporter)

		antlr.ParseTreeWalkerDefault.Walk(listener, psr.Grl())
		assert.True(t, errReporter.HasError(), fmt.Sprintf("Garbage test case %d should have resulted in a parse error", i))
	}
}
