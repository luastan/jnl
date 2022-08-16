/*
Copyright © 2022 Luastan <@luastan>
*/
package main

import (
	"context"
	"flag"
	"github.com/luastan/jnl/cmd"
	"github.com/luastan/jnl/jnl"
)

var (
	validCmds = map[string]struct{}{
		"getHistory": {},
	}
)

func main() {
	flag.Parse()
	command := flag.Args()[0]
	_, isJNLCmd := validCmds[command]
	if isJNLCmd {
		cmd.Execute()
		return
	}

	ex := jnl.NewExecution(context.TODO(), flag.Args())

	err := ex.Start()
	if err != nil {
		jnl.ErrorLogger.Fatalln(err.Error())
	}

	err = ex.Wait()
	if err != nil {
		jnl.ErrorLogger.Println(err.Error())
	}

}
