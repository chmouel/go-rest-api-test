[![GolangCI](https://golangci.com/badges/github.com/chmouel/go-rest-api-test.svg)](https://golangci.com/r/github.com/chmouel/go-rest-api-test)
[![License](https://img.shields.io/github/license/chmouel/go-rest-api-test?)](/LICENSE)

Go rest api tester
==================

Simple HTTP rest api responder, responding whatever you like according to rules. Useful for integration tests.

Used in https://github.com/tektoncd/catalog task testing!

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
