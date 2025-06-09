//nolint:revive // must use this style
package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	// TODO: Setup you testing environment
	// exmaple: Connect to db and setup router etc.

})

var _ = AfterSuite(func() {
	// TODO: remove all testing data...
})

func TestApiTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Workflow Suite Test")
}
