package main

import (
    "os"
    "log"
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/go-sql-driver/mysql"
    "database/sql"
    "strconv"
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

    intId, err := strconv.Atoi(id)
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }

    item, err := queryEntryById(intId)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
        return
    }

    c.IndentedJSON(http.StatusOK, item)
}

func createEntry(c *gin.Context) {
    var newEntry entry

    // Bind received JSON to newEntry
    if err := c.BindJSON(&newEntry); err != nil {
        return
    }

    testEntries = append(testEntries, newEntry)

    id,err := dbAddEntry(newEntry)
    if err != nil {
        log.Panic("Unable to add entry")
    }

    fmt.Println("Added ID: %d", id)

    c.IndentedJSON(http.StatusCreated, gin.H{"message": "Entry created"})
}

func queryEntryById(id int) (entry, error) {
    var ent entry

    row := db.QueryRow("SELECT id,title,description FROM entries WHERE id=?", id)
    if err := row.Scan(&ent.ID, &ent.Title, &ent.Description); err != nil {
        if err == sql.ErrNoRows {
            return ent, fmt.Errorf("queryEntryById [%d]: no such album", id)
        }
        return ent, fmt.Errorf("queryEntryById [%d]: %v", id, err)
    }
    return ent, nil
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

func dbAddEntry(ent entry) (int64, error) {
    result, err := db.Exec("INSERT INTO entries (title, description, image) VALUES (?, ?, '')", ent.Title, ent.Description)
    if err != nil {
        return 0, fmt.Errorf("dbAddEntry: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("dbAddEntry: %v", err)
    }
    return id, nil
}
