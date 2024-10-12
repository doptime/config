
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

[Http]
  CORES = "*"
  Port = 80
  Path = "/"
	JwtSecret = ""
  AutoDataAuth = false
  MaxBufferSize = 10485760
  SUToken = ""

[[APISource]]
	Name    = "doptime"
  UrlBase = "https://api.doptime.cc"
	ApiKey  = "yourapikeyhere"

[Log]
  LogLevel = 1

```