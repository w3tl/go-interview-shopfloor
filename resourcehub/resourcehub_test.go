package resourcehub

import (
	"context"
	"fmt"
	"net/http"
	"shopfloor/client"
	"shopfloor/resource"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceHub(t *testing.T) {
	t.Parallel()

	h := NewServer()

	t.Run("Should add new resource", func(t *testing.T) {
		err := h.AddResource(resource.New("machine 1"))
		require.NoError(t, err)

		assert.Len(t, h.resources, 1)
	})

	t.Run("Should not add already existing resource", func(t *testing.T) {
		err := h.AddResource(resource.New("machine 1"))
		require.ErrorIs(t, err, ErrResourceExists)
	})

	t.Run("Should remove resource", func(t *testing.T) {
		err := h.RemoveResource("machine 1")
		require.NoError(t, err)

		assert.Len(t, h.resources, 0)
	})

	t.Run("Should return an error when delete not existing resource", func(t *testing.T) {
		rescount := len(h.resources)

		err := h.RemoveResource("resource test")
		require.ErrorIs(t, err, ErrResourceNotExists)

		assert.Equal(t, rescount, len(h.resources))
	})

	t.Run("Should get the resource", func(t *testing.T) {
		res := h.GetResource("machine 1")
		require.Nil(t, res)

		err := h.AddResource(resource.New("test resource"))
		require.NoError(t, err)

		res = h.GetResource("test resource")
		require.NotNil(t, res)
	})
}

var srv *http.Server

func StartServer(b *testing.B) {
	if srv == nil {
		hub := NewServer()
		err := hub.AddResource(resource.New("machine_1"))
		require.NoError(b, err)

		srv = hub.Listen(3030)
	}
}

func BenchmarkResourceHub(b *testing.B) {
	b.Cleanup(func() {
		if err := srv.Shutdown(context.Background()); err != nil {
			b.Fatal(err)
		}
	})

	StartServer(b)

	var wg sync.WaitGroup
	for idx := 0; idx < b.N; idx += 1 {
		idx := idx
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()
			err := client.Run(ctx, fmt.Sprintf("client %d", idx+1), "http://localhost:3030", "machine_1")
			if err != nil {
				b.Errorf("Client %d faced with an error", idx+1)
			}
		}()
	}

	wg.Wait()

	if err := srv.Shutdown(context.Background()); err != nil {
		b.Fatal(err)
	}
}
