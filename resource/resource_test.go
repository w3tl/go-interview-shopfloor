package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResource(t *testing.T) {
	t.Parallel()

	t.Run("Should start resource", func(t *testing.T) {
		t.Parallel()

		r := New("machine 1")

		err := r.StartSetup()
		require.NoError(t, err)
		assert.Equal(t, ResourceStatusSetup, r.Status())
	})

	t.Run("Should start resource", func(t *testing.T) {
		t.Parallel()

		r := New("machine 1")

		err := r.StartSetup()
		require.NoError(t, err)

		err = r.Stop()
		require.NoError(t, err)
		assert.Equal(t, ResourceStatusStopped, r.Status())
	})
}
