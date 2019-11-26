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
	BFF         *Person
	Friends     []Person
	Extra       map[string]string
}

type AnotherPerson struct {
	Name   string
	Age    int
	Emails []string
	Other  map[string]string
}

type testCase struct {
	name          string
	instance      interface{}
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
			name:     "person: all normal values",
			instance: Person{},
			values: map[string]interface{}{
				"name":    "Bob",
				"age":     25,
				"emails":  []string{"bob@thebuilder.us", "bobby@notan.org"},
				"dob":     "1999-02-10",
				"bff":     &Person{Name: "Jane"},
				"friends": []Person{Person{Name: "John"}, Person{Name: "Doe"}},
				"extra": map[string]string{
					"gender": "Robot",
				},
			},
			existing: Person{Name: "Bobs"},
			result: Person{
				Name:        "Bob",
				Age:         25,
				Emails:      []string{"bob@thebuilder.us", "bobby@notan.org"},
				BFF:         &Person{Name: "Jane"},
				Friends:     []Person{Person{Name: "John"}, Person{Name: "Doe"}},
				DateOfBirth: "1999-02-10",
				Extra: map[string]string{
					"gender": "Robot",
				},
			},
			expectedError: false,
		},
		testCase{
			name:     "person: all normal + unknown values",
			instance: Person{},
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
			name:     "person: missing values",
			instance: Person{},
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
			name:     "person: override values",
			instance: Person{},
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
		testCase{
			name:     "person: field type mistmatch",
			instance: Person{},
			values: map[string]interface{}{
				"name":   "Bob",
				"age":    "25",
				"emails": "job@test.he",
			},
			existing: Person{Emails: []string{"bobby@oldemail.us"}},
			result: Person{
				Name:   "Bob",
				Emails: []string{"bobby@oldemail.us"},
			},
			expectedError: false,
		},
		testCase{
			name:     "person: instance + existing type mistmatch",
			instance: Person{},
			values: map[string]interface{}{
				"name":   "Bob",
				"age":    25,
				"emails": []string{"somebody@nobody.nb"},
			},
			existing: AnotherPerson{Emails: []string{"bobby@oldemail.us"}},
			result: Person{
				Name:   "Bob",
				Age:    25,
				Emails: []string{"somebody@nobody.nb"},
			},
			expectedError: false,
		},
	}
}

// TestUpdater is the main test function
func (s *UpdaterTestSuite) TestUpdater() {
	for _, tc := range s.testcases {
		s.Run(tc.name, func() {
			updaterFn, err := updater.New(tc.instance)
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

// BenchmarkNewUpdater benchmark the New "Updater" factory function
func BenchmarkNewUpdater(b *testing.B) {
	instance := Person{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		updater.New(instance)
	}
}

// BenchmarkUpdaterBasicTypes benchmark the "Updater" function
// with Basic types
func BenchmarkUpdaterBasicTypes(b *testing.B) {
	updaterFn, err := updater.New(Person{})
	if err != nil {
		b.Error(err)
	}

	existing := &Person{}

	values := map[string]interface{}{
		"name": "Bob",
		"age":  25,
		"dob":  "1999-02-10",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		updaterFn(existing, values)
	}
}

// BenchmarkUpdaterContainerTypes benchmark the "Updater" function
// with Container types
func BenchmarkUpdaterContainerTypes(b *testing.B) {
	updaterFn, err := updater.New(Person{})
	if err != nil {
		b.Error(err)
	}

	existing := &Person{}

	values := map[string]interface{}{
		"emails": []string{"bob@thebuilder.us", "bobby@notan.org"},
		"extra": map[string]string{
			"gender": "Robot",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		updaterFn(existing, values)
	}
}

// BenchmarkUpdaterCompositeTypes benchmark the "Updater" function
// with Composite types
func BenchmarkUpdaterCompositeTypes(b *testing.B) {
	updaterFn, err := updater.New(Person{})
	if err != nil {
		b.Error(err)
	}

	existing := &Person{}

	values := map[string]interface{}{
		"bff":     &Person{Name: "Jane"},
		"friends": []Person{Person{Name: "John"}, Person{Name: "Doe"}},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		updaterFn(existing, values)
	}
}

// BenchmarkUpdaterAllTypes benchmark the "Updater" function
// with All types
func BenchmarkUpdaterAllTypes(b *testing.B) {
	updaterFn, err := updater.New(Person{})
	if err != nil {
		b.Error(err)
	}

	existing := &Person{}

	values := map[string]interface{}{
		"name":    "Bob",
		"age":     25,
		"emails":  []string{"bob@thebuilder.us", "bobby@notan.org"},
		"dob":     "1999-02-10",
		"bff":     &Person{Name: "Jane"},
		"friends": []Person{Person{Name: "John"}, Person{Name: "Doe"}},
		"extra": map[string]string{
			"gender": "Robot",
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		updaterFn(existing, values)
	}
}
