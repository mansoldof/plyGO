package plygo

import (
	"reflect"
	"testing"
)

type TestPerson struct {
	Name   string
	Age    int
	City   string
	Salary float64
}

func TestAtRow_SingleIndex(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	result := From(people).AtRow(1).Collect()

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
	if result[0].Name != "Alice" {
		t.Errorf("Expected Alice (index 1), got %s", result[0].Name)
	}
}

func TestAtRow_MultipleIndices(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
		{"Diana", 28, "Miami", 70000},
	}

	result := From(people).AtRow(1, 3).Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}
	if result[0].Name != "Alice" || result[1].Name != "Charlie" {
		t.Errorf("Expected Alice and Charlie, got %s and %s", result[0].Name, result[1].Name)
	}
}

func TestAtRow_NegativeIndex(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	// -1 should be last element
	result := From(people).AtRow(-1).Collect()

	if len(result) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result))
	}
	if result[0].Name != "Charlie" {
		t.Errorf("Expected Charlie (last element), got %s", result[0].Name)
	}
}

func TestAtRow_FirstAndLast(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	result := From(people).AtRow(1, -1).Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}
	if result[0].Name != "Alice" || result[1].Name != "Charlie" {
		t.Errorf("Expected Alice and Charlie, got %s and %s", result[0].Name, result[1].Name)
	}
}

func TestRowRange(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
		{"Diana", 28, "Miami", 70000},
		{"Eve", 32, "Boston", 85000},
	}

	result := From(people).RowRange(2, 4).Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 results (indices 2-3), got %d", len(result))
	}
	if result[0].Name != "Bob" || result[1].Name != "Charlie" {
		t.Errorf("Expected Bob and Charlie, got %s and %s", result[0].Name, result[1].Name)
	}
}

func TestRowRange_ToEnd(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	result := From(people).RowRange(2, -1).Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 results (from index 2 to end), got %d", len(result))
	}
	if result[0].Name != "Bob" {
		t.Errorf("Expected Bob as first, got %s", result[0].Name)
	}
}

func TestTail(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
		{"Diana", 28, "Miami", 70000},
	}

	result := From(people).Tail(2).Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}
	if result[0].Name != "Charlie" || result[1].Name != "Diana" {
		t.Errorf("Expected Charlie and Diana, got %s and %s", result[0].Name, result[1].Name)
	}
}

func TestSample(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
		{"Diana", 28, "Miami", 70000},
		{"Eve", 32, "Boston", 85000},
	}

	result := From(people).Sample(3).Collect()

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}
}

func TestSlice_Positive(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
		{"Diana", 28, "Miami", 70000},
		{"Eve", 32, "Boston", 85000},
	}

	// Every other element starting from index 1
	result := From(people).Slice(1, 5, 2).Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 results (step=2), got %d", len(result))
	}
	if result[0].Name != "Alice" || result[1].Name != "Charlie" {
		t.Errorf("Expected Alice and Charlie, got %s and %s", result[0].Name, result[1].Name)
	}
}

func TestPositions_Simple(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	positions := From(people).Positions()

	expected := []int{1, 2, 3}
	if !reflect.DeepEqual(positions.Rows, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions.Rows)
	}
}

func TestPositions_AfterFilter(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
		{"Diana", 28, "Miami", 70000},
	}

	positions := From(people).
		Where("Age").GreaterThan(28).
		Positions()

	expected := []int{1, 3} // Alice (index 1) and Charlie (index 3)
	if !reflect.DeepEqual(positions.Rows, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions.Rows)
	}
}

func TestWhich(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	indices := From(people).
		Where("City").Equals("NYC").
		Which()

	expected := []int{1}
	if !reflect.DeepEqual(indices, expected) {
		t.Errorf("Expected indices %v, got %v", expected, indices)
	}
}

func TestFieldNames(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
	}

	fields := From(people).FieldNames()

	expected := []string{"Name", "Age", "City", "Salary"}
	if !reflect.DeepEqual(fields, expected) {
		t.Errorf("Expected fields %v, got %v", expected, fields)
	}
}

func TestFieldCount(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
	}

	count := From(people).FieldCount()

	if count != 4 {
		t.Errorf("Expected 4 fields, got %d", count)
	}
}

func TestAtCol_SingleIndex(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
	}

	result := From(people).AtCol(1).Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(result))
	}

	// Should only have Name field (index 1)
	if len(result[0]) != 1 {
		t.Errorf("Expected 1 column, got %d", len(result[0]))
	}
	if result[0]["Name"] != "Alice" {
		t.Errorf("Expected Alice, got %v", result[0]["Name"])
	}
}

