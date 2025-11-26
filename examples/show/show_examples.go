package main

import (
	"fmt"

	plygo "github.com/mansoldof/plyGO"
)

type Employee struct {
	Name       string
	Age        int
	Department string
	Salary     float64
	Active     bool
}

func main() {
	employees := []Employee{
		{"Alice Johnson", 30, "Engineering", 95000.50, true},
		{"Bob Smith", 45, "Sales", 75000.00, true},
		{"Charlie Brown", 28, "Engineering", 85000.75, true},
		{"Diana Prince", 35, "Marketing", 80000.00, true},
		{"Eve Davis", 42, "Engineering", 110000.25, true},
		{"Frank Miller", 38, "Sales", 72000.50, false},
		{"Grace Lee", 29, "Marketing", 78000.00, true},
		{"Henry Wilson", 50, "Engineering", 125000.00, true},
		{"Ivy Chen", 26, "Sales", 65000.00, true},
		{"Jack Turner", 33, "Marketing", 82000.50, true},
	}

	fmt.Println("=== plyGo Show() Examples ===\n")

	// Example 1: Basic Show
	fmt.Println("1. Basic Show (default style):")
	plygo.From(employees).Show()

	// Example 2: With Title
	fmt.Println("\n2. With Custom Title:")
	plygo.From(employees).
		Show(plygo.WithTitle("Company Employees"))

	// Example 3: Rounded Style
	fmt.Println("\n3. Rounded Box Style:")
	plygo.From(employees).
		AtRow(1, 2, 3).
		Show(plygo.WithStyle("rounded"))

	// Example 4: Double Border Style
	fmt.Println("\n4. Double Border Style:")
	plygo.From(employees).
		AtRow(1, 2, 3).
		Show(plygo.WithStyle("double"))

	// Example 5: Minimal Style
	fmt.Println("\n5. Minimal Style (no borders):")
	plygo.From(employees).
		AtRow(1, 2, 3).
		Show(plygo.WithStyle("minimal"))

	// Example 6: Markdown Style
	fmt.Println("\n6. Markdown Table:")
	plygo.From(employees).
		AtRow(1, 2, 3).
		Show(plygo.WithStyle("markdown"))

	// Example 7: With Row Numbers
	fmt.Println("\n7. With Row Numbers:")
	plygo.From(employees).
		AtRow(1, 2, 3, 4).
		Show(plygo.WithRowNumbers(true))

	// Example 8: Filtered Data with Original Indices
	fmt.Println("\n8. Filtered with Original Indices:")
	plygo.From(employees).
		Where("Department").Equals("Engineering").
		Show(
			plygo.WithTitle("Engineering Team"),
			plygo.WithOriginalIndices(true),
			plygo.WithStyle("rounded"),
		)

	// Example 9: Float Precision
	fmt.Println("\n9. Custom Float Precision (0 decimals):")
	plygo.From(employees).
		AtRow(1, 2, 3).
		Show(plygo.WithFloatPrecision(0))

	// Example 10: Boolean Symbols
	fmt.Println("\n10. Boolean as Symbols:")
	plygo.From(employees).
		AtRow(1, 5, 6).
		Show(plygo.WithBoolStyle("symbols"))

	// Example 11: Select Specific Columns
	fmt.Println("\n11. Selected Columns Only:")
	plygo.From(employees).
		Select("Name", "Department", "Salary").
		Show(plygo.WithStyle("rounded"))

	// Example 12: Complex Pipeline
	fmt.Println("\n12. Complex Filter + Sort:")
	plygo.From(employees).
		Where("Salary").GreaterThan(80000).
		Where("Active").IsTrue().
		OrderBy("Salary").Desc().
		Show(
			plygo.WithTitle("High Earners (Active)"),
			plygo.WithStyle("double"),
			plygo.WithOriginalIndices(true),
		)

	// Example 13: Large Dataset Truncation
	fmt.Println("\n13. Large Dataset (auto-truncation):")
	largeData := make([]Employee, 100)
	for i := 0; i < 100; i++ {
		largeData[i] = Employee{
			Name:       fmt.Sprintf("Employee %d", i+1),
			Age:        25 + (i % 30),
			Department: []string{"Eng", "Sales", "Marketing"}[i%3],
			Salary:     50000 + float64(i*1000),
			Active:     i%3 != 0,
		}
	}
	plygo.From(largeData).Show()

	// Example 14: Show Positions
	fmt.Println("\n14. Show Position Indices:")
	positions := plygo.From(employees).
		Where("Department").Equals("Engineering").
		Where("Salary").GreaterThan(90000).
		Positions()
	plygo.ShowPositions(positions, plygo.WithStyle("rounded"))

	// Example 15: Combined Options
	fmt.Println("\n15. All Options Combined:")
	plygo.From(employees).
		Where("Age").Between(30, 45).
		Show(
			plygo.WithTitle("Employees Age 30-45"),
			plygo.WithStyle("rounded"),
			plygo.WithRowNumbers(true),
			plygo.WithFloatPrecision(2),
			plygo.WithBoolStyle("symbols"),
			plygo.WithMaxRows(5),
		)

	// Example 16: After OrderBy
	fmt.Println("\n16. Sorted by Salary (Top 5):")
	plygo.From(employees).
		OrderBy("Salary").Desc().
		AtRow(1, 2, 3, 4, 5).
		Show(
			plygo.WithTitle("Top 5 Salaries"),
			plygo.WithStyle("double"),
		)

	// Example 17: Column Range
	fmt.Println("\n17. Show with Column Selection:")
	plygo.From(employees).
		AtCol(1, 3, 4).
		Show(plygo.WithTitle("Name, Department & Salary"))

	// Example 18: Compact View
	fmt.Println("\n18. After Transform:")
	plygo.From(employees).
		Transform(func(e Employee) Employee {
			e.Salary *= 1.10 // 10% raise
			return e
		}).
		Where("Salary").GreaterThan(100000).
		Show(
			plygo.WithTitle("After 10% Raise (>100k)"),
			plygo.WithFloatPrecision(0),
		)

	fmt.Println("\n=== End of Examples ===")
}
