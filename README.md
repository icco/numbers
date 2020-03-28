# Numbers

[![Go Report Card](https://goreportcard.com/badge/github.com/icco/numbers)](https://goreportcard.com/report/github.com/icco/numbers)

Pretty straight forward little website.

```ruby
require 'open-uri'
loop { puts open('http://localhost/') {|f| f.read.to_i.chr }; sleep 1 }
```
