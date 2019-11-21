package updater_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.lsl.digital/updater"
)

type Person struct {
	Name        string
	Age         int
	Emails      []string
	DateOfBirth string `json:"dob"`
	Extra       map[string]string
}

type testCase struct {
	name          string
	element       interface{}
	values        map[string]interface{}
	existing      interface{}
	result        interface{}
	expectedError bool
}

// UpdaterTestSuite is the main test suite for package updater
type UpdaterTestSuite struct {
	suite.Suite
	testcases []testCase
}

// SetupTest setups the test suite for future tests
func (s *UpdaterTestSuite) SetupSuite() {
	s.testcases = []testCase{
		testCase{
			name:    "person: all normal values",
			element: Person{},
			values: map[string]interface{}{
				"name":   "Bob",
				"age":    25,
				"emails": []string{"bob@thebuilder.us", "bobby@notan.org"},
				"dob":    "1999-02-10",
				"extra": map[string]string{
					"gender": "Robot",
				},
			},
			existing: Person{},
			result: Person{
				Name:        "Bob",
				Age:         25,
				Emails:      []string{"bob@thebuilder.us", "bobby@notan.org"},
				DateOfBirth: "1999-02-10",
				Extra: map[string]string{
					"gender": "Robot",
				},
			},
			expectedError: false,
		},
		testCase{
			name:    "person: all normal + unknown values",
			element: Person{},
			values: map[string]interface{}{
				"name":   "Bob",
				"age":    25,
				"emails": []string{"bob@thebuilder.us", "bobby@notan.org"},
				"extra": map[string]string{
					"gender": "Class",
				},
				"invalid": true,
			},
			existing: Person{},
			result: Person{
				Name:   "Bob",
				Age:    25,
				Emails: []string{"bob@thebuilder.us", "bobby@notan.org"},
				Extra: map[string]string{
					"gender": "Class",
				},
			},
			expectedError: false,
		},
		testCase{
			name:    "person: missing values",
			element: Person{},
			values: map[string]interface{}{
				"name": "Bob",
				"age":  25,
				"extra": map[string]string{
					"gender": "less",
				},
			},
			existing: Person{Emails: []string{"bobby@oldemail.us"}},
			result: Person{
				Name:   "Bob",
				Age:    25,
				Emails: []string{"bobby@oldemail.us"},
				Extra: map[string]string{
					"gender": "less",
				},
			},
			expectedError: false,
		},
		testCase{
			name:    "person: override values",
			element: Person{},
			values: map[string]interface{}{
				"name":   "Bob",
				"age":    25,
				"emails": []string{"no-reply@lebobby.fr"},
				"extra": map[string]string{
					"gender": "fox",
				},
			},
			existing: Person{Emails: []string{"bobby@oldemail.us"}},
			result: Person{
				Name:   "Bob",
				Age:    25,
				Emails: []string{"no-reply@lebobby.fr"},
				Extra: map[string]string{
					"gender": "fox",
				},
			},
			expectedError: false,
		},
	}
}

// TestUpdater is the main test function
func (s *UpdaterTestSuite) TestUpdater() {
	for _, tc := range s.testcases {
		s.Run(tc.name, func() {
			updaterFn, err := updater.New(tc.element)
			if err != nil && !tc.expectedError {
				s.FailNow("updater.New: %v", err)
			}

			result, err := updaterFn(tc.existing, tc.values)
			if err != nil && !tc.expectedError {
				s.FailNow("updaterFn: %v", err)
			}

			s.Equal(tc.result, result)
		})
	}
}

// TestUpdaterTestSuite is the main entrypoint for UpdaterTestSuite
func TestUpdaterTestSuite(t *testing.T) {
	suite.Run(t, new(UpdaterTestSuite))
}
