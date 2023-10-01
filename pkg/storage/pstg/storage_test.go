package storage

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New()
	//os.Setenv("agrigatordb", "user=postgres password=plazma dbname=agrigatordb sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_Articles(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	posts := []Article{
		{
			Title: "Test Post3",
			Url:   strconv.Itoa(rand.Intn(1_000_000_000)),
		},
	}
	db, err := New()
	if err != nil {
		t.Fatal(err)
	}
	err = db.SaveArticles(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.LastArticles(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}
