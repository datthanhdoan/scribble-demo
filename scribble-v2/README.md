# Scribble v2

Thư viện cơ sở dữ liệu JSON đơn giản cho Go, với các tính năng tự động quản lý ID và thời gian.

## Cài đặt

```bash
go get github.com/sdomino/scribble/v2
```

## Cách sử dụng

### 1. Khởi tạo database

```go
db, err := scribble.New("data")
if err != nil {
    log.Fatal(err)
}
```

### 2. Định nghĩa struct cho dữ liệu

```go
type Post struct {
    Title string `json:"title"`
    Body  string `json:"body"`
}
```

### 3. Các thao tác cơ bản

#### Tạo mới

```go
post := Post{
    Title: "Tiêu đề bài viết",
    Body:  "Nội dung bài viết",
}

err := db.Write("posts", &post)
// ID sẽ tự động được tạo: 1, 2, 3,...
```

#### Đọc

```go
var post Post
err := db.Read("posts", "1", &post)
```

#### Cập nhật

```go
post.Body = "Nội dung mới"
err := db.Update("posts", "1", &post)
// UpdatedAt sẽ tự động được cập nhật
```

#### Xóa

```go
err := db.Delete("posts", "1")
```

## Cấu trúc dữ liệu

Mỗi record được lưu trong một file JSON riêng biệt với cấu trúc:

```json
{
  "record": {
    "id": "1",
    "created_at": "2025-02-06T10:00:00+07:00",
    "updated_at": "2025-02-06T10:00:00+07:00"
  },
  "data": {
    "title": "Tiêu đề bài viết",
    "body": "Nội dung bài viết"
  }
}
```

## Tính năng

1. **Tự động quản lý ID**

   - ID tăng dần từ 1
   - Tự động khôi phục counter khi restart

2. **Tự động quản lý thời gian**

   - `created_at`: Thời điểm tạo record
   - `updated_at`: Thời điểm cập nhật gần nhất

3. **Cấu trúc thư mục**

   ```
   data/
   └── posts/           # Tên collection
       ├── 1.json       # Record với ID = 1
       ├── 2.json       # Record với ID = 2
       └── ...
   ```

## Lưu ý

1. **Hiệu năng**

   - Phù hợp cho ứng dụng nhỏ
   - Không nên dùng cho dữ liệu lớn
   - Không có index hoặc tìm kiếm nâng cao

2. **Backup**

   - Dữ liệu được lưu dưới dạng file JSON
   - Dễ dàng sao lưu và phục hồi
   - Có thể chỉnh sửa trực tiếp file JSON

3. **ID và thời gian**
   - ID là số tăng dần, không thể tùy chỉnh
   - Thời gian tự động cập nhật
   - Định dạng thời gian: RFC3339

## Ví dụ

Xem thêm ví dụ trong file [demo.go](demo.go)
