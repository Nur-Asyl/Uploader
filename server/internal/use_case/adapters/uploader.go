package adapters

import (
	"Manual_Parser/internal/domain/data_excel"
	"Manual_Parser/internal/domain/data_xml"
	"context"
	"github.com/xuri/excelize/v2"
)

type ExcelUploader interface {
	Upload(ctx context.Context, f *excelize.File, req data_excel.RequestExcel, dto map[string]string) (*data_excel.ResponseExcel, error)
}

type XMLUploader interface {
	Upload(ctx context.Context, dto map[string]data_xml.Tag, req data_xml.RequestXML, rows [][]data_xml.TagValue) (*data_xml.ResponseXML, error)
}
