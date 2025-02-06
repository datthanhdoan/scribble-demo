// Package scribble là một cơ sở dữ liệu JSON nhỏ gọn
package scribble

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Version là phiên bản hiện tại của dự án
const Version = "2.0.0"

var (
	ErrMissingCollection = errors.New("missing collection - no place to save record")
)

// Driver là đối tượng tương tác với cơ sở dữ liệu scribble
type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string           // thư mục chứa database
	counter map[string]int64 // bộ đếm cho từng collection
}

// Record là struct cơ bản cho các record
type Record struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// New tạo một database scribble mới tại thư mục chỉ định
func New(dir string) (*Driver, error) {
	dir = filepath.Clean(dir)

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		counter: make(map[string]int64),
	}

	// Khôi phục counter từ các file hiện có
	if err := driver.initializeCounters(); err != nil {
		return nil, err
	}

	if _, err := os.Stat(dir); err == nil {
		return &driver, nil
	}

	return &driver, os.MkdirAll(dir, 0755)
}

// initializeCounters khôi phục bộ đếm từ các file hiện có
func (d *Driver) initializeCounters() error {
	collections, err := ioutil.ReadDir(d.dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, col := range collections {
		if !col.IsDir() {
			continue
		}

		files, err := ioutil.ReadDir(filepath.Join(d.dir, col.Name()))
		if err != nil {
			continue
		}

		var maxID int64
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			// Lấy ID từ tên file (bỏ phần .json)
			idStr := strings.TrimSuffix(f.Name(), ".json")
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				if id > maxID {
					maxID = id
				}
			}
		}
		d.counter[col.Name()] = maxID
	}
	return nil
}

// Write ghi một record vào database trong collection chỉ định
func (d *Driver) Write(collection string, v interface{}) error {
	if collection == "" {
		return ErrMissingCollection
	}

	// Tạo record mới với ID tự động tăng và thời gian
	record := &Record{
		ID:        fmt.Sprintf("%d", d.nextID(collection)),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Ghi record và data vào file
	data := map[string]interface{}{
		"record": record,
		"data":   v,
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, record.ID+".json")
	tmpPath := fnlPath + ".tmp"

	return write(dir, tmpPath, fnlPath, data)
}

// Read đọc một record từ database
func (d *Driver) Read(collection, id string, v interface{}) error {
	if collection == "" {
		return ErrMissingCollection
	}

	record := filepath.Join(d.dir, collection, id)

	// Đọc dữ liệu từ file
	var data struct {
		Record *Record     `json:"record"`
		Data   interface{} `json:"data"`
	}
	data.Data = v

	if err := read(record, &data); err != nil {
		return err
	}

	return nil
}

// Update cập nhật một record trong database
func (d *Driver) Update(collection string, id string, v interface{}) error {
	if collection == "" {
		return ErrMissingCollection
	}

	// Đọc record cũ
	var data struct {
		Record *Record     `json:"record"`
		Data   interface{} `json:"data"`
	}
	record := filepath.Join(d.dir, collection, id)
	if err := read(record, &data); err != nil {
		return err
	}

	// Cập nhật thời gian và data mới
	data.Record.UpdatedAt = time.Now()
	data.Data = v

	// Ghi lại vào file
	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, id+".json")
	tmpPath := fnlPath + ".tmp"

	return write(dir, tmpPath, fnlPath, data)
}

// Delete xóa một record từ database
func (d *Driver) Delete(collection, id string) error {
	if collection == "" {
		return ErrMissingCollection
	}

	mutex := d.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, collection)
	path := filepath.Join(dir, id+".json")

	return os.Remove(path)
}

// nextID tạo ID tiếp theo cho collection
func (d *Driver) nextID(collection string) int64 {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.counter[collection]++
	return d.counter[collection]
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
