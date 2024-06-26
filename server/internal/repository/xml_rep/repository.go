package xml_rep

import (
	"Manual_Parser/internal/domain"
	"Manual_Parser/internal/domain/data_xml"
	"Manual_Parser/pkg/utils"
	"context"
	"database/sql"
	"log/slog"
	"strconv"
	"strings"
)

type XMLRepo struct {
	db *sql.DB
}

func NewXMLRepo(db *sql.DB) *XMLRepo {
	return &XMLRepo{
		db: db,
	}
}

func (r *XMLRepo) Upload(ctx context.Context, dto map[string]data_xml.Tag, req data_xml.RequestXML, rows [][]data_xml.TagValue) (*data_xml.ResponseXML, error) {
	inserted := int64(0)
	total := int64(len(rows))
	queryFields := make([]string, 0)
	queryValues := make([]string, 0)
	queryParams := make([]interface{}, 0)
	valueCount := 1
	failedRows := ""

	slog.Info("Start Uploading")
	for i, row := range rows {
		for _, tag := range row {
			queryField, ok := dto[utils.SaltKey(tag.Tag, tag.Parent)]
			if ok {
				queryFields = append(queryFields, queryField.DB)
				queryValues = append(queryValues, "$"+strconv.Itoa(valueCount))
				valueCount++
				queryParams = append(queryParams, tag.Value)
			}
		}
		if len(queryFields) == 0 {
			slog.Error("Failed to find appropriate tags", "index", i, "row", row)
			return nil, domain.ErrTagsNotFound
		}

		queryString := "INSERT INTO " + req.DBTable + " (" + strings.Join(queryFields, ", ") + ") " + "VALUES (" + strings.Join(queryValues, ", ") + ")"
		_, err := r.db.ExecContext(ctx, queryString, queryParams...)
		if err != nil {
			slog.Error("Failed:", "row", i+1)
			slog.Error("Failed to execute query", "query", queryString, "query params", queryParams, "error", err)
			failedRows += "row: " + strconv.Itoa(i+1) + "\n" + err.Error() + "\n"
			if !domain.IsDataTypeError(err.Error()) {
				return nil, err
			}
		} else {
			inserted++
		}

		queryFields = nil
		queryValues = nil
		queryParams = nil
		valueCount = 1
	}

	slog.Info("Successfully Uploaded!!!", "Total", total, "Inserted", inserted)
	responseXML := &data_xml.ResponseXML{
		Total:      total,
		Inserted:   inserted,
		FailedRows: failedRows,
	}
	return responseXML, nil
}
