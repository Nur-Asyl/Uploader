package excel_uc

import (
	"Manual_Parser/internal/domain/data_json"
	"Manual_Parser/internal/use_case/adapters"
	"context"
	"github.com/xuri/excelize/v2"
	"log/slog"
)

type ExcelUseCase struct {
	uploader adapters.Uploader
}

func NewExcelUseCase(uploader adapters.Uploader) *ExcelUseCase {
	return &ExcelUseCase{
		uploader: uploader,
	}
}

func (uc *ExcelUseCase) Upload(ctx context.Context, f *excelize.File, reqData data_json.RequestData) error {
	slog.Info("Creating DTO")
	dto := make(map[string]string)
	for _, data := range reqData.Fields {
		dto[data.Field] = data.DB
	}
	slog.Info("Created DTO")
	if err := uc.uploader.Upload(ctx, f, reqData, dto); err != nil {
		return err
	}
	return nil
}
