## Filetransfer工具开发

### 1：前提条件

|                     环境                     | 是否必须 |
| :------------------------------------------: | :------: |
|                  Go开发环境                  |    是    |
| 阿里云/腾讯云/华为云/AWS 登任何一个云的AK/SK |    是    |
|                   Bucket桶                   |    是    |
|               Goland / Vs Code               |    否    |
|               Endpoint地域域名               |    是    |

### 2：开发前提

```shell
# 我这里使用阿里云的OSS来开发
1：创建项目工程
go mod init filetransfer

2：下载阿里云OSS的SDK
go get github.com/aliyun/aliyun-oss-go-sdk/oss

3：创建测试目录，测试阿里云sdk的可用性
mkdir example && cd example
# 创建一个main.go的文件，然后开始编写测试内容
```

```go
package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
    // 这个是阿里云oss的endpoint的地址
	endpoint    = "oss-cn-shanghai.aliyuncs.com"  // 这里写自己的OSS桶的地区
    // 阿里云的AK
	accessKeyID = "LTAIxxxxxxxxxxxxAZikw"    // 这里写自己的AK
    // 阿里云的SK
	accessKey   = "k0HxxxxxxxxxxPl2RJopKaC"  // 这里写自己的SK
)

func main() {
	client, err := oss.New(endpoint, accessKeyID, accessKey)  // 创建OSSClient实例。
	if err != nil {
		panic(err)
	}
	lsRes, err := client.ListBuckets()  // 列举存储空间。
	if err != nil {
		panic(err)
	}

	for _, bucket := range lsRes.Buckets {  // 遍历存储空间。
		fmt.Println("Bucket:", bucket.Name)
	}
}
```

```shell
运行命令：go run .\main.go
运行结果：Bucket: layzer

这里拿到了我们bucket的名称，我的名称就叫layzer，所以这足以证明SDK是没有问题的，但是这个其实只是我们在测试程序，我们正式去开发的时候不会这么写的
```

### 3：开始编写程序

```shell
1：创建main.go以及所需要的目录
```

```go
package main

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

// 声明变量（后面走配置文件的方式）
var (
	endpoint     = ""   // 自己测试填写
	accessKey    = ""   // 自己测试填写
	secretSecret = ""   // 自己测试填写
	bucketName   = ""   // 自己测试填写
	uploadFile   = ""   // 自己测试填写
)

func validate() error { // 验证参数
	if endpoint == "" || accessKey == "" || secretSecret == "" {
		return fmt.Errorf("endpoint, accessKey, secretSecret must be set")
	}
	if uploadFile == "" {
		return fmt.Errorf("upload file path required")
	}
	return nil
}

// 实现文件上传
func upload(filepath string) error {
	//1：实例化一个bucket的client
	client, err := oss.New(endpoint, accessKey, secretSecret) // 创建OSSClient实例
	if err != nil {
		return err
	}
	//2：获取bucket对象名称
	bucket, err := client.Bucket(bucketName) // 获取存储空间
	if err != nil {
		return err
	}
	//3：上传文件
	return bucket.PutObjectFromFile(uploadFile, filepath) // 上传文件
}

func main() {
	if err := validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := upload(uploadFile); err != nil { // 上传文件
		fmt.Println("upload file failed, err:", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s upload success", uploadFile)
	}
}
```

```shell
运行命令：go run .\main.go
运行结果：gin.png upload success
```

