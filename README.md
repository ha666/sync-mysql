# sync-mysql
同步mysql数据库

### 配置

#### 应用配置
> page_size看情况配置，推荐是50~500，
> thread_count根据数据库服务器的配置以及网络情况来定

```yaml
  page_size: 20
  thread_count: 4
```

#### 源配置(数据库和kafka二选一)

```yaml
  database:
    name: "abc"
    address: "127.0.0.1"
    port: 3306
    account: "root"
    password: "1234567890"
```

```yaml
  kafka:
    version: "2.6.0"
    addresses:
     - "127.0.0.1:9092"
    topic: "maxwell"
    consumer: "maxwell-consumer"
    databaseName: "test_db"
```

#### 目标配置

```yaml
  databases:
    - name: "ha666db"
      address: "127.0.0.1"
      port: 3306
      account: "root"
      password: "1234567890"
```

```yaml
  kafka:
    - addresses:
        - "127.0.0.1:9092"
      topic: "db-log-2"
```

#### 表名配置

```yaml
mapping:
  tables:
    "t_abcd": "t_cd"
```

#### 配置示例一：从数据库到数据库

```yaml
app:
  page_size: 20
  thread_count: 4
source:
  database:
    name: "abc"
    address: "127.0.0.1"
    port: 3306
    account: "root"
    password: "1234567890"
target:
  databases:
    - name: "ha666db"
      address: "127.0.0.1"
      port: 3306
      account: "root"
      password: "1234567890"
```

#### 配置示例二：从kafka到数据库

> 使用的是maxwell的方案，具体的使用方法可以参考[https://github.com/zendesk/maxwell](https://github.com/zendesk/maxwell)

```yaml
app:
  page_size: 20
  thread_count: 4
source:
  kafka:
    version: "2.6.0"
    addresses:
     - "127.0.0.1:9092"
    topic: "maxwell"
    consumer: "maxwell-consumer"
    databaseName: "test_db"
target:
  databases:
    - name: "ha666db"
      address: "127.0.0.1"
      port: 3306
      account: "root"
      password: "1234567890"
```

#### 注意：
+ 如果源配置和目标配置都有数据库，那数据库配置不能完全相同
+ 如果源配置和目标配置都有kafka，那kafka配置不能完全相同
+ 只同步数据，不同步结构

