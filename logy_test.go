package logx

import "testing"

func init() {
	std.SetJsonFormat(true).SetColor(true)
}

func BenchmarkInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("haha")
	}
}
