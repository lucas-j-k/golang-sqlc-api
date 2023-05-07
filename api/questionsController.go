package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lucas-j-k/go-sqlc-api/sqldb"
)

type QuestionController struct {
	queries *sqldb.Queries
	db      *sql.DB
	context context.Context
}

// Initialise a controller instance
func NewQuestionController(queries *sqldb.Queries, db *sql.DB, context context.Context) *QuestionController {
	return &QuestionController{
		queries: queries,
		db:      db,
		context: context,
	}
}

// ////
// Route Handlers
// ////

// List all question rows
func (controller *QuestionController) ListQuestions(w http.ResponseWriter, r *http.Request) {

	// Retrieve questions from database via SQLC
	questions, err := controller.queries.ListQuestions(controller.context)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	// marshall list of rows to JSON and return via chi render
	render.JSON(w, r, NewQuestionListResponse(questions))
}

// Query for a specific question row based on the id (int32)
func (controller *QuestionController) GetQuestionById(w http.ResponseWriter, r *http.Request) {

	// grab and parse the ID param from the url path
	id := chi.URLParam(r, "id")
	parsed, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	idParam := int32(parsed)

	question, err := controller.queries.GetQuestion(controller.context, idParam)
	if err != nil {
		render.Render(w, r, ErrNotFound(err))
		return
	}

	// List any answers related to the question
	answers, err := controller.queries.ListAnswersForQuestion(controller.context, idParam)
	if err != nil {
		render.Render(w, r, ErrInternalServer(err))
		return
	}

	render.JSON(w, r, NewQuestionResponse(question, answers, true))
}

// Insert a quesion row from POST request
func (controller *QuestionController) InsertQuestion(w http.ResponseWriter, r *http.Request) {

	// Bind the incoming JSON body to our QuestionRequest struct
	payload := &QuestionRequest{}
	if err := render.Decode(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// Write the new question to the DB
	result, err := controller.queries.InsertQuestion(controller.context, payload.Text)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	created, err := result.LastInsertId()
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, SimpleIdResponse{Id: int(created)})
}

// Update a quesion row from PUT request
func (controller *QuestionController) UpdateQuestion(w http.ResponseWriter, r *http.Request) {

	// grab and parse the ID param from the url path
	id := chi.URLParam(r, "id")
	parsed, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	idParam := int32(parsed)

	// Bind the incoming JSON body to our QuestionRequest struct
	payload := &QuestionRequest{}
	if err := render.Decode(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// check the question exists
	_, existingErr := controller.queries.GetQuestion(controller.context, idParam)
	if existingErr != nil {
		render.Render(w, r, ErrNotFound(err))
		return
	}

	// Build the SQL query payload
	sqlPayload := sqldb.UpdateQuestionParams{
		ID:           idParam,
		QuestionText: payload.Text,
	}

	// Write the new question to the DB
	updateErr := controller.queries.UpdateQuestion(controller.context, sqlPayload)
	if updateErr != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, SimpleIdResponse{Id: int(idParam)})
}

// Add an answer to a question
func (controller *QuestionController) InsertAnswer(w http.ResponseWriter, r *http.Request) {

	// get the question ID from the url path
	id := chi.URLParam(r, "id")
	parsed, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	idParam := int32(parsed)

	// Bind the incoming JSON body to our QuestionRequest struct
	payload := &AnswerRequest{}
	if err := render.Decode(r, payload); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// check the question exists
	_, existingErr := controller.queries.GetQuestion(controller.context, idParam)
	if existingErr != nil {
		render.Render(w, r, ErrNotFound(err))
		return
	}

	sqlPayload := sqldb.InsertAnswerParams{
		AnswerText: payload.Text,
		QuestionID: idParam,
	}

	// Write the new question to the DB
	result, err := controller.queries.InsertAnswer(controller.context, sqlPayload)
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	created, err := result.LastInsertId()
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, SimpleIdResponse{Id: int(created)})
}

// Delete a question and all associated answers
func (controller *QuestionController) DeleteQuestion(w http.ResponseWriter, r *http.Request) {

	// get the question ID from the url path
	id := chi.URLParam(r, "id")
	parsed, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	idParam := int32(parsed)

	// check the question exists
	_, existingErr := controller.queries.GetQuestion(controller.context, idParam)
	if existingErr != nil {
		render.Render(w, r, ErrNotFound(err))
		return
	}

	// Delete answers and then question with a transaction
	tx, err := controller.db.Begin()
	if err != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	qtx := controller.queries.WithTx(tx)
	answersErr := qtx.DeleteAnswers(controller.context, idParam)
	if answersErr != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	questionErr := qtx.DeleteQuestion(controller.context, idParam)
	if questionErr != nil {
		render.JSON(w, r, ErrInternalServer(err))
		return
	}

	tx.Commit()

	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}
