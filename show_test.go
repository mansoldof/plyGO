package plygo

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

type ShowTestPerson struct {
	Name   string
	Age    int
	City   string
	Salary float64
	Active bool
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestShow_Basic(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000.50, true},
		{"Bob", 25, "LA", 60000.00, false},
	}

	output := captureOutput(func() {
		From(people).Show()
	})

	if !strings.Contains(output, "Alice") {
		t.Error("Output should contain Alice")
	}
	if !strings.Contains(output, "Bob") {
		t.Error("Output should contain Bob")
	}
	if !strings.Contains(output, "2 rows") {
		t.Error("Output should show row count")
	}
}

func TestShow_WithTitle(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
	}

	output := captureOutput(func() {
		From(people).Show(WithTitle("Test Data"))
	})

	if !strings.Contains(output, "Test Data") {
		t.Error("Output should contain title")
	}
}

func TestShow_StyleRounded(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
	}

	output := captureOutput(func() {
		From(people).Show(WithStyle("rounded"))
	})

	if !strings.Contains(output, "╭") || !strings.Contains(output, "╮") {
		t.Error("Output should use rounded style characters")
	}
}

func TestShow_StyleDouble(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
	}

	output := captureOutput(func() {
		From(people).Show(WithStyle("double"))
	})

	if !strings.Contains(output, "╔") || !strings.Contains(output, "╗") {
		t.Error("Output should use double style characters")
	}
}

func TestShow_StyleMinimal(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
	}

	output := captureOutput(func() {
		From(people).Show(WithStyle("minimal"))
	})

	if !strings.Contains(output, "Alice") {
		t.Error("Output should contain data")
	}
	if strings.Contains(output, "+") || strings.Contains(output, "|") {
		t.Error("Minimal style should not have box borders")
	}
}

func TestShow_StyleMarkdown(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
	}

	output := captureOutput(func() {
		From(people).Show(WithStyle("markdown"))
	})

	if !strings.Contains(output, "|") {
		t.Error("Markdown style should use pipe separators")
	}
}

func TestShow_WithRowNumbers(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
		{"Bob", 25, "LA", 60000, false},
	}

	output := captureOutput(func() {
		From(people).Show(WithRowNumbers(true))
	})

	if !strings.Contains(output, "#") {
		t.Error("Output should have row number header")
	}
}

func TestShow_WithOriginalIndices(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
		{"Bob", 25, "LA", 60000, false},
		{"Charlie", 35, "NYC", 90000, true},
	}

	output := captureOutput(func() {
		From(people).
			Where("Age").GreaterThan(28).
			Show(WithOriginalIndices(true))
	})

	if !strings.Contains(output, "#") {
		t.Error("Output should have index column")
	}
}

func TestShow_FloatPrecision(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000.556, true},
	}

	output := captureOutput(func() {
		From(people).Show(WithFloatPrecision(0))
	})

	if strings.Contains(output, ".") && strings.Contains(output, "75000.") {
		t.Error("Should not show decimals with precision 0")
	}
}

func TestShow_BoolStyleSymbols(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
		{"Bob", 25, "LA", 60000, false},
	}

	output := captureOutput(func() {
		From(people).Show(WithBoolStyle("symbols"))
	})

	if !strings.Contains(output, "✓") || !strings.Contains(output, "✗") {
		t.Error("Should use symbols for booleans")
	}
}

func TestShow_LargeDataset(t *testing.T) {
	people := make([]ShowTestPerson, 100)
	for i := 0; i < 100; i++ {
		people[i] = ShowTestPerson{
			Name:   fmt.Sprintf("Person%d", i+1),
			Age:    20 + i,
			City:   "NYC",
			Salary: 50000.0 + float64(i*1000),
			Active: i%2 == 0,
		}
	}

	output := captureOutput(func() {
		From(people).Show()
	})

	if !strings.Contains(output, "...") {
		t.Error("Large dataset should show truncation indicator")
	}
	if !strings.Contains(output, "100 rows") {
		t.Error("Should show total row count")
	}
	if !strings.Contains(output, "showing") {
		t.Error("Should indicate truncation in footer")
	}
}

func TestShow_EmptyDataset(t *testing.T) {
	people := []ShowTestPerson{}

	output := captureOutput(func() {
		From(people).Show()
	})

	if !strings.Contains(output, "Empty") {
		t.Error("Should show empty dataset message")
	}
}

