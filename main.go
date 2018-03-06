package main

import (
	"fmt"
	"errors"
	"FillBlank/spider"
	"FillBlank/answer"
)

func main() {
	fmt.Println("请输入要查询的作品名称：")
	var bookName string
	fmt.Scan(&bookName)

	fmt.Println("请输入作者：")
	var author string
	fmt.Scan(&author)

	fmt.Println("要我默写哪句话？（格式 - （末尾要有标点）：___,则物与我皆无尽也，___！）")
	var findStr string
	fmt.Scan(&findStr)

	ans, errs := StartSearch(bookName, author, findStr)
	if errs != nil {
		fmt.Println("errrrrrrro:", errs.Error())
		return
	} else {
		showAnswer(ans)
	}

	for {
		fmt.Println("还有么？ (输入 kill 结束)")
		var findS string
		fmt.Scan(&findS)

		if len(findS) == 0 || findS == "kill" {
			break
		}
		ans, err := StartSearch(bookName, author, findS)
		if err != nil {
			fmt.Println("errrrrrrro:", err.Error())
			return
		} else {
			showAnswer(ans)
		}
	}

	fmt.Println("Bye~")
}

func showAnswer(ans string) {
	if len(ans) == 0 {
		fmt.Println("这个答不上来")
	} else {
		fmt.Println("so easy : ", ans)
	}
}

//findStr格式： （末尾要有标点）：_,则物与我皆无尽也，_！
func StartSearch(bookName string, author string, findStr string) (string, error) {
	contents, err := spider.GetContent(bookName,author)
	if err != nil {
		return "", errors.New("搜索失败 --" + err.Error())
	}

	fmt.Println("find content ",contents)

	ans, err := answer.FindTheAnswer(contents, findStr)
	return ans, err
}
