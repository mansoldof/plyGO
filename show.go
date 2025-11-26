package plygo

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode/utf8"
)

type ShowConfig struct {
	maxRows         int
	maxWidth        int
	maxColWidth     int
	style           string
	title           string
	showRowNumbers  bool
	showOriginalIdx bool
	floatPrecision  int
	boolStyle       string
	compact         bool
}

type ShowOption func(*ShowConfig)

func WithMaxRows(n int) ShowOption {
	return func(c *ShowConfig) { c.maxRows = n }
}

func WithMaxWidth(n int) ShowOption {
	return func(c *ShowConfig) { c.maxWidth = n }
}

func WithMaxColWidth(n int) ShowOption {
	return func(c *ShowConfig) { c.maxColWidth = n }
}

func WithStyle(style string) ShowOption {
	return func(c *ShowConfig) { c.style = style }
}

func WithTitle(title string) ShowOption {
	return func(c *ShowConfig) { c.title = title }
}

func WithRowNumbers(show bool) ShowOption {
	return func(c *ShowConfig) { c.showRowNumbers = show }
}

func WithOriginalIndices(show bool) ShowOption {
	return func(c *ShowConfig) { c.showOriginalIdx = show }
}

func WithFloatPrecision(n int) ShowOption {
	return func(c *ShowConfig) { c.floatPrecision = n }
}

func WithBoolStyle(style string) ShowOption {
	return func(c *ShowConfig) { c.boolStyle = style }
}

func WithCompact(compact bool) ShowOption {
	return func(c *ShowConfig) { c.compact = compact }
}

func defaultShowConfig() *ShowConfig {
	return &ShowConfig{
		maxRows:        20,
		maxWidth:       120,
		maxColWidth:    30,
		style:          "simple",
		floatPrecision: 2,
		boolStyle:      "text",
		compact:        false,
	}
}

type tableStyle struct {
	topLeft, topRight, bottomLeft, bottomRight   string
	horizontal, vertical, cross                  string
	topCross, bottomCross, leftCross, rightCross string
	headerSep                                    string
}

func getTableStyle(name string) tableStyle {
	styles := map[string]tableStyle{
		"simple": {
			topLeft: "+", topRight: "+", bottomLeft: "+", bottomRight: "+",
			horizontal: "-", vertical: "|", cross: "+",
			topCross: "+", bottomCross: "+", leftCross: "+", rightCross: "+",
			headerSep: "+",
		},
		"rounded": {
			topLeft: "╭", topRight: "╮", bottomLeft: "╰", bottomRight: "╯",
			horizontal: "─", vertical: "│", cross: "┼",
			topCross: "┬", bottomCross: "┴", leftCross: "├", rightCross: "┤",
			headerSep: "├",
		},
		"double": {
			topLeft: "╔", topRight: "╗", bottomLeft: "╚", bottomRight: "╝",
			horizontal: "═", vertical: "║", cross: "╬",
			topCross: "╦", bottomCross: "╩", leftCross: "╠", rightCross: "╣",
			headerSep: "╠",
		},
		"minimal": {
			topLeft: "", topRight: "", bottomLeft: "", bottomRight: "",
			horizontal: "─", vertical: " ", cross: " ",
			topCross: "", bottomCross: "", leftCross: "", rightCross: "",
			headerSep: "",
		},
		"markdown": {
			topLeft: "|", topRight: "|", bottomLeft: "|", bottomRight: "|",
			horizontal: "-", vertical: "|", cross: "|",
			topCross: "|", bottomCross: "|", leftCross: "|", rightCross: "|",
			headerSep: "|",
		},
	}

	if style, ok := styles[name]; ok {
		return style
	}
	return styles["simple"]
}

func (p *Pipeline[T]) Show(options ...ShowOption) {
	config := defaultShowConfig()
	for _, opt := range options {
		opt(config)
	}

	if len(p.data) == 0 {
		fmt.Println("Empty dataset")
		return
	}

	showTable(p.data, p.originalIndex, config)
}

