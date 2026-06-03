package txt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClip(t *testing.T) {
	t.Run("ASCII", func(t *testing.T) {
		assert.Equal(t, "ASCI", Clip("ASCII", 4))
	})
	t.Run("ShortEnough", func(t *testing.T) {
		assert.Equal(t, "I'm ä lazy BRoWN fox!", Clip("I'm ä lazy BRoWN fox!", 128))
	})
	t.Run("Clip", func(t *testing.T) {
		assert.Equal(t, "I'm ä", Clip("I'm ä lazy BRoWN fox!", 6))
		assert.Equal(t, "I'm ä l", Clip("I'm ä lazy BRoWN fox!", 7))
	})
	t.Run("TrimSpace", func(t *testing.T) {
		assert.Equal(t, "abc", Clip(" abc ty3q5y4y46uy", 4))
	})
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, "", Clip("", -1))
	})
}

func TestClipBytes(t *testing.T) {
	t.Run("ASCII", func(t *testing.T) {
		assert.Equal(t, "ASCI", ClipBytes("ASCII", 4))
	})
	t.Run("ShortEnough", func(t *testing.T) {
		assert.Equal(t, "I'm ä lazy fox!", ClipBytes("I'm ä lazy fox!", 128))
	})
	t.Run("MultiByteNotSplit", func(t *testing.T) {
		// "ä" is 2 bytes; a 5-byte budget over "I'm ä lazy" keeps "I'm" (4 bytes)
		// rather than half of "ä", and a 6-byte budget keeps the whole "I'm ä".
		assert.Equal(t, "I'm", ClipBytes("I'm ä lazy", 5))
		assert.Equal(t, "I'm ä", ClipBytes("I'm ä lazy", 6))
	})
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, "", ClipBytes("", 16))
	})
}

func BenchmarkClipRunesASCII(b *testing.B) {
	s := strings.Repeat("abc def ghi ", 20) // ASCII
	b.ReportAllocs()
	for b.Loop() {
		_ = Clip(s, 50)
	}
}

func BenchmarkClipRunesUTF8(b *testing.B) {
	s := strings.Repeat("Grüße 世", 20) // non-ASCII runes
	b.ReportAllocs()
	for b.Loop() {
		_ = Clip(s, 50)
	}
}

func TestShorten(t *testing.T) {
	t.Run("ShortEnough", func(t *testing.T) {
		assert.Equal(t, "fox!", Shorten("fox!", 6, "..."))
	})
	t.Run("CustomSuffix", func(t *testing.T) {
		assert.Equal(t, "I'm ä...", Shorten("I'm ä lazy BRoWN fox!", 8, "..."))
	})
	t.Run("DefaultSuffix", func(t *testing.T) {
		assert.Equal(t, "I'm…", Shorten("I'm ä lazy BRoWN fox!", 7, ""))
	})
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, "", Shorten("", -1, ""))
	})
}
