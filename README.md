# one-mille

One Million Rows Challenge - High-Performance CSV Import to SQLite in Go

## Description

The "One Million Rows Challenge" tests different approaches to import large CSV datasets (up to 2 million rows) into SQLite as fast as possible. This project demonstrates various optimization techniques including concurrency, batch processing, and thread-safe operations to achieve maximum import speed.

## Performance Challenge Results

The One Million Rows Challenge shows dramatic performance differences depending on implementation:

**Test Hardware:** MacBook Pro (M2 Pro, 16 GB RAM)

| Solution       | Description                             | Time (1M Records) | Speedup       |
| -------------- | --------------------------------------- | ----------------- | ------------- |
| Solution One   | Sequential Single Insert                | ~16.5s            | 1x (Baseline) |
| Solution Two   | 50 Concurrent Workers                   | ~17s              | ~1x           |
| Solution Three | Single Worker with Batch Insert         | ~1.8s             | ~9x           |
| Solution Four  | 10 Concurrent Workers with Batch Insert | ~1.7s             | ~10x          |

## The Four Import Strategies

### Solution One - Sequential Processing

- **Strategy**: One goroutine reads CSV, another executes single INSERT statements
- **Characteristics**: Mutex-protected single transactions
- **Performance**: Slowest solution due to sequential processing
- **Bottleneck**: Individual INSERT statements with separate transactions

### Solution Two - Concurrent Workers

- **Strategy**: 50 parallel worker goroutines process CSV rows
- **Characteristics**: Each worker executes its own INSERT statements
- **Performance**: Significant improvement through parallelization
- **Bottleneck**: Still individual INSERT statements, but parallel

### Solution Three - Batch Processing

- **Strategy**: Single worker collects 1000 records and executes batch INSERTs
- **Characteristics**: Reduced number of transactions through batching
- **Performance**: Dramatic improvement due to fewer database roundtrips
- **Advantage**: Optimal balance between memory usage and performance

### Solution Four - Concurrent Batch Processing (Optimal)

- **Strategy**: 10 parallel workers with 10,000 record batches each
- **Characteristics**: Combines concurrency with batch processing
- **Performance**: Best performance through optimal resource utilization
- **Key Feature**: Mutex-protected batch inserts to avoid lock conflicts

## Installation

```bash
go get github.com/mattn/go-sqlite3
```

## Usage

```bash
# Clone repository
git clone https://github.com/jpierer/one-mille.git
cd one-mille

# Install dependencies
go mod tidy

# Extract the CSV file (required for testing)
unzip customers-1m.csv.zip

# Start the One Million Rows Challenge
go run main.go
```

## Configuration

In `main.go` you can select which CSV file to use:

```go
// const CSV_FILE = "customers-100.csv"     // 100 records
const CSV_FILE = "customers-1m.csv"        // 1 million records
```

The program runs all solutions sequentially:

```go
func (app *App) Run() {
    defer app.db.Close()

    app.TruncateDB()
    app.SolutionOne()

    app.TruncateDB()
    app.SolutionTwo()

    app.TruncateDB()
    app.SolutionThree()

    app.TruncateDB()
    app.SolutionFour()

    // todo add even faster solutions
}
```

## Example Output

```
Starting solution one ...
Solution one done in 16.52167625s
Starting solution two ...
Solution two done in 16.973486334s
Starting solution three ...
Solution three done in 1.7710695s
Starting solution four ...
Solution four done in 1.742708875s
```

## Technical Details

- **Language**: Go 1.24.1
- **Database**: SQLite3 with WAL mode for better concurrent performance
- **CSV Format**: id, name, email, company, city, country, birthday
- **Import Optimizations**:
  - WAL mode for SQLite
  - Prepared statements
  - Batch transactions
  - Buffered channels
  - Worker pool pattern

## Key Learnings from the One Million Rows Challenge

1. **Batch processing** is the critical performance factor for millions of records
2. **Concurrency** provides additional improvements, but with diminishing returns
3. **Mutex synchronization** on batch inserts prevents SQLite lock conflicts
4. **Channel buffering** reduces goroutine blocking with large datasets
5. **WAL mode** significantly improves concurrent write performance

## Support Me

Give a ⭐ if this One Million Rows Challenge project helped you achieve faster CSV imports!

---

_Because I was lazy, I generated this README with Claude_
