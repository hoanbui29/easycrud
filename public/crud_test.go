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

var key int

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

	key = id

	if err != nil {
		t.Errorf("error creating data: %v", err)
		t.FailNow()
	}

	if key == 0 {
		t.Errorf("error creating data: id is 0")
		t.FailNow()
	}

	t.Logf("created data with id: %v", key)
}

func TestDetail(t *testing.T) {
	t.Logf("key: %v", key)
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_LOCAL_TEST"))

	if err != nil {
		t.Fatal("error connecting to database")
	}

	crud := EasyCRUD[TestData, int]{
		db: db,
	}

	value, err := crud.Detail(key)

	if err != nil {
		t.Errorf("error getting data: %v", err)
		t.FailNow()
	}

	t.Logf("data: %v", value)
}

func TestUpdate(t *testing.T) {
	t.Logf("key: %v", key)
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_LOCAL_TEST"))

	if err != nil {
		t.Fatal("error connecting to database")
	}

	crud := EasyCRUD[TestData, int]{
		db: db,
	}

	prev, err := crud.Detail(key)

	if err != nil {
		t.Errorf("error getting data: %v", err)
		t.FailNow()
	}

	prev.Name = "Updated name"

	isSuccess, err := crud.Update(prev)

	if err != nil {
		t.Errorf("error updating data: %v", err)
		t.FailNow()
	}

	if !isSuccess {
		t.Errorf("error updating data: not success")
		t.FailNow()
	}

	t.Logf("data with id %d updated", key)
}

func TestDelete(t *testing.T) {
	t.Logf("key: %v", key)
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_LOCAL_TEST"))

	if err != nil {
		t.Fatal("error connecting to database")
	}

	crud := EasyCRUD[TestData, int]{
		db: db,
	}

	isSuccess, err := crud.Delete(key)

	if err != nil {
		t.Errorf("error deleting data: %v", err)
		t.FailNow()
	}

	if !isSuccess {
		t.Errorf("error deleting data: not success")
		t.FailNow()
	}

	t.Logf("data with id %d deleted", key)
}
