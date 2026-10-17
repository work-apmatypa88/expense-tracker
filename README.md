# expense-tracker

A tiny command-line expense tracker written in Go (standard library only). Expenses are appended to `expenses.csv` in the current directory (override with `EXPENSES_FILE`).

## Usage

Requires Go 1.22+.

```sh
go run . add 12.50 food lunch with team
go run . add 40 transport monthly pass
go run . list
go run . summary
```

Example `summary` output:

```
food              12.50
transport         40.00
total             52.50
```

Build a binary with `go build -o expense-tracker .`
