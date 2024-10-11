
Example config.toml:

```toml title="config.toml"
ConfigUrl = ""

[[Redis]]
  Name = "default"
  Username = "doptime"
  Password = "yourpasswordhere"
  Host = "drangonflydb.local"
  Port = 6379
  DB = 0

[Jwt]
  Secret = ""
  Fields = "*"

[Http]
  CORES = "*"
  Port = 80
  Path = "/"
  Enable = false
  MaxBufferSize = 10485760

[Api]
  ServiceBatchSize = 64

[Data]
  AutoAuth = false

[Setting]
  LogLevel = 1
  SUToken = ""

```