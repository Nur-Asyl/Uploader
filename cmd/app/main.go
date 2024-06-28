package main

import (
	"Manual_Parser/configs"
	"Manual_Parser/internal/delivery/http_v1"
	"Manual_Parser/internal/repository/excel_rep"
	"Manual_Parser/internal/repository/xml_rep"
	"Manual_Parser/internal/use_case/excel_uc"
	"Manual_Parser/internal/use_case/xml_uc"
	"Manual_Parser/pkg/storages"
	"log/slog"
)

func main() {
	slog.Info("Loading config")
	cfg := configs.GetConfig()

	slog.Info("Connecting to the Database")
	storage := storages.Connect(cfg)
	slog.Info("Successfully connected to the Database")

	slog.Info("Setting repositories")
	excelRepo := excel_rep.NewExcelRepo(storage.GetDB())
	xmlRepo := xml_rep.NewXMLRepo(storage.GetDB())

	slog.Info("Setting usecases")
	excelUC := excel_uc.NewExcelUseCase(excelRepo)
	xmlUC := xml_uc.NewXMLUseCase(xmlRepo)

	slog.Info("Setting deliveries")
	deliveryHTTP := http_v1.NewUploadHTTPDelivery(excelUC, xmlUC)

	slog.Info("Running delivery")
	deliveryHTTP.Run(cfg)
}
