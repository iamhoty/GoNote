package main

import (
	"fmt"
	"os"
)

//定义两个变量，一个表示用户id, 一个表示用户密码
var userId int
var userPwd string

func main() {
	// 接受用户的选择
	var key int
	// 判断是否继续聊天
	loop := true
	for loop {
		fmt.Println("----------------欢迎登陆多人聊天系统------------")
		fmt.Println("\t\t 1 登陆聊天室")
		fmt.Println("\t\t 2 注册用户")
		fmt.Println("\t\t 3 退出系统")
		fmt.Println("\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			loop = false
		case 2:
			fmt.Println("注册用户")
			loop = false
		case 3:
			fmt.Println("退出系统")
			//loop = false
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")

		}
	}
	if key == 1 {
		// 用户要登录聊天室
		fmt.Println("请输入账号id:")
		fmt.Scanln(&userId)
		fmt.Println("请输入密码:")
		fmt.Scanln(&userPwd)
		// 登陆函数写在login.go
		login(userId, userPwd)
	} else if key == 2 {
		fmt.Println("进行用户注册的逻辑....")
	}

}
