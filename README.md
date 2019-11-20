  <br>


>**&ensp;&ensp;&ensp;邮箱:mikesen1994@gmail.com   &ensp;&ensp;&ensp;&ensp;  &ensp;&ensp; &ensp;&ensp; vx:965952482**

<br>
&ensp;&ensp;&ensp;本程序是模仿比特币的功能所编写的区块链公链demo,主要应用到了密码学,共识算法,对等网络,区块链防篡改结构等相关知识,并把各个知识点结合到一起,编写成了简单完善的可运行公链demo

<hr>

### 程序特点：

- 基于工作量证明共识算法，数据以区块链的结构进行存储
- 去中心化，运用P2P技术使各个节点之间相对独立
- 主动寻找网络中的对等节点，自动连接并存入本地节点池
- 节点退出时会向全网广播，其余节点动态更新当前可连接节点池
- 挖矿成功节点获得记账权，并向全网广播同步最新区块，其余节点验证通过后存入本地区块链中
- 交易转帐使用UTXO交易模型,支持一次交易存在多笔转账
- 支持中文助记词导入，由助记词生成公私钥密钥对（使用的椭圆曲线算法）
- 交易转账使用私钥进行数字签名，公钥验证,并因为UTXO的结构避免了对于签名的重放攻击问题
- 为未花费UTXO单独建立数据表，优化转账交易速度
- 使用默克尔树生成交易的根hash（当前demo并没有区分区块头与区块体，仅仅是想使用此数据结构练练手）
- 持久化区块链与公私钥信息，存入节点本地数据库中（每个节点拥有自己的独立数据库）
- 自定义挖矿难度值、旷工挖矿奖励值
- 自定义交易池大小，满足指定笔数的交易后才会开始挖矿

<hr>

### 主要模块：
 - 命令调度模块
 - UTXO交易生成模块
 - 密码学加解密模块
 - 区块生成、验证模块
 - 数据持久化模块
 - P2P网络通讯模块
 - 日志输出模块
 <br>
 
