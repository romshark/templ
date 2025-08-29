package main

import (
	"bytes"
	"context"
	"testing"
)

func BenchmarkRender(b *testing.B) {
	ctx := context.Background()
	var buffer bytes.Buffer
	buffer.Grow(16 * 1024)
	for b.Loop() {
		buffer.Reset()
		WithFmt(42).Render(ctx, &buffer)
	}
}
