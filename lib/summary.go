package lib

import (
	"fmt"
	"github.com/jstemmer/go-junit-report/parser"
	"io"
	"strings"
)

type Results map[string][]*Test

const (
	PASS = "pass"
	FAIL = "fail"
	SKIP = "skip"
)

var jsonTestKeys = map[parser.Result]string{
	parser.PASS: PASS,
	parser.FAIL: FAIL,
	parser.SKIP: SKIP,
}

type Test struct {
	Name   string `json:"name"`
	Time   int    `json:"time"`
	Output string `json:"output"`
}

type TestSummary struct {
	TotalTests int     `json:"total_tests"`
	Results    Results `json:"results"`
}

func Parse(r io.Reader) (*TestSummary, error) {
	results := Results{
		PASS: []*Test{},
		FAIL: []*Test{},
		SKIP: []*Test{},
	}

	res, err := parser.Parse(r, "")
	if err != nil {
		return nil, err
	}

	totalTests := 0
	for _, pkg := range res.Packages {
		for _, t := range pkg.Tests {
			key, _ := jsonTestKeys[t.Result]

			className := classNameFromPackageName(pkg.Name)
			jsonTest := &Test{
				Name:   fmt.Sprintf("%s/%s", className, t.Name),
				Time:   t.Time,
				Output: strings.Join(t.Output, "\n"),
			}

			results[key] = append(results[key], jsonTest)
			totalTests += 1
		}
	}

	summary := &TestSummary{
		TotalTests: totalTests,
		Results:    results,
	}

	return summary, nil
}

func classNameFromPackageName(packageName string) string {
	className := packageName
	if idx := strings.LastIndex(className, "/"); idx > -1 && idx < len(packageName) {
		className = packageName[idx + 1:]
	}

	return className
}
