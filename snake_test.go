package updater

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Not using "updater_test" because "toSnakeCase" is not exported

type testCase struct {
	name          string
	input         string
	output        string
	expectedError bool
}

// ToSnakeTestSuite is the snake test suite for package updater
type ToSnakeTestSuite struct {
	suite.Suite
	testcases []testCase
}

// SetupTest setups the test suite for future tests
func (s *ToSnakeTestSuite) SetupSuite() {
	s.testcases = []testCase{
		testCase{
			name:   "capitalized",
			input:  "MustPass",
			output: "must_pass",
		},
		testCase{
			name:   "all lowercase",
			input:  "mustpass",
			output: "mustpass",
		},
		testCase{
			name:   "capitalized with space",
			input:  "really MustPass",
			output: "really _must_pass",
		},
		testCase{
			name:   "capitalized with numbers",
			input:  "Must123Pass456",
			output: "must123_pass456",
		},
		testCase{
			name:   "capitalized with symbols",
			input:  "Must@Pass&",
			output: "must@_pass&",
		},
		testCase{
			name:   "already snake case",
			input:  "must_pass_surely",
			output: "must_pass_surely",
		},
	}
}

// TestToSnakeCase is the main test function
func (s *ToSnakeTestSuite) TestToSnakeCase() {
	for _, tc := range s.testcases {
		s.Run(tc.name, func() {
			output := toSnakeCase(tc.input)
			s.Equal(tc.output, output)
		})
	}
}

// TestToSnakeTestSuite is the main entrypoint for ToSnakeTestSuite
func TestToSnakeTestSuite(t *testing.T) {
	suite.Run(t, new(ToSnakeTestSuite))
}

// BenchmarkToSnakeCase benchmark the "toSnakeCase"  function
func BenchmarkToSnakeCase(b *testing.B) {
	input := "SomeTextToSnake@Case123"

	for i := 0; i < b.N; i++ {
		toSnakeCase(input)
	}
}
