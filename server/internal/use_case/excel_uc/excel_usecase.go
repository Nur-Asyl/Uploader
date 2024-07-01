package excel_uc

import (
	"Manual_Parser/internal/domain/data_excel"
	"Manual_Parser/internal/use_case/adapters"
	"context"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log/slog"
)

type ExcelUseCase struct {
	uploader adapters.ExcelUploader
}

func NewExcelUseCase(uploader adapters.ExcelUploader) *ExcelUseCase {
	return &ExcelUseCase{
		uploader: uploader,
	}
}

func (uc *ExcelUseCase) Upload(ctx context.Context, f *excelize.File, req data_excel.RequestExcel) (*data_excel.ResponseExcel, error) {
	dto := make(map[string]string)
	for _, data := range req.Fields {
		_, ok := dto[data.Field]
		if ok {
			slog.Error("Already exists such field", "field", data.Field)
			return nil, errors.New(fmt.Sprintf("Already exists such field: %+v", data.Field))
		} else {
			dto[data.Field] = data.DB
		}
	}
	res, err := uc.uploader.Upload(ctx, f, req, dto)
	if err != nil {
		return nil, err
	}
	return res, nil
}
