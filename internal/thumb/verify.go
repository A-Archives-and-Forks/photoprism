package thumb

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

// Verify reports an error if the image file cannot be decoded by the active rendering library.
// Used to reject corrupt converter output (e.g. a truncated embedded RAW preview that passes a
// MIME sniff but later fails libvips) so the conversion loop can try the next converter.
func Verify(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("verify: empty filename")
	}

	// Use the loader the thumbnailer will use so the check matches GenerateThumbnails.
	if Library == LibVips {
		VipsInit()

		img, err := vips.LoadImageFromFile(fileName, VipsImportParams())
		if err != nil {
			return err
		}

		img.Close()

		return nil
	}

	_, err := Open(fileName, 1)

	return err
}
