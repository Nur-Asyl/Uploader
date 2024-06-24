package use_case

import (
	"Manual_Parser/internal/domain/data_json"
	"context"
	"github.com/xuri/excelize/v2"
)

type UploadUseCase interface {
	Upload(ctx context.Context, f *excelize.File, reqData data_json.RequestData) error
}
