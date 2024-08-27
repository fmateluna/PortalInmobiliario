package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/net/html"
)

var (
	URL_PORTAL  = "www.portalinmobiliario.com"
	HTTP_PORTAL = "https://" + URL_PORTAL
)

type Client struct {
	address             string
	reportExcelFileName string
	context             ResponseAddress
	Operations          map[string]string
	PropertyTypes       map[string]string
	TimeOut             int
}

type BotConfig struct {
	Timeout int `json:"timeout,omitempty"`
}

func (c *Client) Init(tipoPropiedad string, tipoComercio string) {
	/*
		OPERATION :

			242073 : Arriendo
			242074 : Arriendo temporada
			242075 : Venta propiedad

		PROPERTY_TYPE :

			242060 - Casa
			242062 - Departamento
			242067 - Oficina
			242064 - Estacionamiento
			242065 - Comercial
			242068 - Otros
			242070 - Parcela
			242072 - Tiempo compartido
	*/
	c.Operations = make(map[string]string)

	if IsValueInArray("Arriendo", tipoComercio) {
		c.Operations["Arriendo"] = "242073"
	}

	if IsValueInArray("Arriendo", tipoComercio) {
		c.Operations["Arriendo Temporada"] = "242074"
	}
	if IsValueInArray("Venta", tipoComercio) {
		c.Operations["Venta"] = "242075"
	}

	c.PropertyTypes = make(map[string]string)

	if IsValueInArray("Casa", tipoPropiedad) {
		c.PropertyTypes["Casa"] = "242060"
	}
	if IsValueInArray("Departamento", tipoPropiedad) {
		c.PropertyTypes["Departamento"] = "242062"
	}
	if IsValueInArray("Oficina", tipoPropiedad) {
		c.PropertyTypes["Oficina"] = "242067"
	}
	if IsValueInArray("Estacionamiento", tipoPropiedad) {
		c.PropertyTypes["Estacionamiento"] = "242064"
	}
	if IsValueInArray("Comercial", tipoPropiedad) {
		c.PropertyTypes["Comercial"] = "242065"
	}
	if IsValueInArray("Parcela", tipoPropiedad) {
		c.PropertyTypes["Parcela"] = "242070"
	}
	if IsValueInArray("Otros", tipoPropiedad) {
		c.PropertyTypes["Otros"] = "242068"
	}
	if IsValueInArray("Tiempo compartido", tipoPropiedad) {
		c.PropertyTypes["Tiempo compartido"] = "242072"
	}
}

func (c *Client) LoadConfig() int {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		log.Println("Falta archivo de configuracion")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config BotConfig
	json.Unmarshal(byteValue, &config)
	return config.Timeout
}

func (c *Client) CreateReport(address string, maxResult int, report string, web string, tipoPropiedad string, tipoComercio string) bool {
	c.reportExcelFileName = report
	c.TimeOut = c.LoadConfig()

	defer func() {
		if len(getFinalReport().Properties) > 0 {
			getFinalReport().CreateExcel(c.reportExcelFileName, address)
		}

	}()

	finishPortalInmov := true
	finishYapo := true

	yapoClient := Yapo{}
	if IsValueInArray("yapo", web) {
		finishYapo = false
		go func() {
			yapoClient.Keyword = strings.ToLower(address)
			yapoClient.scrapingProperties(maxResult, tipoComercio, tipoPropiedad)
			finishYapo = yapoClient.Finish
		}()
	}

	if IsValueInArray("portal", web) {
		finishPortalInmov = false
		c.Init(tipoPropiedad, tipoComercio)
		c.address = url.QueryEscape(address)
		c.getContext()
		go func() {
			for indexContext := 0; indexContext < len(c.context); indexContext++ {
				idInvoke := indexContext
				subClient := PortalInmobiliarioWebScraping{}
				subClient.City = c.context[idInvoke].Tree.City
				subClient.State = c.context[idInvoke].Tree.State
				subClient.Label = c.context[idInvoke].Label
				subClient.scrapingProperties(c.Operations, c.PropertyTypes)
				time.Sleep(3 * time.Second)
			}
			finishPortalInmov = true
		}()
	}

	if IsValueInArray("yapo", web) || IsValueInArray("portal", web) {
		bar := progressbar.NewOptions(maxResult,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowCount(),
			progressbar.OptionFullWidth(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]█[reset]",
				SaucerHead:    "[yellow]▒[reset]",
				SaucerPadding: "░",
				BarStart:      "[",
				BarEnd:        "]",
			}))

		tm := time.Now().Add(time.Duration(c.TimeOut) * time.Minute)

		//fmt.Println((finishYapo || finishPortalInmov))

		count := 0
		for len(getFinalReport().Properties) < maxResult {
			size := len(getFinalReport().Properties)
			barStatus := size

			if size > count {
				count = size
				tm = time.Now().Add(time.Duration(c.TimeOut) * time.Minute)
			}

			if finishYapo && finishPortalInmov {
				return true
			}

			bar.Set(barStatus)
			if tm.Before(time.Now()) {
				fmt.Println("\nNo se pudieron obtener mas datos, tiempo de respuesta acabado (", c.TimeOut, " Minutos)")
				return false
			}

		}

	}

	return true
}

