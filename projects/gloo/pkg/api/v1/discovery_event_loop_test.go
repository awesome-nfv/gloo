// Code generated by solo-kit. DO NOT EDIT.

// +build solokit

package v1

import (
	"context"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

var _ = Describe("DiscoveryEventLoop", func() {
	var (
		namespace string
		emitter   DiscoveryEmitter
		err       error
	)

	BeforeEach(func() {

		upstreamClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		upstreamClient, err := NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())

		secretClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		secretClient, err := NewSecretClient(secretClientFactory)
		Expect(err).NotTo(HaveOccurred())

		emitter = NewDiscoveryEmitter(upstreamClient, secretClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = emitter.Upstream().Write(NewUpstream(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = emitter.Secret().Write(NewSecret(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockDiscoverySyncer{}
		el := NewDiscoveryEventLoop(emitter, sync)
		_, err := el.Run([]string{namespace}, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(sync.Synced, 5*time.Second).Should(BeTrue())
	})
})

type mockDiscoverySyncer struct {
	synced bool
	mutex  sync.Mutex
}

func (s *mockDiscoverySyncer) Synced() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.synced
}

func (s *mockDiscoverySyncer) Sync(ctx context.Context, snap *DiscoverySnapshot) error {
	s.mutex.Lock()
	s.synced = true
	s.mutex.Unlock()
	return nil
}
