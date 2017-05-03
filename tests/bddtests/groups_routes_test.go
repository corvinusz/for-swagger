package bddtests

import (
	"math/rand"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/corvinusz/for-swagger/server/groups"
)

var _ = Describe("Test GET /groups", func() {
	Context("Get all groups", func() {
		It("should respond properly", func() {
			var orig, result []groups.Entity
			// get orig
			err := suite.app.Context.Orm.Cols("id").Asc("id").Find(&orig)
			Expect(err).NotTo(HaveOccurred())
			for i := range orig {
				_, err = orig[i].ExtractFrom(suite.app.Context.Orm)
				Expect(err).NotTo(HaveOccurred())
			}
			// get resp
			resp, err := suite.client.R().SetResult(&result).Get("/groups")
			Expect(err).NotTo(HaveOccurred())
			Expect(http.StatusOK).To(Equal(resp.StatusCode()))
			Expect(len(orig)).To(BeNumerically(">=", 4))
			Expect(len(result)).To(Equal(len(orig)))
			Expect(result).To(BeEquivalentTo(orig))
		})
	})
})

var _ = Describe("Test GET /groups?id=", func() {
	Context("with 3 random id", func() {
		It("should respond properly", func() {
			properIds := []uint64{1, 4, 5, 10}
			var err error
			for i := 0; i < 3; i++ {
				id := rand.Int() % 4
				orig := new(groups.Entity)
				result := []groups.Entity{}
				// get orig
				orig.ID = properIds[id]
				_, err = orig.ExtractFrom(suite.app.Context.Orm)
				Expect(err).NotTo(HaveOccurred())
				// get resp
				resp, err := suite.client.R().SetResult(&result).Get("/groups?id=" + strconv.Itoa(int(properIds[id])))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(result)).To(Equal(1))
				Expect(http.StatusOK).To(Equal(resp.StatusCode()))
				Expect(&result[0]).To(BeEquivalentTo(orig))
			}
		})
	})
})
