package bot

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

var lock = &sync.Mutex{}

type FinalReport struct {
	ReporteName string
	Properties  []Propiedad
}

func (fr *FinalReport) CreateExcel(name string, address string) {

	f, err := excelize.OpenFile(name)
	if err != nil {
		f = excelize.NewFile()
	}

	f.DeleteSheet("Sheet1")
	t := time.Now()
	sheet := address + " " + fmt.Sprintf((t.Format("2006-01-02")))
	index := f.NewSheet(sheet)
	f.SetActiveSheet(index)

	f.SetCellValue(sheet, "A1", "Origen Dato")
	f.SetCellValue(sheet, "B1", "URL")
	f.SetCellValue(sheet, "C1", "Comercio")
	f.SetCellValue(sheet, "D1", "Tipo Propiedad")
	f.SetCellValue(sheet, "E1", "Comuna")
	f.SetCellValue(sheet, "F1", "Barrio")
	f.SetCellValue(sheet, "G1", "Direccion")
	f.SetCellValue(sheet, "H1", "Estado")
	f.SetCellValue(sheet, "I1", "Dormitorios")
	f.SetCellValue(sheet, "J1", "Baños")
	f.SetCellValue(sheet, "K1", "Estacionamiento")
	f.SetCellValue(sheet, "L1", "Bodega")
	f.SetCellValue(sheet, "M1", "Metros Cuadrados Totales")
	f.SetCellValue(sheet, "N1", "Metros Cuadradtos Utiles")
	f.SetCellValue(sheet, "O1", "Precio")
	f.SetCellValue(sheet, "P1", "Publicacion")
	f.SetCellValue(sheet, "Q1", "JSON")
	f.SetCellValue(sheet, "R1", "Fecha Data")

	for index, propiedad := range fr.Properties {
		excelIndex := strconv.Itoa(index + 2)
		f.SetCellValue(sheet, "A"+excelIndex, propiedad.Origen)
		f.SetCellValue(sheet, "B"+excelIndex, propiedad.URL)
		f.SetCellValue(sheet, "C"+excelIndex, propiedad.Comercio)
		f.SetCellValue(sheet, "D"+excelIndex, propiedad.TipoPropiedad)
		f.SetCellValue(sheet, "E"+excelIndex, propiedad.Comuna)
		f.SetCellValue(sheet, "F"+excelIndex, propiedad.Barrio)
		f.SetCellValue(sheet, "G"+excelIndex, propiedad.Direccion)
		f.SetCellValue(sheet, "H"+excelIndex, propiedad.Estado)
		f.SetCellValue(sheet, "I"+excelIndex, propiedad.Dormitorios)
		f.SetCellValue(sheet, "J"+excelIndex, propiedad.Banos)
		f.SetCellValue(sheet, "K"+excelIndex, propiedad.Estacionamiento)
		f.SetCellValue(sheet, "L"+excelIndex, propiedad.Bodega)
		f.SetCellValue(sheet, "M"+excelIndex, propiedad.MT2Totales)
		f.SetCellValue(sheet, "N"+excelIndex, propiedad.MT2Utiles)
		f.SetCellValue(sheet, "O"+excelIndex, propiedad.Precio)
		f.SetCellValue(sheet, "P"+excelIndex, propiedad.Fecha)
		f.SetCellValue(sheet, "Q"+excelIndex, propiedad.JSON)
		f.SetCellValue(sheet, "R"+excelIndex, propiedad.FechaData)
	}

	if err := f.SaveAs(name); err != nil {
		log.Fatal(err)
	}
}

var singleReport *FinalReport

func getFinalReport() *FinalReport {
	if singleReport == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleReport == nil {
			singleReport = &FinalReport{}
		}
	}
	return singleReport
}
