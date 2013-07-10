package gohaml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStackHasLenOf0(t *testing.T) {
	s := stack{}
	assert.Equal(t, s.len(), 0)
}

func TestPushIncreasesLenBy1(t *testing.T) {
	s := stack{}
	n := fakeNode{}

	for x := 0; x < 10; x += 1 {
		s.push(&n)
		assert.Equal(t, s.len(), x + 1)
	}
}

func TestPopDecreasesLenToLenMinus1(t *testing.T) {
	s := stack{}
	n := fakeNode{}

	for x := 0; x < 10; x += 1 {
		s.push(&n)
		assert.Equal(t, s.len(), x + 1)
	}

	for x := 0; x < 10; x += 1 {
		s.pop()
		assert.Equal(t, s.len(), 9 - x)
	}
}

func TestPopIsLifo(t *testing.T) {
	s := stack{}
	nodes := []Node{&fakeNode{}, &fakeNode{}, &fakeNode{}}

	for _, node := range nodes {
		s.push(node)
	}

	for i := range nodes {
		n, _ := s.pop()
		assert.Equal(t, n, nodes[len(nodes) - i - 1])
	}

	assert.Equal(t, s.len(), 0)
}

func TestErrorPoppingOnEmptyStack(t *testing.T) {
	s := stack{}
	n, e := s.pop()

	assert.Nil(t, n)
	assert.NotNil(t, e)
}