func (c *Client) getContext() {
	req, err := http.NewRequest("GET", HTTP_PORTAL+"/faceted-search/MLC/RE/searchbox/LOCATION?LOCATION="+c.address, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", URL_PORTAL)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cookie", "_d2id=a32d738f-f00b-4d35-a860-0e1162a12770; _csrf=Gha25KRRypQR4SrN5oYlXEeg; _pi_ga=GA1.2.25651228.1684174812; _pi_ci=25651228.1684174812; _gcl_au=1.1.593197647.1684174813; _fbp=fb.1.1684174813369.686992246; __gads=ID=f5ce8eccbbc02053:T=1684174813:S=ALNI_MZ1LkSSw66pu_pnx0auVhcUOC2V5g; _pi_cbanner_mlc=1; _pi_ga_gid=GA1.2.1202318908.1684550261; _mldataSessionId=dceed0fa-965d-4050-0298-f4dd85bed4a4; __gpi=UID=000009f30c4657ef:T=1684174813:RT=1684550260:S=ALNI_MZ5mySDqSE6MapHSeV0FgJ4xfkq2A; _pi_dc=1")
	req.Header.Set("Device-Memory", "8")
	req.Header.Set("Downlink", "10")
	req.Header.Set("Dpr", "1")
	req.Header.Set("Ect", "4g")
	req.Header.Set("Referer", HTTP_PORTAL)
	req.Header.Set("Rtt", "50")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Viewport-Width", "1058")
	req.Header.Set("X-Csrf-Token", GenerateCsrfToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	responseAddress := ResponseAddress{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&responseAddress)
	if len(responseAddress) == 0 {
		panic("El sitio se jodio!! hay que tomarse esto con calma!")
	}
	c.context = responseAddress

}

type PortalInmobiliarioWebScraping struct {
	Operation    string
	PropertyType string
	City         string
	State        string
	Label        string
}

type Scraping struct {
	CacheKey         string
	OperationCode    string
	PropertyTypeCode string
	Operation        string
	PropertyType     string
	PropertyList     []string
}

func (sc *PortalInmobiliarioWebScraping) scrapingProperties(operations map[string]string, propertyTypes map[string]string) {

	for keyProTy, propertyType := range propertyTypes {
		cacheCity := []string{}
		for keyOperation, operation := range operations {
			cacheKey := sc.City + sc.State

			if len(cacheCity) == 0 || IndexOf(cacheKey, cacheCity) == -1 {
				operationCode := operation
				propertyTypeCode := propertyType
				operation := keyOperation
				propertyType := keyProTy
				scraping := Scraping{}
				scraping.OperationCode = operationCode
				scraping.PropertyTypeCode = propertyTypeCode
				scraping.Operation = operation
				scraping.PropertyType = propertyType
				//log.Println(keyOperation, keyProTy, sc.Label)
				time.Sleep(3 * time.Second)
				go func() {
					url := sc.getURL(scraping.OperationCode, scraping.PropertyTypeCode)
					scraping.getAll(url, 0)
					cacheCity = append(cacheCity, scraping.CacheKey)
				}()
			}
		}
	}
}

func (sc *PortalInmobiliarioWebScraping) getURL(operation string, propertyTpes string) string {

	url := "https://www.portalinmobiliario.com/pi/home/api/search/url?OPERATION=" + operation + "&category=MLC1459&city=" + sc.City + "&state=" + sc.State + "&PROPERTY_TYPE=" + propertyTpes + "&visual_id=PIN"
	//log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", URL_PORTAL)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "_d2id=a32d738f-f00b-4d35-a860-0e1162a12770; _csrf=Gha25KRRypQR4SrN5oYlXEeg; _pi_ga=GA1.2.25651228.1684174812; _pi_ci=25651228.1684174812; _gcl_au=1.1.593197647.1684174813; _fbp=fb.1.1684174813369.686992246; __gads=ID=f5ce8eccbbc02053:T=1684174813:S=ALNI_MZ1LkSSw66pu_pnx0auVhcUOC2V5g; _pi_cbanner_mlc=1; _pi_ga_gid=GA1.2.1202318908.1684550261; __gpi=UID=000009f30c4657ef:T=1684174813:RT=1684550260:S=ALNI_MZ5mySDqSE6MapHSeV0FgJ4xfkq2A; vis-re-search-pi-map-visited=userhide")
	req.Header.Set("Device-Memory", "8")
	req.Header.Set("Downlink", "10")
	req.Header.Set("Dpr", "1")
	req.Header.Set("Ect", "4g")
	req.Header.Set("If-None-Match", "W/\"53-OGFxynMACbqjrIW8fQI4t9+Z0O0\"")
	req.Header.Set("Rtt", "50")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Viewport-Width", "1058")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	urlPropiedades := UrlPropiedades{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&urlPropiedades)

	return urlPropiedades.URL

}

func findHref(node *html.Node) string {
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				return attr.Val
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if href := findHref(child); href != "" {
			return href
		}
	}

	return ""
}

