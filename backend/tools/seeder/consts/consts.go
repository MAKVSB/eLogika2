package consts

import "encoding/json"

var DefaultContent = json.RawMessage(`{
  "type": "doc",
  "content": [{
    "type": "paragraph",
    "attrs": {"textAlign": null},
    "content": [{
      "type": "text",
      "text": "asdasdasdasdasd"
    }]
  }]
}`)