func (s *Selection[T]) Show(options ...ShowOption) {
	config := defaultShowConfig()
	for _, opt := range options {
		opt(config)
	}

	data := s.Collect()
	if len(data) == 0 {
		fmt.Println("Empty dataset")
		return
	}

	showMapTable(data, s.pipeline.originalIndex, config)
}

func ShowPositions(pos PositionIndex, options ...ShowOption) {
	config := defaultShowConfig()
	for _, opt := range options {
		opt(config)
	}

	if len(pos.Rows) == 0 {
		fmt.Println("No positions")
		return
	}

	fmt.Printf("Original Positions: %v\n", pos.Rows)
	if len(pos.Cols) > 0 {
		fmt.Printf("Column Positions: %v\n", pos.Cols)
	}

	style := getTableStyle(config.style)

	headers := []string{"#"}
	if pos.IsMatrix() {
		headers = append(headers, "Row", "Col")
	} else {
		headers = []string{"Position"}
	}

	rows := make([][]string, 0)
	if pos.IsMatrix() {
		for i, row := range pos.Rows {
			for _, col := range pos.Cols {
				rows = append(rows, []string{fmt.Sprintf("%d", i+1), fmt.Sprintf("%d", row), fmt.Sprintf("%d", col)})
			}
		}
	} else {
		for i, row := range pos.Rows {
			rows = append(rows, []string{fmt.Sprintf("%d", i+1), fmt.Sprintf("%d", row)})
		}
	}

	renderTable(headers, rows, style, config)
}

func showTable[T any](data []T, originalIndex []int, config *ShowConfig) {
	if len(data) == 0 {
		return
	}

	headers, rows := extractStructData(data, originalIndex, config)
	style := getTableStyle(config.style)

	if config.title != "" {
		printTitle(config.title, config)
	}

	truncatedRows, omitted := truncateRows(rows, config.maxRows)
	renderTable(headers, truncatedRows, style, config)

	if omitted > 0 {
		fmt.Printf("[%d rows × %d columns] (showing %d of %d)\n",
			len(rows), len(headers), len(truncatedRows), len(rows))
	} else {
		fmt.Printf("[%d rows × %d columns]\n", len(rows), len(headers))
	}
}

func showMapTable(data []map[string]any, originalIndex []int, config *ShowConfig) {
	if len(data) == 0 {
		return
	}

	headers, rows := extractMapData(data, originalIndex, config)
	style := getTableStyle(config.style)

	if config.title != "" {
		printTitle(config.title, config)
	}

	truncatedRows, omitted := truncateRows(rows, config.maxRows)
	renderTable(headers, truncatedRows, style, config)

	if omitted > 0 {
		fmt.Printf("[%d rows × %d columns] (showing %d of %d)\n",
			len(rows), len(headers), len(truncatedRows), len(rows))
	} else {
		fmt.Printf("[%d rows × %d columns]\n", len(rows), len(headers))
	}
}

func extractStructData[T any](data []T, originalIndex []int, config *ShowConfig) ([]string, [][]string) {
	if len(data) == 0 {
		return nil, nil
	}

	val := reflect.ValueOf(data[0])
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	headers := make([]string, 0)
	if config.showRowNumbers || config.showOriginalIdx {
		headers = append(headers, "#")
	}

	if val.Kind() == reflect.Struct {
		typ := val.Type()
		for i := 0; i < typ.NumField(); i++ {
			headers = append(headers, typ.Field(i).Name)
		}
	} else {
		headers = append(headers, "Value")
	}

	rows := make([][]string, 0, len(data))
	for i, item := range data {
		row := make([]string, 0)

		if config.showRowNumbers {
			row = append(row, fmt.Sprintf("%d", i+1))
		} else if config.showOriginalIdx && len(originalIndex) > i {
			row = append(row, fmt.Sprintf("%d", originalIndex[i]))
		}

		itemVal := reflect.ValueOf(item)
		if itemVal.Kind() == reflect.Ptr {
			itemVal = itemVal.Elem()
		}

		if itemVal.Kind() == reflect.Struct {
			for j := 0; j < itemVal.NumField(); j++ {
				field := itemVal.Field(j)
				row = append(row, formatValue(field.Interface(), config))
			}
		} else {
			row = append(row, formatValue(item, config))
		}

		rows = append(rows, row)
	}

	return headers, rows
}

