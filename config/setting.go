package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Prefix string

type Database struct {
	Username string
	Password string
	DBname   string
}

type User struct {
	Username string
	Password string
}

type Setting struct {
	Database Database
	User     User
	Commands []*Command
}

type Command struct {
	Name        string
	Description string
	SQL         string
	Query       bool
	Params      int
	Comments    []string
	Input       []interface{} `json:"-"`
}

func Get() Setting {
	return setting
}

var setting Setting

func LoadSetting(file string) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	byte, _ := ioutil.ReadAll(f)

	if err = json.Unmarshal(byte, &setting); err != nil {
		panic(err)
	} else {
		for _, c := range setting.Commands {
			compileCommand(c)
		}
	}
}

func AddCmd(cmd *Command) (l int) {
	compileCommand(cmd)
	l = len(setting.Commands)
	setting.Commands = append(setting.Commands, cmd)
	return
}

func CheckUser(username, password string) bool {
	if username == setting.User.Username && password == setting.User.Password {
		return true
	} else {
		return false
	}
}

var cmdInput = `
<div class="input-group" style="margin: 5px;">
<span class="input-group-prepend"><a href="#" class="btn btn-default disabled">%s</a></span>
<input type="text" class="form-control" id="execmd-%d" required>
</div>
`

func compileCommand(c *Command) {
	l := len(c.Comments)
	for i := 0; i < c.Params; i++ {
		if i >= l {
			c.Comments = append(c.Comments, "")
		}
		input := fmt.Sprintf(cmdInput, c.Comments[i], i)
		c.Input = append(c.Input, input)
	}
}
