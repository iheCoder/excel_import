# Excel Import Framework

一个功能强大的Go语言Excel数据导入框架，支持一般表格数据和树形结构数据的导入、验证、转换和持久化。

## 主要特性

- **多格式支持**：支持Excel (.xlsx) 和 CSV 文件格式
- **结构化导入**：通过标签系统映射Excel列到Go结构体
- **数据验证**：内置多种格式校验器，支持自定义校验
- **数据转换**：单元格内容格式化和转换
- **树形数据处理**：专门的树形结构导入框架
- **并发处理**：支持并行导入提高性能
- **错误处理**：详细的错误报告和记录
- **进度报告**：导入过程进度跟踪
- **数据库集成**：与GORM无缝对接，支持MySQL等数据库
- **可扩展性**：中间件系统和钩子支持

## 安装

```bash
go get github.com/yourusername/excel_import
```

## 快速开始

### 1. 一般表格数据导入

```go
package main

import (
    "excel_import"
    "excel_import/general_framework"
    "excel_import/utils"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

// 定义数据模型（使用结构体标签映射Excel列）
type Product struct {
    ID       int    `excel:"index:0,chk:on,fcf:int"` // 第0列，开启检查，整数类型
    Name     string `excel:"index:1,fcf:cn"`         // 第1列，中文检查
    Price    float64 `excel:"index:2,fcf:float"`     // 第2列，浮点数检查
    ImageURL string `excel:"index:3,fcf:img"`        // 第3列，图片URL检查
}

// 实现RowModelFactory接口
type ProductFactory struct{}

func (pf *ProductFactory) MinColumnCount() int {
    return 4 // 至少需要4列
}

func (pf *ProductFactory) GetModel() any {
    return &Product{}
}

func main() {
    // 连接数据库
    db, err := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4"), &gorm.Config{})
    if err != nil {
        panic("数据库连接失败")
    }
    
    // 迁移模式
    db.AutoMigrate(&Product{})
    
    // 创建导入框架
    framework := general_framework.NewOneSheetOneModelFramework(
        db,
        &ProductFactory{},
        general_framework.WithControl(&general_framework.DefaultImportControl{}),
    )
    
    // 执行导入
    err = framework.Import("path/to/products.xlsx")
    if err != nil {
        panic("导入失败: " + err.Error())
    }
    
    println("导入成功！")
}
```

### 2. 树形结构数据导入

```go
package main

import (
    "excel_import"
    "excel_import/tree_framework"
    "excel_import/utils"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

// 定义树节点模型
type Category struct {
    ID       int    `excel:"index:0,chk:on,fcf:int"`
    Name     string `excel:"index:1,fcf:cn"`
    ParentID int    `excel:"ctx:parent_id"`
    Level    int
}

// 实现RowModelFactory接口
type CategoryFactory struct{}

func (cf *CategoryFactory) MinColumnCount() int {
    return 2
}

func (cf *CategoryFactory) GetModel() any {
    return &Category{}
}

// 实现节点导入器
type CategoryImporter struct{}

func (ci *CategoryImporter) Import(db *gorm.DB, node *tree_framework.TreeNode, path string) error {
    // 根据TreeNode信息创建Category实例并保存
    category, ok := node.Model.(*Category)
    if !ok {
        return fmt.Errorf("模型类型错误")
    }
    
    // 设置层级信息
    category.Level = node.Level
    
    // 保存到数据库
    return db.Create(category).Error
}

func main() {
    // 连接数据库
    db, err := gorm.Open(mysql.Open("user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4"), &gorm.Config{})
    if err != nil {
        panic("数据库连接失败")
    }
    
    // 迁移模式
    db.AutoMigrate(&Category{})
    
    // 创建树形导入框架
    importer := &CategoryImporter{}
    framework := tree_framework.NewTreeImportStrictOrderFramework(
        db,
        3,              // 树最大深度为3
        2,              // 列数为2
        &CategoryFactory{},
        importer,
    )
    
    // 执行导入
    err = framework.Import("path/to/categories.xlsx")
    if err != nil {
        panic("导入失败: " + err.Error())
    }
    
    println("树形数据导入成功！")
}
```

