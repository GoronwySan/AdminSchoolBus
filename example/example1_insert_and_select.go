package example

import (
	"fmt"
	"log"
	"login/auth"
	"login/config"
	"login/db"
	"login/exception"
)

// 一、数据的插入

type TokenTest struct {
	TokenHash    string `db:"token_hash"`
	TokenRevoked bool   `db:"token_revoked"`
	TokenExpiry  string `db:"token_expiry"`
} // 请注意前面的字段名一定要**大写**，否则会触发**未导出错误**，后面的db名一定要正确，并且是蛇形，如果非蛇形程序会报错

func TestInsert() {
	fmt.Println("Warning, this is just an example, please do not use it in production")
	// 实例化一个token
	token := TokenTest{
		TokenHash:    "1234567890",
		TokenRevoked: false,
		TokenExpiry:  "2020-01-01 00:00:00",
	}

	// 插入的数据库是属于什么身份的，这里是admin
	// table叫什么 以及 你的数据
	index, err := db.Insert(config.RoleAdmin, "tokens", token)

	// 错误处理
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("insert success, index:", index)
	}

}

// 二、数据的查询
// ** 这里我们提供两个函数：
// 由于查询比较复杂，一个函数可能很难涵盖所有需求
// 1、Select函数是安全的查询函数，使用相对复杂，无法实现所有查询
// 2、SelectPrimitive也是安全的查询函数，使用简单，但需要用户自己构建查询语句，并处理语法错误问题
// 两者都只支持一个sql语句，如果你需要多个sql语句，请多次调用
// **在SelectPrimitive中，如果你的语句被过滤，请去db.go文件中的 containsUnsafeSQL 函数中**删除**过滤关键字

// 1、Select函数 使用

// **普通查询
func SelectTest1() {
	// 查询所有未撤销的 Token
	conditionFields := []string{"token_revoked = ?"}
	params := []interface{}{false} // 查询未撤销的 token

	// 你需要一个映射表的结构，然后在这里建一个**结构体_数组**，并通过**结构体指针**方式传参
	var tokens []auth.Token
	// 查询后数据会被填入原本的数组里面

	err := db.Select(
		config.RoleAdmin, // 数据库角色
		"tokens",         // 表名
		&tokens,          // 存放查询结果的结构体数组**指针**
		true,             // 查询所有列
		nil,              // 不指定查询字段，查询所有列
		conditionFields,  // 查询条件
		params,           // 查询条件的参数
		"",               // 无排序
		10,               // 限制返回10条记录
		0,                // 从第0条开始
		"",               // 无分组
		"",               // 无筛选
	)

	if err != nil {
		exception.PrintError(SelectTest1, err)
		log.Fatal("Error during select:", err)
	}
	// ** 数据被填入了原本的数组中
	fmt.Println("Query Result:")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s\n", token.TokenID, token.TokenHash, token.TokenExpiry)
	}
}

// **排序查询
func SelectTest2() {
	// 查询所有未撤销的 Token，按 TokenExpiry 排序
	conditionFields := []string{"token_revoked = ?"}
	params := []interface{}{false} // 查询未撤销的 token

	var tokens []auth.Token

	err := db.Select(
		config.RoleAdmin,   // 数据库角色
		"tokens",           // 表名
		&tokens,            // 存放查询结果的结构体数组指针
		true,               // 查询所有列
		nil,                // 不指定查询字段，查询所有列
		conditionFields,    // 查询条件
		params,             // 查询条件的参数
		"token_expiry ASC", // 按过期时间升序排序
		10,                 // 限制返回10条记录
		0,                  // 从第0条开始
		"",                 // 无分组
		"",                 // 无筛选
	)

	if err != nil {
		exception.PrintError(SelectTest2, err)
		log.Fatal("Error during select:", err)
	}

	fmt.Println("Query Result (Sorted by Expiry):")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s\n", token.TokenID, token.TokenHash, token.TokenExpiry)
	}
}

// ** 多条件
func SelectTest6() {
	// 多条件查询：查询未撤销且过期时间小于 2023-01-01 的 token
	conditionFields := []string{
		"token_revoked = ?", // 条件 1: 未撤销
		"token_expiry < ?",  // 条件 2: 过期时间小于某日期
	}
	params := []interface{}{
		false,                 // token_revoked = false
		"2023-01-01 00:00:00", // token_expiry < "2023-01-01"
	}

	var tokens []auth.Token

	err := db.Select(
		config.RoleAdmin, // 数据库角色
		"tokens",         // 表名
		&tokens,          // 存放查询结果的结构体数组指针
		true,             // 查询所有列
		nil,              // 不指定查询字段，查询所有列
		conditionFields,  // 查询条件
		params,           // 查询条件的参数
		"",               // 无排序
		10,               // 限制返回10条记录
		0,                // 从第0条开始
		"",               // 无分组
		"",               // 无筛选
	)

	if err != nil {
		exception.PrintError(SelectTest6, err)
		log.Fatal("Error during select:", err)
	}

	fmt.Println("Query Result (Multi-condition Query):")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s\n", token.TokenID, token.TokenHash, token.TokenExpiry)
	}
}

