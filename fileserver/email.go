package fileserver

import (
	"crypto/md5"
	"fmt"
	"image/jpeg"
	"log"
	"os"

	goimage "image"

	"github.com/Carbon-X-DAO/QRInvite/image"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// TOOD: send the email
func (server *Server) sendEmail(email string, id uint64) {
	code, err := generateQRCode(id)
	if err != nil {
		log.Printf("failed to generate QR code: %s", err)
		return
	}

	attachment := generateAttachment(server.flyer, code)

	// TODO: make this a buffer of bytes
	att, err := os.Create("./attachment.jpg")
	if err != nil {
		log.Printf("failed to create attachment file: %s", err)
		return
	}

	if err := jpeg.Encode(att, attachment, &jpeg.Options{Quality: jpeg.DefaultQuality}); err != nil {
		log.Printf("failed to encode attachment as JPEG: %s", err)
		return
	}
}

func generateQRCode(id uint64) (goimage.Image, error) {
	hash := md5.Sum([]byte(fmt.Sprintf("%d", id)))
	hashString := string(fmt.Sprintf("%x", hash))

	code, err := qr.Encode(hashString, qr.L, qr.Auto)
	if err != nil {
		return nil, fmt.Errorf("failed to encode id %d as QR code: %w", id, err)
	}

	intsize := 180
	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, intsize, intsize)
	if err != nil {
		return nil, fmt.Errorf("failed to scale QR code: %w", err)
	}

	return code, err
}
func generateAttachment(bottom, top goimage.Image) goimage.Image {
	return image.Layer(bottom, top, 587, -103)
}
