package main

import (
	"FileServer/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	l := log.New(os.Stdout, "-", log.Default().Flags())

	StreamHandler := handlers.NewStreamHandler(l)

	sm := http.NewServeMux()

	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.html")
	})

	sm.Handle("/video", StreamHandler)

	server := http.Server{
		Addr:    ":5000",
		Handler: sm,
	}

	server.ListenAndServe()
}

//gin version
/*
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/video", func(ctx *gin.Context) {
		file, err := os.Open("./files/Cure.mp4")

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
*/
