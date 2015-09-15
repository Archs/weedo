// Package timekey implements a customized Fid to get/set time/mime/size info in fid for seaweedfs
package timekey

import (
	"errors"
	"fmt"
	"mime"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	defaultMimeType = "application/octet-stream"
)

// >>> "%x"%(500*365*24*60*60*1000) # 50 years in nano seconds
// 'e574609f000' # 5.5bytes 44bits

type Fid struct {
	Id     uint32 // volume id
	Key    uint64 // file key for volume
	Cookie uint32 // cookie
}

func NewFid(volId, fullPath string) (*Fid, error) {
	fid := new(Fid)
	if id, err := strconv.ParseUint(volId, 10, 32); err != nil {
		return nil, err
	} else {
		fid.Id = uint32(id)
	}
	if err := fid.InsertKeyAndCookie(fullPath); err != nil {
		return nil, err
	}
	return fid, nil
}

func (f *Fid) VolumeID() string {
	return strconv.Itoa(int(f.Id))
}

// Fid in string form
func (f *Fid) String() string {
	return fmt.Sprintf("%d,%x%08x", f.Id, f.Key, f.Cookie)
}

// Set Fid.Key to current nano seconds since 1970's
// which is monotonouse increased :)
// Let's hope this key would not collides with each other
func (f *Fid) insertKey() {
	f.Key = uint64(time.Now().UnixNano())
}

// Set Fid.Cookie(32 bits) according to the file infomation
// Fid.Cookie contains the mime type info and file size(in KB)
// midx takes the left 10 bits   mask: 0xffc00000
// size tekes the rigth 22 bits  mask: 0x003fffff
func (f *Fid) insertCookie(fullPath string) error {
	// cookie
	mtype := mime.TypeByExtension(path.Ext(fullPath))
	idx, ok := mmap[mtype]
	if !ok {
		idx = mmap[defaultMimeType]
	}
	midx := uint32(idx) << 22
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}
	size := uint32(info.Size() / 1024) // count in KB
	if size == 0 {
		size = 1
	}
	f.Cookie = midx + size
	return nil
}

func (f *Fid) InsertKeyAndCookie(fullPath string) error {
	f.insertKey()
	return f.insertCookie(fullPath)
}

func (f *Fid) Time() time.Time {
	return time.Unix(0, int64(f.Key))
}

// mime type for this fid
func (f *Fid) MimeType() string {
	midx := f.Cookie & 0xffc00000
	midx = midx >> 22
	return mslice[midx]
}

// Size in KB
func (f *Fid) Size() int {
	return int(f.Cookie & 0x003fffff)
}

func ParseFid(s string) (*Fid, error) {
	fid := new(Fid)
	a := strings.Split(s, ",")
	if len(a) != 2 || len(a[1]) <= 8 {
		return nil, errors.New("Fid format invalid")
	}
	// id
	id, err := strconv.ParseUint(a[0], 10, 32)
	if err != nil {
		return nil, err
	}
	fid.Id = uint32(id)
	// key
	index := len(a[1]) - 8
	if fid.Key, err = strconv.ParseUint(a[1][:index], 16, 64); err != nil {
		return nil, err
	}
	// cookie
	if cookie, err := strconv.ParseUint(a[1][index:], 16, 32); err != nil {
		return nil, err
	} else {
		fid.Cookie = uint32(cookie)
	}
	return fid, nil
}
