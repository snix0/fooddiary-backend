package main

import (
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
)

func TestQueryAllEntries(t *testing.T) {
    db, mock, err := sqlmock.New()

    if err != nil {
        t.Fatalf("An error '%s' has occurred while opening the stub DB connection", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{"title", "description", "image"}).AddRow("Beef Wellington", "This is a test description", "images/test.png")

    mock.ExpectQuery("SELECT title,description,image FROM entries").WillReturnRows(rows)

    if _, err = queryAllEntries(db); err != nil {
        t.Errorf("Call to queryAllEntries failed: %s", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}

func TestQueryEntryById(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' has occurred while opening the stub DB connection", err)
    }
    defer db.Close()
    const idx = 1

    rows := sqlmock.NewRows([]string{"title", "description"}).AddRow("Beef Wellington", "This is a test description")

    mock.ExpectQuery("SELECT title,description FROM entries WHERE id=?").WithArgs(idx).WillReturnRows(rows)

    if _, err = queryEntryById(db, idx); err != nil {
        t.Errorf("Call to queryEntryByid failed: %s", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}

func TestQueryAddEntry(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' has occurred while opening the stub DB connection", err)
    }
    defer db.Close()

    ent := entry{"Beef Wellington", "This is a test description", "images/test.png"}

    //rows := sqlmock.NewRows([]string{"title", "description"}).AddRow("Beef Wellington", "This is a test description")

    mock.ExpectExec("INSERT INTO entries").WithArgs("Beef Wellington", "This is a test description", "images/test.png").WillReturnResult(sqlmock.NewResult(1, 1))

    if _, err = queryAddEntry(db, ent); err != nil {
        t.Errorf("Call to queryAddEntry failed: %s", err)
    }

    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }

}
