package main

import (
	"log"
	"net/http"
	"os"

	"FileServer/files"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/video/:id", func(ctx *gin.Context) {

		id := ctx.Param("id")
		movie, err := files.FileList.Search(id)

		log.Println(movie, err)

		if err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}

		file, err := os.Open(movie)

		if err != nil {
			log.Println(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			log.Println(err)
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		http.ServeContent(ctx.Writer, ctx.Request, "", stat.ModTime(), file)
	})

	router.Run(":5000")
}