## 数据验证

框架内置多种数据格式验证器，通过`fcf`标签指定：

- `int`: 整数验证
- `float`: 浮点数验证
- `url`: URL地址验证
- `img`: 图片URL验证
- `cn`: 中文字符验证
- `en`: 英文字符验证
- `pinyin`: 拼音验证
- `hash`: 哈希值验证

示例：
```go
type Product struct {
    ID       int     `excel:"index:0,fcf:int"`
    Price    float64 `excel:"index:1,fcf:float"`
    ImageURL string  `excel:"index:2,fcf:img"`
}
```

## 结构体标签说明

在模型定义中可使用以下标签：

- `index`: 指定Excel中的列索引（从0开始）
- `rewrite`: 是否允许重写该字段（布尔值）
- `chk`: 是否开启该字段的检查
- `ctx`: 标记上下文角色，如`parent_id`用于树形结构
- `fcf`: 格式检查函数，指定数据验证规则
- `id`: 用于标识或链接的ID

示例：
```go
type User struct {
    UserID   int    `excel:"index:0,chk:on,fcf:int"`
    Username string `excel:"index:1,fcf:cn"`
    Email    string `excel:"index:2,fcf:url"`
    Role     string `excel:"index:3,id:role_id"`
}
```

## 高级功能

### 1. 自定义中间件

```go
func MyMiddleware(next general_framework.SectionImportHandler) general_framework.SectionImportHandler {
    return func(tx *gorm.DB, rows [][]string) error {
        // 前置处理
        fmt.Println("导入前处理...")
        
        // 调用下一个处理器
        err := next(tx, rows)
        
        // 后置处理
        fmt.Println("导入后处理...")
        
        return err
    }
}

// 注册中间件
framework := general_framework.NewOneSheetOneModelFramework(
    db,
    &ProductFactory{},
    general_framework.WithMiddlewares([]general_framework.GeneralMiddleware{MyMiddleware}),
)
```

### 2. 正确性检查器

```go
type StockChecker struct {
    Products map[int]int // 产品库存映射
}

func (sc *StockChecker) PreCollect(tx *gorm.DB) error {
    // 预收集数据
    var products []Product
    if err := tx.Find(&products).Error; err != nil {
        return err
    }
    
    sc.Products = make(map[int]int)
    for _, p := range products {
        sc.Products[p.ID] = p.Stock
    }
    
    return nil
}

func (sc *StockChecker) CheckCorrect(tx *gorm.DB) error {
    // 检查逻辑
    for id, stock := range sc.Products {
        if stock < 0 {
            return fmt.Errorf("产品ID %d 库存为负数: %d", id, stock)
        }
    }
    
    return nil
}

// 注册检查器
framework := general_framework.NewOneSheetOneModelFramework(
    db,
    &ProductFactory{},
    general_framework.WithCorrectnessCheckers([]excel_import.CorrectnessChecker{&StockChecker{}}),
)
```

### 3. 进度报告

```go
reporter := util.NewProgressReporter(func(percent float64, message string) {
    fmt.Printf("进度: %.2f%% - %s\n", percent*100, message)
})

framework := general_framework.NewOneSheetOneModelFramework(
    db,
    &ProductFactory{},
    general_framework.WithProgressReporter(reporter),
)
```

## 错误处理

框架提供了详细的错误记录功能：

```go
// 创建错误记录器
recorder := util.NewUnexpectedRecorder("import_failed.csv")

framework := general_framework.NewOneSheetOneModelFramework(
    db,
    &ProductFactory{},
    general_framework.WithRecorder(recorder),
)

// 导入后检查错误
if recorder.HasUnexpected() {
    fmt.Printf("导入过程中有 %d 个错误，详情请查看 import_failed.csv\n", recorder.Count())
}
```

## 注意事项

1. 确保Excel文件的列结构与模型定义匹配
2. 树形结构导入需要明确定义父子关系
3. 大文件导入可能需要调整数据库连接参数
4. 使用事务确保数据一致性
5. 建议在导入前备份数据库

## 许可证

MIT

## 贡献

欢迎提交Pull Request或Issue！
