FROM docker.elastic.co/elasticsearch/elasticsearch:8.15.1

RUN elasticsearch-plugin install analysis-kuromoji
