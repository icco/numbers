# Numbers

Pretty straight forward little website.

```ruby
require 'open-uri'
loop { puts open('http://numbersstation.blue/') {|f| f.read.to_i.chr }; sleep 1 }
```
