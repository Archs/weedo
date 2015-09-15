package main

import (
	"flag"
	"github.com/Archs/weedo"
	"log"
	"os"
	"path/filepath"
)

// -collection string
//       optional collection name
// -debug
//       verbose debug information
// -dir string
//       Upload the whole folder recursively if specified.
// -include string
//       pattens of files to upload, e.g., *.pdf, *.html, ab?d.txt, works togethe
// with -dir
// -maxMB int
//       split files larger than the limit
// -replication string
//       replication type
// -secure.secret string
//       secret to encrypt Json Web Token(JWT)
// -server string
//       SeaweedFS master location (default "localhost:9333")
// -ttl string
//       time to live, e.g.: 1m, 1h, 1d, 1M, 1y
var (
	server      string
	debug       bool
	recursive   bool
	collection  string
	replication string
)

var (
	client *weedo.Client
	fmap   = map[string]string{} // map fid -> filepath
)

func uploadFile(path string) error {
	log.Println("\t", path, "...")
	fid, err := client.AssignUploadTK(path)
	if err != nil {
		return err
	}
	fmap[fid] = path
	return nil
}

func uploadDirectory(dirPath string) error {
	log.Println("Uploading directory:", dirPath, "...")
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Walking into", path, "error:", err)
			return err
		}
		// for dir
		if info.IsDir() {
			err = uploadDirectory(path)
			log.Printf("Uploading dir:%s error:%s", path, err.Error())
			return nil
		}
		// for file
		err = uploadFile(path)
		if err != nil {
			log.Printf("Uploading file:%s error:%s", path, err.Error())
		}
		return nil

	}
	return filepath.Walk(dirPath, walkFn)
}

func main() {
	flag.StringVar(&server, "server", "http://localhost:9333", `SeaweedFS master location`)
	flag.StringVar(&collection, "col", "", `optional collection name`)
	flag.StringVar(&replication, "replication", "", "replication type")
	flag.BoolVar(&debug, "debug", false, "verbose debug information")
	flag.BoolVar(&recursive, "r", false, `upload directory recursivly (default false)`)
	flag.Parse()
	client = weedo.NewClient(server)
	if err := client.Master().Status(); err != nil {
		log.Fatal("invalid client:", err)
	}
	// ok now
	if len(flag.Args()) <= 0 {
		log.Fatalln("no files or directories specified")
	}
	// do upload
	for _, fpath := range flag.Args() {
		info, err := os.Stat(fpath)
		if err != nil {
			log.Printf("Uploading %s error:%s", fpath, err.Error())
			continue
		}
		if info.IsDir() {
			err = uploadDirectory(fpath)
		} else {
			err = uploadFile(fpath)
		}
		if err != nil {
			log.Println("Uploading", fpath, "failed:", err.Error())
		}
	}
}