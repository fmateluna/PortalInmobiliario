package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	bot "portalinmobiliario/bot"
	"strings"
	"time"
)

func main() {
	var address string
	var report string
	var tipoPropiedad string
	var webs string
	var tipoComercio string
	var max int

	flag.StringVar(&address, "d", "La Florida", "Ubicacion")
	flag.StringVar(&report, "r", "Report.xlsx", "nombre reporte")
	flag.IntVar(&max, "m", 60, "Maximo de Registro")
	flag.StringVar(&tipoPropiedad, "p", "Departamento,Casa,", "Tipo Propiedades")
	flag.StringVar(&tipoComercio, "tc", "Arriendo", "Venta/Arriendo")
	flag.StringVar(&webs, "w", "portal", "Yapo o Portal?")
	flag.Parse()
	os.Create(report + "_TMP")

	tipoPropiedad = strings.ReplaceAll(tipoPropiedad, ",,", ",")
	showtipoPropiedad := strings.ReplaceAll(tipoPropiedad, ",", " \n")
	showtipoPropiedad = "\n--------------------------\n" + strings.ReplaceAll(showtipoPropiedad, ",", " \n")

	colored := fmt.Sprintf("\x1b[%dm%s  %s  [Inicio %s]\x1b[0m  \n %s \n", 32, "Inicio deep webscriping ==> ", webs, time.Now(), showtipoPropiedad)
	fmt.Println(colored)

	bot := bot.Client{}
	isOk := bot.CreateReport(address, max, report, webs, tipoPropiedad, tipoComercio)

	os.Remove(report + "_TMP")

	if !isOk {
		fmt.Println("Presione Control + C, para cerrar el proceso")
		consoleReader := bufio.NewReaderSize(os.Stdin, 1)
		fmt.Print(">")
		input, _ := consoleReader.ReadByte()

		ascii := input

		if ascii == 27 || ascii == 3 {

			log.Println("Cerrando...")
			os.Exit(0)
		}
	} else {
		os.Exit(0)
	}

}
