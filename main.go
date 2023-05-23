package main

import (
	"github.com/kirtfieldk/astella/src/structures/conf"
)

func main() {
	var c conf.Conf
	c.GetConf()
	c.BuildApi()
}
