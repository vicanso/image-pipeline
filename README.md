# image-pipeline

`image-pipeline`提供一系列的图片处理功能，如格式转换、缩放、水印添加等，可以通过pipeline的形式串连。提供的功能如下：

- 根据可接受的图片格式自动选择匹配格式，优先为`avif`，其次`webp`
- 根据指定转换格式与质量压缩图片
- 将图片的宽高调整为适合指定的宽高
- 将图片的宽高调整为填满指定的宽高
- 将水印以添加至图片上

## 图片拉取

`image-pipeline`不提供图片的存储功能，其支持以下几类的图片拉取方式：

- HTTP
- 文件目录
- minio
- mongodb的gridfs
- 阿里云oss

### HTTP

HTTP形式的支持多IP节点以及健康检查，配置格式：`http://ip1:port1[,ip2:port2]/ping`，其中`/ping`为对应HTTP服务的健康检查路径，多个ip则以`,`分隔。

```go
func AddHTTPFinder(name, uri string, onStatus ...upstream.StatusListener) error
```

- `name`: HTTP Finder的名称，其名称必须唯一，因为后续pipeline中是通过其名称来指定使用对应的finder，如果名称重复，后面添加的则会覆盖原有的
- `uri`: 源地址，如多个节点用`,`分隔
- `onStatus`: 健康检测变化时的回调，可选

```go
finderName := "myImageHTTPService"
imagepipeline.AddHTTPFinder(finderName, "http://192.168.1.3:8080,192.168.1.6:8080/ping")
```

### File

文件的形式较为简单，只需要指定目录即可，生产环境不建议使用此方式，因为pu如果要使用则建议使用NFS等形式挂载网络存储。

```go
func AddFileFinder(name, basePath string) error 
```

- `name`: 文件Finder的名称，必须唯一
- `basePath`: 文件路径，finder提供的图片拉取均为此目录下

```go
imagepipeline.AddFileFinder("myImages", "/opt/images")
```

### Minio

minio是一款开源的对象存储服务器，兼容亚马逊的S3协议，初始化其finder需要指定连接串。

```go
func AddMinioFinder(name, uri string) error 
```

- `name`: Minio Finder的名称，必须唯一
- `uri`: 连接串，格式如下：`minio://ip:port/?accessKey=key&secretKey=secret`，其中`minio://`并不影响，以`http://`的形式一样可行。`accessKey`与`secretKey`则根据minio的配置填写即可

```go
imagepipeline.AddMinioFinder("myMinioImages", "minio://192.168.1.5:900/?accessKey=key&secretKey=secret")
```

### Gridfs

gridfs是mongodb提供的文件存储，可方便的存储大量文件。

```go
func AddGridFSFinder(name, uri string) error
```

- `name`: Mongodb Gridfs的名称，必须唯一
- `uri`: mongodb的连接串，格式如下：`mongodb://user:pwd@127.0.0.1:27017/admin`

```go
imagepipeline.AddGridFSFinder("myGridfsImages", "mongodb://user:pwd@127.0.0.1:27017/admin")
```

### Aliyun OSS

阿里云的对象存储服务可以方便快捷的存储各类图片。


```go
func AddAliyunOSSFinder(name, uri string) error
```

- `name`: 阿里云oss的名称，必须唯一
- `uri`: 连接串，格式如下：`https://oss-cn-beijing.aliyuncs.com?accessKey=key&secretKey=secret`

```go
imagepipeline.AddGridFSFinder("myOSSImages", "https://oss-cn-beijing.aliyuncs.com?accessKey=key&secretKey=secret")
```

## Pipeline

初始化各类Finder之后，即可以pipeline的形式拼接各类的任务（多个任务以|连接，任务参数以/分隔)，假设所有类型的finder均有初始化(其名称为`类型Finder`，如`httpFinder`)。pipeline在处理任务时，优先按照匹配固定规则，如果都不匹配则以finder的形式来处理。

固定的任务类型如下：`proxy`、`optimize`、`autoOptimize`、`fitResize`、`fillResize`，`optimize`或`autoOptimize`图片压缩转换一般都是作为处理任务。

需要注意，pipeline的任务第一个必须是获取图片数据的，下面是各类任务的描述：

### Proxy

`proxy/https%3A%2F%2Fwww.baidu.com%2Fimg%2FPCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png`，任务描述以`proxy`开头，表示以HTTP形式获取后面URL中的图片

### Optimize

`optimize/192.168.1.1:6002/80/webp`，任务描述以`optimize`开头，第二个参数为[tiny]()的服务地址，它优先以它为key获取env的参数，如为空则直接使用此参数为地址。例如如果设置了TINY_ADDR这个env的值为`192.168.1.1:6002`，则上面的描述可以调整为`optimize/TINY_ADDR/80/webp`。第三个参数`80`表示压缩时选择的质量(可选)，第四个参数`webp`表示转换的图片格式(可选)

### AutoOptimize

`autoOptimize/192.168.1.1:6002/80`，任务描述以`autoOptimize`开头，前三个参数与`optimize`一致。此任务会根据客户端可接受的图片类型选择最优的图片：`avif` -> `webp` -> `原类型`

### FitResize

`fitResize/500/600`，任务描述以`fitResize`开头，后面两个参数为宽、高，此任务会根据指定的宽高调整图片大小

### FillResize

`fillResize/500/600`，任务描述以`fillResize`开头，参数与`fitResize`，只是调整宽高的处理方式不同

### HTTP Finder

`httpFinder/image%2Fbanner.png`，此处假设初始化了一个名为`httpFinder`的http finder。对于http finder，后面的参数则是对应的图片地址，通过此地址获取对应的图片

### File Finder

`fileFinder/image%2Fbanner.png`，此处假设初始化了一个名为`fileFinder`的file finder。对于file finder ，后面的参数则是对应图片的相对地址，通过此地址读取图片


### Minio Finder

`minioFinder/bucketName/objectName`，此处假设初始化了一个名为`minioFinder`的minio finder。对于minio finder，后面的参数对应其bucket与object的名称，通过这两个参数指定图片存储位置

### GridFS Finder

`gridfsFinder/objectID/collection`，此处假设初始化了一个名为`gridfsFinder`的grfidfs finder。对于gridfs finder，第一个参数为图片的objectID，第二个参数为对应的collection(可选，若无此参数则为fs)，读取对应的Object数据

### Aliyun Oss Finder

`aliyunOSSFinder/bucketName/objectName`，此处假设初始化了个名为`aliyunOSSFinder`的阿里云oss finder。它的参数与minio finder一致

## 任务串连

使用pipeline的形式，将相关的图片处理任务串连起来，则可以实现图片的各类优化，一般的处理流程为：

1、拉取图片数据
2、对图片做缩放、水印等处理
3、对图片做压缩格式转换处理

例如下面的为将百度的图标缩小再转换为avif的处理：

```go
// 示例代码忽略了err
jobs, _ := imagepipeline.Parse("proxy/https%3A%2F%2Fwww.baidu.com%2Fimg%2FPCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png|fitResize/100/80|optimize/192.168.1.1:6002/80/avif", "")
img, _ := imagepipeline.Do(context.Background(), nil, jobs...)
fmt.Println(img)
```
