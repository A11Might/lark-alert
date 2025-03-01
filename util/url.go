package util

import (
	"fmt"

	"resty.dev/v3"
)

func GetTextContent(url string) (string, error) {
	c := resty.New()
	defer c.Close()

	res, err := c.R().
		Get(fmt.Sprintf("https://r.jina.ai/%s", url))
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