func extractMapData(data []map[string]any, originalIndex []int, config *ShowConfig) ([]string, [][]string) {
	if len(data) == 0 {
		return nil, nil
	}

	headers := make([]string, 0)
	if config.showRowNumbers || config.showOriginalIdx {
		headers = append(headers, "#")
	}

	fieldOrder := make([]string, 0)
	for key := range data[0] {
		fieldOrder = append(fieldOrder, key)
	}
	sort.Strings(fieldOrder)
	headers = append(headers, fieldOrder...)

	rows := make([][]string, 0, len(data))
	for i, item := range data {
		row := make([]string, 0)

		if config.showRowNumbers {
			row = append(row, fmt.Sprintf("%d", i+1))
		} else if config.showOriginalIdx && len(originalIndex) > i {
			row = append(row, fmt.Sprintf("%d", originalIndex[i]))
		}

		for _, key := range fieldOrder {
			row = append(row, formatValue(item[key], config))
		}

		rows = append(rows, row)
	}

	return headers, rows
}

func formatValue(v any, config *ShowConfig) string {
	if v == nil {
		return "nil"
	}

	switch val := v.(type) {
	case float32, float64:
		if config.floatPrecision == 0 {
			return fmt.Sprintf("%.0f", val)
		}
		return fmt.Sprintf("%.*f", config.floatPrecision, val)
	case bool:
		if config.boolStyle == "symbols" {
			if val {
				return "✓"
			}
			return "✗"
		}
		return fmt.Sprintf("%v", val)
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}

func truncateRows(rows [][]string, maxRows int) ([][]string, int) {
	if len(rows) <= maxRows {
		return rows, 0
	}

	half := maxRows / 2
	result := make([][]string, 0, maxRows+1)

	result = append(result, rows[:half]...)

	if len(rows[0]) > 0 {
		separator := make([]string, len(rows[0]))
		for i := range separator {
			separator[i] = "..."
		}
		result = append(result, separator)
	}

	result = append(result, rows[len(rows)-half:]...)

	return result, len(rows) - maxRows
}

func renderTable(headers []string, rows [][]string, style tableStyle, config *ShowConfig) {
	if len(headers) == 0 {
		return
	}

	colWidths := calculateColumnWidths(headers, rows, config.maxColWidth)

	if style.topLeft == "" {
		renderMinimalTable(headers, rows, colWidths)
		return
	}

	printTopBorder(colWidths, style)
	printRow(headers, colWidths, style, true)
	printHeaderSeparator(colWidths, style)

	for _, row := range rows {
		printRow(row, colWidths, style, false)
	}

	printBottomBorder(colWidths, style)
}

func renderMinimalTable(headers []string, rows [][]string, colWidths []int) {
	printMinimalRow(headers, colWidths, true)

	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w + 2
	}
	fmt.Println(strings.Repeat("─", totalWidth-2))

	for _, row := range rows {
		printMinimalRow(row, colWidths, false)
	}
}

func calculateColumnWidths(headers []string, rows [][]string, maxColWidth int) []int {
	widths := make([]int, len(headers))

	for i, header := range headers {
		widths[i] = utf8.RuneCountInString(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) {
				cellLen := utf8.RuneCountInString(cell)
				if cellLen > widths[i] {
					widths[i] = cellLen
				}
			}
		}
	}

	for i := range widths {
		if widths[i] > maxColWidth {
			widths[i] = maxColWidth
		}
	}

	return widths
}

