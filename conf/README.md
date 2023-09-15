## conf

Parsing yaml, json, toml configuration files to go struct.

<br>

### Example of use

```go
    import "github.com/18721889353/common_pkg/conf"

    // Way 1: No listening profile
	config := &App{}
	err := conf.Parse("test.yml", config)

    // Way 2: Enable listening profile
config := &App{}
	fs := []func(){
		func() {
			fmt.Println("Listening for updates to the configuration file")
		},
	}
	err := conf.Parse("test.yml", config, fs...)
```
