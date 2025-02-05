// Package scribble là một cơ sở dữ liệu JSON nhỏ gọn với các tính năng nâng cấp
package scribble

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jcelliott/lumber"
)

// Version là phiên bản hiện tại của dự án
const Version = "2.0.0"

var (
	ErrMissingResource   = errors.New("missing resource - unable to save record")
	ErrMissingCollection = errors.New("missing collection - no place to save record")
	ErrMissingID         = errors.New("missing ID - record must have an ID")
)

// Record là interface cho các đối tượng có thể lưu trữ
type Record interface {
	GetID() string
}

// Logger là interface cho logger
type Logger interface {
	Fatal(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
}

// Driver là đối tượng tương tác với cơ sở dữ liệu scribble
type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string // thư mục chứa database
	log     Logger // logger
}

// Options dùng để cấu hình scribble
type Options struct {
	Logger // logger được sử dụng
}

// New tạo một database scribble mới tại thư mục chỉ định
func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}
	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating scribble database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)
}

// Write ghi một record vào database trong collection chỉ định
func (d *Driver) Write(collection string, v Record) error {
	if collection == "" {
		return ErrMissingCollection
	}

	id := v.GetID()
	if id == "" {
		return ErrMissingID
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, id+".json")
	tmpPath := fnlPath + ".tmp"

	return write(dir, tmpPath, fnlPath, v)
}

// Read đọc một record từ database
func (d *Driver) Read(collection, id string, v interface{}) error {
	if collection == "" {
		return ErrMissingCollection
	}

	if id == "" {
		return ErrMissingID
	}

	record := filepath.Join(d.dir, collection, id)
	return read(record, v)
}

// ReadAll đọc tất cả records từ một collection
func (d *Driver) ReadAll(collection string) ([][]byte, error) {
	if collection == "" {
		return nil, ErrMissingCollection
	}

	dir := filepath.Join(d.dir, collection)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	return readAll(files, dir)
}

// Update cập nhật một record dựa trên ID
func (d *Driver) Update(collection string, v Record) error {
	return d.Write(collection, v)
}

// Delete xóa một record hoặc toàn bộ collection
func (d *Driver) Delete(collection, id string) error {
	path := filepath.Join(collection, id)
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory named %v", path)
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}

	return nil
}

// BaseRecord là struct cơ bản cho các record
type BaseRecord struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetID trả về ID của record
func (b BaseRecord) GetID() string {
	return b.ID
}

// Private functions

func write(dir, tmpPath, dstPath string, v interface{}) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, dstPath)
}

func read(record string, v interface{}) error {
	b, err := ioutil.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}

func readAll(files []os.FileInfo, dir string) ([][]byte, error) {
	var records [][]byte

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		records = append(records, b)
	}

	return records, nil
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json")
	}
	return
}

func (d *Driver) getOrCreateMutex(collection string) *sync.Mutex {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	m, ok := d.mutexes[collection]

	if !ok {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}

	return m
}
