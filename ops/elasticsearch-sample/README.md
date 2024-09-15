# elasticsearch

## \_cat

https://www.elastic.co/guide/en/elasticsearch/guide/current/_cat_api.html

- ヘルスチェック

```
curl -X GET "localhost:9200/_cat/health?v&pretty"
```

- インデックスチェック

```
curl -X GET "localhost:9200/_cat/indices?v"
curl -X GET "localhost:9200/_cluster/health/sample_index?pretty"
```

- ステータスが yellow の場合
  - レプリカを 0 にする。local の場合は es 環境を 2 つ立てるのが面倒なため（レプリカシャードはプライマリシャードとは別のノードに作る必要があるため、ノードが 1 つだとほぼ確で最初 yellow になる）

```
curl -X PUT "localhost:9200/sample_index_2/_settings" -H 'Content-Type: application/json' -d '{
  "index": {
    "number_of_replicas": 0
  }
}'
```

## mappings

マッピング確認

```
curl -X GET "localhost:9200/sample_index/_mappings"
```
