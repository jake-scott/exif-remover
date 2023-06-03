# Photo EXIF remover

Simple Golang function to be run in GCP as an EventARC handler
configured to trigger on an object write in a bucket.

Requires the OUTPUT_BUCKET environment variable pointing to the
name of the destination bucket.

The service account the function runs as requires object-read
IAM permissions on the input bucket and object-write perms on
the output bucket.

This code takes a very much quick and dirty approach to the job
at hand.  It reads and decodes the entire image into memory as
a Golang generic image and re-encodes it as it writes it back.
This is obviuosly very inefficient -- it would be better not
to decode/encode and just strip the Exif data.. but this is more
work as that would involve writing per image-type code.