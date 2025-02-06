# Scribble Library

Thư viện cơ sở dữ liệu JSON đơn giản cho Go.

## Cấu trúc thư viện

```
scribble/
├── scribble.go    # Mã nguồn chính
└── README.md      # Tài liệu
```

## API Reference

### Types

#### Driver

```go
type Driver struct {
    // private fields
}
```

Driver là đối tượng chính để tương tác với database.

#### Record

```go
type Record struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

Record chứa các trường metadata của mỗi record.

### Functions

#### New

```go
func New(dir string) (*Driver, error)
```

Tạo một database mới tại thư mục chỉ định.

- **Tham số:**
  - `dir`: Đường dẫn thư mục chứa database
- **Trả về:**
  - `*Driver`: Driver để tương tác với database
  - `error`: Lỗi nếu có

#### Write

```go
func (d *Driver) Write(collection string, v interface{}) error
```

Ghi một record mới vào collection.

- **Tham số:**
  - `collection`: Tên collection
  - `v`: Dữ liệu cần ghi (struct hoặc map)
- **Trả về:**
  - `error`: Lỗi nếu có

#### Read

```go
func (d *Driver) Read(collection, id string, v interface{}) error
```

Đọc một record từ collection.

- **Tham số:**
  - `collection`: Tên collection
  - `id`: ID của record
  - `v`: Con trỏ đến biến nhận dữ liệu
- **Trả về:**
  - `error`: Lỗi nếu có

#### Update

```go
func (d *Driver) Update(collection string, id string, v interface{}) error
```

Cập nhật một record trong collection.

- **Tham số:**
  - `collection`: Tên collection
  - `id`: ID của record
  - `v`: Dữ liệu mới
- **Trả về:**
  - `error`: Lỗi nếu có

#### Delete

```go
func (d *Driver) Delete(collection, id string) error
```

Xóa một record từ collection.

- **Tham số:**
  - `collection`: Tên collection
  - `id`: ID của record
- **Trả về:**
  - `error`: Lỗi nếu có

## Errors

```go
var (
    ErrMissingCollection = errors.New("missing collection - no place to save record")
)
```

## Thread Safety

- Mỗi collection có một mutex riêng
- Các thao tác Write/Update/Delete được bảo vệ bởi mutex
- An toàn khi có nhiều goroutine truy cập cùng lúc

## Cấu trúc file

### Thư mục database

```
data/
└── posts/              # Collection
    ├── 1.json         # Record
    ├── 2.json
    └── ...
```

### Format JSON

```json
{
  "record": {
    "id": "1",
    "created_at": "2025-02-06T10:00:00+07:00",
    "updated_at": "2025-02-06T10:00:00+07:00"
  },
  "data": {
    // Dữ liệu người dùng
  }
}
```

## Best Practices

1. **Sử dụng con trỏ**

   ```go
   // Tốt
   db.Write("posts", &post)

   // Không tốt
   db.Write("posts", post)
   ```

2. **Xử lý lỗi**

   ```go
   if err := db.Write("posts", &post); err != nil {
       // Xử lý lỗi
   }
   ```

3. **Đóng/mở file**

   - Thư viện tự động xử lý
   - Không cần đóng/mở thủ công

4. **Collection**
   - Tên collection nên là số nhiều
   - Ví dụ: "posts", "users", "comments"

## Giới hạn

1. **Hiệu năng**

   - Mỗi record một file
   - Không có index
   - Không phù hợp cho dữ liệu lớn

2. **Tìm kiếm**

   - Chỉ tìm theo ID
   - Không hỗ trợ query phức tạp

3. **Concurrency**
   - Mutex cho mỗi collection
   - Có thể block khi nhiều write cùng lúc
