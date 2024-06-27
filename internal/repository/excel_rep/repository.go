package excel_rep

import (
	"Manual_Parser/internal/domain"
	"Manual_Parser/internal/domain/data_excel"
	"context"
	"database/sql"
	"github.com/xuri/excelize/v2"
	"log/slog"
	"strconv"
	"strings"
)

type ExcelRepo struct {
	db *sql.DB
}

func NewExcelRepo(db *sql.DB) *ExcelRepo {
	return &ExcelRepo{db: db}
}

func (r *ExcelRepo) Upload(ctx context.Context, f *excelize.File, reqData data_excel.RequestExcel, dto map[string]string) error {
	slog.Info("Receiving sheet name")
	sheetIndex, err := f.GetSheetIndex(reqData.SheetName)
	if err != nil {
		slog.Error("Failed to get Sheet Name")
		return err
	}
	slog.Info("Sheet", "index", sheetIndex, "name", reqData.SheetName)
	slog.Info("Receiving rows")
	rows, err := f.GetRows(f.GetSheetName(sheetIndex))
	if err != nil {
		slog.Error("Error reading Excel sheet", "error", err)
		return err
	}

	fieldRow := rows[reqData.FieldRow-1]
	slog.Info("fields", "row", fieldRow)

	queryFields := make([]string, 0)
	queryValues := make([]string, 0)
	queryParams := make([]interface{}, 0)
	valueCount := 1

	slog.Info("Start Uploading")
	for i, row := range rows[reqData.DataRow-1:] {
		for j, cell := range row {
			queryField, ok := dto[fieldRow[j]]
			if ok {
				queryFields = append(queryFields, queryField)
				queryValues = append(queryValues, "$"+strconv.Itoa(valueCount))
				valueCount++
				queryParams = append(queryParams, cell)
			}

		}

		if len(queryFields) == 0 {
			return domain.ErrNoFields
		}

		queryString := "INSERT INTO " + reqData.DBTable + " (" + strings.Join(queryFields, ", ") + ") " + "VALUES (" + strings.Join(queryValues, ", ") + ")"
		_, err := r.db.ExecContext(ctx, queryString, queryParams...)
		if err != nil {
			slog.Error("Failed to execute query", "query", queryString, "params", queryParams)
			return err
		}

		queryFields = nil
		queryValues = nil
		queryParams = nil
		valueCount = 1

		slog.Info("INSERTED:", "row", i)
	}

	slog.Info("Successfully Uploaded!!!", "total rows in file", len(rows))
	return nil
}
