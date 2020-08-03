Go rest api tester
==================

Simple HTTP rest api responder, which respondes according to rules.

RULES
=====

* will answer on a `GET` on URL `/repo/foo/bar/issues/1/comments` and reply with a `200` `{"status": 200}` with `content-type: text/json`

```yaml
---
headers:
  method: GET
  path: /repos/{repo:[^/]+/[^/]+}/issues/{issue:[0-9]+}/comments
response:
  status: 200
  # file: post-comment.response.json
  output: '{"status": 200}'
  content-type: text/json
``
