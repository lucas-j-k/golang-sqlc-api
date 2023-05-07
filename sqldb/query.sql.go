// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package sqldb

import (
	"context"
	"database/sql"
	"time"
)

const deleteAnswers = `-- name: DeleteAnswers :exec
DELETE FROM answers
WHERE question_id = ?
`

func (q *Queries) DeleteAnswers(ctx context.Context, questionID int32) error {
	_, err := q.db.ExecContext(ctx, deleteAnswers, questionID)
	return err
}

const deleteQuestion = `-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE id = ?
`

func (q *Queries) DeleteQuestion(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteQuestion, id)
	return err
}

const getQuestion = `-- name: GetQuestion :one
SELECT id, question_text, row_inserted, row_last_updated FROM questions 
WHERE id = ?
`

func (q *Queries) GetQuestion(ctx context.Context, id int32) (Question, error) {
	row := q.db.QueryRowContext(ctx, getQuestion, id)
	var i Question
	err := row.Scan(
		&i.ID,
		&i.QuestionText,
		&i.RowInserted,
		&i.RowLastUpdated,
	)
	return i, err
}

const insertAnswer = `-- name: InsertAnswer :execresult
INSERT INTO answers (answer_text, question_id, row_inserted, row_last_updated) VALUES (?, ?, NOW(), NULL)
`

type InsertAnswerParams struct {
	AnswerText string
	QuestionID int32
}

func (q *Queries) InsertAnswer(ctx context.Context, arg InsertAnswerParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertAnswer, arg.AnswerText, arg.QuestionID)
}

const insertQuestion = `-- name: InsertQuestion :execresult
INSERT INTO questions (question_text, row_inserted, row_last_updated) VALUES (?, NOW(), NULL)
`

func (q *Queries) InsertQuestion(ctx context.Context, questionText string) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertQuestion, questionText)
}

const listAnswersForQuestion = `-- name: ListAnswersForQuestion :many
SELECT id, answer_text, row_inserted, row_last_updated FROM answers
WHERE question_id = ?
`

type ListAnswersForQuestionRow struct {
	ID             int32
	AnswerText     string
	RowInserted    time.Time
	RowLastUpdated sql.NullTime
}

func (q *Queries) ListAnswersForQuestion(ctx context.Context, questionID int32) ([]ListAnswersForQuestionRow, error) {
	rows, err := q.db.QueryContext(ctx, listAnswersForQuestion, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAnswersForQuestionRow
	for rows.Next() {
		var i ListAnswersForQuestionRow
		if err := rows.Scan(
			&i.ID,
			&i.AnswerText,
			&i.RowInserted,
			&i.RowLastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listQuestions = `-- name: ListQuestions :many
SELECT id, question_text, row_inserted, row_last_updated FROM questions
`

func (q *Queries) ListQuestions(ctx context.Context) ([]Question, error) {
	rows, err := q.db.QueryContext(ctx, listQuestions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Question
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.QuestionText,
			&i.RowInserted,
			&i.RowLastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateQuestion = `-- name: UpdateQuestion :exec
UPDATE questions SET question_text = ?, row_last_updated = NOW()
WHERE id = ?
`

type UpdateQuestionParams struct {
	QuestionText string
	ID           int32
}

func (q *Queries) UpdateQuestion(ctx context.Context, arg UpdateQuestionParams) error {
	_, err := q.db.ExecContext(ctx, updateQuestion, arg.QuestionText, arg.ID)
	return err
}
