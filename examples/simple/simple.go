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
	Active bool
}

func main() {
	people := []Person{
		{"Alice", 30, "NYC", 75000, true},
		{"Bob", 25, "LA", 60000, true},
		{"Charlie", 35, "NYC", 90000, false},
		{"Diana", 28, "Chicago", 70000, true},
		{"Eve", 32, "LA", 85000, true},
		{"Frank", 40, "Miami", 95000, true},
	}

	fmt.Println("=== Example 1: Simple AND Filtering ===")
	result1 := plygo.From(people).
		Where("Age").GreaterThan(30).
		Where("Active").IsTrue().
		Collect()
	fmt.Printf("Result: %+v\n\n", result1)

	fmt.Println("=== Example 2: Inline AND ===")
	result2 := plygo.From(people).
		Where("Age").GreaterThan(30).And("Active").IsTrue().
		Collect()
	fmt.Printf("Result: %+v\n\n", result2)

	fmt.Println("=== Example 3: OR Conditions ===")
	result3 := plygo.From(people).
		Where("City").Equals("NYC").Or("City").Equals("LA").
		Collect()
	fmt.Printf("Result: %+v\n\n", result3)

	fmt.Println("=== Example 4: OneOf (cleaner OR) ===")
	result4 := plygo.From(people).
		Where("City").OneOf("NYC", "LA", "Chicago").
		Collect()
	fmt.Printf("Result: %+v\n\n", result4)

	fmt.Println("=== Example 5: Complex AND + OR ===")
	result5 := plygo.From(people).
		Where("Age").GreaterThan(30).
		Where("City").Equals("NYC").Or("City").Equals("LA").
		Where("Active").IsTrue().
		Collect()
	fmt.Printf("Result: %+v\n\n", result5)

	fmt.Println("=== Example 6: WhereSome (OR Groups) ===")
	result6 := plygo.From(people).
		WhereSome(
			plygo.W[Person]("City").Equals("NYC"),
			plygo.W[Person]("City").Equals("LA"),
		).
		Collect()
	fmt.Printf("Result: %+v\n\n", result6)

	fmt.Println("=== Example 7: WhereEvery (AND Groups) ===")
	result7 := plygo.From(people).
		WhereEvery(
			plygo.W[Person]("Age").GreaterThan(30),
			plygo.W[Person]("Active").IsTrue(),
		).
		Collect()
	fmt.Printf("Result: %+v\n\n", result7)

	fmt.Println("=== Example 8: Select Fields ===")
	result8 := plygo.From(people).
		Where("Age").GreaterThan(30).
		Select("Name", "Salary").
		Collect()
	fmt.Printf("Result: %+v\n\n", result8)

	fmt.Println("=== Example 9: OrderBy ===")
	result9 := plygo.From(people).
		OrderBy("Salary").Desc().
		Limit(3).
		Collect()
	fmt.Printf("Result: %+v\n\n", result9)

	fmt.Println("=== Example 10: GroupBy ===")
	result10 := plygo.From(people).
		Where("Active").IsTrue().
		GroupBy("City").
		Count()
	fmt.Printf("Result: %+v\n\n", result10)

	fmt.Println("=== Example 11: GroupBy Sum ===")
	result11 := plygo.From(people).
		GroupBy("City").
		Sum("Salary")
	fmt.Printf("Result: %+v\n\n", result11)

	fmt.Println("=== Example 12: Transform ===")
	result12 := plygo.From(people).
		Where("Active").IsTrue().
		Transform(func(p Person) Person {
			p.Salary *= 1.1
			return p
		}).
		Where("Salary").GreaterThan(80000).
		Collect()
	fmt.Printf("Result: %+v\n\n", result12)

	fmt.Println("=== Example 13: Complex Pipeline ===")
	result13 := plygo.From(people).
		Where("Active").IsTrue().
		Where("Age").Between(25, 35).
		Where("City").OneOf("NYC", "LA").
		OrderBy("Salary").Desc().
		Select("Name", "City", "Salary").
		Collect()
	fmt.Printf("Result: %+v\n\n", result13)
}
