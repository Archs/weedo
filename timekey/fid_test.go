package timekey

import (
	"os"
	"testing"
	"time"
)

var (
	testFile = []string{"mime_data.go", "fid.go"}
)

func TestKey(t *testing.T) {
	now := time.Now()
	t.Log("now:", now)
	fid, err := NewFid("2", testFile[1])
	if err != nil {
		t.Fatal(err)
	}
	t.Log("fid.Time()", fid.Time())
	if now.Sub(fid.Time()).Minutes() > 1.0 {
		t.Fatal("time not match")
	}
}

func TestCookie(t *testing.T) {
	fid, err := NewFid("2", testFile[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Log("fid.MimeType():", fid.MimeType())
	t.Log("fid.Size():", fid.Size(), "KB")
	info, err := os.Stat(testFile[0])
	if err != nil {
		t.Fatal(err)
	}
	t.Log("real file size:", info.Size()/1024)
}
