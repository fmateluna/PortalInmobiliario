package bot

import (
	"fmt"
	"strings"
)

type ResponseAddress []struct {
	ID             string `json:"id,omitempty"`
	Label          string `json:"label,omitempty"`
	Weight         int    `json:"weight,omitempty"`
	DisableIfEmpty bool   `json:"disableIfEmpty,omitempty"`
	Level          string `json:"level,omitempty"`
	Tree           Tree   `json:"tree,omitempty"`
}
type Tree struct {
	State        string `json:"state,omitempty"`
	City         string `json:"city,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
}

type UrlPropiedades struct {
	URL string `json:"url,omitempty"`
}

type Propiedad struct {
	Origen          string //URL	Tipo
	URL             string //URL	Tipo
	Comercio        string //Arriendo o Venta
	Direccion       string
	Comuna          string //Comuna
	Barrio          string //Barrio
	Estado          string //Nuevo/Usado?
	Dormitorios     string //dormitorios
	Banos           string //Baños
	Estacionamiento string //Estacionamiento
	Bodega          string //Bodega
	MT2Totales      string //MT2 Totales
	MT2Utiles       string //MT2 Utiles
	Precio          string //Precio en UF
	Fecha           string //Fecha Publicación
	TipoPropiedad   string //Fecha Publicación
	JSON            string //Fecha Publicación
	FechaData       string //Fecha Publicación
}

func (p *Propiedad) clearString(str string) string {
	return removeAccents(str)
}

func (p *Propiedad) CSV() {
	URL := strings.Split(p.URL, "#")[0]
	fmt.Print(URL, ";")
	fmt.Print(p.clearString(p.Comercio), ";")
	fmt.Print(p.clearString(p.TipoPropiedad), ";")
	fmt.Print(p.clearString(p.Comuna), ";")
	fmt.Print(p.clearString(p.Barrio), ";")
	fmt.Print(p.clearString(p.Direccion), ";")
	fmt.Print(p.clearString(p.Estado), ";")
	fmt.Print(p.clearString(p.Dormitorios), ";")
	fmt.Print(p.clearString(p.Banos), ";")
	fmt.Print(p.clearString(p.Estacionamiento), ";")
	fmt.Print(p.clearString(p.Bodega), ";")
	fmt.Print(p.clearString(p.MT2Totales), ";")
	fmt.Print(p.clearString(p.MT2Utiles), ";")
	fmt.Print(p.clearString(p.Precio), ";")
	fmt.Println(p.clearString(p.Fecha), ";")
}

func (p *Propiedad) Show() {

	//fmt.Println("URL ", p.URL)
	//fmt.Println("Comercio ", p.Comercio)
	//fmt.Println("Comuna ", p.Comuna)
	//fmt.Println("Barrio ", p.Barrio)
	//fmt.Println("Direccion ", p.Direccion)
	//fmt.Println("Estado ", p.Estado)
	//fmt.Println("Dormitorios ", p.Dormitorios)
	//fmt.Println("Banos ", p.Banos)
	//fmt.Println("Estacionamiento ", p.Estacionamiento)
	//fmt.Println("Bodega ", p.Bodega)
	//fmt.Println("MT2Totales ", p.MT2Totales)
	//fmt.Println("MT2Utiles ", p.MT2Utiles)
	//fmt.Println("Precio ", p.Precio)
	//fmt.Println("Fecha ", p.Fecha)
}
