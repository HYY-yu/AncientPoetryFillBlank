package spider

import (
	"fmt"
	"errors"
	"strings"
	"regexp"
)

type SearchResult struct {
	DisplayName string `mapstructure:"display_name"`
	Author      string `mapstructure:"literature_author"`
	Sid         string `mapstructure:"sid"`
}

type Spider interface {
	GetContent(SearchResult) (string, error)
	FindContent(string, string) (SearchResult, error)
}

//根据name author调用不同的Spider去获取数据
func GetContent(name, author string) (content string, err error) {
	spiderMap := MapSpiderManifest()
	for k, v := range spiderMap {
		searchResult, err := v.FindContent(name, author)
		if err != nil {
			fmt.Println(k, " 爬取出了点小问题～", err.Error())
			//换一个Spider
			continue
		}
		tempContent, err := v.GetContent(searchResult)
		if err != nil {
			fmt.Println(k, " 网页解析出了点小问题～", err.Error())
			continue
		}

		content = strings.TrimSpace(tempContent)
		// --- 整理格式 ---
		//去掉\n
		content = strings.Replace(content, "\n", "", -1)
		//去掉-
		content = strings.Replace(content, "—", "", -1)
		//去掉（）里的句子
		regBrackets := regexp.MustCompile("\\(.+?\\)")
		content = regBrackets.ReplaceAllString(content,"")
		// ---   END   ---
		return content, nil
	}

	//所有的Spider都不行
	return "", errors.New("无法爬取内容")
}

func MapSpiderManifest() map[string]Spider {
	//初始化并且注册所有的Spider
	spiderMap := make(map[string]Spider)

	//百度
	baiduSpider := new(BaiduSpider)
	spiderMap["baiduSpider"] = baiduSpider

	//古诗文网
	gushiwenSpider := new(GushiwenSpider)
	spiderMap["gushiwenSpider"] = gushiwenSpider

	return spiderMap
}
