package data_json

type RequestData struct {
	DBTable   string `json:"db_table"`
	SheetName string `json:"sheet_name"`
	FieldRow  int64  `json:"field_row"`
	DataRow   int64  `json:"data_row"`
	Fields    []struct {
		DB    string `json:"db"`
		Field string `json:"field"`
	} `json:"fields"`
}

//type requestData struct {
//	FieldRow int32 `json:"field_row"`
//	DataRow  int32 `json:"data_row"`
//	Fields   struct {
//		Bin               string     `json:"bin,omitempty"`
//		Rnn               string     `json:"rnn,omitempty"`
//		Name              string     `json:"name,omitempty"`
//		FioPayer          string     `json:"fio_payer,omitempty"`
//		FioDirector       string     `json:"fio_director,omitempty"`
//		InnDirector       string     `json:"iin_director,omitempty"`
//		RnnDirector       string     `json:"rnn_director,omitempty"`
//		NumberExamination string     `json:"number_examination,omitempty"`
//		Date              *time.Time `json:"date,omitempty"`
//	} `json:"fields"`
//}

//type requestData struct {
//	FieldRow int           `json:"field_row"`
//	DataRow  int           `json:"data_row"`
//	Fields   []interface{} `json:"fields"`
//}
