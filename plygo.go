package plygo

import (
	"math/rand/v2"
	"reflect"
	"sort"
	"strings"
	"sync"
)

type Pipeline[T any] struct {
	data          []T
	originalIndex []int
	typeCache     *typeCache
}

// typeCache stores reflection metadata for performance
type typeCache struct {
	fieldNames  []string
	fieldMap    map[string]int
	initialized bool
}

var (
	typeCacheLock sync.RWMutex
	typeCaches    = make(map[reflect.Type]*typeCache)
)

type PositionIndex struct {
	Rows []int
	Cols []int
}

func (pi PositionIndex) RowCount() int {
	return len(pi.Rows)
}

func (pi PositionIndex) ColCount() int {
	return len(pi.Cols)
}

func (pi PositionIndex) IsMatrix() bool {
	return len(pi.Cols) > 0
}

func From[T any](data []T) *Pipeline[T] {
	indices := make([]int, len(data))
	for i := range indices {
		indices[i] = i + 1
	}
	p := &Pipeline[T]{
		data:          data,
		originalIndex: indices,
	}
	p.initTypeCache()
	return p
}

func (p *Pipeline[T]) initTypeCache() {
	if len(p.data) == 0 {
		return
	}

	var zero T
	typ := reflect.TypeOf(zero)
	if typ == nil {
		return
	}

	// Try to get cached type info
	typeCacheLock.RLock()
	cache, exists := typeCaches[typ]
	typeCacheLock.RUnlock()

	if exists {
		p.typeCache = cache
		return
	}

	// Create new cache
	typeCacheLock.Lock()
	defer typeCacheLock.Unlock()

	// Double-check after acquiring write lock
	if cache, exists := typeCaches[typ]; exists {
		p.typeCache = cache
		return
	}

	cache = &typeCache{
		fieldMap: make(map[string]int),
	}

	val := reflect.ValueOf(p.data[0])
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		t := val.Type()
		cache.fieldNames = make([]string, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fieldName := t.Field(i).Name
			cache.fieldNames[i] = fieldName
			cache.fieldMap[fieldName] = i
		}
		cache.initialized = true
	}

	typeCaches[typ] = cache
	p.typeCache = cache
}

func (p *Pipeline[T]) AtRow(indices ...int) *Pipeline[T] {
	if len(indices) == 0 {
		return &Pipeline[T]{data: []T{}, originalIndex: []int{}, typeCache: p.typeCache}
	}

	result := make([]T, 0, len(indices))
	resultIdx := make([]int, 0, len(indices))

	for _, idx := range indices {
		pos := p.normalizeIndex(idx)
		if pos >= 0 && pos < len(p.data) {
			result = append(result, p.data[pos])
			resultIdx = append(resultIdx, p.originalIndex[pos])
		}
	}

	return &Pipeline[T]{
		data:          result,
		originalIndex: resultIdx,
		typeCache:     p.typeCache,
	}
}

func (p *Pipeline[T]) RowRange(start, end int) *Pipeline[T] {
	startPos := p.normalizeIndex(start)

	var endPos int
	if end == -1 {
		endPos = len(p.data)
	} else {
		endPos = p.normalizeIndex(end)
	}

	if endPos == -1 {
		endPos = len(p.data)
	}

	if startPos < 0 {
		startPos = 0
	}
	if endPos > len(p.data) {
		endPos = len(p.data)
	}
	if startPos >= endPos {
		return &Pipeline[T]{data: []T{}, originalIndex: []int{}, typeCache: p.typeCache}
	}

	result := make([]T, endPos-startPos)
	copy(result, p.data[startPos:endPos])

	resultIdx := make([]int, endPos-startPos)
	copy(resultIdx, p.originalIndex[startPos:endPos])

	return &Pipeline[T]{
		data:          result,
		originalIndex: resultIdx,
		typeCache:     p.typeCache,
	}
}

func (p *Pipeline[T]) normalizeIndex(idx int) int {
	if idx > 0 {
		return idx - 1
	}
	if idx < 0 {
		return len(p.data) + idx
	}
	return -1
}

