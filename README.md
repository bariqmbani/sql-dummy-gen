# SQL Table Dummy Data Generator

This tool provides a CLI for generating dummy SQL `INSERT` statements based on a provided SQL DDL (Create Table) file. It is designed to help developers quickly populate their local or development databases with realistic-looking test data.

## Main Use Case

The primary use case is when you have a database schema (defined in a `.sql` file) and you need to verify performance, test queries, or just have some data to work with in your UI without manually writing thousands of `INSERT` statements. This tool automates that process by analyzing your table structure and generating compliant dummy data.

**Key Features:**
*   **DDL Parsing**: Automatically analyses your schema to choose appropriate data types.
*   **Batch Processing**: Efficiently handles large datasets by batching `INSERT` statements.
*   **Performance Tracking**: Displays real-time progress and total execution time.
*   **Customizable**: Control row counts, date ranges for audit columns, and output locations.

## Installation

### From Remote (Recommended)

You can install the latest version directly without cloning the repository. The binary will be named `sql-dummy-gen` automatically:

```bash
go install github.com/bariqmbani/sql-dummy-gen@latest
```

### From Source

1.  Clone the repository:
    ```bash
    git clone https://github.com/bariqmbani/sql-dummy-gen.git
    cd sql-dummy-gen
    ```

2.  Build the binary:
    ```bash
    go build -o sql-dummy-gen main.go
    ```

3.  (Optional) Install to your `$GOPATH/bin` locally:
    ```bash
    go install .
    ```

## Usage

Once built, you can run the tool directly from the terminal.

### Basic Command

```bash
./sql-dummy-gen -ddl ./path/to/schema.sql
```

This will parse the provided DDL file and generate an SQL file containing `INSERT` statements (by default, `output-{ddl-filename}.sql`).

### Options

*   `-ddl`: (Required) Path to the SQL file containing `CREATE TABLE` statements.
*   `-num`: Number of rows to generate per table (default: 1).
*   `-output`: Custom output file path. If not provided, it defaults to `output-{ddl-filename}.sql`.
*   `-db`: Target database syntax (currently supports: `mysql`, default: `mysql`).
*   `-created-col`: Name of a timestamp/date column to populate with random dates within a specific range.
*   `-time-range`: The range for the `-created-col` flag in format `YYYY-MM-DD,YYYY-MM-DD`. Defaults to today if not specified.

### Examples

**Generate 1000 rows for proper load testing:**

```bash
./sql-dummy-gen -ddl schema.sql -num 1000
```

**Generate data with a specific audit date range:**

Many tables have a `created_at` or `txn_date` column. You can backfill data over a specific period:

```bash
./sql-dummy-gen -ddl schema.sql -num 500 -created-col created_at -time-range 2023-01-01,2023-12-31
```

## Recent Improvements

This tool has been optimized for performance with large datasets. Recent updates include:
*   **Batching**: `INSERT` statements are now generated in batches (default 1000 rows) to reduce memory overhead and file I/O operations.
*   **Progress Indicators**: For large jobs, the CLI now shows a progress counter.
*   **Dynamic Naming**: Output files are automatically named based on the input DDL if not specified.
