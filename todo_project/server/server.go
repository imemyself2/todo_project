package main

import (
	"log"
	"net/http"
	"os"

	todo "github.com/1set/todotxt"
	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	loginInfo := make(map[string]string)
	todoInfo := make(map[string]todo.TaskList)
	port := os.Getenv("PORT")
	// Test ID
	loginInfo["sample"] = "helloworld123"

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/login", func(c *gin.Context) {
		username := c.Query("username")
		password := c.Query("password")

		if pass, ok := loginInfo[username]; ok {
			if pass == password {
				c.Redirect(http.StatusOK, "/"+username)
			} else {
				c.Status(http.StatusUnauthorized)
			}
		}
	})

	router.GET("/", func(c *gin.Context) {
		username := c.Param("username")
		var erro error
		if todoInfo[username], erro = todo.LoadFromPath("../todo.txt"); erro != nil {
			log.Fatal(erro)
			c.Status(http.StatusInternalServerError)
		} else {

			// os.Chdir("../")
			// cmdOut, _ := exec.Command("./todo", "ls").CombinedOutput()
			// os.Chdir("server")

			// exec.Command("cd server").CombinedOutput()

			todolist := todoInfo[username]

			incomplete := todolist.Filter(todo.FilterNotCompleted)
			incomplete.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc)
			newListIncomplete := incomplete.String()
			type response struct {
				Incomplete string `json:"incomplete"`
				Complete   string `json:"complete"`
			}
			complete := todolist.Filter(todo.FilterCompleted)
			complete.Sort(todo.SortPriorityAsc, todo.SortDueDateAsc, todo.SortCreatedDateAsc)
			newListComplete := complete.String()

			answer := response{
				Incomplete: newListIncomplete,
				Complete:   newListComplete,
			}

			c.JSON(http.StatusOK, answer)
		}

	})

	router.POST("/add", func(c *gin.Context) {

	})

	router.POST("/rm", func(c *gin.Context) {

	})

	router.Run(":" + port)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
