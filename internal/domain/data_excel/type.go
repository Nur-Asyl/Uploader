package data_excel

type RequestExcel struct {
	DBTable   string `json:"db_table"`
	SheetName string `json:"sheet_name"`
	FieldRow  int64  `json:"field_row"`
	DataRow   int64  `json:"data_row"`
	Fields    []struct {
		DB    string `json:"db"`
		Field string `json:"field"`
	} `json:"fields"`
}
