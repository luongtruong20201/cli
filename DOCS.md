# Tài liệu dự án rm-cli

## Tổng quan

`rm-cli` là một thư viện CLI (Command Line Interface) được viết bằng Go, cung cấp một framework đơn giản và mạnh mẽ để xây dựng các ứng dụng dòng lệnh. Dự án này được thiết kế để giúp các nhà phát triển tạo ra các ứng dụng CLI một cách nhanh chóng và dễ dàng.

## Cấu trúc dự án

```
rm-cli/
├── app.go              # Cấu trúc và logic chính của ứng dụng
├── cli.go              # Định nghĩa handler và các biến global (đã comment)
├── command.go          # Cấu trúc và xử lý lệnh con
├── context.go          # Context để truyền dữ liệu giữa các thành phần
├── flag.go             # Định nghĩa và xử lý các flag/option
├── help.go             # Template và logic hiển thị help
├── options.go          # Cấu trúc Options để lưu trữ dữ liệu
├── main/               # Thư mục chứa file main
│   └── main.go         # Entry point của ứng dụng
├── go.mod              # File quản lý dependencies
└── README.md           # Tài liệu cơ bản
```

## Các thành phần chính

### 1. App (app.go)

Cấu trúc `App` là thành phần trung tâm của thư viện:

```go
type App struct {
    Name     string        // Tên ứng dụng
    Usage    string        // Mô tả cách sử dụng
    Version  string        // Phiên bản
    Commands []Command     // Danh sách các lệnh con
    Flags    []Flag        // Danh sách các flag global
    Action   func(context *Context) // Hàm xử lý mặc định
}
```

**Các phương thức chính:**
- `NewApp()`: Tạo một instance mới của App
- `Run(arguments []string)`: Chạy ứng dụng với các tham số đầu vào
- `Command(name string)`: Tìm kiếm lệnh con theo tên
- `hasFlag(flag Flag)`: Kiểm tra xem flag đã tồn tại chưa
- `appendFlag(flag Flag)`: Thêm flag mới vào danh sách

### 2. Command (command.go)

Cấu trúc `Command` đại diện cho một lệnh con:

```go
type Command struct {
    Name        string        // Tên lệnh
    ShortName   string        // Tên viết tắt
    Usage       string        // Mô tả cách sử dụng
    Description string        // Mô tả chi tiết
    Action      func(context *Context) // Hàm xử lý
    Flags       []Flag        // Danh sách flag của lệnh
}
```

**Các phương thức chính:**
- `Run(ctx *Context)`: Thực thi lệnh
- `HasName(name string)`: Kiểm tra tên lệnh (bao gồm cả short name)

### 3. Context (context.go)

Cấu trúc `Context` cung cấp interface để truy cập dữ liệu:

```go
type Context struct {
    App       *App           // Tham chiếu đến App
    flagSet   *flag.FlagSet  // FlagSet của lệnh hiện tại
    globalSet *flag.FlagSet  // FlagSet global
}
```

**Các phương thức chính:**
- `Int(name string)`: Lấy giá trị int từ flag
- `Bool(name string)`: Lấy giá trị bool từ flag
- `String(name string)`: Lấy giá trị string từ flag
- `GlobalInt(name string)`: Lấy giá trị int từ global flag
- `GlobalBool(name string)`: Lấy giá trị bool từ global flag
- `GlobalString(name string)`: Lấy giá trị string từ global flag
- `Args() []string`: Lấy danh sách arguments

### 4. Flag (flag.go)

Thư viện hỗ trợ các loại flag sau:

#### BoolFlag
```go
type BoolFlag struct {
    Name  string
    Usage string
}
```

#### StringFlag
```go
type StringFlag struct {
    Name  string
    Value string
    Usage string
}
```

#### IntFlag
```go
type IntFlag struct {
    Name  string
    Value int
    Usage string
}
```

### 5. Help System (help.go)

Hệ thống help tự động tạo ra các template:

- `AppHelpTemplate`: Template cho help của ứng dụng
- `CommandHelpTemplate`: Template cho help của lệnh con
- `helpCommand`: Lệnh help mặc định

**Các hàm chính:**
- `ShowAppHelp(c *Context)`: Hiển thị help của ứng dụng
- `ShowCommandHelp(c *Context, command string)`: Hiển thị help của lệnh
- `ShowVersion(c *Context)`: Hiển thị phiên bản

### 6. Options (options.go)

Cấu trúc `Options` để lưu trữ dữ liệu dạng key-value:

```go
type Options map[string]interface{}
```

**Các phương thức:**
- `Int(key string)`: Lấy giá trị int
- `String(key string)`: Lấy giá trị string
- `Bool(key string)`: Lấy giá trị bool

## Cách sử dụng

### 1. Tạo ứng dụng cơ bản

```go
package main

import (
    "fmt"
    "os"
    "cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "myapp"
    app.Usage = "Một ứng dụng CLI đơn giản"
    app.Version = "1.0.0"
    
    app.Action = func(c *cli.Context) {
        fmt.Println("Xin chào từ myapp!")
    }
    
    app.Run(os.Args)
}
```

### 2. Thêm flag global

```go
app.Flags = []cli.Flag{
    cli.BoolFlag{
        Name:  "verbose",
        Usage: "Hiển thị thông tin chi tiết",
    },
    cli.StringFlag{
        Name:  "config",
        Value: "config.json",
        Usage: "Đường dẫn đến file cấu hình",
    },
}
```

