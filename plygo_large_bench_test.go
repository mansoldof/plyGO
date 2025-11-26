package plygo

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

// DataFrame represents a realistic data structure with 10 columns
type DataFrame struct {
	ID         int
	Name       string
	Value1     float64
	Value2     float64
	Value3     float64
	Category   string
	Score      float64
	Percentage float64
	Count      int
	Active     bool
}

// Generate normally distributed data
func normalRandom(mean, stdDev float64) float64 {
	// Box-Muller transform for normal distribution
	u1 := rand.Float64()
	u2 := rand.Float64()
	z0 := math.Sqrt(-2.0*math.Log(u1)) * math.Cos(2.0*math.Pi*u2)
	return mean + stdDev*z0
}

// Generate 100,000 rows with normally distributed data
func generateLargeDataFrame(n int) []DataFrame {
	categories := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	names := []string{"Alpha", "Beta", "Gamma", "Delta", "Epsilon", "Zeta", "Eta", "Theta"}

	data := make([]DataFrame, n)
	for i := 0; i < n; i++ {
		data[i] = DataFrame{
			ID:         i + 1,
			Name:       fmt.Sprintf("%s_%d", names[i%len(names)], i),
			Value1:     normalRandom(100.0, 15.0), // mean=100, stddev=15
			Value2:     normalRandom(50.0, 10.0),  // mean=50, stddev=10
			Value3:     normalRandom(200.0, 30.0), // mean=200, stddev=30
			Category:   categories[i%len(categories)],
			Score:      normalRandom(75.0, 12.0),                   // mean=75, stddev=12
			Percentage: math.Abs(normalRandom(50.0, 20.0)),         // mean=50, stddev=20, always positive
			Count:      int(math.Abs(normalRandom(1000.0, 300.0))), // mean=1000, stddev=300
			Active:     i%3 != 0,                                   // ~66% active
		}
	}
	return data
}

// ===== Filtering Benchmarks =====

func BenchmarkLarge_Where_Single(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Value1").GreaterThan(100.0).
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Where_Multiple(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Value1").GreaterThan(90.0).
			Where("Value2").LessThan(60.0).
			Where("Active").IsTrue().
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Where_Complex(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Value1").Between(80.0, 120.0).
			Where("Score").GreaterThan(70.0).
			Where("Category").OneOf("A", "B", "C").
			Where("Active").IsTrue().
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Where_Or(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Category").Equals("A").Or("Category").Equals("B").Or("Category").Equals("C").
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Where_OneOf(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Category").OneOf("A", "B", "C", "D", "E").
			Collect()
		_ = result
	}
}

// ===== Sorting Benchmarks =====

func BenchmarkLarge_OrderBy_Single(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			OrderBy("Value1").Desc().
			Collect()
		_ = result
	}
}

func BenchmarkLarge_OrderBy_Multiple(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			OrderBy("Category").Asc().
			ThenBy("Value1").Desc().
			Collect()
		_ = result
	}
}

func BenchmarkLarge_OrderBy_ThreeFields(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			OrderBy("Category").Asc().
			ThenBy("Active").Desc().
			ThenBy("Score").Desc().
			Collect()
		_ = result
	}
}

// ===== Grouping Benchmarks =====

func BenchmarkLarge_GroupBy_Count(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			GroupBy("Category").
			Count()
		_ = result
	}
}

func BenchmarkLarge_GroupBy_Sum(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			GroupBy("Category").
			Sum("Value1")
		_ = result
	}
}

func BenchmarkLarge_GroupBy_Avg(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			GroupBy("Category").
			Avg("Score")
		_ = result
	}
}

// ===== Complex Pipeline Benchmarks =====

func BenchmarkLarge_Pipeline_FilterSort(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Where("Value1").GreaterThan(100.0).
			OrderBy("Score").Desc().
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Pipeline_FilterSortLimit(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Where("Value1").GreaterThan(90.0).
			OrderBy("Score").Desc().
			Limit(100).
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Pipeline_FilterGroupBy(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			GroupBy("Category").
			Sum("Value1")
		_ = result
	}
}

func BenchmarkLarge_Pipeline_Complex(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Where("Value1").Between(80.0, 120.0).
			Where("Score").GreaterThan(70.0).
			OrderBy("Category").Asc().
			ThenBy("Value1").Desc().
			Limit(1000).
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Pipeline_VeryComplex(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Value1").GreaterThan(85.0).
			Where("Value2").LessThan(60.0).
			Where("Score").Between(60.0, 90.0).
			Where("Category").OneOf("A", "B", "C", "D").
			Where("Active").IsTrue().
			OrderBy("Score").Desc().
			ThenBy("Value1").Asc().
			Limit(500).
			Collect()
		_ = result
	}
}

// ===== Transform Benchmarks =====

func BenchmarkLarge_Transform(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Transform(func(d DataFrame) DataFrame {
				d.Value1 *= 1.1
				d.Score *= 1.05
				return d
			}).
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Transform_AfterFilter(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Transform(func(d DataFrame) DataFrame {
				d.Value1 *= 1.1
				return d
			}).
			Collect()
		_ = result
	}
}

// ===== Selection Benchmarks =====

func BenchmarkLarge_Select(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Select("ID", "Name", "Value1", "Score").
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Select_AfterFilter(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Select("Name", "Value1", "Score").
			Collect()
		_ = result
	}
}

// ===== Distinct Benchmark =====

func BenchmarkLarge_Distinct(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Distinct("Category").
			Collect()
		_ = result
	}
}

// ===== Position-based Benchmarks =====

func BenchmarkLarge_AtRow(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			AtRow(1, 1000, 10000, 50000, 99999).
			Collect()
		_ = result
	}
}

func BenchmarkLarge_RowRange(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			RowRange(1000, 2000).
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Tail(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Tail(1000).
			Collect()
		_ = result
	}
}

// ===== Real-world Scenario Benchmarks =====

func BenchmarkLarge_Scenario_TopPerformers(b *testing.B) {
	// Find top 100 performers with high scores
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Where("Score").GreaterThan(80.0).
			OrderBy("Score").Desc().
			Limit(100).
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Scenario_CategoryAnalysis(b *testing.B) {
	// Analyze categories by average score
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			GroupBy("Category").
			Avg("Score")
		_ = result
	}
}

func BenchmarkLarge_Scenario_OutlierDetection(b *testing.B) {
	// Find outliers (values far from mean)
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Value1").LessThan(70.0).Or("Value1").GreaterThan(130.0).
			Where("Active").IsTrue().
			OrderBy("Value1").Desc().
			Collect()
		_ = result
	}
}

func BenchmarkLarge_Scenario_DetailedReport(b *testing.B) {
	// Generate a detailed filtered report
	data := generateLargeDataFrame(100000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Category").OneOf("A", "B", "C").
			Where("Value1").Between(90.0, 110.0).
			Where("Score").GreaterThan(70.0).
			Where("Active").IsTrue().
			OrderBy("Category").Asc().
			ThenBy("Score").Desc().
			Select("ID", "Name", "Category", "Value1", "Score").
			Collect()
		_ = result
	}
}

// ===== Memory profiling helper =====

func BenchmarkLarge_MemoryProfile_ComplexPipeline(b *testing.B) {
	data := generateLargeDataFrame(100000)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := From(data).
			Where("Active").IsTrue().
			Where("Value1").GreaterThan(90.0).
			Where("Score").Between(60.0, 90.0).
			OrderBy("Category").Asc().
			ThenBy("Score").Desc().
			Limit(1000).
			Collect()
		_ = result
	}
}
