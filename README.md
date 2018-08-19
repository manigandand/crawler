# crawler
A simple web crawler in Go.


Package Manager: `glide`

To install glide: `curl https://glide.sh/get | sh`

Step to run:
1. `glide install`
2. run `./webcrawler`


Available Endpoints:
1. Index page
GET `http://127.0.0.1:8080/`

2. To scrap the site
POST `http://127.0.0.1:8080/crawler/` 

3. To view all the sitemap
GET `http://127.0.0.1:8080/crawler/status/` 


