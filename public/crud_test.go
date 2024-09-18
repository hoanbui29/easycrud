package easycrud

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

type TestData struct {
	_    any    `easycrud:"table=test_data"`
	ID   int    `easycrud:"pkey,column=id"`
	Name string `easycrud:"column=name"`
}

func TestCreate(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_LOCAL_TEST"))

	if err != nil {
		t.Fatal("error connecting to database")
	}

	testData := TestData{
		Name: "User name",
	}
	crud := EasyCRUD[TestData, int]{
		db: db,
	}

	id, err := crud.Create(testData)

	if err != nil {
		t.Errorf("error creating data: %v", err)
		return
	}

	if id == 0 {
		t.Errorf("error creating data: id is 0")
	}

	t.Logf("created data with id: %v", id)
}

func TestDetail(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_LOCAL_TEST"))

	if err != nil {
		t.Fatal("error connecting to database")
	}

	crud := EasyCRUD[TestData, int]{
		db: db,
	}

	value, err := crud.Detail(2)

	if err != nil {
		t.Errorf("error getting data: %v", err)
		return
	}

	t.Logf("data: %v", value)
}
