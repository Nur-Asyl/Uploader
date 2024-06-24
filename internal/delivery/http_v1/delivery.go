package http_v1

import (
	"Manual_Parser/configs"
	"Manual_Parser/internal/domain/data_json"
	"Manual_Parser/internal/use_case"
	"context"
	"encoding/json"
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

type UploadHTTPDelivery struct {
	excelUC use_case.UploadUseCase
}

func NewUploadHTTPDelivery(excelUC use_case.UploadUseCase) *UploadHTTPDelivery {
	return &UploadHTTPDelivery{
		excelUC: excelUC,
	}
}

func (d *UploadHTTPDelivery) UploadExcelFileHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*5)
	defer cancel()

	if r.Method != "POST" {
		slog.Error("Invalid request method")
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	slog.Info("Excel file upload endpoint: START")

	err := r.ParseMultipartForm(1 << 30)
	if err != nil {
		slog.Error("Failed to parse request form excel_rep", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData := r.FormValue("json_data")

	slog.Info(fmt.Sprintf("Request excel_rep: %+v", jsonData))

	var requestData data_json.RequestData
	err = json.NewDecoder(strings.NewReader(jsonData)).Decode(&requestData)
	if err != nil {
		slog.Error("Failed to decode request excel_rep", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("Decoded successfully", "request excel_rep", requestData)

	slog.Info("Form file")
	file, handler, err := r.FormFile("file")
	if err != nil {
		slog.Error("Error retrieving file", "error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	slog.Info("Creating temp file")
	tempFile, err := os.CreateTemp("temp", "excel-*.xlsx")
	if err != nil {
		slog.Error("Error creating temporary file", "error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name()) // to remove temp file that was saved in temp directory
	defer tempFile.Close()

	slog.Info("Copying file")
	_, err = io.Copy(tempFile, file)
	if err != nil {
		slog.Error("Error saving file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Opening file")
	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		slog.Error("Error opening excel file", "error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if err := d.excelUC.Upload(ctx, f, requestData); err != nil {
		slog.Error("Failed to upload file", "error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded and read successfully: %s, Size: %d, Header: %s", handler.Filename, handler.Size, handler.Header)

}

func (d *UploadHTTPDelivery) Run(cfg *configs.Config) {
	addr := fmt.Sprintf(":%s", cfg.Port)
	mux := http.NewServeMux()

	mux.HandleFunc("/upload/excel", d.UploadExcelFileHandler)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("Failed to run http server on port %s Error: %+v", cfg.Port, err)
		}
	}()
	slog.Info("Running http server", "port", cfg.Port)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, os.Interrupt)
	<-quitCh

	slog.Info("Gracefully shutdowned")
}
