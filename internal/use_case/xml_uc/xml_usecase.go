package xml_uc

import (
	"Manual_Parser/internal/domain"
	"Manual_Parser/internal/domain/data_xml"
	"Manual_Parser/internal/use_case/adapters"
	"Manual_Parser/pkg/utils"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type XMLUseCase struct {
	uploader adapters.XMLUploader
}

func NewXMLUseCase(uploader adapters.XMLUploader) *XMLUseCase {
	return &XMLUseCase{uploader: uploader}
}

func (uc *XMLUseCase) Upload(ctx context.Context, rootNode data_xml.Node, req data_xml.RequestXML) error {
	dto := make(map[string]data_xml.Tag)
	for _, tag := range req.Tags {
		key := utils.SaltKey(tag.Tag, tag.Parent)
		_, ok := dto[key]
		if ok {
			slog.Error("Already exists such tag", "tag", tag)
			return errors.New(fmt.Sprintf("Already exists such tag: %+v", tag.Tag))
		} else {
			dto[key] = data_xml.Tag{
				DB:     tag.DB,
				Parent: tag.Parent,
				Tag:    tag.Tag,
			}
		}

	}
	row := make([]data_xml.TagValue, 0)
	rows := make([][]data_xml.TagValue, 0)

	slog.Info("Start getting rows recursion", "DTO", dto)
	uc.getRows(rootNode, dto, &row, &rows)

	if len(rows) == 0 {
		return domain.ErrZeroRows
	}
	if err := uc.uploader.Upload(ctx, dto, req, rows); err != nil {
		return err
	}
	return nil
}

// recursively iterates through xml and appends to the rows the data that appropriate for request's DTO
// Node - current node with field of its name and children, attributes, and tag's value text
// DTO - map[request_tag]db_field
// row - slice of each tag value appropriate to request's DTO
// rows - slice of row
func (uc *XMLUseCase) getRows(node data_xml.Node, dto map[string]data_xml.Tag, row *[]data_xml.TagValue, rows *[][]data_xml.TagValue) {
	// iterate each xml tag for its children tags
	for _, child := range node.Children {
		// check each xml tag if it has children tags
		if len(child.Children) != 0 {
			// Enter child until we meet end child
			uc.getRows(child, dto, row, rows)
			// after leaving the children's we need to check if it was the node that we wanted to get
			_, ok := dto[utils.SaltKey(child.XMLName.Local, node.XMLName.Local)]
			if !ok {
				// check if there was no data appropriate to request's DTO
				if len(*row) != 0 {
					*rows = append(*rows, *row)
					*row = nil
				}
			}
		} else {
			// Add as a row appropriate to our request's DTO
			v, ok := dto[utils.SaltKey(child.XMLName.Local, node.XMLName.Local)]
			if ok {
				// check if parent match to our request's DTO
				if v.Parent == node.XMLName.Local {
					*row = append(*row, data_xml.TagValue{Parent: v.Parent, DB: v.DB, Tag: v.Tag, Value: child.Text})
				}

			}
		}
	}
	// leaving this hell
}
