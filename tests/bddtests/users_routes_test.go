package bddtests

import (
	"math/rand"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/corvinusz/for-swagger/server/users"
)

var _ = Describe("Test GET /users", func() {
	Context("GET all /users", func() {
		It("should respond properly", func() {
			var orig, result []users.Entity
			// get orig
			err := suite.app.C.Orm.Cols("id").Asc("id").Find(&orig)
			Expect(err).NotTo(HaveOccurred())
			for i := range orig {
				err = orig[i].ExtractFrom(suite.app.C.Orm)
				Expect(err).NotTo(HaveOccurred())
				orig[i].Password = ""
			}
			Expect(len(orig)).To(BeNumerically(">=", 8))
			// get resp
			resp, err := suite.client.R().SetResult(&result).Get("/users")
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode()).To(Equal(http.StatusOK))
			// check response with DB data
			Expect(len(result)).To(Equal(len(orig)))
			Expect(result).To(BeEquivalentTo(orig))
		})
	})
	Context("GET one /users/{id}", func() {
		It("should respond properly", func() {
			var err error
			for i := 0; i < 3; i++ {
				id := rand.Int()%7 + 1
				orig := new(users.Entity)
				result := []users.Entity{}
				// get orig
				orig.ID = uint64(id)
				err = orig.ExtractFrom(suite.app.C.Orm)
				Expect(err).NotTo(HaveOccurred())
				orig.Password = ""
				// get resp
				resp, err := suite.client.R().SetResult(&result).Get("/users?id=" + strconv.Itoa(id))
				Expect(err).NotTo(HaveOccurred())
				// check response
				Expect(len(result)).To(Equal(1))
				Expect(result[0].ID).To(Equal(uint64(id)))
				Expect(resp.StatusCode()).To(Equal(http.StatusOK))
				// check reponse data
				Expect(&result[0]).To(BeEquivalentTo(orig))
			}
		})
	})
})

var _ = Describe("Test POST /users", func() {
	Context("POST /users", func() {
		It("should respond properly", func() {
			result := new(users.Entity)
			payload := users.UserInput{
				Login:    "a_test_user_01",
				Password: "a_test_user_01",
				GroupID:  10,
			}
			// make http request
			resp, err := suite.client.R().SetBody(payload).SetResult(result).Post("/users")
			Expect(err).NotTo(HaveOccurred())
			// check response
			Expect(resp.StatusCode()).To(Equal(http.StatusCreated))
			Expect(result.ID).NotTo(BeZero())
			Expect(result.Login).To(Equal(payload.Login))
			Expect(result.Created).NotTo(BeZero())
			Expect(result.Updated).NotTo(BeZero())
			// get original user
			orig := new(users.Entity)
			orig.ID = result.ID
			err = orig.ExtractFrom(suite.app.C.Orm)
			Expect(err).NotTo(HaveOccurred())
			orig.Password = "" // { Password string `json:"-"` }
			// check user data
			Expect(result).To(BeEquivalentTo(orig))
		})
	})
})
