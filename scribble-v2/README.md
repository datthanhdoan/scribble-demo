# Scribble v2 - Demo

Đây là phiên bản nâng cấp của thư viện Scribble, bổ sung thêm nhiều tính năng mới.

## Cách chạy demo

1. Mở terminal trong thư mục này
2. Chạy lệnh:

```bash
go mod tidy
go run demo.go
```

## Tính năng mới

- Sử dụng ID làm khóa chính thay vì Title
- Thêm trường CreatedAt và UpdatedAt để theo dõi thời gian
- Thêm chức năng Update để cập nhật record
- Cải thiện cấu trúc code với interface Record và struct BaseRecord
- Hỗ trợ đa ngôn ngữ trong comment
