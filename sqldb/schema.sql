-- Initialise the database
CREATE DATABASE questions_db;

USE questions_db;

-- Build Initial Table Schemas
CREATE TABLE questions (
    id     INT AUTO_INCREMENT NOT NULL,
    question_text  VARCHAR(2500) NOT NULL,
    row_inserted DATETIME NOT NULL DEFAULT NOW(),
    row_last_updated DATETIME NULL,
    PRIMARY KEY(id)
);

CREATE TABLE answers (
    id INT AUTO_INCREMENT NOT NULL,
    answer_text VARCHAR(2500) NOT NULL,
    question_id INT NOT NULL,
    row_inserted DATETIME NOT NULL DEFAULT NOW(),
    row_last_updated DATETIME NULL,

    CONSTRAINT fk_question
    FOREIGN KEY (question_id) 
    REFERENCES questions(id),
    PRIMARY KEY(id)
);


-- Seed initial dummy data
INSERT INTO questions (question_text,row_inserted,row_last_updated) VALUES ("Question one", NOW(), NULL);
INSERT INTO questions (question_text,row_inserted,row_last_updated) VALUES ("Question two", NOW(), NULL);

INSERT INTO answers (answer_text, question_id, row_inserted,row_last_updated) VALUES ("Answer one to question one", 1, NOW(), NULL);
INSERT INTO answers (answer_text, question_id, row_inserted,row_last_updated) VALUES ("Answer two to question one", 1, NOW(), NULL);
