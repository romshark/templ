package main

import (
	"bytes"
	"context"
	"testing"
	"tst/fork"
	"tst/orig"
)

func BenchmarkRender(b *testing.B) {
	ctx := context.Background()
	var buffer bytes.Buffer
	buffer.Grow(16 * 1024)

	b.Run("orig", func(b *testing.B) {
		for b.Loop() {
			buffer.Reset()
			orig.WithFmt(42).Render(ctx, &buffer)
		}
	})

	b.Run("orig-concat", func(b *testing.B) {
		for b.Loop() {
			buffer.Reset()
			orig.Concat(42).Render(ctx, &buffer)
		}
	})

	b.Run("orig-strbuild", func(b *testing.B) {
		for b.Loop() {
			buffer.Reset()
			orig.Strbuild(42).Render(ctx, &buffer)
		}
	})

	b.Run("fork", func(b *testing.B) {
		for b.Loop() {
			buffer.Reset()
			fork.WithFmt(42).Render(ctx, &buffer)
		}
	})
}
