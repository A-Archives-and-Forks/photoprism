package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/photoprism/photoprism/pkg/media/video"
)

func TestExcludeDefault(t *testing.T) {
	t.Run("ContainsMagicYUVByDefault", func(t *testing.T) {
		assert.NotEmpty(t, DefaultExclude)
		def := video.NewFormats(DefaultExclude)
		assert.True(t, def.Contains(video.CodecMagicYUV))
	})

	t.Run("ConsistentWithExclude", func(t *testing.T) {
		// DefaultExclude is captured at package init from Exclude.String(),
		// so the canonical default must round-trip through NewFormats.
		assert.Equal(t, DefaultExclude, video.NewFormats(DefaultExclude).String())
	})
}

func TestExclude(t *testing.T) {
	// Snapshot and restore the package variable so tests don't leak state.
	saved := Exclude
	t.Cleanup(func() { Exclude = saved })

	t.Run("DefaultBlocksMagicYUV", func(t *testing.T) {
		Exclude = video.NewFormats(DefaultExclude)
		assert.True(t, Exclude.Contains(video.CodecMagicYUV))
		assert.False(t, Exclude.Contains(video.CodecAvc1))
	})

	t.Run("EmptyAllowsEverything", func(t *testing.T) {
		Exclude = video.NewFormats("")
		assert.False(t, Exclude.Contains(video.CodecMagicYUV))
	})

	t.Run("CustomList", func(t *testing.T) {
		Exclude = video.NewFormats("avi, hap")
		assert.True(t, Exclude.Contains("avi"))
		assert.True(t, Exclude.Contains("hap"))
		assert.False(t, Exclude.Contains(video.CodecMagicYUV))
	})
}
