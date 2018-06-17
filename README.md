## 基于deflate的压缩算法，参考了zip

## 参考文章
> - [ZIP压缩算法详细分析及解压实例解释](https://www.cnblogs.com/esingchan/p/3958962.html)
> - [详细图解哈夫曼Huffman编码树](https://blog.csdn.net/FX677588/article/details/70767446)
> - [几种压缩算法实现原理详解](https://blog.csdn.net/ghevinn/article/details/45747465)

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

####部分压缩率测试对比

1406967k 电影文件 重复率很低 压缩率基本差不多(也试了下歌曲 压缩率也很低)

  |压缩方式|压缩后大小(单位:k)|压缩比率(压缩后/压缩前)|
  |:--|:--|:--|
  | zip | 1403910 | 99.7827|
  | rar | 1406345 | 99.9557 |
  | my | 1405906 | 99.9245 |

111632  k 一个测试文件 信息熵较小 重新信息较多
  |压缩方式|压缩后大小(单位:k)|压缩比率(压缩后/压缩前)|
  |:--|:--|:--|
  | zip | 22920 | 20.5317|
  | rar | 16325 | 14.6239 |
  | my | 29726 | 26.6285 |


#### 接下来需要做的地方

> - 大文件压缩肯定还是有问题的 目前的做法是直接把全部文件内容读入内存中 需要优化
> - 压缩时间消耗很长 很吃内存  cpu消耗比较高
> - 对CL（Code Length）还没有进行压缩
> - 代码需要整理下 有部分为了早点看到结果 就是写的比较乱