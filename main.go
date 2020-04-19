package main

import (
    "database/sql"
    "encoding/csv"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    twitterscraper "github.com/n0madic/twitter-scraper"
    "log"
    "os"
    "strings"
    "time"
)

func createdb(){
    database, _ := sql.Open("sqlite3", "./tweets.db")
    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS tweets (id INTEGER PRIMARY KEY, " +
       "tweet_user CHAR(20), tweet_id CHAR(20), hashtags TEXT, html TEXT, ispin BOOLEAN, isretweet BOOLEAN, " +
       "likes int, replies int, retweets int, body TEXT, timestamp BIGINT)")
    statement.Exec()
    statement, _ = database.Prepare("CREATE TABLE IF NOT EXISTS csv_files (id INTEGER PRIMARY KEY, " +
        "csv_filename CHAR(30), csv_last_id int, timestamp BIGINT)")
    statement.Exec()
    //statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
    //statement.Exec("Nic", "Raboy")
    //rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
    //var id int
    //var firstname string
    //var lastname string
    //for rows.Next() {
    //    rows.Scan(&id, &firstname, &lastname)
    //    fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
    //}
}
// https://twitter.com/pcmperu/status/1235350644346761216
func createNewCsvFile(database *sql.DB,rowId int ,twRowId int) {
    rows, err := database.Query("SELECT id,hashtags,body,timestamp,tweet_user FROM tweets where id >?", rowId)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    data := [][]string{}
    data = append(data, []string{"title", "time","content"})
    var id int
    var tweetUser string
    var hashtag string
    var body string
    var twTime int64
    for rows.Next() {
        rows.Scan(&id, &hashtag, &body, &twTime,&tweetUser)
        data = append(data, []string{tweetUser+hashtag , GetTodaysDateTime(twTime), body})
    }
    var filename = "result.csv"
    if err := csvExport(data,filename); err != nil {
        log.Fatal(err)
    }
    // actualizo Tabla csv_last_id
    statement, _ := database.Prepare("INSERT INTO csv_files (csv_filename,csv_last_id,timestamp) values (?,?,?) ")
    statement.Exec(filename,id,twTime)
}

func csvExport(data [][]string,filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, value := range data {
        if err := writer.Write(value); err != nil {
            return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
        }
    }
    return nil
}

func checkLastRowSend(database *sql.DB) bool {
    var rowId int
    var err = database.QueryRow("SELECT csv_last_id FROM csv_files ORDER BY id DESC LIMIT 1").Scan(&rowId)
    if err == nil {
        var twRowId int
        var err2 = database.QueryRow("SELECT id FROM tweets ORDER BY id DESC LIMIT 1").Scan(&twRowId)
        if err2 ==nil {
            if twRowId > rowId {
                createNewCsvFile(database,rowId,twRowId)
                return true
            }
        }else{
            fmt.Println("ERROR de lectura de BASE DE DATOS tweets")
        }
    }else{
        fmt.Println("ERROR de lectura de BASE DE DATOS csv_files")
    }
    return false
}
func GetTodaysDateTime(timestamp int64) string {
    //loc, _ := time.LoadLocation("America/Los_Angeles")
    //current_time := time.Now().In(loc)
    current_time := time.Unix(timestamp, 0)
    return current_time.Format("2006-01-02 15:04:05")
}
func main() {
    //createdb()
    //fmt.Print(GetTodaysDateTime())
    //os.Exit(0)
    tweet_users := []string{"rmapalacios","larepublica_pe","canalN_","diariocorreo","policiaperu","Minsa_Peru","JulianaOxenford","elcomercio_peru","MininterPeru","peru21noticias","pcmperu"}
    //tweet_users := []string{"RRsamanez"}
    database, _ := sql.Open("sqlite3", "./tweets.db")
    statement, _ := database.Prepare("INSERT INTO tweets (tweet_user, " +
        "tweet_id, hashtags, html, ispin, isretweet, likes, replies, retweets, body, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
    for _, tweet_user := range tweet_users {
        fmt.Println(tweet_user)
        for tweet := range twitterscraper.GetTweets(tweet_user, 25) {
            if tweet.Error != nil {
                panic(tweet.Error)
            }
            var row_id string
            var err = database.QueryRow("SELECT id FROM tweets WHERE tweet_id=? AND tweet_user=?",tweet.ID,tweet_user).Scan(&row_id)
            if err == sql.ErrNoRows {
                fmt.Println(tweet.ID)
                s := strings.ReplaceAll(fmt.Sprint(tweet.Hashtags),"[", "")
                ss := strings.ReplaceAll(s,"]", "")
                statement.Exec(tweet_user, tweet.ID,ss, tweet.HTML, tweet.IsPin, tweet.IsRetweet, tweet.Likes, tweet.Replies, tweet.Retweets,
                    tweet.Text, tweet.Timestamp)
            }else{
                //fmt.Println(tweet.ID+" Already Exists")
                fmt.Print(".")
            }
        }
    }

    if checkLastRowSend(database) {
        fmt.Print("Prepare to send new CSV to SumUp")
    }else{
        fmt.Println("Nothing to Update")
    }
}
