package model

import (
	"log"
	"time"

	"github.com/notepad_pro/db"
)

type Todo struct {
	ID        uint
	Task      string
	CreatedAt time.Time
}

type Pagination struct {
	Total           int
	Limit           int
	Current         int
	NumberOfRecords int
}

func init() {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		task TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.DB.Exec(query)
	if err != nil {
		log.Fatalf("error creating table: %v", err)
	}

	insertDemo := `
		INSERT INTO todos (task) VALUES
		('turn on pc'),
		('fire up terminal'),
		('open nvim'),
		('code till heart’s content');
	`
	_, err = db.DB.Exec(insertDemo)

	if err != nil {
		log.Fatalf("error creating table: %v", err)
	}
}

func GetList(page int) ([]*Todo, *Pagination) {
	paginate := &Pagination{}
	countQuery := `Select count(*) as count from todos`
	err := db.DB.QueryRow(countQuery).Scan(&paginate.NumberOfRecords)
	if err != nil {
		log.Fatalf("getting error fetching count for pagination : %v", err)
	}
	paginate.Limit = 10
	if paginate.NumberOfRecords < ((paginate.Current - 1) * paginate.Limit) {
		log.Fatalf("need a valid page number")
	}
	paginate.Current = page
	paginate.Total = (paginate.NumberOfRecords + paginate.Limit - 1) / paginate.Limit

	query := `SELECT id,task,created_at from todos Limit 10 offset ?`

	rows, err := db.DB.Query(query, paginate.Current)
	if err != nil {
		log.Fatalf("getting error fetching todo : %v", err)
	}
	defer rows.Close()

	var todo []*Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Task, &t.CreatedAt); err != nil {
			log.Fatal(err)
		}
		todo = append(todo, &t)
	}
	return todo, paginate
}

func CreateTodo(todo *Todo) error {
	query := `INSERT INTO todos (task,created_at) VALUES
		(?);`
	_, err := db.DB.Exec(query, todo.Task)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTodo(id int) error {
	query := `DELETE FROM todos WHERE id = ?;`
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTodo(id int, todo *Todo) error {
	query := `UPDATE TABLE todos SET task = ? WHERE id = ?`
	_, err := db.DB.Exec(query, todo.Task, id)
	if err != nil {
		return err
	}
	return nil
}