func printTitle(title string, config *ShowConfig) {
	fmt.Println()
	fmt.Printf("%s%s\n", strings.Repeat(" ", 8), title)
}

func printTopBorder(widths []int, style tableStyle) {
	if style.topLeft == "" {
		return
	}

	fmt.Print(style.topLeft)
	for i, width := range widths {
		fmt.Print(strings.Repeat(style.horizontal, width+2))
		if i < len(widths)-1 {
			fmt.Print(style.topCross)
		}
	}
	fmt.Println(style.topRight)
}

func printBottomBorder(widths []int, style tableStyle) {
	if style.bottomLeft == "" {
		return
	}

	fmt.Print(style.bottomLeft)
	for i, width := range widths {
		fmt.Print(strings.Repeat(style.horizontal, width+2))
		if i < len(widths)-1 {
			fmt.Print(style.bottomCross)
		}
	}
	fmt.Println(style.bottomRight)
}

func printHeaderSeparator(widths []int, style tableStyle) {
	if style.headerSep == "" {
		return
	}

	fmt.Print(style.headerSep)
	for i, width := range widths {
		fmt.Print(strings.Repeat(style.horizontal, width+2))
		if i < len(widths)-1 {
			fmt.Print(style.cross)
		}
	}
	if style.headerSep == "|" {
		fmt.Println(style.headerSep)
	} else {
		fmt.Println(style.rightCross)
	}
}

func printRow(cells []string, widths []int, style tableStyle, isHeader bool) {
	fmt.Print(style.vertical)

	for i, width := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}

		cell = truncateString(cell, width)

		aligned := alignCell(cell, width, isHeader || isNumeric(cell))
		fmt.Printf(" %s ", aligned)
		fmt.Print(style.vertical)
	}
	fmt.Println()
}

func printMinimalRow(cells []string, widths []int, isHeader bool) {
	for i, width := range widths {
		cell := ""
		if i < len(cells) {
			cell = cells[i]
		}

		cell = truncateString(cell, width)
		aligned := alignCell(cell, width, isHeader || isNumeric(cell))
		fmt.Printf("%s", aligned)

		if i < len(widths)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()
}

func alignCell(cell string, width int, rightAlign bool) string {
	cellLen := utf8.RuneCountInString(cell)
	padding := width - cellLen

	if padding <= 0 {
		return cell
	}

	if rightAlign {
		return strings.Repeat(" ", padding) + cell
	}
	return cell + strings.Repeat(" ", padding)
}

func truncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}

	if maxLen <= 3 {
		return strings.Repeat(".", maxLen)
	}

	runes := []rune(s)
	return string(runes[:maxLen-3]) + "..."
}

func isNumeric(s string) bool {
	if s == "" || s == "..." {
		return false
	}

	for _, r := range s {
		if r >= '0' && r <= '9' || r == '.' || r == '-' || r == '+' {
			continue
		}
		return false
	}
	return true
}

func (c *Condition[T]) Show(options ...ShowOption) {
	filtered := c.Collect()
	config := defaultShowConfig()
	for _, opt := range options {
		opt(config)
	}

	if len(filtered) == 0 {
		fmt.Println("Empty dataset")
		return
	}

	indices := c.Positions().Rows
	showTable(filtered, indices, config)
}

func (s *Sorter[T]) Show(options ...ShowOption) {
	sorted := s.Collect()
	config := defaultShowConfig()
	for _, opt := range options {
		opt(config)
	}

	if len(sorted) == 0 {
		fmt.Println("Empty dataset")
		return
	}

	showTable(sorted, s.pipeline.originalIndex, config)
}

func (s *Sorter[T]) AtRow(indices ...int) *Pipeline[T] {
	sorted := s.Collect()
	return From(sorted).AtRow(indices...)
}

func (s *Sorter[T]) AtCol(indices ...int) *Selection[T] {
	sorted := s.Collect()
	return From(sorted).AtCol(indices...)
}
