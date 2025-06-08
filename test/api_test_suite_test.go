//nolint:revive // must use this style
package api_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = BeforeSuite(func() {
	// parsing env

})

var _ = AfterSuite(func() {

})

func TestApiTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Workflow Suite Test")
}
