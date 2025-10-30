# plygo

A lightweight, type-safe data manipulation library for Go inspired by dplyr and R.

## Features

-   Fluent, chainable API
-   Type-safe operations with generics
-   Position-based data access (R-style with 1-based indexing)
-   **Elegant table visualization with Show()**
-   Clean, readable syntax
-   Zero dependencies
-   High performance

## Installation

``` bash
go get github.com/mansoldof/plyGO
```

## Quick Start

``` go
package main

import "github.com/mansoldof/plyGO"

type Person struct {
    Name   string
    Age    int
    City   string
    Salary float64
}

func main() {

    people := []Person{
        {"Frank", 43, "FL", 88000},
        {"Alice", 30, "NYC", 75000},
        {"Bob", 25, "LA", 60000},
        {"Charlie", 35, "NYC", 90000},
    }

    // Display with elegant formatting
    plygo.From(people).Show()

    // Filter and show
    plygo.From(people).
        Where("Age").GreaterThan(30).
        Show(plygo.WithTitle("Age > 30"))

}
```

<div align="center"> 
<img src="https://github.com/user-attachments/assets/294a8e81-b999-4101-bf7f-efd9da26c54e" alt="plyGO example" width="293" style="border-radius:12px; box-shadow:0 4px 18px rgba(0,0,0,0.2); margin-top:1rem;"/>
<p><em>Console output</em></p>
</div>

## Table Visualization

### Basic Usage

``` go
    // Display with elegant formatting
    plygo.From(people).Show(plygo.WithStyle("rounded"))

    // Filter and show
    plygo.From(people).
        Where("Salary").LessThan(80000).
        Show(plygo.WithTitle("Salary < $80,000"),
            plygo.WithStyle("rounded"))
```

<div align="center"> 
<img src="https://github.com/user-attachments/assets/baf3e38e-f5bb-417c-b20f-4a354b584b0e" alt="plyGO example" width="284" style="border-radius:12px; box-shadow:0 4px 18px rgba(0,0,0,0.2); margin-top:1rem;"/>
<p><em>Console output</em></p>
</div>

## Integration with goTableView (Compact Version)

