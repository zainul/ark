package xs3

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUploadS3(t *testing.T) {
	s3 := NewS3Upload("bucker", "/s3/main", "accessKey", "secretKey", "south-east-asia", 2)

	err := s3.UploadFile(time.Now().Format(time.RFC3339Nano), []byte("some file byte"))

	assert.NotEqual(t, nil, err)
}
