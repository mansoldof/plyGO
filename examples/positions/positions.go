package main

import (
	"fmt"

	plygo "github.com/mansoldof/plyGO"
)

type Person struct {
	Name   string
	Age    int
	City   string
	Salary float64
}

func main() {
	people := []Person{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "NYC", 90000},
		{"Diana", 28, "Chicago", 70000},
		{"Eve", 32, "LA", 85000},
		{"Frank", 40, "Miami", 95000},
	}

	fmt.Println("=== plyGo Position-Based Data Retrieval Examples ===\n")

	// Example 1: AtRow - Get specific rows by position (1-based indexing)
	fmt.Println("1. AtRow - Get first and third person:")
	result1 := plygo.From(people).AtRow(1, 3).Collect()
	for _, p := range result1 {
		fmt.Printf("   %s, %d years\n", p.Name, p.Age)
	}

	// Example 2: AtRow with negative indices (Python-style)
	fmt.Println("\n2. AtRow with negative index - Get last person:")
	result2 := plygo.From(people).AtRow(-1).Collect()
	fmt.Printf("   %s, %d years\n", result2[0].Name, result2[0].Age)

	// Example 3: RowRange - Get range of rows
	fmt.Println("\n3. RowRange - Get persons 2 through 4:")
	result3 := plygo.From(people).RowRange(2, 5).Collect()
	for _, p := range result3 {
		fmt.Printf("   %s, %s\n", p.Name, p.City)
	}

	// Example 4: Tail - Get last n elements
	fmt.Println("\n4. Tail - Get last 3 people:")
	result4 := plygo.From(people).Tail(3).Collect()
	for _, p := range result4 {
		fmt.Printf("   %s\n", p.Name)
	}

	// Example 5: AtCol - Select columns by position
	fmt.Println("\n5. AtCol - Get columns 1 and 3 (Name and City):")
	result5 := plygo.From(people).AtCol(1, 3).Collect()
	for _, row := range result5 {
		fmt.Printf("   %s from %s\n", row["Name"], row["City"])
	}

	// Example 6: AtRow + AtCol - Matrix-style access
	fmt.Println("\n6. AtRow + AtCol - Get specific rows and columns:")
	result6 := plygo.From(people).
		AtRow(1, 3, 5).
		AtCol(1, 4).
		Collect()
	for _, row := range result6 {
		fmt.Printf("   %s: $%.0f\n", row["Name"], row["Salary"])
	}

	// Example 7: Positions - Get original indices after filtering
	fmt.Println("\n7. Positions - Find indices of people over 30:")
	positions := plygo.From(people).
		Where("Age").GreaterThan(30).
		Positions()
	fmt.Printf("   Original indices: %v\n", positions.Rows)

	// Example 8: Which - R-style index retrieval
	fmt.Println("\n8. Which - Find indices where City is LA:")
	indices := plygo.From(people).
		Where("City").Equals("LA").
		Which()
	fmt.Printf("   Indices: %v\n", indices)

	// Example 9: Matrix positions (rows and columns)
	fmt.Println("\n9. Matrix positions (rows and columns):")
	matrixPos := plygo.From(people).
		AtRow(1, 3, 5).
		AtCol(1, 2, 4).
		Positions()
	fmt.Printf("   Row indices: %v\n", matrixPos.Rows)
	fmt.Printf("   Col indices: %v\n", matrixPos.Cols)
	fmt.Printf("   Is matrix: %v\n", matrixPos.IsMatrix())

	// Example 10: Slice - Python-style slicing with step
	fmt.Println("\n10. Slice - Every other person (step=2):")
	result10 := plygo.From(people).Slice(1, -1, 2).Collect()
	for _, p := range result10 {
		fmt.Printf("    %s\n", p.Name)
	}

	// Example 11: Field introspection
	fmt.Println("\n11. Field Names and Count:")
	fieldNames := plygo.From(people).FieldNames()
	fmt.Printf("    Fields: %v\n", fieldNames)
	fmt.Printf("    Count: %d\n", plygo.From(people).FieldCount())

	// Example 12: Using positions to re-access original data
	fmt.Println("\n12. Get positions then use them:")
	highAgeIndices := plygo.From(people).
		Where("Age").GreaterThan(30).
		Which()
	fmt.Printf("    High age indices: %v\n", highAgeIndices)

	// Re-access using those positions
	selected := plygo.From(people).AtRow(highAgeIndices...).Collect()
	fmt.Printf("    Re-accessed %d people:\n", len(selected))
	for _, p := range selected {
		fmt.Printf("      - %s\n", p.Name)
	}

	fmt.Println("\n=== End of Examples ===")
}
