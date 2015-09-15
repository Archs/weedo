package fidcookies

import (
	"mime"
	"os"
	"path"
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

// Fid in string form
func (f *Fid) String() string {
	return fmt.Sprintf("%d,%x%8x", f.Id, f.Key, f.Cookie)
}

// Set Fid.Key to current nano seconds since 1970's
// which is monotonouse increased :)
// Let's hope this key would not collides with each other
func (f *Fid) SetTimeKey() {
	f.Key = time.Now().UnixNano()
}

// Set Fid.Cookie(32 bits) according to the file infomation
// Fid.Cookie contains the mime type info and file size(in KB)
// midx takes the left 10 bits
// size tekes the rigth 22 bits
func (f *Fid) SetCookie(fullPath string) error {
	mtype := mime.TypeByExtension(path.Ext(fullPathm))
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

func ParseFid(fid string) (*Fid, error) {
	a := strings.Split(s, ",")
	if len(a) != 2 || len(a[1]) <= 8 {
		return fid, errors.New("Fid format invalid")
	}
	if fid.Id, err = strconv.ParseUint(a[0], 10, 32); err != nil {
		return
	}
	index := len(a[1]) - 8
	if fid.Key, err = strconv.ParseUint(a[1][:index], 16, 64); err != nil {
		return
	}
	if fid.Cookie, err = strconv.ParseUint(a[1][index:], 16, 32); err != nil {
		return
	}

	return
}
