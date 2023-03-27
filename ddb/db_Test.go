package ddb

import (
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	tt := TaskBean{
		Name:        "aaa",
		Url:         "http://www.google.com",
		Create_time: time.Now().String(),
	}
	InsertDB(tt)
}
