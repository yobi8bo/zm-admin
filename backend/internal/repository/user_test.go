package repository

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestIsDuplicateEntry(t *testing.T) {
	if !isDuplicateEntry(&mysql.MySQLError{Number: 1062}) {
		t.Fatal("MySQL duplicate-entry error was not recognized")
	}
	if isDuplicateEntry(&mysql.MySQLError{Number: 1045}) {
		t.Fatal("non-duplicate MySQL error was recognized as duplicate")
	}
	if isDuplicateEntry(errors.New("other error")) {
		t.Fatal("generic error was recognized as duplicate")
	}
}