func (p *Pipeline[T]) Tail(n int) *Pipeline[T] {
	if n <= 0 {
		return &Pipeline[T]{data: []T{}, originalIndex: []int{}, typeCache: p.typeCache}
	}
	if n >= len(p.data) {
		return p
	}

	start := len(p.data) - n
	result := make([]T, n)
	copy(result, p.data[start:])

	resultIdx := make([]int, n)
	copy(resultIdx, p.originalIndex[start:])

	return &Pipeline[T]{
		data:          result,
		originalIndex: resultIdx,
		typeCache:     p.typeCache,
	}
}

func (p *Pipeline[T]) Sample(n int) *Pipeline[T] {
	if n <= 0 || len(p.data) == 0 {
		return &Pipeline[T]{data: []T{}, originalIndex: []int{}, typeCache: p.typeCache}
	}
	if n >= len(p.data) {
		return p
	}

	indices := rand.Perm(len(p.data))[:n]
	sort.Ints(indices)

	result := make([]T, n)
	resultIdx := make([]int, n)

	for i, idx := range indices {
		result[i] = p.data[idx]
		resultIdx[i] = p.originalIndex[idx]
	}

	return &Pipeline[T]{
		data:          result,
		originalIndex: resultIdx,
		typeCache:     p.typeCache,
	}
}

func (p *Pipeline[T]) Slice(start, end, step int) *Pipeline[T] {
	if step == 0 {
		step = 1
	}

	startPos := p.normalizeIndex(start)

	var endPos int
	if end == -1 {
		endPos = len(p.data)
	} else {
		endPos = p.normalizeIndex(end)
	}

	if endPos == -1 {
		endPos = len(p.data)
	}

	if startPos < 0 {
		startPos = 0
	}
	if endPos > len(p.data) {
		endPos = len(p.data)
	}

	result := make([]T, 0)
	resultIdx := make([]int, 0)

	if step > 0 {
		for i := startPos; i < endPos; i += step {
			result = append(result, p.data[i])
			resultIdx = append(resultIdx, p.originalIndex[i])
		}
	} else {
		for i := endPos - 1; i >= startPos; i += step {
			result = append(result, p.data[i])
			resultIdx = append(resultIdx, p.originalIndex[i])
		}
	}

	return &Pipeline[T]{
		data:          result,
		originalIndex: resultIdx,
	}
}

func (p *Pipeline[T]) Positions() PositionIndex {
	return PositionIndex{
		Rows: p.originalIndex,
		Cols: []int{},
	}
}

func (p *Pipeline[T]) Which() []int {
	return p.originalIndex
}

func (p *Pipeline[T]) FieldNames() []string {
	if len(p.data) == 0 {
		return []string{}
	}

	// Use cache if available
	if p.typeCache != nil && p.typeCache.initialized {
		return p.typeCache.fieldNames
	}

	// Fallback for maps and other types
	val := reflect.ValueOf(p.data[0])
	if val.Kind() == reflect.Map {
		keys := val.MapKeys()
		result := make([]string, len(keys))
		for i, key := range keys {
			result[i] = key.String()
		}
		sort.Strings(result)
		return result
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return []string{}
	}

	typ := val.Type()
	result := make([]string, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		result[i] = typ.Field(i).Name
	}

	return result
}

func (p *Pipeline[T]) FieldCount() int {
	return len(p.FieldNames())
}

func (p *Pipeline[T]) AtCol(indices ...int) *Selection[T] {
	if len(indices) == 0 {
		return &Selection[T]{pipeline: p, fields: []string{}}
	}

	fieldNames := p.FieldNames()
	selectedFields := make([]string, 0, len(indices))

	for _, idx := range indices {
		pos := p.normalizeIndex(idx)
		if pos >= 0 && pos < len(fieldNames) {
			selectedFields = append(selectedFields, fieldNames[pos])
		}
	}

	return &Selection[T]{
		pipeline: p,
		fields:   selectedFields,
	}
}

func (p *Pipeline[T]) ColRange(start, end int) *Selection[T] {
	fieldNames := p.FieldNames()

	startPos := p.normalizeIndex(start)

	var endPos int
	if end == -1 {
		endPos = len(fieldNames)
	} else {
		endPos = p.normalizeIndex(end)
	}

	if endPos == -1 {
		endPos = len(fieldNames)
	}

	if startPos < 0 {
		startPos = 0
	}
	if endPos > len(fieldNames) {
		endPos = len(fieldNames)
	}
	if startPos >= endPos {
		return &Selection[T]{pipeline: p, fields: []string{}}
	}

	selectedFields := make([]string, endPos-startPos)
	copy(selectedFields, fieldNames[startPos:endPos])

	return &Selection[T]{
		pipeline: p,
		fields:   selectedFields,
	}
}

