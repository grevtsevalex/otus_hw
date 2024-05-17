package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkStats(b *testing.B) {
	b.Helper()
	b.StopTimer()

	r, _ := zip.OpenReader("testdata/users.dat.zip")
	defer r.Close()

	data, _ := r.File[0].Open()
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		GetDomainStat(data, "biz")
		b.StopTimer()
	}
}