func TestShow_AfterFilter(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
		{"Bob", 25, "LA", 60000, false},
		{"Charlie", 35, "NYC", 90000, true},
	}

	output := captureOutput(func() {
		From(people).
			Where("Age").GreaterThan(28).
			Show()
	})

	if !strings.Contains(output, "Alice") {
		t.Error("Should contain filtered results")
	}
	if strings.Contains(output, "Bob") {
		t.Error("Should not contain filtered out items")
	}
}

func TestShow_AfterSelect(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
		{"Bob", 25, "LA", 60000, false},
	}

	output := captureOutput(func() {
		From(people).
			Select("Name", "City").
			Show()
	})

	if !strings.Contains(output, "Alice") || !strings.Contains(output, "NYC") {
		t.Error("Should show selected fields")
	}
	if strings.Contains(output, "75000") {
		t.Error("Should not show non-selected fields")
	}
}

func TestShow_ComplexPipeline(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000, true},
		{"Bob", 25, "LA", 60000, true},
		{"Charlie", 35, "NYC", 90000, true},
		{"Diana", 28, "Chicago", 70000, true},
	}

	output := captureOutput(func() {
		From(people).
			Where("Age").GreaterThan(26).
			OrderBy("Salary").Desc().
			Show(WithTitle("High Earners"), WithStyle("rounded"))
	})

	if !strings.Contains(output, "High Earners") {
		t.Error("Should show title")
	}
	if !strings.Contains(output, "Charlie") {
		t.Error("Should contain filtered and sorted results")
	}
}

func TestShow_LongStrings(t *testing.T) {
	people := []ShowTestPerson{
		{"This is a very long name that should be truncated", 30, "NYC", 75000, true},
	}

	output := captureOutput(func() {
		From(people).Show(WithMaxColWidth(20))
	})

	if !strings.Contains(output, "...") {
		t.Error("Long strings should be truncated with ...")
	}
}

func TestShowPositions_Basic(t *testing.T) {
	pos := PositionIndex{
		Rows: []int{1, 3, 5},
	}

	output := captureOutput(func() {
		ShowPositions(pos)
	})

	if !strings.Contains(output, "Original Positions") {
		t.Error("Should show positions label")
	}
	if !strings.Contains(output, "1") || !strings.Contains(output, "3") || !strings.Contains(output, "5") {
		t.Error("Should show all positions")
	}
}

func TestShowPositions_Matrix(t *testing.T) {
	pos := PositionIndex{
		Rows: []int{1, 2},
		Cols: []int{1, 3},
	}

	output := captureOutput(func() {
		ShowPositions(pos)
	})

	if !strings.Contains(output, "Column Positions") {
		t.Error("Should show column positions for matrix")
	}
}

func TestShowPositions_Empty(t *testing.T) {
	pos := PositionIndex{
		Rows: []int{},
	}

	output := captureOutput(func() {
		ShowPositions(pos)
	})

	if !strings.Contains(output, "No positions") {
		t.Error("Should show no positions message")
	}
}

func TestShow_MaxRows(t *testing.T) {
	people := make([]ShowTestPerson, 50)
	for i := 0; i < 50; i++ {
		people[i] = ShowTestPerson{
			Name: fmt.Sprintf("Person%d", i+1),
			Age:  20 + i,
		}
	}

	output := captureOutput(func() {
		From(people).Show(WithMaxRows(10))
	})

	lines := strings.Split(output, "\n")
	dataLines := 0
	for _, line := range lines {
		if strings.Contains(line, "Person") {
			dataLines++
		}
	}

	if dataLines > 11 {
		t.Errorf("Should show at most 10 rows + separator, got %d lines", dataLines)
	}
}

func TestShow_CombinedOptions(t *testing.T) {
	people := []ShowTestPerson{
		{"Alice", 30, "NYC", 75000.123, true},
		{"Bob", 25, "LA", 60000.456, false},
	}

	output := captureOutput(func() {
		From(people).Show(
			WithTitle("Test"),
			WithStyle("rounded"),
			WithRowNumbers(true),
			WithFloatPrecision(1),
			WithBoolStyle("symbols"),
		)
	})

	if !strings.Contains(output, "Test") {
		t.Error("Should have title")
	}
	if !strings.Contains(output, "╭") {
		t.Error("Should use rounded style")
	}
	if !strings.Contains(output, "#") {
		t.Error("Should have row numbers")
	}
}

func TestShow_SimpleTypes(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	output := captureOutput(func() {
		From(numbers).Show()
	})

	if !strings.Contains(output, "1") || !strings.Contains(output, "5") {
		t.Error("Should show simple type values")
	}
}
