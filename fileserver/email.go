package fileserver

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"time"

	goimage "image"

	"github.com/Carbon-X-DAO/QRInvite/image"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

var body string = `
<html>
<body>
<h1>Gracias por participar.</h1>

	Te has ganado un premio.

	Que puedes reclamar
	<ol>
	<li> En la Carroza durante marcha. </li>
	<li> En la tarima de el evento después de la marcha en el parque luces. </li>
	</ol>

	Movimiento Cannabico Colombiano.
</body>
</html>
`

// TOOD: send the email
func (server *Server) sendEmail(email string, id uint64) error {
	code, err := generateQRCode(id)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %w", err)
	}

	attachment := generateAttachment(server.flyer, code)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, attachment, &jpeg.Options{Quality: 100}); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	subject := `Movimiento Cannabico Colombiano Premio 11 de Diciembre 2021`
	msg := server.mg.NewMessage("noreply@cieloverde.io", subject, "", email)
	msg.SetHtml(body)

	msg.AddReaderAttachment("boleto.jpg", ioutil.NopCloser(&buf))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, _, err = server.mg.Send(ctx, msg); err != nil {
		return fmt.Errorf("failed to send message to recipient: %w", err)
	}

	return nil
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
