package main

import (
	"fmt"
)

func main() {
	var t int
	fmt.Scan(&t) // 读取测试用例的数量

	for i := 0; i < t; i++ {
		var a, b int
		fmt.Scan(&a, &b) // 读取两个整数
		fmt.Println(a + b) // 输出它们的和
	}
}