func (sc *Scraping) getAll(url string, page int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", URL_PORTAL)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "_d2id=a32d738f-f00b-4d35-a860-0e1162a12770; _csrf=Gha25KRRypQR4SrN5oYlXEeg; _pi_ga=GA1.2.25651228.1684174812; _pi_ci=25651228.1684174812; _gcl_au=1.1.593197647.1684174813; _fbp=fb.1.1684174813369.686992246; __gads=ID=f5ce8eccbbc02053:T=1684174813:S=ALNI_MZ1LkSSw66pu_pnx0auVhcUOC2V5g; _pi_cbanner_mlc=1; _pi_ga_gid=GA1.2.1202318908.1684550261; __gpi=UID=000009f30c4657ef:T=1684174813:RT=1684550260:S=ALNI_MZ5mySDqSE6MapHSeV0FgJ4xfkq2A; vis-re-search-pi-map-visited=userhide; _mldataSessionId=5f1698de-a06e-43cb-054a-354e388fe2f6")
	req.Header.Set("Device-Memory", "8")
	req.Header.Set("Downlink", "10")
	req.Header.Set("Dpr", "1")
	req.Header.Set("Ect", "4g")
	req.Header.Set("Rtt", "50")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Viewport-Width", "1058")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if resp != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		htmlPropiedades := string(body)
		p := strings.NewReader(htmlPropiedades)
		doc, _ := goquery.NewDocumentFromReader(p)

		doc.Find(".andes-card--padding-16").Each(func(i int, el *goquery.Selection) {
			resHTML, err := el.Html()
			if err == nil {
				sp := strings.NewReader(resHTML)
				link, _ := goquery.NewDocumentFromReader(sp)

				link.Find(".andes-carousel-snapped__container--arrows-visible").Each(func(i int, link *goquery.Selection) {

					propiedad, _ := link.Html()

					doc, err := html.Parse(strings.NewReader(propiedad))
					if err != nil {
						fmt.Println("Error parsing HTML:", err)
						return
					}

					direccion := findHref(doc)

					if strings.Index(direccion, "MLC") > -1 {
						//go func() {
						propiedad := sc.getPropriedad(direccion)
						propiedad.JSON = "N/A"
						getFinalReport().Properties = append(getFinalReport().Properties, propiedad)
						//	}()
					}
				})

			}

		})

		//log.Println("nice!")
		if page == 0 {
			doc.Find(".andes-pagination__page-count").Each(func(i int, el *goquery.Selection) {
				paginas := el.Text()
				paginas = strings.Split(paginas, " ")[1]
				//log.Println("Tiene ", paginas, " Para navegar")
				_page, err := strconv.Atoi(paginas)
				if err != nil {
					page = 0
				} else {
					page = int(_page)
					for indexPage := 1; indexPage < page; indexPage++ {
						next := "/_Desde_" + strconv.Itoa(50*indexPage+1) + "_NoIndex_True"
						//log.Println(url + next)
						sc.getAll(url+next, indexPage)
					}
				}
			})
		}
	}
}
func (sc *Scraping) getAllX(url string, page int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", URL_PORTAL)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "_d2id=a32d738f-f00b-4d35-a860-0e1162a12770; _csrf=Gha25KRRypQR4SrN5oYlXEeg; _pi_ga=GA1.2.25651228.1684174812; _pi_ci=25651228.1684174812; _gcl_au=1.1.593197647.1684174813; _fbp=fb.1.1684174813369.686992246; __gads=ID=f5ce8eccbbc02053:T=1684174813:S=ALNI_MZ1LkSSw66pu_pnx0auVhcUOC2V5g; _pi_cbanner_mlc=1; _pi_ga_gid=GA1.2.1202318908.1684550261; __gpi=UID=000009f30c4657ef:T=1684174813:RT=1684550260:S=ALNI_MZ5mySDqSE6MapHSeV0FgJ4xfkq2A; vis-re-search-pi-map-visited=userhide; _mldataSessionId=5f1698de-a06e-43cb-054a-354e388fe2f6")
	req.Header.Set("Device-Memory", "8")
	req.Header.Set("Downlink", "10")
	req.Header.Set("Dpr", "1")
	req.Header.Set("Ect", "4g")
	req.Header.Set("Rtt", "50")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Viewport-Width", "1058")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if resp != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		htmlPropiedades := string(body)
		p := strings.NewReader(htmlPropiedades)
		doc, _ := goquery.NewDocumentFromReader(p)

		//pattern := regexp.MustCompile(`<div class="andes-card andes-card--flat andes-card--default ui-search-result ui-search-result--res andes-card--padding-default">(.*?)<\/div>`)
		pattern := regexp.MustCompile(`<div class="andes-card ui-search-result ui-search-result--res andes-card--flat andes-card--padding-16" id="(.*?)">(.*?)<\/div>`)

		split := pattern.Split(htmlPropiedades, -1)
		for i, propiedad := range split {
			if propiedad != "" && i > 0 {
				regex := `href="(.*?)""`
				patternAtrib := regexp.MustCompile(regex)
				splitAtrib := patternAtrib.Split(propiedad, -1)
				direccion := splitAtrib[0]

				if strings.Index(direccion, "MLC") > -1 && len(strings.Split(direccion, "\"")) > 3 {
					direccion = strings.Split(direccion, "\"")[3]
					//go func() {
					propiedad := sc.getPropriedad(direccion)
					propiedad.JSON = "N/A"
					getFinalReport().Properties = append(getFinalReport().Properties, propiedad)
					//	}()
				}
			}
		}
		//log.Println("nice!")
		if page == 0 {
			doc.Find(".andes-pagination__page-count").Each(func(i int, el *goquery.Selection) {
				paginas := el.Text()
				paginas = strings.Split(paginas, " ")[1]
				//log.Println("Tiene ", paginas, " Para navegar")
				_page, err := strconv.Atoi(paginas)
				if err != nil {
					page = 0
				} else {
					page = int(_page)
					for indexPage := 1; indexPage < page; indexPage++ {
						next := "/_Desde_" + strconv.Itoa(50*indexPage+1) + "_NoIndex_True"
						//log.Println(url + next)
						sc.getAll(url+next, indexPage)
					}
				}
			})
		}
	}
}

