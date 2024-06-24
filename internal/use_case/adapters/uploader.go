package adapters

import (
	"Manual_Parser/internal/domain/data_json"
	"context"
	"github.com/xuri/excelize/v2"
)

type Uploader interface {
	Upload(ctx context.Context, f *excelize.File, req data_json.RequestData, dto map[string]string) error
}
