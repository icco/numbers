# Numbers

Pretty straight forward little website.

```ruby
loop { open('http://numbersstation.blue/') {|f| f.read.to_i.chr }; sleep 1 }
```
