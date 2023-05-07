package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lucas-j-k/go-sqlc-api/api"
	"github.com/lucas-j-k/go-sqlc-api/sqldb"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()

	// initialize env vars and SQL connection
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	port := viper.Get("PORT")

	err := sqldb.Connect()

	if err != nil {
		log.Panic("Unable to connect to MYSQL")
	}

	// initialize SQLC and http controller
	queries := sqldb.New(sqldb.DB)
	questionsController := api.NewQuestionController(queries, sqldb.DB, ctx)

	// setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// route handlers
	r.Route("/questions", func(r chi.Router) {
		r.Get("/{id}", questionsController.GetQuestionById)
		r.Put("/{id}", questionsController.UpdateQuestion)
		r.Delete("/{id}", questionsController.DeleteQuestion)
		r.Post("/{id}/answers", questionsController.InsertAnswer)
		r.Get("/", questionsController.ListQuestions)
		r.Post("/", questionsController.InsertQuestion)
	})

	// healthcheck
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong"))
	})

	fmt.Printf("Server running on port [%v]\n\n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}
