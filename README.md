# tantowi/gini

[![Go Reference](https://pkg.go.dev/badge/github.com/tantowi/gini.svg)](https://pkg.go.dev/github.com/tantowi/gini)
  ![GitHub tag](https://img.shields.io/github/v/tag/tantowi/gini?label=version)

<br>

Module `tantowi/gini` is a Go module to read INI format configuration file
<br>

[![Go](https://github.com/tantowi/gini/actions/workflows/action-on-push.yml/badge.svg)](https://github.com/tantowi/gini/actions/workflows/action-on-push.yml)
<br>

## HOW TO USE

If we have configuration file in INI format like this :

```ini
[setting]
color=red
width=700
height=450

[server]
host=10.10.20.20
port=3344
```

First, import the library

```go
import (
   "github.com/tantowi/gini"
)
```

Then load the file :

```go
ini, err := gini.LoadFile("/etc/myapp/config.ini")
if err != nil {
	fmt.Println(err.Error())
	return
}
```

And read the value :

```go
color  := ini.Read("setting", "color")
width  := ini.Read("setting", "width")
height := ini.Read("setting", "height")

host := ini.Read("server", "host")
port := ini.Read("server", "port")
```

<br><br>

## Copyright

Copyright &copy; 2020 Tantowi Mustofa, ttw@tantowi.com<br>
Licensed under [MIT License](LICENSE)
