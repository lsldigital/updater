package updater_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// UpdaterTestSuite is the main test suite for package updater
type UpdaterTestSuite struct {
	suite.Suite
}

// SetupTest setups the test suite for future tests
func (suite *UpdaterTestSuite) SetupTest() {
	// TODO: implement
	suite.Fail("not implemented yet")
}

// TestUpdater is the main test function
func (suite *UpdaterTestSuite) TestUpdater() {
	// TODO: implement
	suite.FailNow("not implemented yet")
}

// TestUpdaterTestSuite is the main entrypoint for UpdaterTestSuite
func TestUpdaterTestSuite(t *testing.T) {
	suite.Run(t, new(UpdaterTestSuite))
}
