# Large Dataset Performance Benchmarks

## Dataset Specifications

- **Rows:** 100,000
- **Columns:** 10 (ID, Name, Value1, Value2, Value3, Category, Score, Percentage, Count, Active)
- **Data Distribution:** Normally distributed numeric values (using Box-Muller transform)
- **Categories:** 10 distinct values (A-J)
- **Active Rate:** ~66% of records

## Benchmark Results (100,000 Rows)

### Filtering Operations

| Operation | Time (ms) | Memory (MB) | Allocations |
|-----------|-----------|-------------|-------------|
| **Where Single** | 23.5 | 19.6 | 300,007 |
| **Where Multiple** (3 filters) | 70.9 | 88.2 | 713,061 |
| **Where Complex** (4 filters) | 77.2 | 87.6 | 756,015 |
| **Where Or** (3 conditions) | 44.5 | 32.7 | 720,011 |
| **Where OneOf** (5 values) | 37.3 | 19.6 | 300,008 |

**Key Insights:**
- Single filter on 100K rows: **23.5ms** - excellent performance
- Multiple filters scale linearly
- OneOf is more efficient than multiple Or conditions

### Sorting Operations

| Operation | Time (ms) | Memory (MB) | Allocations |
|-----------|-----------|-------------|-------------|
| **OrderBy Single Field** | 520.1 | 357.4 | 3,489,643 |
| **OrderBy Two Fields** | 975.4 | 645.6 | 6,491,119 |
| **OrderBy Three Fields** | 1,569.1 | 912.9 | 9,276,124 |

**Key Insights:**
- Sorting is the most expensive operation
- Multi-field sorts scale with number of fields
- Each additional sort field roughly doubles the cost

### Grouping Operations

| Operation | Time (ms) | Memory (MB) | Allocations |
|-----------|-----------|-------------|-------------|
| **GroupBy Count** | 17.6 | 10.4 | 100,007 |
| **GroupBy Sum** | 30.2 | 20.0 | 200,007 |
| **GroupBy Avg** | 35.0 | 20.0 | 200,013 |

**Key Insights:**
- GroupBy operations are very efficient
- Count is fastest (17.6ms for 100K rows)
- Avg slightly slower due to division calculations

### Complex Pipelines

| Operation | Time (ms) | Memory (MB) | Allocations |
|-----------|-----------|-------------|-------------|
| **Filter + Sort** | 230.3 | 161.2 | 1,570,423 |
| **Filter + Sort + Limit** | 327.7 | 225.1 | 2,142,134 |
| **Filter + GroupBy** | 53.5 | 47.1 | 433,350 |
| **Complex** (filter + sort + limit) | 386.4 | 287.9 | 2,729,227 |
| **Very Complex** (5 filters + 2 sorts + limit) | 172.7 | 166.8 | 1,419,493 |

**Key Insights:**
- Complex pipelines under **400ms** for 100K rows
- Filter + GroupBy much faster than Filter + Sort
- Limiting results doesn't reduce sort cost (full dataset sorted first)

### Transform Operations

| Operation | Time (ms) | Memory (MB) | Allocations |
|-----------|-----------|-------------|-------------|
| **Transform** (all rows) | 3.4 | 11.2 | 4 |
| **Transform After Filter** | 36.7 | 41.3 | 300,014 |

**Key Insights:**
- Transform is extremely fast: **3.4ms** for 100K rows
- Only 4 allocations for full dataset transform!
- Most efficient operation per row

### Selection Operations

| Operation | Time (ms) | Memory (MB) | Allocations |
|-----------|-----------|-------------|-------------|
| **Select** (4 columns) | 50.2 | 44.8 | 300,003 |
| **Select After Filter** | 61.5 | 63.7 | 500,013 |

### Position-Based Operations

| Operation | Time (µs) | Memory (KB) | Allocations |
|-----------|-----------|-------------|-------------|
| **AtRow** (5 specific rows) | 190 | 784 | 5 |
| **RowRange** (1,000 rows) | 179 | 888 | 5 |
| **Tail** (1,000 rows) | 172 | 880 | 3 |

**Key Insights:**
- Position-based operations are **extremely fast** (sub-millisecond)
- Constant-time access regardless of dataset size

### Real-World Scenarios

| Scenario | Description | Time (ms) | Memory (MB) |
|----------|-------------|-----------|-------------|
| **Top Performers** | Filter active + high scores, sort, limit 100 | 166.9 | 118.8 |
| **Category Analysis** | Filter + GroupBy average score | 56.8 | 47.1 |
| **Outlier Detection** | Find values >2σ from mean | 55.2 | 43.0 |
| **Detailed Report** | Multi-filter + multi-sort + select columns | 102.3 | 68.8 |

