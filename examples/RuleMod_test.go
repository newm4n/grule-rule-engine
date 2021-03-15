package examples

import (
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	RuleA = `rule One {
	when true 
	then OK();
}`

	RuleB = `rule Two {
	when true 
	then OK();
}`
)

func TestRuleAdd(t *testing.T) {
	knowledgeLibrary := ast.NewKnowledgeLibrary()
	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	byteArr := pkg.NewBytesResource([]byte(RuleA))
	err := ruleBuilder.BuildRuleFromResource("Test", "0.0.1", byteArr)
	assert.NoError(t, err)

	knowledgeBase := knowledgeLibrary.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.Equal(t, 1, len(knowledgeBase.RuleEntries))
	_, ok := knowledgeBase.RuleEntries["One"]
	assert.True(t, ok)
	_, ok = knowledgeBase.RuleEntries["Two"]
	assert.False(t, ok)

	byteArr = pkg.NewBytesResource([]byte(RuleB))
	err = ruleBuilder.BuildRuleFromResource("Test", "0.0.1", byteArr)
	assert.NoError(t, err)

	knowledgeBase = knowledgeLibrary.NewKnowledgeBaseInstance("Test", "0.0.1")
	assert.Equal(t, 2, len(knowledgeBase.RuleEntries))
	_, ok = knowledgeBase.RuleEntries["One"]
	assert.True(t, ok)
	_, ok = knowledgeBase.RuleEntries["Two"]
	assert.True(t, ok)
}
