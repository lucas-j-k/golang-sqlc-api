package api

import (
	"time"

	"github.com/lucas-j-k/go-sqlc-api/sqldb"
)

// ////
// Data Types
// ////

// New Question POST Body
// Decode the HTTP body JSON to this struct shape
type QuestionRequest struct {
	Text string `json:"text"`
}

// New Answer POST Body
// Decode the HTTP body JSON to this struct shape
type AnswerRequest struct {
	Text string `json:"text"`
}

// Question API response object
// DB result rows are mapped to this shape, which is then marshalled to JSON using the defined tags
type QuestionResponse struct {
	Id        int32             `json:"id"`
	Text      string            `json:"text"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt sqldb.NullTime    `json:"updatedAt"`
	Answers   *[]AnswerResponse `json:"answers,omitempty"` // pointer and omitempty allows this field to be excluded from list responses
}

// Answer response object
// Answer rows from DB mapped into this struct, then marshalled into JSON for the response
type AnswerResponse struct {
	Id        int32          `json:"id"`
	Text      string         `json:"text"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt sqldb.NullTime `json:"updatedAt"`
}

// ////
// Mappers
// ////

// Map Question DB result row into our response struct
func NewQuestionResponse(question sqldb.Question, answers []sqldb.ListAnswersForQuestionRow, includeAnswers bool) *QuestionResponse {
	resp := &QuestionResponse{
		Id:        question.ID,
		Text:      question.QuestionText,
		CreatedAt: question.RowInserted,
		UpdatedAt: sqldb.NullTime(question.RowLastUpdated),
	}

	mappedAnswers := []AnswerResponse{}

	// only map and attach the answers array if includeAnswers = true
	if includeAnswers == false {
		resp.Answers = nil
		return resp
	}

	mappedAnswers = NewAnswersListResponse(answers)
	resp.Answers = &mappedAnswers

	return resp
}

// Map list of multiple questions
func NewQuestionListResponse(questions []sqldb.Question) []*QuestionResponse {
	list := []*QuestionResponse{}
	for _, question := range questions {
		list = append(list, NewQuestionResponse(question, []sqldb.ListAnswersForQuestionRow{}, false))
	}
	return list
}

// Map single answer from DB struct to response struct
func NewAnswerResponse(answer sqldb.ListAnswersForQuestionRow) AnswerResponse {
	resp := AnswerResponse{
		Id:        answer.ID,
		Text:      answer.AnswerText,
		CreatedAt: answer.RowInserted,
		UpdatedAt: sqldb.NullTime(answer.RowLastUpdated),
	}
	return resp
}

// Map a list of multiple answers
func NewAnswersListResponse(answers []sqldb.ListAnswersForQuestionRow) []AnswerResponse {
	list := []AnswerResponse{}
	for _, answer := range answers {
		list = append(list, NewAnswerResponse(answer))
	}
	return list
}
