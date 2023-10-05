package yqlib

import (
	"container/list"
	"fmt"
)

func isKeyOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- isKeyOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		results.PushBack(createBooleanCandidate(candidate, candidate.IsMapKey))
	}

	return context.ChildContext(results), nil
}

func getKeyOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- getKeyOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		if candidate.Key != nil {
			results.PushBack(candidate.Key)
		}
	}

	return context.ChildContext(results), nil

}

func keysOperator(d *dataTreeNavigator, context Context, expressionNode *ExpressionNode) (Context, error) {
	log.Debugf("-- keysOperator")

	var results = list.New()

	for el := context.MatchingNodes.Front(); el != nil; el = el.Next() {
		candidate := el.Value.(*CandidateNode)

		var targetNode *CandidateNode
		if candidate.Kind == MappingNode {
			targetNode = getMapKeys(candidate)
		} else if candidate.Kind == SequenceNode {
			targetNode = getIndices(candidate)
		} else {
			return Context{}, fmt.Errorf("Cannot get keys of %v, keys only works for maps and arrays", candidate.Tag)
		}

		// result := candidate.CreateReplacement(targetNode)
		results.PushBack(targetNode)
	}

	return context.ChildContext(results), nil
}

func getMapKeys(node *CandidateNode) *CandidateNode {
	contents := make([]*CandidateNode, 0)
	for index := 0; index < len(node.Content); index = index + 2 {
		contents = append(contents, node.Content[index])
	}
	return &CandidateNode{Kind: SequenceNode, Tag: "!!seq", Content: contents}
}

func getIndices(node *CandidateNode) *CandidateNode {
	var contents = make([]*CandidateNode, len(node.Content))

	for index := range node.Content {
		contents[index] = &CandidateNode{
			Kind:  ScalarNode,
			Tag:   "!!int",
			Value: fmt.Sprintf("%v", index),
		}
	}

	return &CandidateNode{Kind: SequenceNode, Tag: "!!seq", Content: contents}
}