func (p *Pipeline[T]) Where(field string) *Condition[T] {
	return &Condition[T]{
		pipeline: p,
		field:    field,
		filters:  make([]filter[T], 0),
	}
}

func (p *Pipeline[T]) WhereSome(conditions ...*ConditionGroup[T]) *Pipeline[T] {
	if len(conditions) == 0 {
		return p
	}

	result := make([]T, 0)
	resultIdx := make([]int, 0)
	for i, item := range p.data {
		for _, cond := range conditions {
			if cond.evaluate(item) {
				result = append(result, item)
				resultIdx = append(resultIdx, p.originalIndex[i])
				break
			}
		}
	}
	return &Pipeline[T]{data: result, originalIndex: resultIdx}
}

func (p *Pipeline[T]) WhereEvery(conditions ...*ConditionGroup[T]) *Pipeline[T] {
	if len(conditions) == 0 {
		return p
	}

	result := make([]T, 0)
	resultIdx := make([]int, 0)
	for i, item := range p.data {
		match := true
		for _, cond := range conditions {
			if !cond.evaluate(item) {
				match = false
				break
			}
		}
		if match {
			result = append(result, item)
			resultIdx = append(resultIdx, p.originalIndex[i])
		}
	}
	return &Pipeline[T]{data: result, originalIndex: resultIdx}
}

func (p *Pipeline[T]) Select(fields ...string) *Selection[T] {
	return &Selection[T]{
		pipeline: p,
		fields:   fields,
	}
}

func (p *Pipeline[T]) OrderBy(field string) *Sorter[T] {
	return &Sorter[T]{
		pipeline: p,
		sorts:    []sortField{{field: field, desc: false}},
	}
}

func (p *Pipeline[T]) GroupBy(field string) *Grouping[T] {
	return &Grouping[T]{
		pipeline: p,
		field:    field,
	}
}

func (p *Pipeline[T]) Transform(fn func(T) T) *Pipeline[T] {
	result := make([]T, len(p.data))
	resultIdx := make([]int, len(p.data))
	copy(resultIdx, p.originalIndex)
	for i, item := range p.data {
		result[i] = fn(item)
	}
	return &Pipeline[T]{data: result, originalIndex: resultIdx}
}

func (p *Pipeline[T]) Limit(n int) *Pipeline[T] {
	if n >= len(p.data) {
		return p
	}
	return &Pipeline[T]{data: p.data[:n], originalIndex: p.originalIndex[:n]}
}

func (p *Pipeline[T]) Skip(n int) *Pipeline[T] {
	if n >= len(p.data) {
		return &Pipeline[T]{data: []T{}, originalIndex: []int{}}
	}
	return &Pipeline[T]{data: p.data[n:], originalIndex: p.originalIndex[n:]}
}

func (p *Pipeline[T]) Distinct(field string) *Pipeline[T] {
	seen := make(map[any]bool)
	result := make([]T, 0)
	resultIdx := make([]int, 0)

	for i, item := range p.data {
		val := getFieldValue(item, field)
		key := valueKey(val)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
			resultIdx = append(resultIdx, p.originalIndex[i])
		}
	}
	return &Pipeline[T]{data: result, originalIndex: resultIdx}
}

func (p *Pipeline[T]) Collect() []T {
	return p.data
}

func (p *Pipeline[T]) First() (T, bool) {
	var zero T
	if len(p.data) == 0 {
		return zero, false
	}
	return p.data[0], true
}

func (p *Pipeline[T]) Last() (T, bool) {
	var zero T
	if len(p.data) == 0 {
		return zero, false
	}
	return p.data[len(p.data)-1], true
}

func (p *Pipeline[T]) Count() int {
	return len(p.data)
}

type Condition[T any] struct {
	pipeline *Pipeline[T]
	field    string
	filters  []filter[T]
	orMode   bool
}

