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
    "path/filepath"
    "github.com/google/uuid"
)

type entry struct {
    Title string `json:"title"`
    Description string `json:"description"`
    Image string `json:"image"`
}

type Env struct {
    db *sql.DB
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

    fmt.Println("Attempting to connect to database")

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    if pingErr := db.Ping(); pingErr != nil {
        log.Fatal(pingErr)
    }

    fmt.Println("Connected to database")

    env := &Env{db: db}

    // Set up routes for API
    router := gin.Default()
    router.Use(CORSMiddleware())
    router.GET("/", env.getAllEntries)
    router.GET("/entries/:id", env.getEntryById)
    router.POST("/submit", env.createEntry)
    router.Static("/images", "images/")

    router.Run(":80")
}

// For testing only
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func (e *Env) getAllEntries(c *gin.Context) {
    allEntries, err := queryAllEntries(e.db)
    if err != nil {
        log.Panic("Unable to fetch all entries: ", err)
        return
    }

    c.IndentedJSON(http.StatusOK, allEntries)
}

func (e *Env) getEntryById(c *gin.Context) {
    id := c.Param("id")

    intId, err := strconv.Atoi(id)
    if err != nil {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid id"})
        return
    }

    item, err := queryEntryById(e.db, intId)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "entry not found"})
        return
    }

    c.IndentedJSON(http.StatusOK, item)
}

func (e *Env) createEntry(c *gin.Context) {
    var newEntry entry

    file, err := c.FormFile("file")
    if err != nil {
        fmt.Println(err)
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{ "message": "Error uploading image" })
        return
    }

    extension := filepath.Ext(file.Filename)
    newFileName := uuid.New().String() + extension

    fmt.Printf("New filename being saved to %s", newFileName)

    if err := c.SaveUploadedFile(file, "images/" + newFileName); err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{ "message": "Unable to save the file" })
        return
    }

    // Bind received JSON to newEntry
    fmt.Printf("%s | %s\n", c.PostForm("title"), c.PostForm("description"))

    newEntry = entry{Title: c.PostForm("title"), Description: c.PostForm("description"), Image: "images/" + newFileName}

    id,err := queryAddEntry(e.db, newEntry)
    if err != nil {
        log.Panic("Unable to add entry")
    }

    fmt.Printf("Added ID: %d\n", id)

    c.IndentedJSON(http.StatusCreated, gin.H{"message": "Entry created"})
}

func queryEntryById(db *sql.DB, id int) (entry, error) {
    var ent entry

    row := db.QueryRow("SELECT title,description FROM entries WHERE id=?", id)
    if err := row.Scan(&ent.Title, &ent.Description); err != nil {
        if err == sql.ErrNoRows {
            return ent, fmt.Errorf("queryEntryById [%d]: no such album", id)
        }
        return ent, fmt.Errorf("queryEntryById [%d]: %v", id, err)
    }
    return ent, nil
}

func queryAllEntries(db *sql.DB) ([]entry, error) {
    var entries []entry

    rows, err := db.Query("SELECT title,description,image FROM entries")
    if err != nil {
        return nil, fmt.Errorf("queryAllEntries: %v", err)
    }

    defer rows.Close()

    for rows.Next() {
        var ent entry
        if err := rows.Scan(&ent.Title, &ent.Description, &ent.Image); err != nil {
            return nil, fmt.Errorf("queryAllEntries: %v", err)
        }
        entries = append(entries, ent)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("queryAllEntries: %v", err)
    }

    return entries, nil
}

func queryAddEntry(db *sql.DB, ent entry) (int64, error) {
    result, err := db.Exec("INSERT INTO entries (title, description, image) VALUES (?, ?, ?)", ent.Title, ent.Description, ent.Image)
    if err != nil {
        return 0, fmt.Errorf("queryAddEntry: %v", err)
    }
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("queryAddEntry: %v", err)
    }
    return id, nil
}
