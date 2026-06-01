package es

// ItemMapping ES 物品索引映射（支持 IK 中文分词，需安装 analysis-ik 插件）
const ItemMapping = `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "ik_smart_analyzer": {
          "type": "custom",
          "tokenizer": "ik_smart"
        },
        "ik_max_analyzer": {
          "type": "custom",
          "tokenizer": "ik_max_word"
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "user_id": { "type": "long" },
      "title": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart",
        "fields": {
          "keyword": { "type": "keyword" },
          "suggest": {
            "type": "completion",
            "analyzer": "ik_max_word"
          }
        }
      },
      "description": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_smart"
      },
      "category_id": { "type": "integer" },
      "category_name": { "type": "keyword" },
      "estimated_value": { "type": "float" },
      "condition": { "type": "keyword" },
      "images": { "type": "keyword" },
      "status": { "type": "keyword" },
      "location": { "type": "keyword" },
      "view_count": { "type": "integer" },
      "want_items": { "type": "text", "analyzer": "ik_smart" },
      "user_nickname": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}`

// ItemMappingFallback 降级映射：当 IK 插件不可用时使用 ES 内置分析器。
// 中文用 standard 分析器（单字切分），效果不如 IK，但保证系统可用。
const ItemMappingFallback = `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "user_id": { "type": "long" },
      "title": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" },
          "suggest": { "type": "completion" }
        }
      },
      "description": { "type": "text", "analyzer": "standard" },
      "category_id": { "type": "integer" },
      "category_name": { "type": "keyword" },
      "estimated_value": { "type": "float" },
      "condition": { "type": "keyword" },
      "images": { "type": "keyword" },
      "status": { "type": "keyword" },
      "location": { "type": "keyword" },
      "view_count": { "type": "integer" },
      "want_items": { "type": "text", "analyzer": "standard" },
      "user_nickname": { "type": "keyword" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}`
