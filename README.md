# Scribble Demo - Bài tập môn Ứng dụng phân tán

Dự án này là một phần của bài tập môn Hệ thống phân tán, tập trung vào việc nghiên cứu và phát triển thêm tính năng cho thư viện [Scribble](https://github.com/sdomino/scribble) - một cơ sở dữ liệu JSON đơn giản.

## Cấu trúc dự án

Dự án được chia thành 2 phần chính:

### 1. Scribble v1 (Phiên bản gốc)

- Nằm trong thư mục `scribble-v1/`
- Triển khai cơ bản của một cơ sở dữ liệu JSON
- Sử dụng Title làm khóa chính
- Hỗ trợ các thao tác CRUD cơ bản

### 2. Scribble v2 (Phiên bản nâng cấp)

- Nằm trong thư mục `scribble-v2/`
- Bổ sung nhiều tính năng mới:
  - Sử dụng ID làm khóa chính
  - Thêm trường CreatedAt và UpdatedAt
  - Thêm chức năng Update
  - Cải thiện cấu trúc code với interface Record
  - Hỗ trợ đa ngôn ngữ trong comment

## Cách chạy demo

### Demo phiên bản gốc (v1)

```bash
cd scribble-v1
go mod tidy
go run demo.go
```

### Demo phiên bản nâng cấp (v2)

```bash
cd scribble-v2
go mod tidy
go run demo.go
```

## So sánh hai phiên bản

### Phiên bản gốc (v1)

1. **Cấu trúc dữ liệu**

   - Sử dụng Title làm khóa chính
   - Không có cấu trúc chuẩn cho các record
   - Không theo dõi thời gian

2. **Chức năng**
   - Write: Ghi dữ liệu với Title làm key
   - Read: Đọc dữ liệu theo Title
   - ReadAll: Đọc tất cả record trong collection
   - Delete: Xóa record hoặc collection

### Phiên bản nâng cấp (v2)

1. **Cấu trúc dữ liệu**

   - Sử dụng ID làm khóa chính
   - Có BaseRecord làm cấu trúc chuẩn
   - Theo dõi thời gian tạo/cập nhật

2. **Chức năng**
   - Write: Ghi dữ liệu với ID làm key
   - Read: Đọc dữ liệu theo ID
   - ReadAll: Đọc tất cả record trong collection
   - Update: Cập nhật record theo ID
   - Delete: Xóa record hoặc collection

## Đóng góp

Các tính năng đã phát triển thêm:

1. Thêm cơ chế quản lý ID và thời gian
2. Cải thiện cấu trúc code với interface và struct chuẩn
3. Thêm chức năng Update
4. Hỗ trợ đa ngôn ngữ trong comment

## Công nghệ sử dụng

- Ngôn ngữ: Go 1.23.0
