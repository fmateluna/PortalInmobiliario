package bot

import (
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func IndexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int, specialChar string) string {
	b := make([]byte, length)
	charsetGen := charset + specialChar
	for i := range b {
		b[i] = charsetGen[seededRand.Intn(len(charsetGen))]
	}
	return string(b)
}

// ayA1sslL-G2OY3Ir8doLcGK-mo6v2KK_tOxi
// ah2aZxdS-rmKzCU3PUe8fLEHyi9_9QY6PDA4
func GenerateCsrfToken() string {
	return "a" + generateRandomString(7, "") + "-" + generateRandomString(27, "_-")
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)

	result = strings.ReplaceAll(result, "Â", "")
	result = strings.ReplaceAll(result, "ñ", "n")
	result = strings.ReplaceAll(result, "Ñ", "N")
	return result
}

func LastDateURL(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "http2.mlstatic.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("If-Modified-Since", "Tue Aug 30 14:46:31 UTC 2022")
	req.Header.Set("If-None-Match", "\"3822679362\"")
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
		// handle err
	}

	defer resp.Body.Close()

	dateString := resp.Header.Get("Last-Modified")
	layout := "Mon Jan 02 15:04:05 MST 2006"
	parsedTime, err := time.Parse(layout, dateString)
	if err != nil {
		return "N/A"
	}

	formattedDate := parsedTime.Format("02/01/2006")
	//fmt.Println(formattedDate)
	return formattedDate
}

func IsValueInArray(value string, values string) bool {
	values = RemoveLastComma(values)
	array := strings.Split(values, ",")
	for _, valueIn := range array {
		if strings.ToUpper(value) == strings.ToUpper(valueIn) {
			return true
		}
	}
	return false
}

func RemoveLastComma(text string) string {
	if text[len(text)-1] == ',' {
		return text[:len(text)-1]
	}
	return text
}
