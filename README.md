# Gomato

 Simple scrapper written in Go to scrape the restaurants details from zomato web. The scrapper use [Colly](https://github.com/gocolly/colly) framework to scrape.
 
## Run
    go mod tidy
    go run main.go

## Config Parameters
    Depth - can be used to define the nested depth level for crawling
    
    StartPoint - Starting point of the crawler

    CsvFilePath - Filename or the filepath of the generating csv
    
    