你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# Content Addressable

Package contentaddressable contains tools for writing content addressable files.
Files are written to a temporary location, and only renamed to the final
location after the file's OID (Object ID) has been verified.

```go
filename := "path/to/01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b"
file, err := contentaddressable.NewFile(filename)
if err != nil {
  panic(err)
}
defer file.Close()

file.Oid // 01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b

written, err := io.Copy(file, someReader)

if err == nil {
// Move file to final location if OID is verified.
  err = file.Accept()
}

if err != nil {
  panic(err)
}
```

See the [godocs](http://godoc.org/github.com/technoweenie/go-contentaddressable)
for details.

## Installation

    $ go get github.com/technoweenie/go-contentaddressable

Then import it:

    import "github.com/technoweenie/go-contentaddressable"

## Note on Patches/Pull Requests

1. Fork the project on GitHub.
2. Make your feature addition or bug fix.
3. Add tests for it. This is important so I don't break it in a future version
   unintentionally.
4. Commit, do not mess with version or history. (if you want to have
   your own version, that is fine but bump version in a commit by itself I can
   ignore when I pull)
5. Send me a pull request. Bonus points for topic branches.
