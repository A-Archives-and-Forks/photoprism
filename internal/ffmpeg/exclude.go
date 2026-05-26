package ffmpeg

import (
	"github.com/photoprism/photoprism/pkg/media/video"
)

// Exclude contains the video container and codec formats that should not be processed by FFmpeg.
var Exclude = video.NewFormats(video.CodecMagicYUV)

// DefaultExclude contains the video container and codec formats that should not be processed by default.
var DefaultExclude = Exclude.String()
