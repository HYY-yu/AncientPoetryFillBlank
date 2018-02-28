package answer

import (
	"regexp"
	"strings"
	"errors"
)

type Find struct {
	BlankId     int
	PreString   string // 前导诗句
	BlankString string // 本空
	PostString  string // 后置诗句
	BlankNum    int
	BlankFinish bool //是否是完全块
}

//find the answer by contents.
//findSrc : blank is ___
func FindTheAnswer(contents string, findSrc string) (string, error) {
	finds := make([]Find, 0, 0)
	findSplits, findPun := SplitByPunctuation(findSrc)
	if len(findSplits) != len(findPun) {
		return "", errors.New("标点符号要输入完整，才是一句话。")
	}

	var j = 0
	for i := range findSplits {
		if findSplits[i] == BLANK_STRING {
			var newFind Find
			newFind.BlankId = j
			newFind.BlankNum = i
			//前一句
			if i > 0 {
				newFind.PreString = findSplits[i-1]
			}
			//后一句
			if i < len(findSplits)-1 {
				newFind.PostString = findSplits[i+1]
			}

			finds = append(finds, newFind)
			j++
		}
	}

	//构造答案
	makeAnswer(contents, finds)

	//Return答案
	var returnStr string
	for i := range findSplits {
		if findSplits[i] == BLANK_STRING {
			returnStr += "<ans>" + GetFindsByJ(finds, i).BlankString + "</ans>"
		} else {
			returnStr += findSplits[i]
		}
		//标点
		returnStr += findPun[i]
	}
	return returnStr, nil
}

func GetFindsByJ(finds []Find, j int) (Find) {
	for i := range finds {
		if finds[i].BlankNum == j {
			return finds[i]
		}
	}
	return Find{}
}

func makeAnswer(contents string, finds []Find) {
	for i := range finds {
		//自主块：可自主推导出答案的Find
		//完全块：PreString BlankString PostString 均有的Find
		//推导所有自主块的答案 - 让它们变成完全块
		if checkStringNotBlank(finds[i].PreString) || checkStringNotBlank(finds[i].PostString) {
			rootFind := &finds[i]
			if checkStringNotBlank(rootFind.PreString) {
				makeWithPreContent(contents, rootFind)
			} else {
				makeWithPostContent(contents, rootFind)
			}
		}
	}

	// 遍历完全块，找它们左右是否有非自主块。给它们添加答案
	for i := range finds {
		if finds[i].BlankFinish {
			root := i
			//前向搜索
			if root > 0 {
				if !finds[root-1].BlankFinish {
					for i := root - 1; i >= 0; i-- {
						if !finds[i].BlankFinish {
							finds[i].BlankString = finds[i+1].PreString
							finds[i].PostString = finds[i+1].BlankString

							//find finds[i].PreString
							makeWithPostContent(contents, &finds[i])
						} else {
							break
						}
					}
				}
			}
			//后向搜索
			if root < len(finds)-1 {
				if !finds[root+1].BlankFinish {
					for i := root + 1; i < len(finds); i++ {
						if !finds[i].BlankFinish {
							finds[i].BlankString = finds[i-1].PostString
							finds[i].PreString = finds[i-1].BlankString

							//find finds[i].PostString
							makeWithPreContent(contents, &finds[i])
						} else {
							break
						}
					}
				}
			}
		}
	}
}

func checkStringNotBlank(check string) bool {
	return len(check) > 0 && check != BLANK_STRING
}

func makeWithPreContent(contents string, newFind *Find) {
	allC, _ := SplitByPunctuation(contents)
	for l := range allC {
		if allC[l] == newFind.PreString && l < len(allC)-1 {
			newFind.BlankString = allC[l+1]
			if l < len(allC)-2 {
				newFind.PostString = allC[l+2]
			}
			newFind.BlankFinish = true
		}
	}
}

func makeWithPostContent(contents string, newFind *Find) {
	allC, _ := SplitByPunctuation(contents)
	for l := range allC {
		if allC[l] == newFind.PostString && l > 0 {
			newFind.BlankString = allC[l-1]
			if l-1 > 0 {
				newFind.PreString = allC[l-2]
			}
			newFind.BlankFinish = true
		}
	}
}

func SplitByPunctuation(s string) ([]string, []string) {
	regPunctuation, _ := regexp.Compile(`[,.，。?、？！!;；]`)
	//匹配标点符号，保存下来。
	toPun := regPunctuation.FindAllString(s, -1)
	x := regPunctuation.ReplaceAllString(s, "@")
	result := strings.Split(x, "@")
	//如果最后一个为空，去掉
	if len(result[len(result)-1]) == 0 {
		result = result[:len(result)-1]
	}
	return result, toPun
}
