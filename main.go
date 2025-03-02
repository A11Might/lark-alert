package main

import (
	"fmt"
	"log"
	"time"

	"github.com/A11Might/lark-alert/model"
	"github.com/A11Might/lark-alert/util"
	"github.com/dustin/go-humanize"
	"github.com/samber/lo"
	"resty.dev/v3"
)

type AlgoliaResponse struct {
	Hits []*Hit `json:"hits"`
}

type Hit struct {
	Title       string `json:"title"`
	CreatedAt   string `json:"created_at"`
	URL         string `json:"url"`
	Author      string `json:"author"`
	Points      int    `json:"points"`
	ObjectID    string `json:"objectID"`
	NumComments int    `json:"num_comments"`
}

func main() {
	log.Default().Println("start fetching headlines")
	endTime := time.Now().Unix()
	startTime := endTime - 25*60*60

	c := resty.New()
	defer c.Close()

	res, err := c.R().
		// SetDebug(true).
		SetQueryString(fmt.Sprintf("page=0&hitsPerPage=10&numericFilters=created_at_i>%d,created_at_i<%d", startTime, endTime)).
		SetResult(&AlgoliaResponse{}).
		Get("https://hn.algolia.com/api/v1/search")
	if err != nil {
		log.Default().Println(err)
		return
	}
	top10objs := lo.Map((*res.Result().(*AlgoliaResponse)).Hits, func(hit *Hit, index int) *Hit {
		if hit.URL == "" {
			hit.URL = fmt.Sprintf("https://news.ycombinator.com/item?id=%s", hit.ObjectID)
		}
		return hit
	})

	list := lo.Map(top10objs, func(hit *Hit, idx int) map[string]string {
		log.Default().Printf("start fetch [%s]\n", hit.Title)
		urlContent, err := util.GetTextContent(hit.URL)
		if err != nil {
			log.Default().Println(err)
			return nil
		}
		summary, err := util.CallOpenAIAPI(model.PromptOneSentenceSummary, urlContent)
		if err != nil {
			log.Default().Println(err)
			return nil
		}

		createdAt, err := time.Parse(time.RFC3339, hit.CreatedAt)
		if err != nil {
			log.Default().Println(err)
			return nil
		}

		timeAgo := humanize.Time(createdAt)

		title := fmt.Sprintf("[%s](%s)", hit.Title, hit.URL)
		status := fmt.Sprintf("%d points by [%s](https://news.ycombinator.com/user?id=%s) %s | [%d comments](https://news.ycombinator.com/item?id=%s)",
			hit.Points, hit.Author, hit.Author, timeAgo, hit.NumComments, hit.ObjectID)

		log.Default().Printf("end fetch [%s]\n", hit.Title)
		return map[string]string{
			"title":        title,
			"list_content": summary,
			"status":       status,
		}
	})
	util.SendCard(util.BuildTemplateCard(list))
	log.Default().Println("fetching headlines end")
}
