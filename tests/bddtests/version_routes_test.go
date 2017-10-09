package bddtests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/corvinusz/for-swagger/server/version"
)

var _ = Describe("Test /version", func() {
	Context("GET /version", func() {
		It("should respond properly", func() {
			result := new(version.GetVersionBody)
			resp, err := suite.client.R().SetResult(result).Get("/version")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(200))
			Expect(result.Result).To(Equal("OK"))
			Expect(result.Version).To(Equal(suite.app.C.Config.Version))
			Expect(result.ServerTime).NotTo(BeZero())
		})
	})
})
