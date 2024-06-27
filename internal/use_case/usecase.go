package use_case

import (
	"Manual_Parser/internal/domain/data_excel"
	"Manual_Parser/internal/domain/data_xml"
	"context"
	"github.com/xuri/excelize/v2"
)

type ExcelUseCase interface {
	Upload(ctx context.Context, f *excelize.File, reqData data_excel.RequestExcel) error
}

type XMLUseCase interface {
	Upload(ctx context.Context, rootNode data_xml.Node, reqData data_xml.RequestXML) error
}
