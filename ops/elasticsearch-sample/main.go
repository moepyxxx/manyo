package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/kuromojitokenizationmode"
)

var (
	es *elasticsearch.TypedClient
	err error
)

const (
	Address = "http://localhost:9200"
)

func init() {
    // // client
    // es, err := elasticsearch.NewClient(elasticsearch.Config{
    //     Addresses: []string{"http://localhost:9200"},
    // })
    // if err != nil {
    //     log.Fatalf("Error creating the client: %s", err)
    // }

    // typed client
    es, err = elasticsearch.NewTypedClient(elasticsearch.Config{
        Addresses: []string{Address},
    })
    if err != nil {
        log.Fatalf("Error creating the client: %s", err)
    }
}

type SampleIndexData struct {
	Waka string `json:"waka"`
	Butate string `json:"butate"`
	Author string `json:"author"`
	Bangou int `json:"bangou"`
	Maki int `json:"maki"`
	Yomikudashi string `json:"yomikudashi"`
}

func main() {
	ctx := context.Background()
	indexName := "sample_index"

	// err := createIndex(ctx, indexName)
	// if err != nil {
	// 	log.Fatalf("Error createIndex: %s", err)
	// }

	// data := &SampleIndexData{
	// 	Waka: "あかねさす紫野行き標野行き野守は見ずや君が袖振る",
	// }

	// err = index(ctx, indexName, data)
	// if err != nil {
	// 	log.Fatalf("Error index: %s", err)
	// }
	
	err = deleteIndex(ctx, indexName)

	// err = query(ctx, indexName)
}

func boolPtr(b bool) *bool {
    return &b
}

func stringPtr(s string) *string {
	return &s
}

func createIndex(ctx context.Context, indexName string) error {
	_, err = es.Indices.Create(indexName).Request(&create.Request{
		Settings: &types.IndexSettings{
			IndexSettings: map[string]json.RawMessage{
				"refresh_interval": json.RawMessage(`"30s"`),
				"max_result_window": json.RawMessage(`"5000"`),
				"number_of_replicas": json.RawMessage(`"1"`),
			},
			Analysis: &types.IndexSettingsAnalysis{
				Tokenizer: map[string]types.Tokenizer{
					"kuromoji_user_dict": types.KuromojiTokenizer{
						Type: "kuromoji_tokenizer",
						Mode: kuromojitokenizationmode.Extended,
						DiscardPunctuation: boolPtr(false),
						UserDictionary: stringPtr("dic/kogo.txt"),
					},
				},
				Analyzer: map[string]types.Analyzer{
					"my_analyzer": types.CustomAnalyzer{
						Type: "custom",
						Tokenizer: "kuromoji_user_dict",
					},
				},
			},
		},
		Mappings: &types.TypeMapping{
			Properties: map[string]types.Property{
				"waka": types.TextProperty{
					Analyzer: stringPtr("my_analyzer"),
				},
			},
		},
	}).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func deleteIndex(ctx context.Context, indexName string) error {
	_, err = es.Indices.Delete(indexName).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func query(ctx context.Context, indexName string) error {
	// ageLte := 40.0
	// ageLteC := (*types.Float64)(&ageLte)

	// ageGte := 13.0
	// ageGteC := (*types.Float64)(&ageGte)

	pageStart := 0
	size := 50

	req := &search.Request{
		// クエリを書く
		Query: &types.Query{},
		// ページのスタート地点
		From: &pageStart,
		// 返す数
		Size: &size,
		// ソート指定
		// Sort: []types.SortCombinations{
		// 	types.SortOptions{SortOptions: map[string]types.FieldSort{
		// 		"created_at": {Order: &sortorder.Desc},
		// 	}},
		// },
	}
	res, err := es.Search().
		Index(indexName).
		Request(req).
		Do(ctx)

	if err != nil {
		log.Fatalf("Error query: %s", err)
	}

	// total数を出す
	fmt.Println(res.Hits.Total)

	ds := []*SampleIndexData{}
	for _, hit := range res.Hits.Hits {
		var d *SampleIndexData
		if err := json.Unmarshal(hit.Source_, &d); err != nil {
			log.Fatalf("Error decoding: %s", err)
		}
		ds = append(ds, d)
	}

	// データを出す
	fmt.Println(ds)

	return nil
}

func index(ctx context.Context, indexName string, data *SampleIndexData) error {
	_, err = es.Index(indexName).Request(data).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func delete(ctx context.Context, indexName string, id string) error {
	_, err = es.Delete(indexName, id).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func get(ctx context.Context, indexName string, id string) error {
	res, err := es.Get(indexName, id).Do(ctx)
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}