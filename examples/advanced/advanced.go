package main

import (
	"fmt"

	plygo "github.com/mansoldof/plyGO"
)

type Employee struct {
	ID         int
	Name       string
	Department string
	Salary     float64
	YearsExp   int
	Active     bool
}

func main() {
	employees := []Employee{
		{1, "Alice Johnson", "Engineering", 95000, 5, true},
		{2, "Bob Smith", "Marketing", 70000, 3, true},
		{3, "Charlie Brown", "Engineering", 120000, 8, true},
		{4, "Diana Prince", "Sales", 85000, 4, false},
		{5, "Eve Adams", "Engineering", 78000, 2, true},
		{6, "Frank Miller", "Marketing", 92000, 6, true},
		{7, "Grace Lee", "Sales", 110000, 7, true},
		{8, "Henry Ford", "Engineering", 88000, 4, true},
	}

	fmt.Println("=== Advanced Example 1: Multi-Level Filtering ===")
	result1 := plygo.From(employees).
		Where("Active").IsTrue().
		Where("Department").Equals("Engineering").Or("Department").Equals("Sales").
		Where("Salary").GreaterThan(80000).
		OrderBy("Salary").Desc().
		Collect()

	for _, emp := range result1 {
		fmt.Printf("  %s (%s): $%.0f\n", emp.Name, emp.Department, emp.Salary)
	}

	fmt.Println("\n=== Advanced Example 2: Department Statistics ===")
	avgSalaries := plygo.From(employees).
		Where("Active").IsTrue().
		GroupBy("Department").
		Avg("Salary")

	fmt.Println("  Average Salaries by Department:")
	for dept, avg := range avgSalaries {
		fmt.Printf("    %s: $%.2f\n", dept, avg)
	}

	fmt.Println("\n=== Advanced Example 3: Top Earners Per Department ===")
	maxSalaries := plygo.From(employees).
		GroupBy("Department").
		Max("Salary")

	fmt.Println("  Highest Salaries by Department:")
	for dept, max := range maxSalaries {
		fmt.Printf("    %s: $%.0f\n", dept, max)
	}

	fmt.Println("\n=== Advanced Example 4: Experience-Based Filtering ===")
	result4 := plygo.From(employees).
		Where("YearsExp").Between(3, 6).
		Where("Active").IsTrue().
		OrderBy("YearsExp").Asc().
		ThenBy("Salary").Desc().
		Select("Name", "Department", "YearsExp", "Salary").
		Collect()

	fmt.Println("  Employees with 3-6 years experience:")
	for _, emp := range result4 {
		fmt.Printf("    %s (%s): %d years, $%.0f\n",
			emp["Name"], emp["Department"], emp["YearsExp"], emp["Salary"])
	}

	fmt.Println("\n=== Advanced Example 5: Salary Adjustments ===")
	adjusted := plygo.From(employees).
		Where("Active").IsTrue().
		Where("Salary").LessThan(90000).
		Transform(func(e Employee) Employee {
			e.Salary *= 1.15
			return e
		}).
		OrderBy("Salary").Desc().
		Collect()

	fmt.Println("  After 15% raise (for active employees < $90k):")
	for _, emp := range adjusted {
		fmt.Printf("    %s: $%.2f\n", emp.Name, emp.Salary)
	}

	fmt.Println("\n=== Advanced Example 6: Complex Grouping with WhereSome ===")
	result6 := plygo.From(employees).
		WhereSome(
			plygo.WhereEvery(
				plygo.W[Employee]("Department").Equals("Engineering"),
				plygo.W[Employee]("Salary").GreaterThan(90000),
			),
			plygo.WhereEvery(
				plygo.W[Employee]("Department").Equals("Sales"),
				plygo.W[Employee]("YearsExp").GreaterThan(5),
			),
		).
		OrderBy("Salary").Desc().
		Collect()

	fmt.Println("  Engineers earning >$90k OR Sales with >5 years exp:")
	for _, emp := range result6 {
		fmt.Printf("    %s (%s): $%.0f, %d years\n",
			emp.Name, emp.Department, emp.Salary, emp.YearsExp)
	}

	fmt.Println("\n=== Advanced Example 7: Pagination ===")
	page1 := plygo.From(employees).
		Where("Active").IsTrue().
		OrderBy("Salary").Desc().
		Limit(3).
		Collect()

	page2 := plygo.From(employees).
		Where("Active").IsTrue().
		OrderBy("Salary").Desc().
		Skip(3).
		Limit(3).
		Collect()

	fmt.Println("  Page 1 (Top 3):")
	for _, emp := range page1 {
		fmt.Printf("    %s: $%.0f\n", emp.Name, emp.Salary)
	}

	fmt.Println("  Page 2 (Next 3):")
	for _, emp := range page2 {
		fmt.Printf("    %s: $%.0f\n", emp.Name, emp.Salary)
	}

	fmt.Println("\n=== Advanced Example 8: Department Count ===")
	deptCounts := plygo.From(employees).
		Where("Active").IsTrue().
		GroupBy("Department").
		Count()

	fmt.Println("  Active Employees per Department:")
	for dept, count := range deptCounts {
		fmt.Printf("    %s: %d employees\n", dept, count)
	}

	fmt.Println("\n=== Advanced Example 9: String Operations ===")
	result9 := plygo.From(employees).
		Where("Name").StartsWith("A").
		Or("Name").StartsWith("E").
		Collect()

	fmt.Println("  Names starting with 'A' or 'E':")
	for _, emp := range result9 {
		fmt.Printf("    %s (%s)\n", emp.Name, emp.Department)
	}

	fmt.Println("\n=== Advanced Example 10: Unique Departments ===")
	unique := plygo.From(employees).
		Where("Active").IsTrue().
		Distinct("Department").
		Collect()

	fmt.Println("  Unique Active Departments:")
	for _, emp := range unique {
		fmt.Printf("    %s\n", emp.Department)
	}
}
