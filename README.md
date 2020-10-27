# sync-mysql
同步mysql数据库

### 配置

#### 应用配置

```yaml
  page_size: 20
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
    version: "0.10.2.1"
    addresses:
      - "127.0.0.1:9092"
    topics:
      - "db-log-1"
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
  kafkas:
    - version: "0.10.2.1"
      addresses:
        - "127.0.0.1:9092"
      topic:
        - "db-log-2"
      consumer: "consumer-obj"
```

#### 配置示例一：从数据库到数据库

```yaml
app:
  page_size: 20
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

```yaml
app:
  page_size: 20
source:
  kafka:
    version: "0.10.2.1"
    addresses:
      - "127.0.0.1:9092"
    topics:
      - "db-log-1"
    consumer: "consumer-obj"
target:
  databases:
    - name: "ha666db"
      address: "127.0.0.1"
      port: 3306
      account: "root"
      password: "1234567890"
```

#### 配置示例三：从数据库到kafka

```yaml
app:
  page_size: 20
source:
  database:
    name: "abc"
    address: "127.0.0.1"
    port: 3306
    account: "root"
    password: "1234567890"
target:
  kafkas:
    - version: "0.10.2.1"
      addresses:
        - "127.0.0.1:9092"
      topics:
        - "db-log-2"
```

#### 配置示例四：从数据库到数据库和kafka

```yaml
app:
  page_size: 20
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
  kafkas:
    - version: "0.10.2.1"
      addresses:
        - "127.0.0.1:9092"
      topics:
        - "db-log-2"
```


#### 注意：
+ 如果源配置和目标配置都有数据库，那数据库配置不能完全相同
+ 如果源配置和目标配置都有kafka，那kafka配置不能完全相同
+ 只同步数据，不同步结构

