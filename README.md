# Tweets Scrap
this was tested with GO 1.14
```
go get -u github.com/n0madic/twitter-scraper
go get github.com/mattn/go-sqlite3
go run main.go
```
Before to RUN the fist time uncomment the Database creation Function
```
func main() {
    createdb()
    os.Exit(0)
```
To change the list of Tweets Users change this Attay
```
tweet_users := []string{"rmapalacios","larepublica_pe","canalN_","diariocorreo","policiaperu","Minsa_Peru","JulianaOxenford","elcomercio_peru","MininterPeru","peru21noticias","pcmperu"}
```

finally
```
go build main.go
```
