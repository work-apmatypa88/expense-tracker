// Command expense-tracker records expenses in a CSV file and summarises them by category.
package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func dbPath() string {
	if p := os.Getenv("EXPENSES_FILE"); p != "" {
		return p
	}
	return "expenses.csv"
}

// load returns rows as [date, category, amount, note].
func load() ([][]string, error) {
	f, err := os.Open(dbPath())
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return csv.NewReader(f).ReadAll()
}

func add(args []string) error {
	if len(args) < 2 {
		return errors.New("usage: expense-tracker add <amount> <category> [note...]")
	}
	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil || amount <= 0 {
		return fmt.Errorf("invalid amount %q", args[0])
	}
	f, err := os.OpenFile(dbPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	row := []string{time.Now().Format("2006-01-02"), strings.ToLower(args[1]), strconv.FormatFloat(amount, 'f', 2, 64), strings.Join(args[2:], " ")}
	if err := w.Write(row); err != nil {
		return err
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	fmt.Printf("added %s %s\n", row[2], row[1])
	return nil
}

func list() error {
	rows, err := load()
	if err != nil {
		return err
	}
	if len(rows) == 0 {
		fmt.Println("no expenses yet")
	}
	for _, r := range rows {
		fmt.Printf("%s  %-12s %10s  %s\n", r[0], r[1], r[2], r[3])
	}
	return nil
}

func summary() error {
	rows, err := load()
	if err != nil {
		return err
	}
	totals := map[string]float64{}
	var grand float64
	for _, r := range rows {
		v, err := strconv.ParseFloat(r[2], 64)
		if err != nil {
			return fmt.Errorf("bad amount %q in %s", r[2], dbPath())
		}
		totals[r[1]] += v
		grand += v
	}
	cats := make([]string, 0, len(totals))
	for c := range totals {
		cats = append(cats, c)
	}
	sort.Strings(cats)
	for _, c := range cats {
		fmt.Printf("%-12s %10.2f\n", c, totals[c])
	}
	fmt.Printf("%-12s %10.2f\n", "total", grand)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: expense-tracker add|list|summary")
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "add":
		err = add(os.Args[2:])
	case "list":
		err = list()
	case "summary":
		err = summary()
	default:
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
