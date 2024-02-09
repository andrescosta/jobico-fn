package main

import (
	"fmt"

	"github.com/andrescosta/goico/pkg/yamlutil"
	"github.com/andrescosta/jobico/internal/test"
)

func main() {
	p := test.NewTestPackage(test.SchemaRefIds{SchameRef: "sch1", SchemaRefOk: "sch1_ok", SchemaRefError: "sch1_error"}, "run1")
	s, err := yamlutil.Marshal(&p)
	if err != nil {
		print(err)
	}
	fmt.Println(*s)
}