func (sc *Scraping) getPropriedad(url string) Propiedad {

	propiedad := Propiedad{}
	propiedad.URL = url
	propiedad.Origen = "Portal Inmobiliario"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", URL_PORTAL)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cookie", "_d2id=a32d738f-f00b-4d35-a860-0e1162a12770; _csrf=Gha25KRRypQR4SrN5oYlXEeg; _pi_ga=GA1.2.25651228.1684174812; _pi_ci=25651228.1684174812; _gcl_au=1.1.593197647.1684174813; _fbp=fb.1.1684174813369.686992246; __gads=ID=f5ce8eccbbc02053:T=1684174813:S=ALNI_MZ1LkSSw66pu_pnx0auVhcUOC2V5g; _pi_cbanner_mlc=1; _pi_ga_gid=GA1.2.1202318908.1684550261; _mldataSessionId=dceed0fa-965d-4050-0298-f4dd85bed4a4; __gpi=UID=000009f30c4657ef:T=1684174813:RT=1684550260:S=ALNI_MZ5mySDqSE6MapHSeV0FgJ4xfkq2A; _pi_dc=1")
	req.Header.Set("Device-Memory", "8")
	req.Header.Set("Downlink", "10")
	req.Header.Set("Dpr", "1")
	req.Header.Set("Ect", "4g")
	req.Header.Set("Referer", HTTP_PORTAL)
	req.Header.Set("Rtt", "50")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.Header.Set("Viewport-Width", "1058")
	req.Header.Set("X-Csrf-Token", GenerateCsrfToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	s := string(body)

	p := strings.NewReader(s)
	doc, _ := goquery.NewDocumentFromReader(p)

	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})

	atributos := make(map[string]string)
	keys := []string{}

	doc.Find(".andes-table__header--left").Each(func(i int, el *goquery.Selection) {
		key := el.Text()
		atributos[key] = ""
		keys = append(keys, key)
	})

	doc.Find(".andes-table__column--value").Each(func(i int, el *goquery.Selection) {
		value := el.Text()
		atributos[keys[i]] = value
	})

	doc.Find(".ui-pdp-image").Each(func(i int, el *goquery.Selection) {
		band, ok := el.Attr("src")
		if ok {

			if strings.HasPrefix(band, "http") && (propiedad.FechaData == "N/A" || propiedad.FechaData == "") {
				propiedad.FechaData = LastDateURL(band)
			}

		}
	})

	doc.Find(".ui-pdp-color--GRAY").Each(func(i int, el *goquery.Selection) {
		if i == 0 {
			propiedad.Fecha = el.Text()
		}
	})
	/*
		doc.Find(".ui-pdp-subtitle").Each(func(i int, el *goquery.Selection) {
			propiedad.Comercio = el.Text()
		})
	*/
	propiedad.Comercio = sc.Operation
	propiedad.TipoPropiedad = sc.PropertyType

	doc.Find(".ui-vip-location__subtitle").Each(func(i int, el *goquery.Selection) {
		propiedad.Direccion = el.Text()
	})

	doc.Find(".andes-money-amount__fraction").Each(func(i int, el *goquery.Selection) {
		propiedad.Precio += el.Text() + " "
	})

	doc.Find(".andes-breadcrumb__link").Each(func(i int, el *goquery.Selection) {
		/*
			if i == 0 {
				propiedad.TipoPropiedad += " " + el.Text()
			}
		*/
		if i == 2 {
			propiedad.Estado += el.Text()
		}
		if i == 3 || i == 4 {
			propiedad.Comuna += el.Text() + " "
		}
		/*
			if i == 1 {
				propiedad.Comercio = el.Text()
			}
		*/

		if i == 5 {
			propiedad.Barrio = el.Text()
		}

	})

	for k, v := range atributos {

		if k == "Dormitorios" {
			propiedad.Dormitorios = v
		}
		if k == "Baños" {
			propiedad.Banos = v
		}
		if k == "Estacionamientos" {
			propiedad.Estacionamiento = v
		}

		if k == "Bodegas" {
			propiedad.Bodega = v
		}
		if k == "Superficie total" {
			propiedad.MT2Totales = v
		}
		if k == "Superficie útil" {
			propiedad.MT2Utiles = v
		}
		if k == "Precio" {
			propiedad.Precio = v
		}
		if k == "Fecha" {
			propiedad.Fecha = v
		}
	}

	//propiedad.CSV(
	return propiedad
}
