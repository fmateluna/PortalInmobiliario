package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Yapo struct {
	Query     string //{"category":[1220,1240,1260]}
	Keyword   string //ñuñoa
	MaxLimit  int
	Finish    bool
	Comercio  string
	Propiedad string
}

func (yp *Yapo) scrapingProperties(limit int, tipoComercio string, tipoPropiedad string) {
	yp.Finish = false
	yp.MaxLimit = limit
	yp.Comercio = tipoComercio
	yp.Propiedad = tipoPropiedad
	yp.MakeQuery()
	yp.Find()
	yp.Finish = true
}

func (yp *Yapo) MakeQuery() {

	comercial := ""
	//"1220" = venta
	if IsValueInArray("venta", yp.Comercio) {
		comercial = comercial + "1220,"
	}

	//"1240" = arriendo
	if IsValueInArray("arriendo", yp.Comercio) {
		comercial = comercial + "1240,"
	}

	//"1260" =
	if IsValueInArray("Arriendo de temporada", yp.Comercio) {
		comercial = comercial + "1260,"
	}

	comercial = RemoveLastComma(comercial)

	tiposDeVenta := "{\"category\":[" + comercial + "]}"

	yp.Query = url.QueryEscape(tiposDeVenta)
}

func (yp *Yapo) Find() {

	max := yp.StepByStep(0)
	if max < yp.MaxLimit {
		yp.MaxLimit = max
	}
	steps := max / 10
	for step := 1; step < steps; step++ {
		subStep := step
		go func() {
			max = yp.StepByStep(subStep)
		}()
		time.Sleep(2 * time.Second)
		if ((step+1)*10) > max || step == 0 {
			return
		}
	}
}

func (yp *Yapo) StepByStep(step int) int {
	url := "https://public-api.yapo.cl/buyers/search?page=" + strconv.Itoa(step) + "&limit=10&query=" + yp.Query + "&keyword=" + url.QueryEscape(yp.Keyword)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0
	}
	req.Header.Set("Authority", "public-api.yapo.cl")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Origin", "https://new.yapo.cl")
	req.Header.Set("Referer", "https://new.yapo.cl/")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Chref", "WEB")
	req.Header.Set("X-Cmref", "client")
	req.Header.Set("X-Commerce", "Yapo")
	req.Header.Set("X-Country", "CL")
	req.Header.Set("X-Domain", "Buyer")
	req.Header.Set("X-Rhsref", "new.yapo.cl")
	req.Header.Set("X-Txref", "a7b0d681-053e-4ae8-a965-a0800ada2420")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}
	body, _ := ioutil.ReadAll(resp.Body)
	yapoResponse := YapoResponse{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&yapoResponse)

	if len(yapoResponse.Ads) == 0 {
		return 0
	}

	for _, yapo := range yapoResponse.Ads {
		if strings.ToLower(yapo.Location.CommuneName) == strings.ToLower(yp.Keyword) && (IsValueInArray(yapo.EstateType, yp.Propiedad)) {
			propiedad := Propiedad{}
			propiedad.TipoPropiedad = yapo.EstateType
			propiedad.Direccion = yapo.Subject
			propiedad.Comuna = yapo.Location.RegionName
			propiedad.Barrio = yapo.Location.CommuneName
			propiedad.Precio = strconv.Itoa(yapo.Price) + " " + yapo.Currency
			propiedad.MT2Totales = strconv.Itoa(yapo.Size)
			propiedad.MT2Utiles = strconv.Itoa(yapo.UtilSize)
			propiedad.Banos = strconv.Itoa(yapo.Bathrooms)
			propiedad.Dormitorios = strconv.Itoa(yapo.Rooms)
			propiedad.Fecha = yapo.ListTime

			propiedad.Comercio = yapo.Category.Name
			urlYapo := yp.MakeURL(yapo)
			propiedad.URL = urlYapo
			propiedad.Origen = "Yapo"
			details, urlJson := yp.getFeatures(yapo.ListID)

			propiedad.FechaData = yapo.ListTime

			propiedad.JSON = urlJson
			propiedad.Estacionamiento = "0"
			propiedad.Bodega = "0"
			for _, feature := range details.Features {
				if feature.Label == "garages" {
					propiedad.Estacionamiento = feature.Value
				}

				if feature.Label == "others" {
					for _, prop := range strings.Split(feature.Value, ",") {
						if prop == "Bodega" {
							propiedad.Bodega = "1"
						}
					}
				}
			}
			getFinalReport().Properties = append(getFinalReport().Properties, propiedad)
		}
	}
	defer resp.Body.Close()
	return yapoResponse.Pagination.TotalRecords
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func (yp *Yapo) getFeatures(id int) (FeaturesResponse, string) {
	url := "https://public-api.yapo.cl/buyers/items/" + strconv.Itoa(id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "public-api.yapo.cl")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Origin", "https://new.yapo.cl")
	req.Header.Set("Referer", "https://new.yapo.cl/")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("X-Chref", "WEB")
	req.Header.Set("X-Cmref", "client")
	req.Header.Set("X-Commerce", "Yapo")
	req.Header.Set("X-Country", "CL")
	req.Header.Set("X-Domain", "Buyer")
	req.Header.Set("X-Rhsref", "new.yapo.cl")
	req.Header.Set("X-Txref", "a7b0d681-053e-4ae8-a965-a0800ada2420")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	features := FeaturesResponse{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&features)
	return features, url
}

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func (yp *Yapo) MakeURL(yapo Ads) string {
	ref := "https://new.yapo.cl/buscar?buscar=" + yapo.Subject
	url := "https://new.yapo.cl/" + yapo.Category.ParentName + "/" + strings.ReplaceAll(clearString(yapo.Subject), " ", "-") + "_" + strconv.Itoa(yapo.ListID)
	url = strings.ToLower(url)

	if yp.is200Status(url, yapo) {
		return url
	} else {
		return ref
	}
}

