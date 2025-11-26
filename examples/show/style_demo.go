package main

import plygo "github.com/mansoldof/plyGO"

type Sample struct {
	ID     int
	Name   string
	Value  float64
	Active bool
}

func main() {
	data := []Sample{
		{1, "Alice", 100.50, true},
		{2, "Bob", 200.75, false},
		{3, "Charlie", 150.25, true},
	}

	println("╔════════════════════════════════════════════╗")
	println("║  plyGo Show() - Style Comparison Demo     ║")
	println("╚════════════════════════════════════════════╝")
	println()

	println("1. SIMPLE STYLE (default)")
	println("   Classic box-drawing characters")
	println()
	plygo.From(data).Show()

	println("\n2. ROUNDED STYLE")
	println("   Modern rounded corners")
	println()
	plygo.From(data).Show(plygo.WithStyle("rounded"))

	println("\n3. DOUBLE STYLE")
	println("   Bold double-line borders")
	println()
	plygo.From(data).Show(plygo.WithStyle("double"))

	println("\n4. MINIMAL STYLE")
	println("   Clean, no borders")
	println()
	plygo.From(data).Show(plygo.WithStyle("minimal"))

	println("\n5. MARKDOWN STYLE")
	println("   Perfect for documentation")
	println()
	plygo.From(data).Show(plygo.WithStyle("markdown"))

	println("\n6. WITH ALL OPTIONS")
	println("   Title + Row numbers + Symbols + Custom precision")
	println()
	plygo.From(data).Show(
		plygo.WithTitle("Sample Data Report"),
		plygo.WithStyle("double"),
		plygo.WithRowNumbers(true),
		plygo.WithBoolStyle("symbols"),
		plygo.WithFloatPrecision(1),
	)

	println("\n╔════════════════════════════════════════════╗")
	println("║  All styles work with filters & sorts!    ║")
	println("╚════════════════════════════════════════════╝")
}
