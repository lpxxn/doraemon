package internal

import (
	"fmt"
	"image"
	"log"

	"github.com/skip2/go-qrcode"
)

// RenderQRString as a QR code
func RenderQRString(s string) {
	q, err := qrcode.New(s, qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q.ToSmallString(false))
}

// RenderQRImage returns a QR code as an image.Image
func RenderQRImage(s string) image.Image {
	q, err := qrcode.New(s, qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}
	return q.Image(256)
}
