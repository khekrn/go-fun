package caching

import "testing"

var fa int

func BenchmarkLinkedListTraversal(b *testing.B) {
	var a int
	for i := 0; i < b.N; i++ {
		a = LinkedListTraversal()
	}
	fa = a
}

func BenchmarkColumnTraverse(b *testing.B) {
	var a int
	for i := 0; i < b.N; i++ {
		a = ColumnTraverse()
	}
	fa = a
}

func BenchmarkRowTraverse(b *testing.B) {
	var a int
	for i := 0; i < b.N; i++ {
		a = RowTraverse()
	}
	fa = a
}
