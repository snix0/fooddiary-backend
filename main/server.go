package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type entry struct {
    ID string `json:"id"`
    Title string `json:"title"`
    Description string `json:"description"`
    // Image TODO
}

var testEntries = []entry{
    { ID: "1", Title: "Beef Bourgignon", Description: "Julia's Finest" },
    { ID: "2", Title: "Chicken Rice", Description: "Singapore on a plate" },
    { ID: "3", Title: "Egg Tarts", Description: "Sweet custard goodness" },
}

func main() {
    router := gin.Default()
    router.GET("/", getAllEntries)
    router.GET("/entries/:id", getEntryById)
    router.POST("/submit", createEntry)

    router.Run("localhost:3000")
}

func getAllEntries(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, testEntries)
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
