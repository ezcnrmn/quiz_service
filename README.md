a simple quiz service that provides the following APIs:
- .../quiz
  GET: returns list of all quizzes
  POST: gets a new quiz and saves it
  PATCH: edits existing quiz
  DELETE: deletes quiz
- .../quiz/:id
  GET: returns quiz data
  POST: checks quiz answers

  how to run:
  1. build main.go
  2. run .exe
  you can specify address using flagg --addr (default addres is 127.0.0.1:3001)

all quizzes are stored in the "quizzes_storage" folder next to the .exe file
