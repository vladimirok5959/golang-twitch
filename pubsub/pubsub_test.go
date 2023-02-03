package pubsub_test

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/vladimirok5959/golang-twitch/pubsub"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PubSub", func() {
	var ctx = context.Background()

	Context("PubSub", func() {
		var ps *pubsub.PubSub

		BeforeEach(func() {
			ps = pubsub.NewWithURL(url.URL{Scheme: "ws", Host: "example.com", Path: ""})
		})

		AfterEach(func() {
			ps.Close()
		})

		Context("Listen", func() {
			It("create new connection for each 50 topics", func() {
				Expect(len(ps.Connections)).To(Equal(0))

				for i := 1; i <= 45; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 1, i)
				}
				Expect(len(ps.Connections)).To(Equal(1))

				for i := 1; i <= 5; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 1, i)
				}
				Expect(len(ps.Connections)).To(Equal(1))

				for i := 1; i <= 50; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 2, i)
				}
				Expect(len(ps.Connections)).To(Equal(2))

				for i := 1; i <= 50; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 3, i)
				}
				Expect(len(ps.Connections)).To(Equal(3))
			})
		})

		Context("Unlisten", func() {
			It("remove connection without topics", func() {
				Expect(len(ps.Connections)).To(Equal(0))

				for i := 1; i <= 50; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 1, i)
				}
				Expect(len(ps.Connections)).To(Equal(1))

				ps.Listen(ctx, "community-points-channel-v1", 2, 1)
				Expect(len(ps.Connections)).To(Equal(2))

				ps.Unlisten(ctx, "community-points-channel-v1", 2, 1)
				Expect(len(ps.Connections)).To(Equal(1))

				for i := 1; i <= 50; i++ {
					ps.Unlisten(ctx, "community-points-channel-v1", 1, i)
				}
				Expect(len(ps.Connections)).To(Equal(0))
			})
		})

		Context("Topics", func() {
			It("return topics", func() {
				Expect(len(ps.Connections)).To(Equal(0))

				for i := 1; i <= 50; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 1, i)
				}
				Expect(len(ps.Connections)).To(Equal(1))

				ps.Listen(ctx, "community-points-channel-v1", 2, 1)
				Expect(len(ps.Connections)).To(Equal(2))

				Expect(ps.Topics()).To(ContainElements(
					"community-points-channel-v1.2.1",
				))
			})
		})

		Context("HasTopic", func() {
			It("checks topics", func() {
				Expect(len(ps.Connections)).To(Equal(0))

				ps.Listen(ctx, "community-points-channel-v1", 1)
				Expect(ps.HasTopic("unknown")).To(BeFalse())
				Expect(ps.HasTopic("community-points-channel-v1", 1)).To(BeTrue())
			})
		})

		Context("TopicsCount", func() {
			It("return topics count", func() {
				Expect(ps.TopicsCount()).To(Equal(0))
				for i := 1; i <= 50; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 1, i)
				}
				Expect(ps.TopicsCount()).To(Equal(50))

				for i := 1; i <= 5; i++ {
					ps.Listen(ctx, "community-points-channel-v1", 2, i)
				}
				Expect(ps.TopicsCount()).To(Equal(55))
			})
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

	Context("Connection", func() {
		var c *pubsub.Connection

		BeforeEach(func() {
			c = pubsub.NewConnection(url.URL{Scheme: "ws", Host: "example.com", Path: ""})
		})

		AfterEach(func() {
			c.Close()
		})

		Context("AddTopic", func() {
			It("add new topic", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))

				c.AddTopic("community-points-channel-v1.1")
				Expect(c.TopicsCount()).To(Equal(1))
				Expect(c.Topics()).To(ContainElements(
					"community-points-channel-v1.1",
				))

				c.AddTopic("community-points-channel-v1.2")
				Expect(c.TopicsCount()).To(Equal(2))
				Expect(c.Topics()).To(ContainElements(
					"community-points-channel-v1.1",
					"community-points-channel-v1.2",
				))
			})

			It("not add the same topics", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))

				c.AddTopic("community-points-channel-v1.1")
				Expect(c.TopicsCount()).To(Equal(1))
				Expect(c.Topics()).To(ContainElements(
					"community-points-channel-v1.1",
				))

				c.AddTopic("community-points-channel-v1.1")
				Expect(c.TopicsCount()).To(Equal(1))
				Expect(c.Topics()).To(ContainElements(
					"community-points-channel-v1.1",
				))
			})

			It("add not more than 50 topics", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))

				for i := 1; i <= 60; i++ {
					c.AddTopic(fmt.Sprintf("community-points-channel-v1.%d", i))
				}

				Expect(c.TopicsCount()).To(Equal(50))
			})
		})

		Context("RemoveTopic", func() {
			It("remove topic", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))

				c.AddTopic("community-points-channel-v1.1")
				Expect(c.TopicsCount()).To(Equal(1))
				Expect(c.Topics()).To(ContainElements(
					"community-points-channel-v1.1",
				))

				c.RemoveTopic("community-points-channel-v1.1")
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))
			})
		})

		Context("RemoveAllTopics", func() {
			It("remove all topic", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))

				c.AddTopic("community-points-channel-v1.1")
				c.AddTopic("community-points-channel-v1.2")
				Expect(c.TopicsCount()).To(Equal(2))
				Expect(c.Topics()).To(ContainElements(
					"community-points-channel-v1.1",
					"community-points-channel-v1.2",
				))

				c.RemoveAllTopics()
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))
			})
		})

		Context("Topics", func() {
			It("return topics", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))

				c.AddTopic("community-points-channel-v1.1")
				Expect(c.TopicsCount()).To(Equal(1))
				Expect(c.Topics()).To(ContainElements(
					"community-points-channel-v1.1",
				))
			})
		})

		Context("HasTopic", func() {
			It("checks topics", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				Expect(c.Topics()).To(Equal([]string{}))

				c.AddTopic("community-points-channel-v1.1")
				Expect(c.HasTopic("unknown")).To(BeFalse())
				Expect(c.HasTopic("community-points-channel-v1.1")).To(BeTrue())
			})
		})

		Context("TopicsCount", func() {
			It("return topics count", func() {
				Expect(c.TopicsCount()).To(Equal(0))
				c.AddTopic("community-points-channel-v1.1")
				Expect(c.TopicsCount()).To(Equal(1))
				c.AddTopic("community-points-channel-v1.2")
				Expect(c.TopicsCount()).To(Equal(2))
			})
		})
	})
})

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PubSub")
}
