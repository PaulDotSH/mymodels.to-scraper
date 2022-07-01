package scraper

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/antchfx/htmlquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DoRequest(method, url, reqdata, cookie string) []byte {
	client := &http.Client{}

	// Declare post data
	PostData := strings.NewReader(reqdata)

	// Declare HTTP Method and Url
	req, err := http.NewRequest(method, url, PostData)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.53 Safari/537.36")

	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	//req.Header.Set("accept-encoding", "gzip, deflate, br")

	req.Header.Set("accept-language", "en-US,en;q=0.9")

	req.Header.Set("sec-fetch-site", "same-origin")
	//text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9

	req.Header.Set("x-requested-with", "XMLHttpRequest")

	req.Header.Set("accept", "*/*")

	if err != nil {
		fmt.Println(err)
	}
	// Set cookie
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	//fmt.Println(resp.StatusCode)
	// Read response
	data, err := ioutil.ReadAll(resp.Body)
	// error handle
	if err != nil {
		fmt.Printf("error = %s \n", err)
	}
	return data
}

//Normal request instead of XHTTP thingy
func DoRequest2(method, url, reqdata, cookie string) []byte {
	client := &http.Client{}

	// Declare post data
	PostData := strings.NewReader(reqdata)

	// Declare HTTP Method and Url
	req, err := http.NewRequest(method, url, PostData)

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.53 Safari/537.36")

	req.Header.Set("content-type", "text/html; charset=UTF-8")

	//req.Header.Set("accept-encoding", "gzip, deflate, br")

	req.Header.Set("accept-language", "en-US,en;q=0.9")

	req.Header.Set("sec-fetch-site", "same-origin")

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

	if err != nil {
		fmt.Println(err)
	}
	// Set cookie
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	//fmt.Println(resp.StatusCode)
	// Read response
	data, err := ioutil.ReadAll(resp.Body)
	// error handle
	if err != nil {
		fmt.Printf("error = %s \n", err)
	}
	return data
}

func removeEverythingBefore(str, rem string) string {
	slices := strings.Split(str, rem)
	Len := len(slices)
	if Len == 1 {
		return ""
	}
	return slices[Len-1]
}

func downloadFile(filepath string, url string) (err error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func DownloadModelVideos(url, cookie string) {
	modelName := getModelName(url)
	for i := 0; true; i += 1 {
		resp := string(DoRequest("POST", "https://mymodels.to/m/pixei", "t=videos&p="+strconv.Itoa(i), "PHPSESSID=6k074qvo1ifbmuljvfrprcldqi"))
		if strings.HasPrefix(resp, "{\"videos\":null,") {
			return
		}
		json := strings.Split(resp, "\"html\"")[0] + "\"html\":\"\"}"
		fmt.Println(json)

		jsonParsed, err := gabs.ParseJSON([]byte(json))
		if err != nil {
			fmt.Println("1___")
			fmt.Println(json)
			panic(err)
		}
		v, err := jsonParsed.S("videos").Children()
		if err != nil {
			fmt.Println("2___")
			fmt.Println(jsonParsed)
			panic(err)
		}
		// https://mymodels.to/assets/images/gallery/pixei/0gvt4q4oyapbncrqk3p6i_source.mp4
		for _, child := range v {
			fileName := strings.ReplaceAll(child.S("name").String(), "\"", "")
			downloadUrl := "https://mymodels.to/assets/images/gallery/" + modelName + "/" + strings.ReplaceAll(child.S("name").String(), "\"", "")

			fmt.Println("Downloading", downloadUrl)
			err := downloadFile("./models/"+modelName+"/videos/"+fileName, downloadUrl)
			if err != nil {
				fmt.Println("Error downloading file", downloadUrl, " e:", err)
			}

		}
	}
}

func ScrapeModel(url, cookie string) {
	DownloadModelPhotos(url, cookie)
	DownloadModelVideos(url, cookie)
	DownloadModelVideos2(url, cookie)
}

func DownloadModelVideos2(url, cookie string) {
	modelName := getModelName(url)
	resp := string(DoRequest2("GET", url, "", cookie))
	doc, err := htmlquery.Parse(strings.NewReader(resp))
	if err != nil {
		fmt.Println(err)
	}
	list := htmlquery.Find(doc, "//li/a/img[@class=\"img-responsive\"]")
	fmt.Println(len(list))
	//models := make([]string, len(list))
	for _, n := range list {
		downloadUrl := strings.ReplaceAll(htmlquery.SelectAttr(n, "src"), ".png", ".mp4")
		fileName := removeEverythingBefore(downloadUrl, "/")
		fmt.Println("Downloading", downloadUrl)
		err := downloadFile("./models/"+modelName+"/videos/"+fileName, downloadUrl)
		if err != nil {
			fmt.Println("Error downloading file", downloadUrl, " e:", err)
		}
	}
}

func DownloadModelPhotos(url, cookie string) {
	modelName := getModelName(url)

	//os.Mkdir("./models", 0755)
	os.MkdirAll("./models/"+modelName+"/images/", 0755)
	os.MkdirAll("./models/"+modelName+"/videos/", 0755)

	value := true
	for i := 0; value; i += 1 {
		resp := string(DoRequest("POST", url, "t=photos&p="+strconv.Itoa(i), cookie))
		//Since this site is retarded, the value doesn't always work

		if strings.HasSuffix(resp, "\"photos\":null}") {
			return
		}

		json := strings.Split(resp, "\"html\"")[0] + "\"html\":\"\"}"

		jsonParsed, err := gabs.ParseJSON([]byte(json))
		if err != nil {
			fmt.Println("1___")
			fmt.Println(json)
			panic(err)
		}
		v, err := jsonParsed.S("photos").Children()
		if err != nil {
			//fmt.Println("2___")
			//fmt.Println(jsonParsed)
			return
		}
		for _, child := range v {
			fileName := strings.ReplaceAll(child.S("name").String(), "\"", "")
			downloadUrl := "https://mymodels.to/assets/images/gallery/" + modelName + "/" + strings.ReplaceAll(child.S("name").String(), "\"", "")

			fmt.Println("Downloading", downloadUrl)
			err := downloadFile("./models/"+modelName+"/images/"+fileName, downloadUrl)
			if err != nil {
				fmt.Println("Error downloading file", downloadUrl, " e:", err)
			}
		}

		value = jsonParsed.Path("has_more").String() == "true"
	}

	value = true

}

func GetAllModelNames() []string {
	models := GetAllModelsUrl()
	for i := 0; i < len(models); i++ {
		models[i] = removeEverythingBefore(models[i], "/")
	}
	return models
}

func getModelName(url string) string {
	return removeEverythingBefore(url, "/")
}

func GetAllModelsUrl() []string {
	resp := string(DoRequest("POST", "https://mymodels.to/admin/search-model", "action=search&searchKey=", "PHPSESSID=6k074qvo1ifbmuljvfrprcldqi"))
	doc, err := htmlquery.Parse(strings.NewReader(resp))
	if err != nil {
		fmt.Println(err)
	}
	list := htmlquery.Find(doc, "//a[@class=\"btn btn-primary btn-md\"]")
	models := make([]string, len(list))
	for i, n := range list {
		models[i] = htmlquery.SelectAttr(n, "href")
	}
	return models
}
