package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYouTubeID(t *testing.T) {
	valid := []string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://youtube.com/watch?v=dQw4w9WgXcQ&t=42s",
		"https://m.youtube.com/watch?list=PL123&v=dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ?si=abc",
		"https://www.youtube.com/shorts/dQw4w9WgXcQ",
		"https://www.youtube.com/embed/dQw4w9WgXcQ",
		"https://www.youtube.com/live/dQw4w9WgXcQ/",
		"https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ",
		"https://WWW.YouTube.com/watch?v=dQw4w9WgXcQ",
	}
	for _, link := range valid {
		id, ok := YouTubeID(link)
		assert.True(t, ok, link)
		assert.Equal(t, "dQw4w9WgXcQ", id, link)
		assert.Equal(t, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", canonicalVideoURL(link), link)
	}

	invalid := []string{
		"",
		"dQw4w9WgXcQ",
		"http://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://vimeo.com/123456789",
		"https://youtube.com.evil.com/watch?v=dQw4w9WgXcQ",
		"https://evil.com/?https://youtu.be/dQw4w9WgXcQ",
		"https://user:senha@www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://www.youtube.com:8443/watch?v=dQw4w9WgXcQ",
		"https://www.youtube.com/watch?v=curto",
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ\"><script>",
		"https://www.youtube.com/channel/dQw4w9WgXcQ",
		"https://www.youtube.com/",
		"javascript:alert(1)",
	}
	for _, link := range invalid {
		_, ok := YouTubeID(link)
		assert.False(t, ok, link)
		assert.Empty(t, canonicalVideoURL(link), link)
	}
}