func TestAtCol_MultipleIndices(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
	}

	result := From(people).AtCol(1, 3).Collect()

	// Should have Name (index 1) and City (index 3) fields
	if len(result[0]) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(result[0]))
	}
	if result[0]["Name"] != "Alice" || result[0]["City"] != "NYC" {
		t.Errorf("Expected Name and City fields")
	}
}

func TestColRange(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
	}

	result := From(people).ColRange(2, 4).Collect()

	// Should have Age (index 2) and City (index 3)
	if len(result[0]) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(result[0]))
	}
	if _, hasAge := result[0]["Age"]; !hasAge {
		t.Error("Expected Age field")
	}
	if _, hasCity := result[0]["City"]; !hasCity {
		t.Error("Expected City field")
	}
}

func TestAtRowAndAtCol_Combined(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	result := From(people).
		AtRow(1, 3).
		AtCol(1, 2).
		Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(result))
	}
	if len(result[0]) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(result[0]))
	}
	if result[0]["Name"] != "Alice" {
		t.Errorf("Expected Alice, got %v", result[0]["Name"])
	}
}

func TestPositions_WithColumns(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
	}

	positions := From(people).
		AtRow(1).
		AtCol(1, 3).
		Positions()

	expectedRows := []int{1}
	expectedCols := []int{1, 3}

	if !reflect.DeepEqual(positions.Rows, expectedRows) {
		t.Errorf("Expected row positions %v, got %v", expectedRows, positions.Rows)
	}
	if !reflect.DeepEqual(positions.Cols, expectedCols) {
		t.Errorf("Expected col positions %v, got %v", expectedCols, positions.Cols)
	}
	if !positions.IsMatrix() {
		t.Error("Expected IsMatrix() to be true")
	}
}

func TestSelection_AtRow(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	result := From(people).
		Select("Name", "City").
		AtRow(1, 3).
		Collect()

	if len(result) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(result))
	}
	if result[0]["Name"] != "Alice" {
		t.Errorf("Expected Alice, got %v", result[0]["Name"])
	}
}

func TestComplexPipeline_WithPositions(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "NYC", 90000},
		{"Diana", 28, "Chicago", 70000},
		{"Eve", 32, "LA", 85000},
	}

	positions := From(people).
		Where("Age").GreaterThan(27).
		Where("City").OneOf("NYC", "LA").
		Positions()

	// Should match: Alice (1), Charlie (3), Eve (5)
	expected := []int{1, 3, 5}
	if !reflect.DeepEqual(positions.Rows, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions.Rows)
	}
}

func TestPositionTracking_ThroughTransform(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
		{"Charlie", 35, "Chicago", 90000},
	}

	positions := From(people).
		Transform(func(p TestPerson) TestPerson {
			p.Salary *= 1.1
			return p
		}).
		Where("Salary").GreaterThan(80000).
		Positions()

	// After 10% raise: Alice=82500, Bob=66000, Charlie=99000
	// Charlie should be at original index 3
	expected := []int{1, 3} // Alice and Charlie
	if !reflect.DeepEqual(positions.Rows, expected) {
		t.Errorf("Expected positions %v, got %v", expected, positions.Rows)
	}
}

func TestAtRow_OutOfBounds(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
	}

	result := From(people).AtRow(10).Collect()

	if len(result) != 0 {
		t.Errorf("Expected 0 results for out of bounds index, got %d", len(result))
	}
}

func TestAtRow_ZeroIndex(t *testing.T) {
	people := []TestPerson{
		{"Alice", 30, "NYC", 75000},
		{"Bob", 25, "LA", 60000},
	}

	// Index 0 should be invalid (1-based indexing)
	result := From(people).AtRow(0).Collect()

	if len(result) != 0 {
		t.Errorf("Expected 0 results for index 0, got %d", len(result))
	}
}

func TestPositionIndex_Methods(t *testing.T) {
	pi := PositionIndex{
		Rows: []int{1, 2, 3},
		Cols: []int{1, 2},
	}

	if pi.RowCount() != 3 {
		t.Errorf("Expected RowCount 3, got %d", pi.RowCount())
	}
	if pi.ColCount() != 2 {
		t.Errorf("Expected ColCount 2, got %d", pi.ColCount())
	}
	if !pi.IsMatrix() {
		t.Error("Expected IsMatrix to be true")
	}
}
