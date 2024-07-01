package http_v1

import (
	"Manual_Parser/configs"
	"Manual_Parser/internal/domain"
	"Manual_Parser/internal/domain/data_excel"
	"Manual_Parser/internal/domain/data_xml"
	"Manual_Parser/internal/middleware/cors_mw"
	"Manual_Parser/internal/use_case"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	FRONTEND_FORM_JSON = "json_data"
	FRONTEND_FORM_FILE = "file"
)

type UploadHTTPDelivery struct {
	excelUC use_case.ExcelUseCase
	xmlUC   use_case.XMLUseCase
}

func NewUploadHTTPDelivery(excelUC use_case.ExcelUseCase, xmlUC use_case.XMLUseCase) *UploadHTTPDelivery {
	return &UploadHTTPDelivery{
		excelUC: excelUC,
		xmlUC:   xmlUC,
	}
}

func (d *UploadHTTPDelivery) UploadExcelFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*5)
	defer cancel()

	if r.Method != http.MethodPost {
		slog.Error("Invalid request method")
		http.Error(w, domain.ErrInvalidRequestMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(1 << 30)
	if err != nil {
		slog.Error("Failed to parse request form excel_rep", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData := r.FormValue(FRONTEND_FORM_JSON)

	var requestData data_excel.RequestExcel
	err = json.NewDecoder(strings.NewReader(jsonData)).Decode(&requestData)
	if err != nil {
		slog.Error("Failed to decode request excel_json", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formFile, fheader, err := r.FormFile(FRONTEND_FORM_FILE)
	if err != nil {
		slog.Error("Error retrieving file", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	slog.Info("Creating temp file")
	tempFile, err := os.CreateTemp("temp", "excel-*.xlsx")
	if err != nil {
		slog.Error("Error creating temporary file", "error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name()) // to remove temp file that was saved in temp directory
	defer tempFile.Close()

	slog.Info("Copying file to temporary")
	_, err = io.Copy(tempFile, formFile)
	if err != nil {
		slog.Error("Error saving file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Opening file")
	file, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		slog.Error("Error opening excel file", "error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	res, err := d.excelUC.Upload(ctx, file, requestData)
	if err != nil {
		slog.Error("Failed to upload file", "error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("File uploaded and read successfully", "File", fheader.Filename, "Size", fheader.Size, "Header", fheader.Header)

	err = json.NewEncoder(w).Encode(*res)
	if err != nil {
		slog.Error("Failed to encode excel response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (d *UploadHTTPDelivery) UploadXMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*5)
	defer cancel()

	if r.Method != http.MethodPost {
		slog.Error("Invalid request method")
		http.Error(w, domain.ErrInvalidRequestMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(1 << 30); err != nil {
		slog.Error("Failed to parse multipart form", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData := r.FormValue(FRONTEND_FORM_JSON)

	var req data_xml.RequestXML
	if err := json.NewDecoder(strings.NewReader(jsonData)).Decode(&req); err != nil {
		slog.Error("Failed to decode xml json", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formFile, fheader, err := r.FormFile(FRONTEND_FORM_FILE)
	if err != nil {
		slog.Error("Error retrieving file", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	slog.Info("Creating temporary file")
	tempFile, err := os.CreateTemp("temp", "xml-*.xml")
	if err != nil {
		slog.Error("Error creating temporary file", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	slog.Info("Copying file to temporary")
	_, err = io.Copy(tempFile, formFile)
	if err != nil {
		slog.Error("Failed to copy temp file", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Opening file")
	file, err := os.Open(tempFile.Name())
	if err != nil {
		slog.Error("Failed to open temporary file", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	slog.Info("Decoding xml file")
	rootExist := false
	dec := xml.NewDecoder(file)
	for t, err := dec.Token(); t != nil; t, err = dec.Token() {
		if err != nil {
			slog.Error("Failed to decode token", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch se := t.(type) {
		case xml.StartElement:

			if strings.Trim(se.Name.Local, " ") == strings.Trim(req.Root, " ") {
				rootExist = true
				var rootNode data_xml.Node
				if err = dec.DecodeElement(&rootNode, &se); err != nil {
					slog.Error("Failed to decode xml", "error", err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				res, err := d.xmlUC.Upload(ctx, rootNode, req)
				if err != nil {
					slog.Error("Failed to upload xml", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = json.NewEncoder(w).Encode(*res)
				if err != nil {
					slog.Error("Failed to encode xml response", "error", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				break
			}
		}
	}

	if !rootExist {
		slog.Error("Failed to find root", "root", req.Root)
		http.Error(w, errors.New("Failed to find root").Error(), http.StatusBadRequest)
		return
	}

	slog.Info("XML uploaded successfully", "File", fheader.Filename, "Size", fheader.Size, "Header", fheader.Header)
	w.WriteHeader(http.StatusOK)
}

func (d *UploadHTTPDelivery) Run(cfg *configs.Config) {
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	mux := http.NewServeMux()

	mux.HandleFunc("/upload/excel", cors_mw.CORS(d.UploadExcelFileHandler))
	mux.HandleFunc("/upload/xml", cors_mw.CORS(d.UploadXMLHandler))

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("Failed to run http server on port %s Error: %+v", cfg.Server.Port, err)
		}
	}()
	slog.Info("Running http server", "port", cfg.Server.Port)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, os.Interrupt)
	<-quitCh

	slog.Info("Gracefully shutdowned")
}
