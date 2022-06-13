# go-strchecker
## 背景
在大型项目开发过程中，经常会遇到打印大量日志，输出信息和在源码中写注释的情况。对于软件开发来说，我们一般都是打印输出英文的日志（主要考虑软件在各种环境下的兼容性，如果打印中文日志可能会出现乱码，另外英文日志更容易搜索，更容易后续做国际化），但是对于我们中国人来说，很容易就把中文全角的中文标点符号一不注意就写到日志中了。不过源码中的注释因为是完全面向开发者的，不会面向客户，所以如果研发团队全是中国人，那么代码注释用中文就更有效率。
在实际开发过程中，确实就发现了打印日志中包含了中文标点的情况，但是如果我们直接用中文标点在IDE中进行全文搜索，就好发现大量的代码注释使用中文标点，而到底哪里是日志打印时的中文标点，哪里是注释中的中文标点，根本看不出来。于是我参考golangci-lint的代码扫描检查功能，写了一个Go源码中字符串规范检查的lint工具：strchecker。
## strchecker介绍
strchecker可以扫描某个文件夹或者该文件夹下的所有子文件夹中的go代码，并对其中的go代码进行语法分析，构建语法树，找到其中申明的常量、变量、函数参数、返回值、赋值、case语句等场景下的字符串string类型，然后对这些字符串进行正则匹配。系统默认的正则匹配方式是只有ASCII字符才是合法字符，只要超过一个字节的（比如中文、中文标点等都是多字节的）就会被匹配到，而匹配到的字符串就算是非法字符串，并最终将这些非法字符串打印出来。
下面举一个示例：
1. 安装strchecker
```
go install github.com/studyzy/go-strchecker/cmd/strchecker@latest
```
2. 找到我们要进行扫描的文件夹，这里就以go-strchecker/testdata/ 这个文件夹为例，进行非法字符串扫描。
```
strchecker $GOPATH/src/github.com/studyzy/go-strchecker/testdata
```
3. 输出结果如下：
```
0 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/call.go:9:60 has invalid string: "！"
1 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/call.go:10:11 has invalid string: "a！b"
2 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/call.go:11:5 has invalid string: "aa！"
3 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/call.go:12:40 has invalid string: "bb！"
4 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/main.go:10:30 has invalid string: "not found！"
5 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/main.go:12:17 has invalid string: "no，data！"
6 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/main.go:15:14 has invalid string: "Hello，World！"
7 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/main.go:16:12 has invalid string: "Current time："
8 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/main.go:19:15 has invalid string: "한국어"
9 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/main.go:20:15 has invalid string: "にほんご"
10 /Users/devinzeng/go/src/github.com/studyzy/go-strchecker/testdata/main.go:22:14 has invalid string: ":) 😁😁😁"
```
4. 如果是在Goland这样的IDE中的，那么我们就可以直接点击源码路径，定位到具体的代码位置了。

5. 然后对其中的字符串进行修复，将其中的中文标点替换成英文标点。
6. 如果我们有一些特殊的要求，而不是只允许ASCII码表中的内容才是合法内容，比如我们允许中英文，但是不允许日文、韩文等，那么怎么办？于是我在参数中预置了ASCII表允许和ASCII+中文+中文标点允许这两种常用的匹配类型。如果我们想允许ASCII和中文，那么命令是：
```
strchecker -invalid-type=1 ./testdata/...
```
【注意：这里最后的...表示testdata目录下的所有子文件和子文件夹，会递归的扫描，当然因为我们testdata没有子文件夹，所以这个...加或者不加都是一样的。】
输入结果为：
```
0 testdata/main.go:19:15 has invalid string: "한국어"
1 testdata/main.go:20:15 has invalid string: "にほんご"
2 testdata/main.go:22:14 has invalid string: ":) 😁😁😁"
```
## 结论
strchecker是一个专门用于扫描Golang源码中字符串是否包含特定正则表达式的Lint工具。使用strchecker可以快速找到Go源码中字符串中隐藏的中文标点、非中英文字符等，很适合用于国人在大型go项目中扫描日志输出或者其他字符串定义时不小心出现的中文标点的情况。

当然，如果本身项目的源码中连注释都不允许用中文和中文标点，那么就直接用IDE的search功能即可，本工具是不扫描源码中注释的内容的。