![image](https://img2022.cnblogs.com/blog/2222036/202211/2222036-20221121010510226-1528234844.png)

```shell
这样我们的小工具其实就可以正常使用了，但是我们要知道，工具嘛肯定不能一直以main.go的方式运行，我们肯定得编译成包，然后以二进制的方式运行，所以我们少不了对工具有些配置，比如配置我们上面变量上的一些东西按照配置文件的方式传递给我们的程序。
```

### 4：读取配置文件

```shell
以前看过我Go学习笔记的应该有印象，里面有一个包叫做flag，它可以使我们的工具变成一个cli，那么我们下面就来看看怎么把我们这个main.go变成一个可以和用户交互的CLI文件上传工具，我们再次完善一下代码
```

```go
package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

// 声明变量
var (
	endpoint     = "oss-cn-shanghai.aliyuncs.com"
	accessKey    = "LTAI5txxxxxxxxxxxxZikw"
	secretSecret = "kxxxxxxxxxxxxxxxxxxxxpKaC"
	bucketName   = "layzer"
	uploadFile   = ""
)

// 使需要上传的文件参数化
func loadParams() {
	flag.StringVar(&uploadFile, "f", "", "set oss upload file")
	flag.Parse()
}

func validate() error { // 验证参数
	if endpoint == "" || accessKey == "" || secretSecret == "" {
		return fmt.Errorf("endpoint, accessKey, secretSecret must be set")
	}
	if uploadFile == "" {
		return fmt.Errorf("upload file path required")
	}
	return nil
}

// 实现文件上传
func upload(filepath string) error {
	//1：实例化一个bucket的client
	client, err := oss.New(endpoint, accessKey, secretSecret) // 创建OSSClient实例
	if err != nil {
		return err
	}
	//2：获取bucket对象名称
	bucket, err := client.Bucket(bucketName) // 获取存储空间
	if err != nil {
		return err
	}
	//3：上传文件
	return bucket.PutObjectFromFile(uploadFile, filepath) // 上传文件
}

func main() {
	// 加载参数
	loadParams()
	// 验证参数
	if err := validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 上传文件
	if err := upload(uploadFile); err != nil { // 上传文件
		fmt.Println("upload file failed, err:", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s upload success", uploadFile)
	}
}
```

```shell
运行命令：go run .\main.go -f main.go
运行结果：main.go upload success
解释：这里使用flag解析-f的参数传递给了程序，程序去读取传递过来的参数赋值给&uploadFile，然后就等于uploadFile的值变成了用户自定义，这个时候我们后面就可以随意跟我们的参数了，当然还是按照规范，-f就是 --file指定上传的文件名称，当然了，这个其实还是不完美的，因为用的人不同，我们要传递的ep，ak，sk，bucketname都是不同的，所以我们可以使用两种方式来解决这个问题，

1：读取变量来解决这个问题
2：使用参数化传递配置到程序内

帮助测试：go run main.go -h
测试结果：Usage of C:\Users\ADMINI~1\AppData\Local\Temp\go-build2287442688\b001\exe\main.exe:
  -f string
        set oss upload file
        
        
我们最后修改一下配置文件的获取方式然后就可以打包测试了
```

```go
package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

// 声明变量
var (
	endpoint     = os.Getenv("OSS_ENDPOINT")   // OSS endpoint
	accessKey    = os.Getenv("OSS_ACCESS_KEY")	// OSS accessKey
	secretSecret = os.Getenv("OSS_SECRET_KEY")	// OSS secretKey
	bucketName   = os.Getenv("OSS_BUCKET_NAME")  // OSS bucketName
	uploadFile   = ""
)

// 使需要上传的文件参数化
func loadParams() {
	flag.StringVar(&uploadFile, "f", "", "set oss upload file")
	flag.Parse()
}

func validate() error { // 验证参数
	if endpoint == "" || accessKey == "" || secretSecret == "" {
		return fmt.Errorf("endpoint, accessKey, secretSecret must be set")
	}
	if uploadFile == "" {
		return fmt.Errorf("upload file path required")
	}
	return nil
}

// 实现文件上传
func upload(filepath string) error {
	//1：实例化一个bucket的client
	client, err := oss.New(endpoint, accessKey, secretSecret) // 创建OSSClient实例
	if err != nil {
		return err
	}
	//2：获取bucket对象名称
	bucket, err := client.Bucket(bucketName) // 获取存储空间
	if err != nil {
		return err
	}
	//3：上传文件
	return bucket.PutObjectFromFile(uploadFile, filepath) // 上传文件
}

func main() {
	// 加载参数
	loadParams()
	// 验证参数
	if err := validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 上传文件
	if err := upload(uploadFile); err != nil { // 上传文件
		fmt.Println("upload file failed, err:", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s upload success", uploadFile)
	}
}

```



```shell
# 打包程序并测试
打包Linux须知：set GOOS=linux
打包命令：go build -o oss-cli main.go

# 将程序上传到Linux
[root@cdk-server ~]# chmod +x oss-cli 
[root@cdk-server ~]# ./oss-cli -h
Usage of ./oss-cli:
  -f string
    	set oss upload file
# 因为前面我们配置的是获取环境变量，所以我们需要提前配置好环境变量

[root@cdk-server ~]# cat /etc/profile | tail -n 4
export OSS_ENDPOINT="oss-cn-shanghai.aliyuncs.com"
export OSS_ACCESS_KEY="LTAI5txxxxxxxxxxxxxxZikw"
export OSS_SECRET_KEY="k0Hxxxxxxxxxxxxxxxxx2RJopKaC"
export OSS_BUCKET_NAME="layzer"

[root@cdk-server ~]# source /etc/profile

# 创建测试文件并使用程序上传

[root@cdk-server ~]# touch myfile
[root@cdk-server ~]# ./oss-cli -f myfile 
myfile upload success
```

![image](https://img2022.cnblogs.com/blog/2222036/202211/2222036-20221121024615597-965118487.png)

### 5：增加文件下载连接

```shell
上面我们看到的是文件上传的配置，那么，既然我们上传了文件，肯定是提供别人下载的，那么我们是否可以顺便打印出文件的下载链接呢？这个是毋庸置疑的，当然是可以的，那么我们来看看如何打印文件的下载链接，其实下载链接是可以直接拼接的，我们就简单的拼接一下就OK了
```

```go
package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

// 声明变量
var (
	endpoint     = os.Getenv("OSS_ENDPOINT")    // OSS endpoint
	accessKey    = os.Getenv("OSS_ACCESS_KEY")  // OSS accessKey
	secretSecret = os.Getenv("OSS_SECRET_KEY")  // OSS secretKey
	bucketName   = os.Getenv("OSS_BUCKET_NAME") // OSS bucketName
	uploadFile   = ""
)

// 使需要上传的文件参数化
func loadParams() {
	flag.StringVar(&uploadFile, "f", "", "set oss upload file")
	flag.Parse()
}

func validate() error { // 验证参数
	if endpoint == "" || accessKey == "" || secretSecret == "" {
		return fmt.Errorf("endpoint, accessKey, secretSecret must be set")
	}
	if uploadFile == "" {
		return fmt.Errorf("upload file path required")
	}
	return nil
}

// 实现文件上传
func upload(filepath string) error {
	//1：实例化一个bucket的client
	client, err := oss.New(endpoint, accessKey, secretSecret) // 创建OSSClient实例
	if err != nil {
		return err
	}
	//2：获取bucket对象名称
	bucket, err := client.Bucket(bucketName) // 获取存储空间
	if err != nil {
		return err
	}
	//3：上传文件
	return bucket.PutObjectFromFile(uploadFile, filepath) // 上传文件
}

func main() {
	// 加载参数
	loadParams()
	// 验证参数
	if err := validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// 上传文件
	if err := upload(uploadFile); err != nil { // 上传文件
		fmt.Println("upload file failed, err:", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s upload success", uploadFile)
		fmt.Println("下载地址：", "https"+"://"+bucketName+"."+endpoint+"/"+uploadFile)  // 这里拼接一下就OK了
	}
}
```

```shell
测试单文件和目录

[root@cdk-server ~]# ./oss-cli -f aliyun 
aliyun Upload Success！！！
Download URL: https://layzer.oss-cn-shanghai.aliyuncs.com/aliyun

[root@cdk-server ~]# ./oss-cli -f 1/2/3/aliyun 
1/2/3/aliyun Upload Success！！！
Download URL: https://layzer.oss-cn-shanghai.aliyuncs.com/1/2/3/aliyun
```

