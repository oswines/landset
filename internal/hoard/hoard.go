package hoard

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
	api "github.com/oswines/landset"
)

const makeTable string = `
CREATE TABLE IF NOT EXISTS hoard (
id INTEGER NOT NULL PRIMARY KEY,
inning TEXT,
created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`
const updateTringger string = `
CREATE TRIGGER IF NOT EXISTS update_timestamp
AFTER UPDATE ON hoard
FOR EACH ROW
BEGIN
    UPDATE hoard
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;
`

const file string = "hoard.db"

type Hoard struct {
	db *sql.DB
}

func NewHoard() (*Hoard, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(makeTable); err != nil {
		return nil, err
	}
	if _, err := db.Exec(updateTringger); err != nil {
		return nil, err
	}
	return &Hoard{
		db: db,
	}, nil
}

func (c *Hoard) Insert(inlay api.Inlay) (int, error) {
	res, err := c.db.Exec("INSERT INTO hoard (inning) VALUES(?);", inlay.Inning)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	log.Printf("Added %v as %d", inlay, id)
	return int(id), nil
}

var ErrIDNotFound = errors.New("ID not found")

func (c *Hoard) Fetch(id int) (api.Inlay, error) {
	log.Printf("Getting %d", id)

	row := c.db.QueryRow("SELECT id, inning, created_at, updated_at FROM hoard WHERE id=?", id)

	inlay := api.Inlay{}
	var err error
	if err = row.Scan(&inlay.ID, &inlay.Inning, &inlay.CreatedAt, &inlay.UpdatedAt); err == sql.ErrNoRows {
		log.Printf("ID not found")
		return api.Inlay{}, ErrIDNotFound
	}
	return inlay, err
}