####  命令调度模块
&ensp;&ensp;&ensp;  启动程序后，控制台捕捉用户输入信息，通过对用户的输入解析出命令以及跟随在命令后的值。根据不同命令对程序进行相关操作
<br>
####  UTXO交易生成模块
&ensp;&ensp;&ensp;交易转账模块基于UTXO模型，但并没有引入比特币脚本，脚本处直接使用数字签名的字节数组进行替代。当用户A转账给用户B时，需要用户A使用私钥对"**输入**" (包含了用户A所拥有的"**输出**"交易hash、索引等信息）进行数字签名，生成交易后发送给其他节点，其他节点则使用用户A的公钥对其进行签名验证。</br>&ensp;&ensp;&ensp;由于UTXO的特殊结构，天然的避免了重放攻击，并不需要像以太坊账户系统一样添加nonce值，但是为了避免UTXO的重复计算问题，在上一笔转账未打包进区块之前暂不支持同一地址的再次转账</br>&ensp;&ensp;&ensp;支持一笔交易多笔转账，并为了优化转账查询速度创建了UTXO数据表专门用于存储所有区块链中未花费的**输出**。</br>&ensp;&ensp;&ensp;想要了解更多关于UTXO相关，建议参考[这篇文章](https://draveness.me/utxo-account-models)
<br>
####  密码学加解密模块
   1. 单向散列函数：sha256  ripemd-160
   
    主要用于将整体区块通过计算转换为固定长度的字符串,方便进行数据校验
   2. 编解码算法：base58 
   
    由于私钥原始长度过长不利于记忆，使用base58编码对私钥、地址进行可视化编码
   3. 非对称加密：椭圆曲线算法（crypto/elliptic p256）
   
    通过助记词文本提取7对中文词语作为种子，通过使用椭圆曲线算法生成公私钥密钥对，私钥用于对交易数据进行数字签名，公钥对签名进行验证来确保发起人身份。
    公钥通过一系列运算生成地址，地址用于查询余额，以及接收转账Token
&ensp;&ensp;&ensp;地址生成规则如下：
    
 - 通过椭圆曲线算法生成公钥
 - 对公钥进行sha256散列和ripemd160散列,获得publickeyHash
 - 在publickeyHash前面加上version(版本)字节数组获得versionPublickeyHash
 - 对versionPublickeyHash进行两次sha256散列并取前４位字节，获得tailfHash
 - 将tailfHash拼接到versionPublickeyHash后面，获得公钥的最终Hash即finalHash
  - 最后将finalHash进行Base58编码得到比特币地址

	
> 曾经有个疑问，为何比特币生成地址要这么麻烦，既然非对称加密只拥有公钥是无法倒推出私钥的，为何不直接使用公钥当地址，而是对公钥进行hash多次来取得地址，直到最近看了篇文章才明白，该文章提到量子计算机是可以破解椭圆曲线加密的，其可以通过公钥快速寻找到私钥信息。但是量子计算机很难逆转Hash算法(或者说需要2的80次方个步骤来破解Hash)，所以你的比特币放在一个未支付过的地址中(根据UTXO交易模型，输出存的是公钥Hash而不是公钥,这同样解释了为何UTXO输入存的是公钥而输出存的是公钥Hash)是相当安全的。也就是说已有花费的地址在面对量子计算机面前是不安全的，没有花费的地址有较强的抗量子性。
####  区块生成、验证
&ensp;&ensp;&ensp;  基于POW共识算法生成区块，首先根据难度值（可在配置文件里定义）来定义挖矿难度（一串大数），通过调用go自身的随机数包crypto/rand来不断的变换随机数nonce(上个版本用的nonce值自身累加的方法,但是分叉的概率太大)，不断哈希区块自身来使最终计算出来的区块自身hash值小于当前定义的挖矿难度则获得出块权利。</br>&ensp;&ensp;&ensp;  出块节点可获得奖励代币并拥有记账权，出块后像全网进行广播。其余P2P节点收到区块后首先对区块自身hash进行验证，其次检验区块里的prehash与本地的前区块hash是否一致，最后存入本地数据库中。
<br>
####  数据持久化模块
&ensp;&ensp;&ensp; 持久化层基于KV型数据库blot多封装了一层,主要接口为put、view、delete。每次调用接口会单独打开、关闭数据库的句柄，所以不会出现被其他线程占用的情况。数据库分别建立了三个表 BlockBucket（用于存放区块的详细信息）、AddrBucket（用于存放本地钱包数据）、UTXOBucket（用于存放未消费UTXO数据）
<br>
####  P2P网络通讯模块
&ensp;&ensp; 使用适合局域网寻址的mdns技术，由于所使用的包在windows下存在找不到网络的bug，所以本程序建议在linux/mac下运行</br>              &ensp;&ensp; 节点启动后会自动在局域网中寻找其他对等节点，发现后会存放在节点池中（存于内存），节点之间相互通讯的数据前十二个字节默认为命令，根据命令不同来对本地的区块链相关信息进行反馈 </br>&ensp;&ensp;&ensp;主要运行原理为分发区块与收到交易后的挖矿：</br>
</br>&ensp;&ensp;获取区块流程：
 1. 互相对比区块高度
 2. 获取缺失的区块hash
 3. 通过区块hash来接收缺失的整个区块
 4. 区块验证，存入数据库

&ensp;&ensp;挖矿流程：
 1. 通过某个节点发送交易数据到全网节点
 2. 节点接收到交易，对交易进行签名验证,余额验证
 3. 验证通过后存入交易池，满足交易池大小后开始挖矿
 4. 挖矿成功，全网广播区块高度
 5. 发送区块到其他节点
 6. 其他节点进行区块验证，存入数据库
####  日志输出模块
&ensp;&ensp;&ensp; 使用自制的log包，程序启动后会默认在当前目录下(可在配置文件设置)生成log+端口号的日志文件，所有程序产生的debug信息都会打印到此日志文件中，建议开启一个窗口进行实时监听以方便观察节点之间的交互，以及区块生成的详细步骤

【日志包特点】：

   - 支持定向输出日志到指定文件
   - 支持一键隐藏调试信息
   - 支持彩色打印（windows/linux/mac均支持）
   - 显示输出日志的类名、函数/方法名
 
<br>
<hr>

主要使用的工具包
---------------------------
包     | 用途
-------- | -----
[github.com/boltdb/bolt](https://github.com/boltdb/bolt)| k,v型数据库
[github.com/spf13/viper](https://github.com/spf13/viper)  | 配置文件读取工具
[github.com/golang/crypto](https://github.com/golang/crypto)| 密码学相关工具
[github.com/libp2p/go-libp2p](https://github.com/libp2p/go-libp2p)  | ipfs旗下的p2p通讯工具
[github.com/corgi-kx/logcustom](https://github.com/corgi-kx/logcustom)| 日志输出工具

 <hr>
 
### 程序运行教程：

**1.下载后编译**

本demo建议在linux/mac下运行，否则会出现助记词乱码，找不到对等网络的问题

```shell
 git clone https://github.com/corgi-kx/blockchain_golang.git
```
```shell
 go build -mod=vendor -o chain main.go
```
<br>

**2.打开多个窗口**

为了简化操作，在同一台电脑中启动不同端口来模拟P2P节点（三个窗口用于启动程序，三个窗口用于实时查看日志）
>实机操作时，如果出现找不到其他节点情况可能是防火墙问题，请关闭防火墙后在试

![在这里插入图片描述](https://img-blog.csdnimg.cn/20191118103707708.png)

<br>

**3.修改配置文件**
  
  主要修改本地监听ip，本地监听端口。其他的默认即可</br>
  不建议调小难度阀值，避免产生区块分叉情况，demo暂未对区块分叉做处理
```shell
 vi config.yaml
```
```yaml
blockchain:
  #挖矿难度值,越大越难挖
  mine_difficulty_value: 24
  #挖矿奖励代币数量
  token_reward_num: 25
  #交易池大小(满足多少条交易才开始进行挖矿)
  trade_pool_length: 2
  #日志存放路径
  log_path: "./"
  #中文助记词种子路径
  chinese_mnemonic_path: "./chinese_mnemonic_world.txt"
network:
  #本地监听IP
  listen_host: "192.168.0.164"
  #本地监听端口
  listen_port: "9000"
  #节点组唯一标识名称(如果节点间名称不同会找不到网络)
  rendezvous_string: "meetme"
  #网络传输流的协议id(如果节点间id不同发送不了数据)
  protocol_id: "/chain/1.1.0"

```

<br>

**4.启动节点,创建钱包,生成创世区块**

启动节点1
```shell
 ./chain
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191118101305498.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

通过命令,先生成三个钱包地址

```
> generateWallet
助记词： ["肺段","生地","齿槽","几维","中葡","芒鱼","光华"]
私钥： 6HrLjHE4Qm31dZFGjemwNLZM3iqnxoSUqKb5VtEKbWzh
地址： 12BwtcVWimms9rrKxxoCev68woGyMYS4sk
> generateWallet
助记词： ["扭伤","剪创","肌病","下陷","广发","浊音","斜疝"]
私钥： 7yBRSB46q8ZeEiYbwZDSvKzzsh1MYAygeo2i689uEMAf
地址： 1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS
> generateWallet
助记词： ["心室","缺缸","瓣胃","黑茶","份额","张铜","回游"]
私钥： 872CCeLS8bDrC7bdSoFrgUSWm57eqTdypEhKbErYC9xi
地址： 1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD
```
生成创世区块(赋予第一个地址100Tokens)
```
> genesis -a 12BwtcVWimms9rrKxxoCev68woGyMYS4sk -v 100
已成生成创世区块
```

日志1实时查看日志(可以看到挖矿过程)
```shell
 tail -f log9000.txt 
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191118144251486.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

<br>

**5.同步区块**

节点2,节点3依次修改配置文件的端口号为9001,9002,启动这两个节点来同步创世区块</br>
这时节点1的日志监测到网络中存在的其他节点
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191118145703154.png)节点2,节点3 启动后会自动同步创世区块
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191118145752942.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

<br>

**6.进行转帐操作**

每个节点设置挖矿奖励地址(也可以不设置,不设置的情况下,节点挖到矿后不会产生奖励)</br>
节点1设置挖矿奖励地址:
```
> setRewardAddr -a 12BwtcVWimms9rrKxxoCev68woGyMYS4sk
已设置地址12BwtcVWimms9rrKxxoCev68woGyMYS4sk为挖矿奖励地址！
```
节点2设置挖矿奖励地址:
```
> setRewardAddr -a 1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS
已设置地址1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS为挖矿奖励地址！
```
节点3设置挖矿奖励地址:
```
> setRewardAddr -a 1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD
已设置地址1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD为挖矿奖励地址！
```
节点1进行转帐操作(创世地址像其他两个地址每个转帐10Tokens)
```
> transfer -from ["12BwtcVWimms9rrKxxoCev68woGyMYS4sk","12BwtcVWimms9rrKxxoCev68woGyMYS4sk"] -to ["1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS","1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD"] -amount [10,10]
已执行转帐命令
```
![在这里插入图片描述](https://img-blog.csdnimg.cn/2019111815314125.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzM1OTExMTg0,size_16,color_FFFFFF,t_70)

<br>

**7.查看余额**

三个节点中,由节点2挖到区块,里所应当节点2获得挖矿奖励25Tokens
![在这里插入图片描述](https://img-blog.csdnimg.cn/20191118153547470.png)此时在任意节点敲入`getBalance`查看余额命令,可以查看三个地址的余额信息
```
> getBalance -a 12BwtcVWimms9rrKxxoCev68woGyMYS4sk
地址:12BwtcVWimms9rrKxxoCev68woGyMYS4sk的余额为：80
> getBalance -a 1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS
地址:1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS的余额为：35
> getBalance -a 1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD
地址:1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD的余额为：10
```

<br>

**8.查看区块详细信息**

任意节点输入`printAllBlock`命令查看区块信息

区块1为创世区块,只有赋予`12BwtcVWimms9rrKxxoCev68woGyMYS4sk`的100UTXO输出

可以看到区块2:</br>
第一笔交易,地址`12BwtcVWimms9rrKxxoCev68woGyMYS4sk`先花掉创世区块额度为100的UTXO,给自身生成一个90UTXO,给地址`1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS` 生成10UTXO</br>
第二笔交易地址`12BwtcVWimms9rrKxxoCev68woGyMYS4sk`使用第一笔交易输出的90额度的UTXO,给自身生成一个80UTXO 以及给地址`1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD` 生成10UTXO</br>
第三笔交易为挖矿奖励交易,所以只有输出,没有输入,给地址`1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS` 生成25UTXO(在配置文件中设置的25奖励额度)

```
> printAllBlock                                   
========================================================================================================
本块hash         00000008acfb9a8dcf3bb923f4eb6f2ddfc27dcaff861ea6848a9074ca46d85b
        ------------------------------交易数据------------------------------
         本次交易id:  988ecbe7f374855aa94addb873f22960cf43646bdaeb562533f3e683478270db
          tx_input：
                        交易id:  bb717bd6717c8cae3829875187b97f256859277ad4a52ac57cdbc132895ca154
                        索引:    0
                        签名信息:    8c8b0628ceadebbc9e97b490a40a23494d3f8286f1af045f1e1f18d529c49a90afa194799182c264ee15871b5dd35c773e5dd46427fc8e2c268356ce09f6b60b
                        公钥:    8e0f1fe7d6177f11027818663048392cee8952cefcf1ceeec8edc84e176f46cedd338575f709b412eeab904d7027056354038f8aef7a1940f45264f7116ba793
                        地址:    12BwtcVWimms9rrKxxoCev68woGyMYS4sk
          tx_output：
                        金额:    90    
                        公钥Hash:    0d0a1aeb1baf838828a54ac97b09524f0b0c3210    
                        地址:    12BwtcVWimms9rrKxxoCev68woGyMYS4sk
                        ---------------
                        金额:    10    
                        公钥Hash:    6eb2d1846217aa089dfa26e3147b767e1de0b08d    
                        地址:    1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS
         本次交易id:  443b4a4f04204bd8ed2bdfcc096642a27457c27aa47c2ee81486d7440b059521
          tx_input：
                        交易id:  988ecbe7f374855aa94addb873f22960cf43646bdaeb562533f3e683478270db
                        索引:    0
                        签名信息:    2a064297227ba07c7ea92eebb1d43f3fe4dfbd6c7e78be8ec2d1d30e20fa51500c8bdb591c11908a877aeef61b4c64f9a851cc44af441cbe6893e1b80e42032c
                        公钥:    8e0f1fe7d6177f11027818663048392cee8952cefcf1ceeec8edc84e176f46cedd338575f709b412eeab904d7027056354038f8aef7a1940f45264f7116ba793
                        地址:    12BwtcVWimms9rrKxxoCev68woGyMYS4sk
          tx_output：
                        金额:    80    
                        公钥Hash:    0d0a1aeb1baf838828a54ac97b09524f0b0c3210    
                        地址:    12BwtcVWimms9rrKxxoCev68woGyMYS4sk
                        ---------------
                        金额:    10    
                        公钥Hash:    8fa79c32a067830be3b16ade637d370e1d1e6e0d    
                        地址:    1E6aRBxfncAsypUnjGxPJYbR4JQ3gZ6hHD
         本次交易id:  2420c67272ab7832d6148a36a6b38166862d12e265f184439e6ab2e606b01245
          tx_input：
          tx_output：
                        金额:    25    
                        公钥Hash:    6eb2d1846217aa089dfa26e3147b767e1de0b08d    
                        地址:    1B6KYdABXZDwq8xGTbdDknpHBo11CkihxS
        --------------------------------------------------------------------
时间戳           2019-11-18 03:23:57 PM
区块高度         2
随机数           2808567053068705071
上一个块hash     0000007d7b7c7b540d9d1b0d1d06b6936e1bc613f6ab7de1ae0275cdaef4e4a4
========================================================================================================
本块hash         0000007d7b7c7b540d9d1b0d1d06b6936e1bc613f6ab7de1ae0275cdaef4e4a4
        ------------------------------交易数据------------------------------
         本次交易id:  bb717bd6717c8cae3829875187b97f256859277ad4a52ac57cdbc132895ca154
          tx_input：
                        交易id:  
                        索引:    -1
                        签名信息:    
                        公钥:    
                        地址:    
          tx_output：
                        金额:    100    
                        公钥Hash:    0d0a1aeb1baf838828a54ac97b09524f0b0c3210    
                        地址:    12BwtcVWimms9rrKxxoCev68woGyMYS4sk
        --------------------------------------------------------------------
时间戳           2019-11-18 10:43:41 AM
区块高度         1
随机数           8604076799988393002
上一个块hash     0000000000000000000000000000000000000000000000000000000000000000
========================================================================================================

```

<br>

**9.其他**

你也可以在节点2,节点3发起转帐,不过首先需要通过助记词导入钱包信息,例子如下:
```
> importMnword -m ["扭伤","剪创","肌病","下陷","广发","浊音","斜疝"]
```
<br>
更多功能请自行发掘 :)

 <br>
 