func (yp *Yapo) is200Status(url string, yapo Ads) bool {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "new.yapo.cl")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "_gcl_au=1.1.1934210061.1684174815; _ga=GA1.1.8331166.1684174815; _fbp=fb.1.1684174815986.638703869; _cc_id=a6cc530903c5ba1ed9d53d9080a15398; panoramaId_expiry=1686877999455; panoramaId=3e6df0d3ee71985d696e52ef5cd616d539386c4b4495077858474428b259bca6; panoramaIdType=panoIndiv; _pbjs_userid_consent_data=3524755945110770; _pubcid=2d3b331b-3cde-45b4-a3f2-aed89e71fa17; pbjs-unifiedid=%7B%22TDID%22%3A%2204c19b11-bb2f-4daa-8b11-7a7c779ca0ae%22%2C%22TDID_LOOKUP%22%3A%22TRUE%22%2C%22TDID_CREATED_AT%22%3A%222023-05-09T01%3A13%3A35%22%7D; ln_or=eyI0NjEzOTg2IjoiZCJ9; _hjSessionUser_3357621=eyJpZCI6ImQ2OGUyMDM5LWE5MjItNTk1ZS05ZWFmLWZiY2U1NDZmMWI5MyIsImNyZWF0ZWQiOjE2ODYyNzM4NTM5MDUsImV4aXN0aW5nIjp0cnVlfQ==; __gsas=ID=11031057b1fadd28:T=1686718151:RT=1686718151:S=ALNI_Mbu9ak7vcZj6iisio8vFnlhH1bbcQ; ___iat_vis=F81B9E7F7176DB68.56d7c72814cea24ebfd8b7a5fd4f4d8a.1686745400372.8797646842a63dd96cfb5c7c9777590e.BEBRZJIEAA.11111111.1.0; _ga_T1TKGF3G8X=GS1.1.1686745114.8.1.1686745456.60.0.0; __gads=ID=59c0b59dd893318a-22093132a7b400e9:T=1684174814:RT=1686745456:S=ALNI_MbMBBMyr_N6AkEnCA5mgUBf2DeOnw; __gpi=UID=000009f30c3d45b0:T=1684174814:RT=1686745456:S=ALNI_MZ2WdMmWeJsleHwRpCwmPiJo3evYg; cto_bundle=nqyDOF81U284R1IzRGFNUmxiTnhKN0RnQ29ma0hqN1dENTlsV0UzN2RkYUpnWEo3Q28lMkZHSFVqTzYyY1M3bUl2ZFZ5cTVJQzJUZSUyQkxDNEJXV3pwRGglMkJ5JTJCem5RZkhIN0V1VCUyQlZZSiUyRnJSTiUyQmtCRHIweTlWSFFIa2x6RiUyQjBvUjZ1MVh5Qk44QVFFa3F5Z1ljN25ubzM3RVBOOTN3JTNEJTNE; cto_bidid=57oi-F9FZGwyWVNsNlRab2NTaFZobDhsc3hQcWdxZUZWc1U4MjQ4YkVqbW1Ed3dCZFJidFI3OCUyRmFKUzVmQWU5S0ljJTJGWmNqQTZYJTJCV2ZXbTh4bXBYRzN2ZFd6QyUyRldnJTJGcUhzTFRab3RwOVM3bkdZTUklM0Q")
	req.Header.Set("If-Modified-Since", "Wed, 14 Jun 2023 04:39:50 GMT")
	req.Header.Set("Sec-Ch-Ua", "\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true

}
