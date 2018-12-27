package appsearch

import (
	"github.com/neal1991/gshark/logger"
	"github.com/neal1991/gshark/models"
	"github.com/neal1991/gshark/vars"
	"time"
	"github.com/gocolly/colly"
	"fmt"
	"strings"
	"strconv"
)

func ScheduleTasks(duration time.Duration) {
	for {
		logger.Log.Infof("Complete the scan of APP, start to sleep %v seconds", duration*time.Second)
		time.Sleep(duration * time.Second)
	}
}

func GenerateSearchCodeTask() (map[int][]models.Rule, error) {
	result := make(map[int][]models.Rule)
	rules, err := models.GetValidRules()
	ruleNum := len(rules)
	batch := ruleNum / vars.SearchNum

	for i := 0; i < batch; i++ {
		result[i] = rules[vars.SearchNum*i : vars.SearchNum*(i+1)]
	}

	if ruleNum%vars.SearchNum != 0 {
		result[batch] = rules[vars.SearchNum*batch : ruleNum]
	}
	return result, err
}

func SaveResults(results []*models.APPSearchResult) {
	for _, result := range results {
		result.Insert()
	}
}

func SearchForApp(rule models.Rule) []*models.APPSearchResult {
	appSearchResults := make([]*models.APPSearchResult, 0)
	if rule.Caption == "HUAWEI" {
		baseUrl := "http://appstore.huawei.com"
		url := baseUrl + "/search/" + rule.Pattern
		for i := 1; i < 99; i++ {
			c := colly.NewCollector()
			linkUrl := url + "/" + strconv.Itoa(i)
			c.OnHTML("body", func(e *colly.HTMLElement) {
				hasNext, appSearchResult := saveAppSearchResult(baseUrl, e)
				appSearchResults = append(appSearchResults, appSearchResult)
				if !hasNext {
					return
				}
				fmt.Println(linkUrl)
			})
			c.Visit(linkUrl)
		}
		// todo
		// other app market
	} else {

	}
	return appSearchResults
}

func saveAppSearchResult(baseUrl string, e *colly.HTMLElement)  (bool, *models.APPSearchResult){
	appSearchResult := new(models.APPSearchResult)
	count := 0
	var hasNext bool
	e.ForEach(".list-game-app.dotline-btn.nofloat", func(i int, element *colly.HTMLElement) {
		var title = element.ChildText(".title")
		var content = element.ChildText(".content")
		var deployDate = strings.Replace(element.ChildText(".date"),
			"发布时间： ", "", -1)
		var appUrl = baseUrl + element.ChildAttr(".title a", "href")
		appSearchResult.Name = &title
		appSearchResult.Description = &content
		appSearchResult.DeployDate = &deployDate
		appSearchResult.APPUrl = &appUrl
		appSearchResult.Status = 0
		fmt.Println(*appSearchResult.Name)
		fmt.Println(*appSearchResult.Description)
		fmt.Println(*appSearchResult.DeployDate)
		fmt.Println(*appSearchResult.APPUrl)
		count++
	})
	if count < 2 {
		hasNext = false
	} else {
		hasNext = true
	}
	return hasNext, appSearchResult
}
