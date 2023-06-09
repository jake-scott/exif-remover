package exifstrip

//package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"

	"image/gif"
	"image/jpeg"
	"image/png"

	_ "golang.org/x/image/webp"
)

var config = struct {
	outputBucket string
}{}

func init() {
	doInit()
	functions.CloudEvent("StripExif", stripExifEvent)
}

// StorageObjectData contains metadata of the Cloud Storage object.
type StorageObjectData struct {
	Bucket         string    `json:"bucket,omitempty"`
	Name           string    `json:"name,omitempty"`
	Metageneration int64     `json:"metageneration,string,omitempty"`
	TimeCreated    time.Time `json:"timeCreated,omitempty"`
	Updated        time.Time `json:"updated,omitempty"`
}

func doInit() {
	if env, ok := os.LookupEnv("OUTPUT_BUCKET"); ok {
		config.outputBucket = env
	} else {
		panic("Could not find OUTPUT_BUCKET env var")
	}

	log.Printf("Using config: %+v", config)
}

func idLogger(id string) func(format string, v ...any) {
	return func(format string, v ...any) {
		format = "[evt %s] " + format
		args := []any{id}
		args = append(args, v...)
		log.Printf(format, args...)
	}
}

func stripExifEvent(ctx context.Context, e event.Event) error {
	id := e.ID()
	log := idLogger(id)
	log("Processing new event (type: %s)", e.Type())

	var data StorageObjectData
	if err := e.DataAs(&data); err != nil {
		return fmt.Errorf("[%s] BAD EVENT: %e", id, err)
	}

	log("Bucket: %s, Object: %s", data.Bucket, data.Name)
	return StripExif(ctx, id, data.Bucket, data.Name)
}

func StripExif(ctx context.Context, id, inBucket, object string) error {
	log := idLogger(id)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		return fmt.Errorf("[%s] GCS error: %w", id, err)
	}

	rc, err := client.Bucket(inBucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("[%s] opening object for reading: %w", id, err)
	}
	defer rc.Close()
	m, _, err := image.Decode(rc)
	if err != nil {
		return fmt.Errorf("[%s] decoding object: %w", id, err)
	}

	fileType := strings.ToLower(filepath.Ext(object))
	outName := object

	if fileType == ".webp" {
		outName = object[0:len(object)-len(fileType)] + ".jpg"
	}

	wc := client.Bucket(config.outputBucket).Object(outName).NewWriter(ctx)
	defer wc.Close()

	switch fileType {
	case ".jpeg", ".jpg":
		log("Writing JPEG data")
		err = jpeg.Encode(wc, m, nil)
	case ".gif":
		log("Writing GIF data")
		err = gif.Encode(wc, m, nil)
	case ".png":
		log("Writing PNG data")
		err = png.Encode(wc, m)
	case ".webp":
		log("Writing WEBP as JPG data %s", outName)
		err = jpeg.Encode(wc, m, nil)
	default:
		err = fmt.Errorf("Unsupported file type %s", fileType)
	}

	if err != nil {
		return fmt.Errorf("[%s] opening object for writing: %w", id, err)
	}

	return nil
}
