package lib

import (
	"io"
	"strings"

	"github.com/improbable-io/go-junit-report/parser"

	"io/ioutil"
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
	PackageName string `json:"package_name"`
	TestName    string `json:"test_name"`
	Time        int    `json:"time"`
	Output      string `json:"output"`
}

type TestSummary struct {
	TotalTests  int     `json:"total_tests"`
	BuildErrors string  `json:"build_errors"`
	Results     Results `json:"results"`
}

func Parse(stdoutReader io.Reader, stderrReader io.Reader) (*TestSummary, error) {
	results := Results{
		PASS: []*Test{},
		FAIL: []*Test{},
		SKIP: []*Test{},
	}

	res, err := parser.Parse(stdoutReader, "")
	if err != nil {
		return nil, err
	}

	totalTests := 0
	for _, pkg := range res.Packages {
		for _, t := range pkg.Tests {
			key, _ := jsonTestKeys[t.Result]

			jsonTest := &Test{
				PackageName: pkg.Name,
				TestName:    t.Name,
				Time:        t.Time,
				Output:      strings.Join(t.Output, "\n"),
			}

			results[key] = append(results[key], jsonTest)
			totalTests += 1
		}
	}

	buildErrorBytes, err := ioutil.ReadAll(stderrReader)
	if err != nil {
		return nil, err
	}

	summary := &TestSummary{
		TotalTests:  totalTests,
		Results:     results,
		BuildErrors: string(buildErrorBytes),
	}

	return summary, nil
}
