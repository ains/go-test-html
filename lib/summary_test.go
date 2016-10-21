package lib

import (
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	stdOutFile string
	stdErrFile string
	summary    *TestSummary
}

var testCases = []TestCase{
	{
		stdOutFile: "01-stdout",
		stdErrFile: "01-stderr",
		summary: &TestSummary{
			TotalTests: 3,
			BuildErrors: "",
			Results: Results{
				PASS: []*Test{
					{
						PackageName:"cobra/doc",
						TestName:"TestBashCompletions",
						Time: 10,
						Output: "",
					},
					{
						PackageName:"cobra/doc",
						TestName:"TestBashCompletionHiddenFlag",
						Time: 0,
						Output: "",
					},
					{
						PackageName:"cobra/doc",
						TestName:"TestBashCompletionDeprecatedFlag",
						Time: 0,
						Output: "",
					},
				},
				FAIL: []*Test{},
				SKIP: []*Test{},
			},
		},
	},
	{
		stdOutFile: "02-stdout",
		stdErrFile: "02-stderr",
		summary: &TestSummary{
			TotalTests: 3,
			BuildErrors: "",
			Results: Results{
				PASS: []*Test{
					{
						PackageName:"cobra/doc",
						TestName:"TestGenManNoGenTag",
						Time: 0,
						Output: "",
					},
				},
				FAIL: []*Test{
					{
						PackageName:"cobra/doc",
						TestName:"TestGenManDoc",
						Time: 0,
						Output: "cmd_test.go:144: Line: 59 Unexpected response.",
					},
				},
				SKIP: []*Test{
					{
						PackageName:"cobra/doc",
						TestName:"TestGenMdNoTag",
						Time: 0,
						Output: "md_docs_test.go:72: yolo",
					},
				},
			},
		},
	},
	{
		stdOutFile: "03-stdout",
		stdErrFile: "03-stderr",
		summary: &TestSummary{
			TotalTests: 1,
			BuildErrors: "build-error",
			Results: Results{
				PASS: []*Test{},
				FAIL: []*Test{
					{
						PackageName:"cobra/doc",
						TestName:"TestGenManDoc",
						Time: 0,
						Output: "cmd_test.go:144: Line: 59 Unexpected response.",
					},
				},
				SKIP: []*Test{},
			},
		},
	},
}

func TestSummaryParser(t *testing.T) {
	for _, testCase := range testCases {
		t.Logf("Running: (%s, %s)", testCase.stdOutFile, testCase.stdErrFile)

		stdoutFile, err := os.Open("tests/" + testCase.stdOutFile)
		require.NoError(t, err)

		stdoutErr, err := os.Open("tests/" + testCase.stdErrFile)
		require.NoError(t, err)

		actual, err := Parse(stdoutFile, stdoutErr)
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.EqualValues(t, testCase.summary, actual)
	}
}