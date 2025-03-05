# 初始化 Git 仓库
git init

# 创建 go.mod
go mod init github.com/kenneth-wang/go-demo

# 创建 .gitignore
cat <<EOF > .gitignore
# Go build artifacts
*.exe
*.out
*.test

# Go modules
/go/bin/
/go/pkg/
/go/src/
vendor/

# Logs & temp files
*.log
*.swp
EOF

# 创建 README.md
cat <<EOF > README.md
# Go Demo Repository
This repository contains various Go demos categorized by topic.
Each demo is independent and can be executed separately.
EOF

# 创建目录结构
mkdir -p basic/hello-world
mkdir -p basic/variables
mkdir -p basic/functions
mkdir -p concurrency/goroutines
mkdir -p concurrency/channels
mkdir -p web/http-server
mkdir -p web/gin-demo
mkdir -p database/mysql-demo
mkdir -p database/postgres-demo
mkdir -p cli/cobra-demo
mkdir -p testing/unit-tests
mkdir -p testing/benchmarks
mkdir -p advanced/reflection
mkdir -p advanced/generics
mkdir -p tools/logger
mkdir -p tools/config-loader

# 创建 Hello World 示例
cat <<EOF > basic/hello-world/main.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go Demo!")
}
EOF

cat <<EOF > basic/hello-world/README.md
# Hello World Demo
This is a simple Hello World example in Go.
Run it with:
\`\`\`sh
go run main.go
\`\`\`
EOF

# 创建 Goroutine 示例
cat <<EOF > concurrency/goroutines/main.go
package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from Goroutine!")
}

func main() {
	go sayHello()
	time.Sleep(1 * time.Second)
	fmt.Println("Main function finished!")
}
EOF

cat <<EOF > concurrency/goroutines/README.md
# Goroutine Demo
This demo shows how to create and run a simple goroutine.
\`\`\`sh
go run main.go
\`\`\`
EOF

# 创建 HTTP Server 示例
cat <<EOF > web/http-server/main.go
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, this is a simple HTTP server in Go!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
EOF

cat <<EOF > web/http-server/README.md
# HTTP Server Demo
This demo creates a simple HTTP server in Go.
\`\`\`sh
go run main.go
\`\`\`
Then open \`http://localhost:8080\` in your browser.
EOF

# 创建 Gin Web 框架示例
cat <<EOF > web/gin-demo/main.go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Gin!"})
	})
	r.Run(":8080")
}
EOF

cat <<EOF > web/gin-demo/README.md
# Gin Framework Demo
This demo shows how to create a web server using Gin.
\`\`\`sh
go run main.go
\`\`\`
Then open \`http://localhost:8080\` in your browser.
EOF

# 创建数据库 MySQL 连接示例
cat <<EOF > database/mysql-demo/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/testdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Connected to MySQL!")
}
EOF

cat <<EOF > database/mysql-demo/README.md
# MySQL Connection Demo
This demo connects to a MySQL database.
\`\`\`sh
go run main.go
\`\`\`
Make sure MySQL is running and replace the credentials accordingly.
EOF

# 创建 Cobra CLI 工具示例
cat <<EOF > cli/cobra-demo/main.go
package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "CLI App Example",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello from CLI!")
		},
	}

	rootCmd.Execute()
}
EOF

cat <<EOF > cli/cobra-demo/README.md
# Cobra CLI Demo
This demo shows how to create a CLI tool using Cobra.
\`\`\`sh
go run main.go
\`\`\`
EOF

# 提交到 Git
git add .
git commit -m "Initial commit: created structured Go demo repository"

# 输出成功信息
echo "✅ Go Demo Repository initialized successfully!"
echo "👉 Now you can start coding in ~/Projects/go-demo"
