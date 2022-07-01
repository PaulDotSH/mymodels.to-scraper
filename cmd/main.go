package main

import "mymodels.to_scraper/scraper"

func main() {
	//modelUrls := scraper.GetAllModelsUrl()
	//for _, url := range modelUrls {
	//	scraper.ScrapeModel(url, "PHPSESSID=6k074qvo1ifbmuljvfrprcldqi")
	//}

	scraper.ScrapeModel("https://mymodels.to/m/miakhalifa", "PHPSESSID=6k074qvo1ifbmuljvfrprcldqi")
}