## Performance Characteristics

### Time Complexity (100K rows)

| Operation Type | Time Range | Notes |
|----------------|------------|-------|
| **Position Access** | <200 µs | AtRow, RowRange, Tail - O(k) where k=selected rows |
| **Transform** | 3-4 ms | O(n) - highly optimized |
| **GroupBy** | 17-35 ms | O(n) - hash-based grouping |
| **Simple Filter** | 23-40 ms | O(n) - linear scan |
| **Multiple Filters** | 70-80 ms | O(n*f) where f=filter count |
| **Selection** | 50-60 ms | O(n*c) where c=column count |
| **Sort Single** | 520 ms | O(n log n) |
| **Sort Multiple** | 1-1.6 sec | O(n log n * f) where f=field count |

### Memory Efficiency

- **Filtering:** ~200-300 bytes per row retained
- **Sorting:** ~3.5 KB per row (includes temp buffers)
- **GroupBy:** ~100-200 bytes per row
- **Transform:** Minimal overhead (~112 bytes per row)
- **Position Ops:** Negligible (<10 bytes per row)

### Throughput Metrics (100,000 rows)

| Operation | Rows/Second | Rows/ms |
|-----------|-------------|---------|
| Transform | 29.3 million | 29,300 |
| GroupBy Count | 5.7 million | 5,700 |
| Simple Filter | 4.3 million | 4,300 |
| Multiple Filters | 1.4 million | 1,400 |
| Sort Single | 192 thousand | 192 |

## Recommendations for Large Datasets

### ✅ Best Practices

1. **Use OneOf instead of multiple Or conditions**
   - OneOf: 37.3ms vs Or chain: 44.5ms (19% faster)

2. **Filter before sorting**
   - Reduces dataset size before expensive sort operation
   - Filter 50% → Sort is 2x faster

3. **Limit early when possible**
   - Use Where before OrderBy for better performance
   - GroupBy operations don't benefit from Limit

4. **Transform is highly efficient**
   - Can safely transform 100K rows in 3.4ms
   - Only 4 allocations for entire operation

5. **Use GroupBy for aggregations**
   - Much faster than filter + manual counting
   - 17.6ms for counting across 100K rows

### ⚠️ Performance Considerations

1. **Sorting is expensive**
   - Single field sort: 520ms
   - Consider filtering first to reduce dataset
   - Multi-field sorts scale multiplicatively

2. **Selection creates new structures**
   - Select creates map[string]any for each row
   - Use Select only when necessary
   - 50ms overhead for 100K rows

3. **Memory allocation patterns**
   - Filtering pre-allocates with estimates
   - Sorting requires temporary buffers
   - GroupBy uses hash maps (efficient)

## Comparison with Other Operations

```
100,000 rows processing time:

Transform:           ████ 3.4 ms
GroupBy Count:       ██████████████████ 17.6 ms
Single Filter:       ████████████████████████ 23.5 ms
OneOf Filter:        █████████████████████████████████████ 37.3 ms
Select:              ██████████████████████████████████████████████████ 50.2 ms
Multiple Filters:    ███████████████████████████████████████████████████████████████████████ 70.9 ms
Sort Single:         ████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████████ 520.1 ms
```

## Scalability Analysis

Based on the benchmarks, here's how performance scales:

| Rows | Simple Filter | Sort | GroupBy Count | Complex Pipeline |
|------|---------------|------|---------------|------------------|
| 100 | 24 µs | 167 µs | ~165 µs | ~3 ms |
| 1,000 | 232 µs | 1.9 ms | ~160 µs | ~59 ms |
| 10,000 | 2.8 ms | 18.7 ms | ~170 µs | N/A |
| 100,000 | 23.5 ms | 520 ms | 17.6 ms | 386 ms |

**Observations:**
- Filtering scales linearly: O(n)
- GroupBy is nearly constant for hash operations: O(1) lookup
- Sorting follows O(n log n) pattern
- Complex pipelines scale with their most expensive operation

## Conclusion

The plyGO library performs excellently on large datasets (100K rows):

- ✅ **Fast filtering:** 23.5ms for single condition
- ✅ **Efficient aggregation:** 17.6ms for GroupBy
- ✅ **Ultra-fast transforms:** 3.4ms with only 4 allocations
- ✅ **Instant position access:** <200µs for any row subset
- ⚠️ **Sorting is expensive:** 520ms for single field (expected for comparison sorts)

The library maintains its fluent, elegant API while delivering strong performance for data analysis tasks on datasets up to 100,000+ rows.
