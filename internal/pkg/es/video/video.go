package video

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"StreamCore/internal/pkg/domain"
	"StreamCore/internal/pkg/es/model"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type VideoElastic interface {
	AddVideo(ctx context.Context, video *domain.Video) error
	RemoveVideo(ctx context.Context, vid uint) error
	SearchVideo(ctx context.Context, query *domain.VideoQuery) ([]uint, int64, error)
}

func (c *VideoEsClient) CreateIndex(ctx context.Context) error {
	// esapi create an index with mapping
	req := &esapi.IndicesCreateRequest{
		Index: c.indexName,
		Body:  strings.NewReader(mappingVid),
	}
	_, err := req.Do(ctx, c.es)
	if err != nil {
		return fmt.Errorf("VideoEsClient.CreateIndex failed: %w", err)
	}
	return nil
}

func (c *VideoEsClient) AddVideo(ctx context.Context, video *domain.Video) error {
	videoDoc := model.VideoEs{
		Title:       video.Title,
		Description: video.Description,
		AuthorId:    video.AuthorId,
	}

	_, err := c.es.Index(c.indexName).
		Id(strconv.FormatInt(int64(video.Id), 10)).
		Document(videoDoc).
		Do(ctx)
	if err != nil {
		return fmt.Errorf("VideoEsClient.AddVideo failed: %w", err)
	}
	return nil
}

func (c *VideoEsClient) RemoveVideo(ctx context.Context, vid uint) error {
	vidstr := strconv.FormatInt(int64(vid), 10)
	_, err := c.es.Delete(c.indexName, vidstr).Do(ctx)
	if err != nil {
		return fmt.Errorf("VideoEsClient.RemoveVideo failed: %w", err)
	}
	return nil
}

func (c *VideoEsClient) SearchVideo(ctx context.Context, query *domain.VideoQuery) ([]uint, int64, error) {
	resp, err := c.es.Search().Index(c.indexName).Request(&search.Request{
		Query: c.buildQuery(query),
	}).Do(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("VideoEsClient.SearchVideo failed: %w", err)
	}

	hcount := resp.Hits.Total.Value
	hits := make([]uint, 0)
	for _, hit := range resp.Hits.Hits {
		id, _ := strconv.ParseUint(*hit.Id_, 10, 64)
		hits = append(hits, uint(id))
	}
	return hits, hcount, nil
}

var mappingVid = `{
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "author_id": {
        "type": "keyword"
      },
      "title": {
        "type": "text",
        "analyzer": "standard"
      },
      "description": {
        "type": "text",
        "analyzer": "standard"
      },
      "username": {
        "type": "text",
        "analyzer": "standard"
      },
      "published_at": {
        "type": "date"
      }
    }
  }
}`

func (c *VideoEsClient) buildQuery(query *domain.VideoQuery) *types.Query {
	/* like:
	{
	  "query": {
	    "bool": {
	      "must": [
	        {
	          "match": {
	            "title": "xxx"
	          }
	        },
	        {
	          "match": {
	            "description": "xxxx"
	          }
	        }
	      ],
	      "filter": [
	        { "term": { "author_id": "10" } },
			{
			  "range": {
			    "published_at": {
				  "gte": "1990-01-01T08:00:00Z"
				  "lte": "2030-12-31T23:59:59Z"
			    }
			  }
			}
	      ]
	    }
	  }
	}*/
	bq := &types.BoolQuery{}
	anyCondition := false

	// must - match query (title)
	if query.TitleMatches != "" {
		bq.Must = append(bq.Must, types.Query{
			Match: map[string]types.MatchQuery{
				"title": {Query: query.TitleMatches},
			},
		})
		anyCondition = true
	}

	// must - match query (description)
	if query.DescMatches != "" {
		bq.Must = append(bq.Must, types.Query{
			Match: map[string]types.MatchQuery{
				"description": {Query: query.DescMatches},
			},
		})
		anyCondition = true
	}

	// must - match query (username)
	if query.UsernameMatches != nil {
		bq.Must = append(bq.Must, types.Query{
			Match: map[string]types.MatchQuery{
				"username": {Query: *query.UsernameMatches},
			},
		})
		anyCondition = true
	}

	// filter - term query (author_id)
	if query.AuthorIdIsExact != nil {
		authorId := strconv.FormatInt(int64(*query.AuthorIdIsExact), 10)
		bq.Filter = append(bq.Filter, types.Query{
			Term: map[string]types.TermQuery{
				"author_id": {Value: authorId},
			},
		})
		anyCondition = true
	}

	// filter - range query (from & to time)
	rq := types.DateRangeQuery{}
	timeRange := false
	if query.FromDate != nil {
		rq.Gte = query.FromDate
		timeRange = true
		anyCondition = true
	}
	if query.ToDate != nil {
		rq.Lte = query.ToDate
		timeRange = true
		anyCondition = true
	}
	if timeRange {
		bq.Filter = append(bq.Filter, types.Query{
			Range: map[string]types.RangeQuery{
				"published_at": rq,
			},
		})
	}

	if anyCondition { // any condition -> take boolQuery
		return &types.Query{
			Bool: bq,
		}
	} else { // no condition -> take a matchAllQuery
		return &types.Query{
			MatchAll: &types.MatchAllQuery{},
		}
	}
}

type VideoEsClient struct {
	es        *elastic.TypedClient
	indexName string
}

func NewVideoElastic(es *elastic.TypedClient) VideoElastic {
	return &VideoEsClient{
		es:        es,
		indexName: "videos",
	}
}
