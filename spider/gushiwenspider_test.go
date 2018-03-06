package spider

import (
	"testing"
	"fmt"
)

func TestGushiwenSpider_FindContent(t *testing.T) {
	g := &GushiwenSpider{}
	gotSearchResult, err := g.FindContent("南", "辛弃疾")
	g.GetContent(gotSearchResult)
	if err != nil{
		fmt.Println(err.Error())
	}
}
