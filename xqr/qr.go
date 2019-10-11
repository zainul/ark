package xqr

import (
	"fmt"
	"path"

	qrcode "github.com/skip2/go-qrcode"
)

// Code is struct for creating / initializing qrcode.
type Code struct {
	Content   string // qrcode content
	Size      int    // qrcode size. Qrcode is always square
	Filename  string // for accomodate full file path qrcode we want to upload
	imgReview string // private field for accomodate link on first step
}

// Minimal size of qrcode is 300px according to Tokopedia upload tools
// Return error if size less than 300
func validateQrCodeSpec(qrCode *Code) error {
	if qrCode.Size < 300 {
		return fmt.Errorf("Code: Too small Code size")
	}

	return nil
}

// NewQR set the qr
func NewQR(content string, size int) Code {
	return Code{
		Content: content,
		Size:    size,
	}
}

// GenerateQrCodeImage is function to generate qrcode image.
// Filename is full filepath. Make sure path is available
// to write.
func (qrCode *Code) GenerateQrCodeImage(filename string) error {
	err := validateQrCodeSpec(qrCode)
	if err != nil {
		return err
	}

	// Clean filepath
	filename = path.Clean(filename)

	err = qrcode.WriteFile(qrCode.Content, qrcode.Medium, qrCode.Size, filename)
	if err != nil {
		return err
	}

	qrCode.Filename = filename
	return nil
}

// GenerateQrCodeImageByte is function to generate qrcode image resulting in byte array.
func (qrCode *Code) GenerateQrCodeImageByte() ([]byte, error) {
	err := validateQrCodeSpec(qrCode)
	if err != nil {
		return nil, err
	}

	png, err := qrcode.Encode(qrCode.Content, qrcode.Medium, qrCode.Size)
	if err != nil {
		return nil, err
	}

	return png, nil
}
