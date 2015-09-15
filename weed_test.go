// weed_test.go
package weedo

import (
	"os"
	"testing"
)

var (
	client   = NewClient("localhost:9333", "localhost:8888")
	filename = "hello.txt"
)

func TestAssign(t *testing.T) {
	fid, err := client.Master().Assign()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("assign", fid)

	fid, err = client.Master().AssignN(3)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("assign 3", fid)
}

func TestFid(t *testing.T) {
	fid, err := client.Master().Assign()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("assign", fid)
	f, err := ParseFid(fid)
	if err != nil {
		t.Fatal(err)
	}
	if f.String() != fid {
		t.Error("Fid.String() not match")
	}
	t.Log("Fid.String()", f.String())
}

func TestAssginUpload(t *testing.T) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	fid, size, err := client.AssignUpload(filename, "text/plain", file)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("assign upload", filename, fid, size)
}

func TestMasterSubmit(t *testing.T) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	fid, size, err := client.Master().Submit(filename, "text/plain", file)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("master submit", filename, fid, size)
}

func TestGetUrl(t *testing.T) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	fid, _, err := client.Master().Submit(filename, "text/plain", file)
	if err != nil {
		t.Fatal(err)
	}
	publicUrl, url, err := client.GetUrl(fid)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("publicUrl:", publicUrl, "url:", url)
}

func TestDelete(t *testing.T) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	fid, size, err := client.Master().Submit(filename, "text/plain", file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("submit", fid, size)
	if err := client.Delete(fid, 1); err != nil {
		t.Fatal(err)
	}
	t.Log(fid, "deleted")
}

func TestDeleteN(t *testing.T) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	fid, size, err := client.Master().Submit(filename, "text/plain", file)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("submit", fid, size)
	if err := client.Delete(fid, 3); err != nil {
		t.Fatal(err)
	}
	t.Log(fid, "deleted")
}

func TestFilerUpload(t *testing.T) {
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Filer("localhost:8888").Upload("text/world.txt", "text/plain", file); err != nil {
		t.Fatal(err)
	}
}

func TestFilerDelete(t *testing.T) {
	if err := client.Filer("localhost:8888").Delete("text/"); err != nil {
		t.Fatal(err)
	}
}

func TestFilerDir(t *testing.T) {
	dir, err := client.Filer("localhost:8888").Dir("/")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(dir)
}
