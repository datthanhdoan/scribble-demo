# Scribble Demo

Đây là demo cho thư viện Scribble - một cơ sở dữ liệu JSON đơn giản cho Go.

## Cấu trúc thư mục

```
scribble-demo/
├── scribble-v1/     # Phiên bản cũ (tham khảo)
└── scribble-v2/     # Phiên bản mới với các cải tiến
```

## Tính năng chính của Scribble v2

- Lưu trữ dữ liệu dưới dạng file JSON
- Tự động tạo ID và quản lý thời gian
- Thread-safe với mutex
- API đơn giản, dễ sử dụng
- Không cần cài đặt database server

## Hướng dẫn chạy demo

1. Chuyển vào thư mục scribble-v2:

```bash
cd scribble-v2
```

2. Chạy demo:

```bash
go run demo.go
```

Demo sẽ thể hiện các thao tác cơ bản:

- Tạo bài viết mới
- Đọc bài viết
- Cập nhật bài viết
- Xóa bài viết

## Tài liệu chi tiết

Xem thêm:

- [Hướng dẫn sử dụng Scribble v2](scribble-v2/README.md)
- [Mã nguồn thư viện](scribble-v2/library/scribble/)
