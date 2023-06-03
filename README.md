# Photo EXIF remover

Simple Golang function to be run in GCP as an EventARC handler
configured to trigger on an object write in a bucket.

Requires the OUTPUT_BUCKET environment variable pointing to the
name of the destination bucket.

The service account the function runs as requires object-read
IAM permissions on the input bucket and object-write perms on
the output bucket.

