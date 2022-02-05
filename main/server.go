package main

import (
    "os"
    "log"
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/go-sql-driver/mysql"
    "database/sql"
)

type entry struct {
    ID string `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
    // Image TODO
}

var db *sql.DB

var testEntries = []entry{
    { ID: "1", Title: "Beef Bourgignon", Description: "Julia's Finest" },
    { ID: "2", Title: "Chicken Rice", Description: "Singapore on a plate" },
    { ID: "3", Title: "Egg Tarts", Description: "Sweet custard goodness" },
}

func main() {
    // Establish database connection
    cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "fdmysql:3306",
        DBName: "fooddiary",
    }

    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    if pingErr := db.Ping(); pingErr != nil {
        log.Fatal(pingErr)
    }

    fmt.Println("Connected to database")

    // Set up routes for API
    router := gin.Default()
    router.GET("/", getAllEntries)
    router.GET("/entries/:id", getEntryById)
    router.POST("/submit", createEntry)

    router.Run("localhost:3000")
}

func getAllEntries(c *gin.Context) {
    allEntries, err := queryAllEntries()
    if err != nil {
        log.Fatal("Unable to fetch all entries: %v", err)
        return
    }

    c.IndentedJSON(http.StatusOK, allEntries)
}

func getEntryById(c *gin.Context) {
    id := c.Param("id")

    for _,item := range testEntries {
        if item.ID == id {
            c.IndentedJSON(http.StatusOK, item)
            return
        }
    }

    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
}

func createEntry(c *gin.Context) {
    var newEntry entry

    // Bind received JSON to newEntry
    if err := c.BindJSON(&newEntry); err != nil {
        return
    }

    testEntries = append(testEntries, newEntry)
    c.IndentedJSON(http.StatusCreated, testEntries)
}

func queryAllEntries() ([]entry, error) {
    var entries []entry

    rows, err := db.Query("SELECT id,title,description FROM entries")
    if err != nil {
        return nil, fmt.Errorf("queryAllEntries: %v", err)
    }

    defer rows.Close()

    for rows.Next() {
        var ent entry
        if err := rows.Scan(&ent.ID, &ent.Title, &ent.Description); err != nil {
            return nil, fmt.Errorf("queryAllEntries: %v", err)
        }
        entries = append(entries, ent)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("queryAllEntries: %v", err)
    }

    return entries, nil
}
