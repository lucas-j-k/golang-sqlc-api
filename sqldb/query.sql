-- name: ListQuestions :many
SELECT id, question_text, row_inserted, row_last_updated FROM questions;

-- name: GetQuestion :one
SELECT id, question_text, row_inserted, row_last_updated FROM questions 
WHERE id = ?;

-- name: InsertQuestion :execresult
INSERT INTO questions (question_text, row_inserted, row_last_updated) VALUES (?, NOW(), NULL);
SELECT id FROM questions WHERE id = LAST_INSERT_ID();

-- name: UpdateQuestion :exec
UPDATE questions SET question_text = ?, row_last_updated = NOW()
WHERE id = ?;

-- name: ListAnswersForQuestion :many
SELECT id, answer_text, row_inserted, row_last_updated FROM answers
WHERE question_id = ?;

-- name: DeleteAnswers :exec
DELETE FROM answers
WHERE question_id = ?;

-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE id = ?;

-- name: InsertAnswer :execresult
INSERT INTO answers (answer_text, question_id, row_inserted, row_last_updated) VALUES (?, ?, NOW(), NULL);
SELECT id FROM answers WHERE id = LAST_INSERT_ID();