type filter[T any] struct {
	fn     func(T) bool
	orMode bool
}

func (c *Condition[T]) Equals(value any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			return compareEqual(getFieldValue(item, c.field), value)
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) NotEquals(value any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			return !compareEqual(getFieldValue(item, c.field), value)
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) GreaterThan(value any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			return compareNumeric(getFieldValue(item, c.field), value, ">")
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) GreaterOrEqual(value any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			return compareNumeric(getFieldValue(item, c.field), value, ">=")
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) LessThan(value any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			return compareNumeric(getFieldValue(item, c.field), value, "<")
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) LessOrEqual(value any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			return compareNumeric(getFieldValue(item, c.field), value, "<=")
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) Between(min, max any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			val := getFieldValue(item, c.field)
			return compareNumeric(val, min, ">=") && compareNumeric(val, max, "<=")
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) OneOf(values ...any) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			val := getFieldValue(item, c.field)
			for _, v := range values {
				if compareEqual(val, v) {
					return true
				}
			}
			return false
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) Contains(substr string) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			val := getFieldValue(item, c.field)
			if str, ok := val.(string); ok {
				return strings.Contains(str, substr)
			}
			return false
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) StartsWith(prefix string) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			val := getFieldValue(item, c.field)
			if str, ok := val.(string); ok {
				return strings.HasPrefix(str, prefix)
			}
			return false
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) EndsWith(suffix string) *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			val := getFieldValue(item, c.field)
			if str, ok := val.(string); ok {
				return strings.HasSuffix(str, suffix)
			}
			return false
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) IsTrue() *Condition[T] {
	return c.Equals(true)
}

func (c *Condition[T]) IsFalse() *Condition[T] {
	return c.Equals(false)
}

