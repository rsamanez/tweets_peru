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
finally
```
go build main.go
```
