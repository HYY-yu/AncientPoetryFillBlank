package spider

import (
	"net/http"
	"net/url"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type GushiwenSpider struct {
}

func (*GushiwenSpider) GetContent(searchResult SearchResult) (string, error) {
	client := &http.Client{
	}

	u, _ := url.Parse(searchResult.Sid)
	q := u.Query()
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36`)
	req.Header.Add("DNT", "1")
	req.Header.Add("Host", "www.gushiwen.org")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Add("Referer", "http://www.gushiwen.org/shiwen/")

	resp, errDo := client.Do(req)
	if errDo != nil || resp.StatusCode != 200 {
		return "", errors.New("无法连接古诗词网 " + errDo.Error())
	}

	//解析HTML
	docm, errDoc := goquery.NewDocumentFromResponse(resp)
	if errDoc != nil {
		return "", errors.New("解析DOM错误 " + errDo.Error())
	}

	divContent := docm.Find("div .contson")
	content := divContent.First().Text()
	return content, nil
}

func (*GushiwenSpider) FindContent(name string, author string) (searchResult SearchResult, err error) {
	baseUrl := "http://www.gushiwen.org/shiwen2017/ajaxSearchSo.aspx"
	client := &http.Client{
	}

	u, _ := url.Parse(baseUrl)
	q := u.Query()
	q.Set("valueKey", name)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36`)
	req.Header.Add("DNT", "1")
	req.Header.Add("Host", "www.gushiwen.org")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8")
	req.Header.Add("Referer", "http://www.gushiwen.org/shiwen/")

	resp, errDo := client.Do(req)
	if errDo != nil || resp.StatusCode != 200 {
		err = errors.New("无法连接古诗词网 " + errDo.Error())
		return
	}

	//解析HTML
	docm, errDoc := goquery.NewDocumentFromResponse(resp)
	if errDoc != nil {
		err = errors.New("解析DOM错误 " + errDoc.Error())
		return
	}

	divMain := docm.Find(".main")
	//只找诗文部分
	divMain.EachWithBreak(func(i int, selection *goquery.Selection) bool {
		divLeft := selection.Find(".mleft")
		if divLeft.Text() == "诗文" {
			divRight := selection.Find(".mright")
			//遍历mright子节点
			divRight.Children().EachWithBreak(func(i int, selection *goquery.Selection) bool {
				oneListItem := selection.Text()
				findName, findAuthor := splitTheItem(oneListItem)
				if len(findName) == 0 {
					//不是这条 跳过
					return true
				}

				if findAuthor == author {
					//找到a标签的href属性
					a := selection.Find("a")
					href, exist := a.Attr("href")
					if exist {
						searchResult.Sid = href
						searchResult.DisplayName = findName
						searchResult.Author = findAuthor
						return false
					}
					return true
				}
				return true
			})
		}
		return true
	})

	if len(searchResult.Sid) == 0 {
		return searchResult, errors.New(" 未搜索到结果 -- " + name + " -- ")
	}

	return
}
func splitTheItem(s string) (string, string) {
	temp := strings.Split(s, " - ")
	if len(temp) != 2 {
		return "", ""
	}
	return strings.TrimSpace(temp[0]), strings.TrimSpace(temp[1])
}
