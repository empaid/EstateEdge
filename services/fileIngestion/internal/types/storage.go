package types

import (
	"context"

	"github.com/empaid/estateedge/services/common/genproto/fileIngestion"
)

type Storage interface {
	ReturnPreSignedUploadURL(ctx context.Context, file *fileIngestion.File) (url string, err error)
}
