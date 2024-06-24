package main

import (
	"Manual_Parser/configs"
	"Manual_Parser/internal/delivery/http_v1"
	"Manual_Parser/internal/repository/excel_rep"
	"Manual_Parser/internal/use_case/excel_uc"
	"Manual_Parser/pkg/storages"
	"log"
)

func main() {
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration")
	}

	storage, err := storages.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Database", err)
	}
	log.Println("Successfully connected to Database")

	excelRepo := excel_rep.NewExcelRepo(storage.GetDB())
	excelUC := excel_uc.NewExcelUseCase(excelRepo)

	deliveryHTTP := http_v1.NewUploadHTTPDelivery(excelUC)
	deliveryHTTP.Run(cfg)
}