### 3. Thêm lệnh con

```go
app.Commands = []cli.Command{
    {
        Name:      "build",
        ShortName: "b",
        Usage:     "Build ứng dụng",
        Action: func(c *cli.Context) {
            fmt.Println("Đang clear
            clear
            clear
             ứng dụng...")
        },
    },
    {
        Name:      "deploy",
        ShortName: "d",
        Usage:     "Deploy ứng dụng",
        Action: func(c *cli.Context) {
            fmt.Println("Đang deploy ứng dụng...")
        },
    },
}
```

### 4. Xử lý flag trong Action

```go
app.Action = func(c *cli.Context) {
    if c.GlobalBool("verbose") {
        fmt.Println("Chế độ verbose được bật")
    }
    
    configPath := c.GlobalString("config")
    fmt.Printf("Sử dụng file cấu hình: %s\n", configPath)
    
    args := c.Args()
    if len(args) > 0 {
        fmt.Printf("Arguments: %v\n", args)
    }
}
```

## Tính năng

### 1. Tự động tạo help
- Help được tạo tự động cho ứng dụng và các lệnh con
- Hỗ trợ lệnh `help` hoặc `-h` để hiển thị thông tin trợ giúp

### 2. Quản lý phiên bản
- Tự động thêm flag `--version` để hiển thị phiên bản
- Hỗ trợ lệnh `version` để kiểm tra phiên bản

### 3. Xử lý lỗi
- Hiển thị thông báo lỗi rõ ràng khi sử dụng sai
- Tự động hiển thị help khi có lỗi

### 4. Hỗ trợ nhiều loại flag
- Bool flags: `--verbose`, `-v`
- String flags: `--config=path`, `-c=path`
- Int flags: `--port=8080`, `-p=8080`

## Ví dụ hoàn chỉnh

```go
package main

import (
    "fmt"
    "os"
    "cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "greet"
    app.Usage = "Chào hỏi người dùng"
    app.Version = "1.0.0"
    
    // Thêm flag global
    app.Flags = []cli.Flag{
        cli.BoolFlag{
            Name:  "formal",
            Usage: "Sử dụng cách chào trang trọng",
        },
        cli.StringFlag{
            Name:  "name",
            Value: "World",
            Usage: "Tên người cần chào",
        },
    }
    
    // Thêm lệnh con
    app.Commands = []cli.Command{
        {
            Name:      "hello",
            ShortName: "h",
            Usage:     "Chào hỏi thân thiện",
            Action: func(c *cli.Context) {
                name := c.GlobalString("name")
                if c.GlobalBool("formal") {
                    fmt.Printf("Xin chào, %s!\n", name)
                } else {
                    fmt.Printf("Chào %s!\n", name)
                }
            },
        },
        {
            Name:      "goodbye",
            ShortName: "g",
            Usage:     "Chào tạm biệt",
            Action: func(c *cli.Context) {
                name := c.GlobalString("name")
                fmt.Printf("Tạm biệt %s!\n", name)
            },
        },
    }
    
    // Action mặc định
    app.Action = func(c *cli.Context) {
        fmt.Println("Sử dụng 'greet help' để xem các lệnh có sẵn")
    }
    
    app.Run(os.Args)
}
```

## Build và chạy ứng dụng

### Build dự án

1. **Đảm bảo có Go 1.24.1 trở lên:**
```bash
go version
```

2. **Build ứng dụng:**
```bash
# Từ thư mục gốc của dự án
go build -o greet main/main.go
```

3. **Kiểm tra file đã được tạo:**
```bash
ls -la greet
```

### Chạy ứng dụng

```bash
# Hiển thị help
./greet --help

# Chào hỏi thân thiện
./greet hello

# Chào hỏi trang trọng
./greet --formal hello

# Chào với tên cụ thể
./greet --name="Alice" hello

# Chào tạm biệt
./greet --name="Bob" goodbye

# Hiển thị phiên bản
./greet --version
```

### Các lệnh build khác

```bash
# Build với tên file khác
go build -o myapp main/main.go

# Build và chạy trực tiếp
go run main/main.go --help

# Build cho các platform khác
GOOS=windows GOARCH=amd64 go build -o greet.exe main/main.go
GOOS=darwin GOARCH=amd64 go build -o greet-mac main/main.go
```

## Lưu ý

1. **Trạng thái dự án**: Theo README.md, dự án này đang trong giai đoạn phát triển (Work in Progress) và chưa sẵn sàng để release.

2. **Tương thích**: Dự án sử dụng Go 1.24.1 và có thể cần cập nhật để tương thích với các phiên bản Go mới hơn.

3. **Testing**: Dự án có đầy đủ các file test để đảm bảo chất lượng code.

4. **Mở rộng**: Có thể dễ dàng mở rộng thêm các loại flag mới hoặc tính năng khác.

## Kết luận

`rm-cli` là một thư viện CLI đơn giản nhưng mạnh mẽ, cung cấp tất cả các tính năng cần thiết để xây dựng ứng dụng dòng lệnh trong Go. Với cấu trúc rõ ràng và API dễ sử dụng, nó giúp các nhà phát triển tập trung vào logic nghiệp vụ thay vì xử lý các chi tiết kỹ thuật của CLI.
