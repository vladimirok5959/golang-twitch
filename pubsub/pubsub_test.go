package pubsub_test

import (
	"net/url"
	"testing"

	"github.com/vladimirok5959/golang-twitch/pubsub"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PubSub", func() {
	var ps *pubsub.PubSub

	BeforeEach(func() {
		ps = pubsub.NewWithURL(url.URL{Scheme: "ws", Host: "example.com", Path: ""})
	})

	Context("Topic", func() {
		It("generate correct topic", func() {
			Expect(ps.Topic("channel-bits-events-v1.123")).To(Equal("channel-bits-events-v1.123"))
			Expect(ps.Topic("channel-bits-events-v1", 123)).To(Equal("channel-bits-events-v1.123"))
			Expect(ps.Topic("channel-bits-events-v1", "123")).To(Equal("channel-bits-events-v1.123"))
			Expect(ps.Topic("channel-bits-events-v1", 123, 456)).To(Equal("channel-bits-events-v1.123.456"))
			Expect(ps.Topic("channel-bits-events-v1", 123, "456")).To(Equal("channel-bits-events-v1.123.456"))
			Expect(ps.Topic("channel-bits-events-v1", "123", 456)).To(Equal("channel-bits-events-v1.123.456"))
			Expect(ps.Topic("channel-bits-events-v1", "123", "456")).To(Equal("channel-bits-events-v1.123.456"))
			Expect(ps.Topic("channel-bits-events-v1", 123, 456, 789)).To(Equal("channel-bits-events-v1.123.456.789"))
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PubSub")
}
