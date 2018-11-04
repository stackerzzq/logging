package main

import (
	"github.com/gosets/logging"
)

func main() {
	log:=logging.New("")
	log.Debugf("hello world %s", "aaaa")
	log.Infof("hello world %s", "bbbb")
	log.Warnf("hello world %s", "cccc")
	log.Errorf("hello world %s", "dddd")
	log.Panicf("hello world %s", "eeee")
	log.Fatalf("hello world %s", "ffff")
}