// ** 范围查询
func SelectTest7() {
	// 查询过期时间在两个日期范围内的 token
	conditionFields := []string{
		"token_expiry BETWEEN ? AND ?", // 条件: 过期时间在某个范围内
	}
	params := []interface{}{
		"2022-01-01 00:00:00", // 从 2022-01-01 开始
		"2023-01-01 00:00:00", // 到 2023-01-01
	}

	var tokens []auth.Token

	err := db.Select(
		config.RoleAdmin, // 数据库角色
		"tokens",         // 表名
		&tokens,          // 存放查询结果的结构体数组指针
		true,             // 查询所有列
		nil,              // 不指定查询字段，查询所有列
		conditionFields,  // 查询条件
		params,           // 查询条件的参数
		"",               // 无排序
		10,               // 限制返回10条记录
		0,                // 从第0条开始
		"",               // 无分组
		"",               // 无筛选
	)

	if err != nil {
		exception.PrintError(SelectTest7, err)
		log.Fatal("Error during select:", err)
	}

	fmt.Println("Query Result (Range Query):")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s\n", token.TokenID, token.TokenHash, token.TokenExpiry)
	}
}

// **条件组合查询 (AND + OR)
func SelectTest9() {
	// 查询条件：未撤销的 token 或过期时间早于某日期的 token
	conditionFields := []string{
		"(token_revoked = ? OR token_expiry < ?)", // 使用 OR 组合条件
	}
	params := []interface{}{
		false,                 // token_revoked = false
		"2023-01-01 00:00:00", // token_expiry < "2023-01-01"
	}

	var tokens []auth.Token

	err := db.Select(
		config.RoleAdmin, // 数据库角色
		"tokens",         // 表名
		&tokens,          // 存放查询结果的结构体数组指针
		true,             // 查询所有列
		nil,              // 不指定查询字段，查询所有列
		conditionFields,  // 查询条件
		params,           // 查询条件的参数
		"",               // 无排序
		10,               // 限制返回10条记录
		0,                // 从第0条开始
		"",               // 无分组
		"",               // 无筛选
	)

	if err != nil {
		exception.PrintError(SelectTest9, err)
		log.Fatal("Error during select:", err)
	}

	fmt.Println("Query Result (AND + OR):")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s\n", token.TokenID, token.TokenHash, token.TokenExpiry)
	}
}

// ** 模糊查询
func SelectTest10() {
	// 查询 TokenHash 含有特定子串的记录
	conditionFields := []string{
		"token_hash LIKE ?", // 模糊查询条件
	}
	params := []interface{}{
		"%abc%", // 查询 TokenHash 中含有 "abc" 的记录
	}

	var tokens []auth.Token

	err := db.Select(
		config.RoleAdmin, // 数据库角色
		"tokens",         // 表名
		&tokens,          // 存放查询结果的结构体数组指针
		true,             // 查询所有列
		nil,              // 不指定查询字段，查询所有列
		conditionFields,  // 查询条件
		params,           // 查询条件的参数
		"",               // 无排序
		10,               // 限制返回10条记录
		0,                // 从第0条开始
		"",               // 无分组
		"",               // 无筛选
	)

	if err != nil {
		exception.PrintError(SelectTest10, err)
		log.Fatal("Error during select:", err)
	}

	fmt.Println("Query Result (LIKE Query):")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s\n", token.TokenID, token.TokenHash, token.TokenExpiry)
	}
}

// ** 多表查询
func SelectTest8() {
	type tem struct {
		UserID        string `db:"user_id"`
		UserRegisDate string `db:"user_registry_date"`
	}

	var tems []tem

	params := []interface{}{"1"}

	// 检测role和user_id是否匹配，正式检查
	err := db.Select(config.RoleAdmin, "usersPass p, usersInfo i", &tems,
		false, []string{"p.user_id", "user_registry_date"}, []string{"p.user_id = i.user_id", "p.user_id = ?"},
		params, "", 9999, 0, "", "")

	if err != nil {
		exception.PrintError(SelectTest8, err)
		log.Fatal("Error during select:", err)
	}

	for _, tem := range tems {
		fmt.Printf("User ID: %s, Registry Date: %s\n", tem.UserID, tem.UserRegisDate)
	}
}

// 2、SelectPrimitive函数 使用

func SelectTest11() {
	// 构建 SQL 查询，查询所有 Token
	sqlQuery := "SELECT token_id, token_hash, token_revoked, token_expiry FROM tokens"
	params := []interface{}{} // 无条件的查询

	// 存放查询结果的结构体数组
	var tokens []auth.Token

	// 执行查询并将结果存入 tokens 数组
	err := db.SelectPrimitive(config.RoleAdmin, sqlQuery, params, &tokens)
	if err != nil {
		exception.PrintError(SelectTest11, err)
		log.Fatal("Error during select:", err)
	}

	fmt.Println("Query Result (All Tokens):")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s, Revoked: %v\n",
			token.TokenID, token.TokenHash, token.TokenExpiry, token.TokenRevoked)
	}
}

func SelectTest13() {
	// 构建 SQL 查询，查询 TokenID 在多个值中的 Token
	sqlQuery := "SELECT token_id, token_hash, token_revoked, token_expiry FROM tokens WHERE token_id = (?) AND token_expiry < (?)"
	params := []interface{}{"12345", "2023-01-01 00:00:00"} // 条件：token_id 在这几个值中

	// 存放查询结果的结构体数组
	var tokens []auth.Token

	// 执行查询并将结果存入 tokens 数组
	err := db.SelectPrimitive(config.RoleAdmin, sqlQuery, params, &tokens)
	if err != nil {
		exception.PrintError(SelectTest13, err)
		log.Fatal("Error during select:", err)
	}

	fmt.Println("Query Result (IN Query):")
	for _, token := range tokens {
		fmt.Printf("Token ID: %s, Hash: %s, Expiry: %s, Revoked: %v\n",
			token.TokenID, token.TokenHash, token.TokenExpiry, token.TokenRevoked)
	}
}