func (c *Condition[T]) IsNull() *Condition[T] {
	c.filters = append(c.filters, filter[T]{
		fn: func(item T) bool {
			return getFieldValue(item, c.field) == nil
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *Condition[T]) And(field string) *Condition[T] {
	c.field = field
	c.orMode = false
	return c
}

func (c *Condition[T]) Or(field string) *Condition[T] {
	if len(c.filters) > 0 {
		c.filters[len(c.filters)-1].orMode = false
	}
	c.field = field
	c.orMode = true
	return c
}

func (c *Condition[T]) Where(field string) *Condition[T] {
	filtered, indices := c.executeWithIndices()
	return &Condition[T]{
		pipeline: &Pipeline[T]{data: filtered, originalIndex: indices, typeCache: c.pipeline.typeCache},
		field:    field,
		filters:  make([]filter[T], 0),
	}
}

func (c *Condition[T]) Select(fields ...string) *Selection[T] {
	filtered, indices := c.executeWithIndices()
	return &Selection[T]{
		pipeline: &Pipeline[T]{data: filtered, originalIndex: indices, typeCache: c.pipeline.typeCache},
		fields:   fields,
	}
}

func (c *Condition[T]) OrderBy(field string) *Sorter[T] {
	filtered, indices := c.executeWithIndices()
	return &Sorter[T]{
		pipeline: &Pipeline[T]{data: filtered, originalIndex: indices, typeCache: c.pipeline.typeCache},
		sorts:    []sortField{{field: field, desc: false}},
	}
}

func (c *Condition[T]) GroupBy(field string) *Grouping[T] {
	filtered, indices := c.executeWithIndices()
	return &Grouping[T]{
		pipeline: &Pipeline[T]{data: filtered, originalIndex: indices, typeCache: c.pipeline.typeCache},
		field:    field,
	}
}

func (c *Condition[T]) Transform(fn func(T) T) *Pipeline[T] {
	filtered, indices := c.executeWithIndices()
	p := &Pipeline[T]{data: filtered, originalIndex: indices, typeCache: c.pipeline.typeCache}
	return p.Transform(fn)
}

func (c *Condition[T]) Limit(n int) *Pipeline[T] {
	filtered, indices := c.executeWithIndices()
	p := &Pipeline[T]{data: filtered, originalIndex: indices, typeCache: c.pipeline.typeCache}
	return p.Limit(n)
}

func (c *Condition[T]) Distinct(field string) *Pipeline[T] {
	filtered, indices := c.executeWithIndices()
	p := &Pipeline[T]{data: filtered, originalIndex: indices, typeCache: c.pipeline.typeCache}
	return p.Distinct(field)
}

func (c *Condition[T]) Collect() []T {
	return c.execute()
}

func (c *Condition[T]) Positions() PositionIndex {
	// Execute filters to get matching items
	filtered := c.execute()

	if len(filtered) == 0 {
		return PositionIndex{Rows: []int{}, Cols: []int{}}
	}

	// Build a map of filtered items for matching
	filteredMap := make(map[any]bool)
	for _, item := range filtered {
		filteredMap[item] = true
	}

	// Find original indices of filtered items
	indices := make([]int, 0, len(filtered))
	for i, item := range c.pipeline.data {
		if filteredMap[item] {
			indices = append(indices, c.pipeline.originalIndex[i])
			delete(filteredMap, item) // Remove to handle duplicates correctly
		}
	}

	return PositionIndex{Rows: indices, Cols: []int{}}
}

func (c *Condition[T]) Which() []int {
	return c.Positions().Rows
}

func (c *Condition[T]) execute() []T {
	result, _ := c.executeWithIndices()
	return result
}

func (c *Condition[T]) executeWithIndices() ([]T, []int) {
	if len(c.filters) == 0 {
		return c.pipeline.data, c.pipeline.originalIndex
	}

	// Pre-allocate with estimated capacity
	estimatedSize := len(c.pipeline.data) / 2
	if estimatedSize < 16 {
		estimatedSize = 16
	}
	result := make([]T, 0, estimatedSize)
	indices := make([]int, 0, estimatedSize)

	for i, item := range c.pipeline.data {
		if c.evaluate(item) {
			result = append(result, item)
			indices = append(indices, c.pipeline.originalIndex[i])
		}
	}
	return result, indices
}

func (c *Condition[T]) evaluate(item T) bool {
	if len(c.filters) == 0 {
		return true
	}

	groups := make([][]filter[T], 0)
	currentGroup := make([]filter[T], 0)

	for i := 0; i < len(c.filters); i++ {
		f := c.filters[i]

		if f.orMode {
			if len(currentGroup) == 0 && i > 0 {
				currentGroup = append(currentGroup, c.filters[i-1])
				groups = groups[:len(groups)-1]
			}
			currentGroup = append(currentGroup, f)

			if i == len(c.filters)-1 || !c.filters[i+1].orMode {
				groups = append(groups, currentGroup)
				currentGroup = make([]filter[T], 0)
			}
		} else {
			if len(currentGroup) == 0 {
				groups = append(groups, []filter[T]{f})
			}
		}
	}

	for _, group := range groups {
		if len(group) == 1 {
			if !group[0].fn(item) {
				return false
			}
		} else {
			hasMatch := false
			for _, f := range group {
				if f.fn(item) {
					hasMatch = true
					break
				}
			}
			if !hasMatch {
				return false
			}
		}
	}

	return true
}

type ConditionGroup[T any] struct {
	field   string
	filters []filter[T]
	orMode  bool
}

func W[T any](field string) *ConditionGroup[T] {
	return &ConditionGroup[T]{
		field:   field,
		filters: make([]filter[T], 0),
	}
}

func WhereEvery[T any](conditions ...*ConditionGroup[T]) *ConditionGroup[T] {
	return &ConditionGroup[T]{
		filters: flattenConditions(conditions, false),
	}
}

func WhereSome[T any](conditions ...*ConditionGroup[T]) *ConditionGroup[T] {
	return &ConditionGroup[T]{
		filters: flattenConditions(conditions, true),
	}
}

func flattenConditions[T any](conditions []*ConditionGroup[T], orMode bool) []filter[T] {
	result := make([]filter[T], 0)
	for _, cond := range conditions {
		for _, f := range cond.filters {
			f.orMode = orMode
			result = append(result, f)
		}
	}
	return result
}

func (cg *ConditionGroup[T]) Equals(value any) *ConditionGroup[T] {
	cg.filters = append(cg.filters, filter[T]{
		fn: func(item T) bool {
			return compareEqual(getFieldValue(item, cg.field), value)
		},
		orMode: cg.orMode,
	})
	return cg
}

func (cg *ConditionGroup[T]) GreaterThan(value any) *ConditionGroup[T] {
	cg.filters = append(cg.filters, filter[T]{
		fn: func(item T) bool {
			return compareNumeric(getFieldValue(item, cg.field), value, ">")
		},
		orMode: cg.orMode,
	})
	return cg
}

func (cg *ConditionGroup[T]) LessThan(value any) *ConditionGroup[T] {
	cg.filters = append(cg.filters, filter[T]{
		fn: func(item T) bool {
			return compareNumeric(getFieldValue(item, cg.field), value, "<")
		},
		orMode: cg.orMode,
	})
	return cg
}

func (cg *ConditionGroup[T]) IsTrue() *ConditionGroup[T] {
	return cg.Equals(true)
}

func (cg *ConditionGroup[T]) evaluate(item T) bool {
	for _, f := range cg.filters {
		if !f.fn(item) {
			return false
		}
	}
	return true
}

type Selection[T any] struct {
	pipeline *Pipeline[T]
	fields   []string
}

func (s *Selection[T]) Where(field string) *ConditionMap {
	selected := s.execute()
	return &ConditionMap{
		pipeline: &Pipeline[map[string]any]{data: selected},
		field:    field,
		filters:  make([]filter[map[string]any], 0),
	}
}

func (s *Selection[T]) OrderBy(field string) *SorterMap {
	selected := s.execute()
	return &SorterMap{
		pipeline: &Pipeline[map[string]any]{data: selected},
		sorts:    []sortField{{field: field, desc: false}},
	}
}

func (s *Selection[T]) GroupBy(field string) *GroupingMap {
	selected := s.execute()
	return &GroupingMap{
		pipeline: &Pipeline[map[string]any]{data: selected},
		field:    field,
	}
}

func (s *Selection[T]) Collect() []map[string]any {
	return s.execute()
}

func (s *Selection[T]) Positions() PositionIndex {
	allFields := s.pipeline.FieldNames()
	colIndices := make([]int, 0, len(s.fields))

	for _, field := range s.fields {
		for i, f := range allFields {
			if f == field {
				colIndices = append(colIndices, i+1)
				break
			}
		}
	}

	return PositionIndex{
		Rows: s.pipeline.originalIndex,
		Cols: colIndices,
	}
}

func (s *Selection[T]) AtRow(indices ...int) *Selection[T] {
	return &Selection[T]{
		pipeline: s.pipeline.AtRow(indices...),
		fields:   s.fields,
	}
}

func (s *Selection[T]) RowRange(start, end int) *Selection[T] {
	return &Selection[T]{
		pipeline: s.pipeline.RowRange(start, end),
		fields:   s.fields,
	}
}

func (s *Selection[T]) execute() []map[string]any {
	result := make([]map[string]any, len(s.pipeline.data))

	for i, item := range s.pipeline.data {
		row := make(map[string]any)
		v := reflect.ValueOf(item)

		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.Struct {
			for _, fieldName := range s.fields {
				field := v.FieldByName(fieldName)
				if field.IsValid() && field.CanInterface() {
					row[fieldName] = field.Interface()
				}
			}
		}

		result[i] = row
	}

	return result
}

type ConditionMap struct {
	pipeline *Pipeline[map[string]any]
	field    string
	filters  []filter[map[string]any]
	orMode   bool
}

func (c *ConditionMap) Equals(value any) *ConditionMap {
	c.filters = append(c.filters, filter[map[string]any]{
		fn: func(item map[string]any) bool {
			return compareEqual(item[c.field], value)
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *ConditionMap) GreaterThan(value any) *ConditionMap {
	c.filters = append(c.filters, filter[map[string]any]{
		fn: func(item map[string]any) bool {
			return compareNumeric(item[c.field], value, ">")
		},
		orMode: c.orMode,
	})
	c.orMode = false
	return c
}

func (c *ConditionMap) Or(field string) *ConditionMap {
	c.field = field
	c.orMode = true
	return c
}

func (c *ConditionMap) Collect() []map[string]any {
	return c.execute()
}

func (c *ConditionMap) execute() []map[string]any {
	if len(c.filters) == 0 {
		return c.pipeline.data
	}

	result := make([]map[string]any, 0)
	for _, item := range c.pipeline.data {
		match := true
		for _, f := range c.filters {
			if !f.fn(item) {
				match = false
				break
			}
		}
		if match {
			result = append(result, item)
		}
	}
	return result
}

type sortField struct {
	field string
	desc  bool
}

type Sorter[T any] struct {
	pipeline *Pipeline[T]
	sorts    []sortField
}

func (s *Sorter[T]) Desc() *Sorter[T] {
	if len(s.sorts) > 0 {
		s.sorts[len(s.sorts)-1].desc = true
	}
	return s
}

func (s *Sorter[T]) Asc() *Sorter[T] {
	if len(s.sorts) > 0 {
		s.sorts[len(s.sorts)-1].desc = false
	}
	return s
}

func (s *Sorter[T]) ThenBy(field string) *Sorter[T] {
	s.sorts = append(s.sorts, sortField{field: field, desc: false})
	return s
}

func (s *Sorter[T]) Where(field string) *Condition[T] {
	sorted, indices := s.executeWithIndices()
	return &Condition[T]{
		pipeline: &Pipeline[T]{data: sorted, originalIndex: indices, typeCache: s.pipeline.typeCache},
		field:    field,
		filters:  make([]filter[T], 0),
	}
}

func (s *Sorter[T]) Select(fields ...string) *Selection[T] {
	sorted, indices := s.executeWithIndices()
	return &Selection[T]{
		pipeline: &Pipeline[T]{data: sorted, originalIndex: indices, typeCache: s.pipeline.typeCache},
		fields:   fields,
	}
}

func (s *Sorter[T]) Limit(n int) *Pipeline[T] {
	sorted, indices := s.executeWithIndices()
	p := &Pipeline[T]{data: sorted, originalIndex: indices, typeCache: s.pipeline.typeCache}
	return p.Limit(n)
}

func (s *Sorter[T]) Skip(n int) *Pipeline[T] {
	sorted, indices := s.executeWithIndices()
	p := &Pipeline[T]{data: sorted, originalIndex: indices, typeCache: s.pipeline.typeCache}
	return p.Skip(n)
}

func (s *Sorter[T]) Collect() []T {
	return s.execute()
}

func (s *Sorter[T]) execute() []T {
	result, _ := s.executeWithIndices()
	return result
}

func (s *Sorter[T]) executeWithIndices() ([]T, []int) {
	if len(s.sorts) == 0 {
		return s.pipeline.data, s.pipeline.originalIndex
	}

	n := len(s.pipeline.data)
	result := make([]T, n)
	copy(result, s.pipeline.data)

	indices := make([]int, n)
	copy(indices, s.pipeline.originalIndex)

	// Create index array for sorting
	sortIndices := make([]int, n)
	for i := range sortIndices {
		sortIndices[i] = i
	}

	// Sort the index array
	sort.Slice(sortIndices, func(i, j int) bool {
		for _, sf := range s.sorts {
			vi := getFieldValue(result[sortIndices[i]], sf.field)
			vj := getFieldValue(result[sortIndices[j]], sf.field)

			cmp := compareValues(vi, vj)
			if cmp != 0 {
				if sf.desc {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return false
	})

	// Reorder result and indices based on sorted indices
	sortedResult := make([]T, n)
	sortedIndices := make([]int, n)
	for i, idx := range sortIndices {
		sortedResult[i] = result[idx]
		sortedIndices[i] = indices[idx]
	}

	return sortedResult, sortedIndices
}

type SorterMap struct {
	pipeline *Pipeline[map[string]any]
	sorts    []sortField
}

func (s *SorterMap) Desc() *SorterMap {
	if len(s.sorts) > 0 {
		s.sorts[len(s.sorts)-1].desc = true
	}
	return s
}

func (s *SorterMap) ThenBy(field string) *SorterMap {
	s.sorts = append(s.sorts, sortField{field: field, desc: false})
	return s
}

func (s *SorterMap) Collect() []map[string]any {
	return s.execute()
}

func (s *SorterMap) execute() []map[string]any {
	if len(s.sorts) == 0 {
		return s.pipeline.data
	}

	result := make([]map[string]any, len(s.pipeline.data))
	copy(result, s.pipeline.data)

	sort.Slice(result, func(i, j int) bool {
		for _, sf := range s.sorts {
			vi := result[i][sf.field]
			vj := result[j][sf.field]

			cmp := compareValues(vi, vj)
			if cmp != 0 {
				if sf.desc {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return false
	})

	return result
}

type Grouping[T any] struct {
	pipeline *Pipeline[T]
	field    string
}

func (g *Grouping[T]) Count() map[any]int {
	result := make(map[any]int)

	for _, item := range g.pipeline.data {
		key := valueKey(getFieldValue(item, g.field))
		result[key]++
	}

	return result
}

func (g *Grouping[T]) Sum(sumField string) map[any]float64 {
	result := make(map[any]float64)

	for _, item := range g.pipeline.data {
		key := valueKey(getFieldValue(item, g.field))
		val := toFloat64(getFieldValue(item, sumField))
		result[key] += val
	}

	return result
}

func (g *Grouping[T]) Avg(avgField string) map[any]float64 {
	sums := make(map[any]float64)
	counts := make(map[any]int)

	for _, item := range g.pipeline.data {
		key := valueKey(getFieldValue(item, g.field))
		val := toFloat64(getFieldValue(item, avgField))
		sums[key] += val
		counts[key]++
	}

	result := make(map[any]float64)
	for key, sum := range sums {
		result[key] = sum / float64(counts[key])
	}

	return result
}

func (g *Grouping[T]) Min(minField string) map[any]any {
	result := make(map[any]any)

	for _, item := range g.pipeline.data {
		key := valueKey(getFieldValue(item, g.field))
		val := getFieldValue(item, minField)

		if existing, ok := result[key]; !ok || compareValues(val, existing) < 0 {
			result[key] = val
		}
	}

	return result
}

func (g *Grouping[T]) Max(maxField string) map[any]any {
	result := make(map[any]any)

	for _, item := range g.pipeline.data {
		key := valueKey(getFieldValue(item, g.field))
		val := getFieldValue(item, maxField)

		if existing, ok := result[key]; !ok || compareValues(val, existing) > 0 {
			result[key] = val
		}
	}

	return result
}

type GroupingMap struct {
	pipeline *Pipeline[map[string]any]
	field    string
}

func (g *GroupingMap) Count() map[any]int {
	result := make(map[any]int)

	for _, item := range g.pipeline.data {
		key := valueKey(item[g.field])
		result[key]++
	}

	return result
}

func (g *GroupingMap) Sum(sumField string) map[any]float64 {
	result := make(map[any]float64)

	for _, item := range g.pipeline.data {
		key := valueKey(item[g.field])
		val := toFloat64(item[sumField])
		result[key] += val
	}

	return result
}

func getFieldValue(item any, fieldName string) any {
	v := reflect.ValueOf(item)

	if v.Kind() == reflect.Map {
		mapVal := v.MapIndex(reflect.ValueOf(fieldName))
		if mapVal.IsValid() {
			return mapVal.Interface()
		}
		return nil
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	field := v.FieldByName(fieldName)
	if !field.IsValid() || !field.CanInterface() {
		return nil
	}

	return field.Interface()
}

func compareEqual(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}
	return reflect.DeepEqual(a, b)
}

func compareNumeric(a, b any, op string) bool {
	aVal := toFloat64(a)
	bVal := toFloat64(b)

	switch op {
	case ">":
		return aVal > bVal
	case ">=":
		return aVal >= bVal
	case "<":
		return aVal < bVal
	case "<=":
		return aVal <= bVal
	}
	return false
}

func compareValues(a, b any) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	aVal := toFloat64(a)
	bVal := toFloat64(b)

	if aVal < bVal {
		return -1
	}
	if aVal > bVal {
		return 1
	}

	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			if sa < sb {
				return -1
			}
			if sa > sb {
				return 1
			}
		}
	}

	return 0
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int8:
		return float64(val)
	case int16:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case uint:
		return float64(val)
	case uint8:
		return float64(val)
	case uint16:
		return float64(val)
	case uint32:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return val
	}
	return 0
}

func valueKey(v any) any {
	if v == nil {
		return "<nil>"
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Map {
		return "<complex>"
	}
	return v
}
