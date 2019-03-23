## 基于deflate的压缩算法，参考了zip（主要是想尝试下自己写个压缩工具） [![travis-ci.org](https://img.shields.io/travis/whalecold/zlip/master.svg)](https://travis-ci.org/whalecold/compress)

## 参考文章
> - [ZIP压缩算法详细分析及解压实例解释](https://www.cnblogs.com/esingchan/p/3958962.html)
> - [详细图解哈夫曼Huffman编码树](https://blog.csdn.net/FX677588/article/details/70767446)
> - [几种压缩算法实现原理详解](https://blog.csdn.net/ghevinn/article/details/45747465)

### 参数
> - 使用参数 -d(true/false) true 表示解压 false 表示压缩  
> - -source=file 被压缩或者是被解压的文件名
> - -dest=file 压缩后或者是解压后的文件名

### 2018-06-11  demo版

> - 把huffman算法基本完成了，但是压缩了不是很高，简单看下问题，是因为序列化huffman树的占用了很大的空间 这里有很大的优化空间（毕竟这只是算法的一部分）
> - 对于中文的压缩率没有英文高  因为这个是基于字节去建立huffman树的 中文编码中的byte分母要比英文大很多
> - 顺手修复了最后一位信息错误的bug

### 2018-06-14 提交

> - 到现在为止增加了lz77算法 把huffman树转成了deflate树
> - deflate树信息转字节流和反转基本跑通测试了
> - `todo` distance树好了 literal/length树还没好 但基本上也就是时间问题啦
> - 增加了几个测试用例
> - 添加对字节流的编解码

### 2018-06-17 version 1.0.0

> - 基本算是一个阶段了 distance树也好了 基本按照先lz77算法，然后按照deflate算法进行压缩了

### 2018-06-19

>- 增加了多线程压缩和解压 减少了压缩耗时

### 2018-06-22

>- 优化内存消耗
>- 优化了一波压缩耗时 现在在同等条件下压缩时间是rar的一半 但是压缩效率和内存没有rar做的好

#### 部分压缩率测试对比(my1 是经过一版优化之后的数据)

 __测试平台__:_windows_  __cpu__:_i5-4590 3.30GHz_  __version__:_go1.10_

1406967k 电影文件 重复率很低 压缩率基本差不多(也试了下歌曲 压缩率也很低)

|压缩方式|压缩后大小(单位:k)|压缩比率(压缩后/压缩前)|压缩耗时(ms)|解压耗时(ms)|内存(M)
|:--|:--|:--|:--|:--|:--|
| zip | 1403910 | 99.7827|35370|10020|10|
| rar | 1406345 | 99.9557 |102910|16800|105|
| my | 1405906 | 99.9245 |很久就是了|很久就是了|2300|
| my1 | 1405664 | 99.9245 |129335|65219|2300|
| my2 | 1405664 | 99.9245 |110115|xxx|2300|
| my3 | 1405664 | 99.9245 |100654|40073|250|
| my4 | 1405664 | 99.9245 |52191|27901|250|


  `rar的内存占用一直稳定在100m左右，而zip则只占用了10m，好厉害！我的内存在2g左右(已经优化到250m左右)，差距不是一般的大。
  综合看下来zip除了在压缩效率方面比rar稍微差了点，在时间和内存消耗方面都优于rar，感觉zip要好一些。
  当然最后的使用还是看具体的需求了。`

111632  k 一个测试文件 信息熵较小 重复信息较多

|压缩方式|压缩后大小(单位:k)|压缩比率(压缩后/压缩前)|压缩耗时(ms)|解压耗时(ms)|
|:--|:--|:--|:--|:--|
| zip | 22920 | 20.5317|待测|待测|
| rar | 16325 | 14.6239 |待测|待测|
| my | 29726 | 26.6285 |没测|没测|
| my1 | 26989 | 24.1767 |2596|1395|
| my4 | 27565 | 24.6927 |1347|625|这里牺牲了一定的压缩率


#### 需要优化的地方
- [ ] 对CL（Code Length）还没有进行压缩(这部分需要对huffman树进行剪枝 暂时先不做了)
- [ ] 目前对一个被压缩对象小于2个重复字段的压缩会出现bug

#### 已经优化好
- [x] 大文件压缩肯定还是有问题的 目前的做法是直接把全部文件内容读入内存中`(优化：把文件分成很多小块进行压缩，
块大小可以配置，配置参数是LZ77_ChunkSize，解压的时候再逐块解压然后在组装)`
- [x] 压缩时间消耗很长 `(优化：采用多协程压缩，因为把文件分块了，所以按照每个块一个协程进行压缩协程任务之间是没有竞争的，
不会因为锁竞争导致效率降低)`
- [x] 运行内存消耗过大`(优化：之前那一版是一个块一个协程，在文件很大的时候会占用很多的内存，然后我又思考了一下，
协程是在线程上面进行了一层封装，所以在不改进算法的前提下影响压缩时间的本质还是cpu个数。假设cpu的个数为n, 
在设置runtime.GOMAXPROCS(n)之后，启动100个协程和n个协程进行压缩时间是差不多的,n个协程还能节约部分系统协程调度的开销(测试
下来在压缩1.4g电影文件的时候时间消耗减少了10s,这个算是意外之喜吧!)，
开启100个协程还有更多的内存开销（这里指的是在压缩文件过程中所消耗的内存），所以我改进的方法就是只创建n个协程和一个调度器协程，
这几个协程去调度器里申请需要的信息进行压缩 改进过后再压缩大文件的时候内存显著降低。)`

#### 其他优化
> - 把原码与huffman编码的映射由map改成了数组下标索引，在这个算法里原码的值都在300以内，
所以不需要很大的数组，也不会造成空间浪费。在压缩1.4g电影文件的时候 能节约20s左右的时间

#### 分块压缩块的大小(LZ77_ChunkSize)对压缩效率有一定的影响 需要取一个适当的值 目前取得是5M 没有进行很深的研究 

####TODO
- [ ]  完善测试用例
- [ ]  支持现在的zip
- [ ]  支持压缩文件夹
- [ ]  完善ci
