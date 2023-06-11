package main

import (
	"context"
	"log"

	"github.cm/jake-scott/exif-remover/exifstrip"
)

func main() {
	if err := exifstrip.StripExif(context.Background(), "abc", "photos-in.poptart.org", "16154941227038.png.webp"); err != nil {
		log.Println(err)
	}
}