plyGO integrates seamlessly with [goTableView](https://github.com/mansoldof/goTableView) for GUI visualization. Transform your data with plyGO's fluent API, then display it in a native Windows table viewer.

### Example

``` go
package main

import (
    "github.com/mansoldof/plyGO"
    "github.com/mansoldof/goTableView"
)

type Person struct {
    Name   string
    Age    int
    City   string
    Salary float64
}

func main() {
    people := []Person{
        {"Frank", 43, "FL", 88000},
        {"Alice", 30, "NYC", 75000},
        {"Bob", 25, "LA", 60000},
        {"Charlie", 35, "NYC", 90000},
    }

	// Display initial table
	gotableview.FromStructs("Initial table", people).Show(false)

    // Filter, sort, and display in GUI
    result := plygo.From(people).
        Where("Age").GreaterThan(30).
        OrderBy("Salary").Desc().
        Collect()

    gotableview.FromStructs("High Earners", result).Show()
}
```

<div align="center"> <img src="https://github.com/user-attachments/assets/8a30c562-474b-4a2a-95a2-aa177ed03102" alt="plyGO with goTableView basic example" width="600" 
                         style="border-radius:12px; box-shadow:0 4px 18px rgba(0,0,0,0.2); margin-top:1rem;">
<p><em>Rendered GUI table — built directly from Go code.</em></p> 
</div>

Install both:

``` bash
go get github.com/mansoldof/plyGO
go get github.com/mansoldof/goTableView
```

**Note**: goTableView is Windows-only. For terminal display, use plyGO's `Show()` method.

## Table Styles

``` go
// Rounded borders
plygo.From(people).Show(plygo.WithStyle("rounded"))

// Double borders
plygo.From(people).Show(plygo.WithStyle("double"))

// Minimal (no borders)
plygo.From(people).Show(plygo.WithStyle("minimal"))

// Markdown format
plygo.From(people).Show(plygo.WithStyle("markdown"))
```

## Customization Options

``` go
// With title
plygo.From(people).Show(plygo.WithTitle("Employee Data"))

// With row numbers
plygo.From(people).Show(plygo.WithRowNumbers(true))

// Show original indices after filtering
plygo.From(people).
    Where("Age").GreaterThan(30).
    Show(plygo.WithOriginalIndices(true))

// Custom float precision
plygo.From(people).Show(plygo.WithFloatPrecision(0))

// Boolean as symbols (✓/✗)
plygo.From(people).Show(plygo.WithBoolStyle("symbols"))

// Limit displayed rows
plygo.From(people).Show(plygo.WithMaxRows(10))

// Combined options
plygo.From(people).Show(
    plygo.WithTitle("Report"),
    plygo.WithStyle("rounded"),
    plygo.WithRowNumbers(true),
    plygo.WithFloatPrecision(2),
)
```

### Smart Features

-   **Auto-truncation**: Large datasets automatically show first/last rows
-   **Type detection**: Numbers right-aligned, strings left-aligned
-   **Long text handling**: Truncates with "..." when needed
-   **Empty handling**: Graceful display for empty datasets
-   **Works everywhere**: With filters, sorts, selections, and all plyGo operations

### Show Throughout Pipeline

``` go
// Display at any point in the pipeline
plygo.From(people).
    Show(plygo.WithTitle("Original")).
    Where("Age").GreaterThan(30).
    Show(plygo.WithTitle("Filtered")).
    OrderBy("Salary").Desc().
    Show(plygo.WithTitle("Sorted"))
```

## Position-Based Operations

### Row Access

``` go
// Get specific rows (1-based indexing)
plygo.From(people).AtRow(1, 3, 5).Collect()

// Negative indices (Python-style)
plygo.From(people).AtRow(-1).Collect()        // Last element
plygo.From(people).AtRow(1, -1).Collect()     // First and last

// Range of rows
plygo.From(people).RowRange(2, 5).Collect()   // Rows 2-4
plygo.From(people).RowRange(3, -1).Collect()  // From row 3 to end

// Tail - last n rows
plygo.From(people).Tail(3).Collect()

// Slice with step
plygo.From(people).Slice(1, -1, 2).Collect()  // Every other row
```

### Column Access

``` go
// Get specific columns by position
plygo.From(people).AtCol(1, 3).Collect()      // Name and City

// Range of columns
plygo.From(people).ColRange(2, 4).Collect()   // Age and City

// Field introspection
fields := plygo.From(people).FieldNames()     // ["Name", "Age", "City", "Salary"]
count := plygo.From(people).FieldCount()       // 4
```

### Position Tracking

``` go
// Get original indices after filtering
positions := plygo.From(people).
    Where("Age").GreaterThan(30).
    Positions()

// Display positions elegantly
plygo.ShowPositions(positions)
```

## Filtering

``` go
// AND (implicit)
Where("Age").GreaterThan(30).
Where("Active").IsTrue()

// OR
Where("City").Equals("NYC").Or("City").Equals("LA")

// OneOf (cleaner OR)
Where("City").OneOf("NYC", "LA", "Chicago")

// Between
Where("Age").Between(25, 35)

// String operations
Where("Name").Contains("Smith")
Where("Name").StartsWith("A")
Where("Name").EndsWith("son")
```

## Grouping & Aggregation

``` go
GroupBy("City").Count()
GroupBy("City").Sum("Salary")
GroupBy("City").Avg("Age")
GroupBy("City").Min("Salary")
GroupBy("City").Max("Salary")
```

## Sorting

``` go
OrderBy("Salary").Desc()
OrderBy("Age").Asc().ThenBy("Salary").Desc()
```

## Transformation

``` go
Transform(func(p Person) Person {
    p.Salary *= 1.1
    return p
})
```

## Examples

-   `examples/simple/simple.go` - Basic filtering and operations
-   `examples/positions/positions.go` - Position-based operations
-   **`examples/show/show_examples.go` - Table visualization examples**

## Show() Options Reference

| Option | Description | Example |
|--------------------|-------------------------------|----------------------|
| `WithStyle(style)` | Table style: "simple", "rounded", "double", "minimal", "markdown" | `WithStyle("rounded")` |
| `WithTitle(title)` | Add title above table | `WithTitle("Report")` |
| `WithRowNumbers(bool)` | Show row numbers | `WithRowNumbers(true)` |
| `WithOriginalIndices(bool)` | Show original positions after filtering | `WithOriginalIndices(true)` |
| `WithFloatPrecision(n)` | Decimal places for floats | `WithFloatPrecision(0)` |
| `WithBoolStyle(style)` | Boolean display: "text" or "symbols" | `WithBoolStyle("symbols")` |
| `WithMaxRows(n)` | Maximum rows to display | `WithMaxRows(20)` |
| `WithMaxWidth(n)` | Maximum total width | `WithMaxWidth(120)` |
| `WithMaxColWidth(n)` | Maximum column width | `WithMaxColWidth(30)` |

## License

MIT
