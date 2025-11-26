package plygo

import (
	"fmt"
	"testing"
)

type BenchPerson struct {
	Name   string
	Age    int
	City   string
	Salary float64
	Active bool
}

func generateBenchData(n int) []BenchPerson {
	cities := []string{"NYC", "LA", "Chicago", "Miami", "Boston"}
	names := []string{"Alice", "Bob", "Charlie", "Diana", "Eve", "Frank", "Grace", "Henry"}

	data := make([]BenchPerson, n)
	for i := 0; i < n; i++ {
		data[i] = BenchPerson{
			Name:   fmt.Sprintf("%s%d", names[i%len(names)], i),
			Age:    25 + (i % 40),
			City:   cities[i%len(cities)],
			Salary: float64(50000 + (i%50)*1000),
			Active: i%3 != 0,
		}
	}
	return data
}

// Benchmark basic filtering
func BenchmarkWhere_Simple(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Where("Age").GreaterThan(30).Collect()
	}
}

func BenchmarkWhere_Multiple(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).
			Where("Age").GreaterThan(30).
			Where("Active").IsTrue().
			Where("Salary").LessThan(80000).
			Collect()
	}
}

func BenchmarkWhere_Or(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).
			Where("City").Equals("NYC").Or("City").Equals("LA").
			Collect()
	}
}

func BenchmarkWhere_OneOf(b *testing.B) {
	data := generateBenchData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).
			Where("City").OneOf("NYC", "LA", "Chicago").
			Collect()
	}
}

func BenchmarkWhere_StringOps(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).
			Where("Name").StartsWith("Alice").
			Collect()
	}
}

// Benchmark sorting
func BenchmarkOrderBy_Single(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).OrderBy("Age").Desc().Collect()
	}
}

func BenchmarkOrderBy_Multiple(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).
			OrderBy("City").Asc().
			ThenBy("Age").Desc().
			Collect()
	}
}

// Benchmark grouping
func BenchmarkGroupBy_Count(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).GroupBy("City").Count()
	}
}

func BenchmarkGroupBy_Sum(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).GroupBy("City").Sum("Salary")
	}
}

// Benchmark complex pipelines
func BenchmarkPipeline_Complex(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Where("Age").GreaterThan(25).
			OrderBy("Salary").Desc().
			Collect()
		_ = result
	}
}

// Benchmark field operations (reflection-heavy)
func BenchmarkFieldNames(b *testing.B) {
	data := generateBenchData(100)
	p := From(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.FieldNames()
	}
}

// Benchmark distinct
func BenchmarkDistinct(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Distinct("City").Collect()
	}
}

// Benchmark transform
func BenchmarkTransform(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Transform(func(p BenchPerson) BenchPerson {
			p.Salary *= 1.1
			return p
		}).Collect()
	}
}

// Benchmark select (creates maps)
func BenchmarkSelect(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Select("Name", "Salary").Collect()
	}
}

// Benchmark position operations
func BenchmarkAtRow(b *testing.B) {
	data := generateBenchData(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).AtRow(1, 10, 100, 500, 999).Collect()
	}
}

// Different data sizes
func BenchmarkWhere_Size100(b *testing.B) {
	data := generateBenchData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Where("Age").GreaterThan(30).Collect()
	}
}

func BenchmarkWhere_Size10000(b *testing.B) {
	data := generateBenchData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).Where("Age").GreaterThan(30).Collect()
	}
}

func BenchmarkOrderBy_Size100(b *testing.B) {
	data := generateBenchData(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).OrderBy("Age").Desc().Collect()
	}
}

func BenchmarkOrderBy_Size10000(b *testing.B) {
	data := generateBenchData(10000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).OrderBy("Age").Desc().Collect()
	}
}
