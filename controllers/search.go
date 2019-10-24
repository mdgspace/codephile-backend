package controllers

import (
	"context"
	"encoding/json"
	search "github.com/mdg-iitr/Codephile/services/elastic"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
)

// @Title Search
// @Description Endpoint to search users
// @Security token_auth read:user
// @Param	count		query 	string	true		"No of search objects to be returned"
// @Param	query		query 	string	true		"Search query"
// @Success 200 {object} []models.User
// @Failure 403 "search query is too small"
// @router /search [get]
func (u *UserController) Search() {
	query := u.GetString("query")
	if len(query) < 4 {
		u.Ctx.ResponseWriter.WriteHeader(403)
		u.Data["json"] = "search query is too small"
		u.ServeJSON()
		return
	}
	count := u.GetString("count")
	c, err := strconv.Atoi(count)
	//Default query response size
	if err != nil {
		c = 15
	}
	pq := elastic.NewQueryStringQuery("*" + query + "*").Field("username").Field("handle.codechef").Field("handle.spoj").Field("handle.codeforces").Field("handle.hackerrank").Fuzziness("4")
	q := elastic.NewMultiMatchQuery(query, "username", "handle.codechef", "handle.spoj", "handle.codeforces", "handle.hackerrank").Fuzziness("4")
	bq := elastic.NewBoolQuery().Should(q, pq)
	client := search.GetElasticClient()
	result, err := client.Search().Index("codephile").
		Pretty(false).Query(bq).Size(c).
		Do(context.Background())
	if err != nil {
		log.Println(err.Error())
	}
	var results []interface{}
	for _, hit := range result.Hits.Hits {
		var result interface{}
		err := json.Unmarshal(hit.Source, &result)
		if err != nil {
			log.Println(err.Error())
		}
		results = append(results, result)
	}
	u.Data["json"] = results
	u.ServeJSON()
}
