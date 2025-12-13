package main

import (
	"fmt"
	"log"

	"github.com/luxingwen/sgin/pkg/app"
	"github.com/luxingwen/sgin/pkg/config"
)

// MyExtra is the demo config structure used in this example
type MyExtra struct {
	FeatureFlag bool `mapstructure:"feature_flag"`
	Limit       int  `mapstructure:"limit"`
}

func main() {
	// Initialize config from example config file
	config.InitConfigWithFile("examples/withextra/config.yaml")

	// ---方式 A：手动解码，再注入 App---
	var manual MyExtra
	if err := config.UnmarshalKey("my_extra2", &manual); err != nil {
		log.Printf("manual UnmarshalKey failed: %v\n", err)
	} else {
		// inject decoded struct into App
		a := app.NewAppWithOptions(app.WithExtra("manual_extra", "my_extra2", &manual))
		if v, ok := app.GetExtraAs[*MyExtra](a, "manual_extra"); ok {
			fmt.Printf("manual injected extra: %+v\n", *v)
		} else {
			fmt.Println("manual_extra not found or wrong type")
		}
	}

	// ---方式 B：直接把一个空指针传给 WithExtra，让它内部去解码---
	a2 := app.NewAppWithOptions(app.WithExtra("auto_extra", "my_extra2", &MyExtra{}))
	if v2, ok := app.GetExtraAs[*MyExtra](a2, "auto_extra"); ok {
		fmt.Printf("auto injected extra: %+v\n", *v2)
	} else {
		fmt.Println("auto_extra not found or wrong type")
	}

	fmt.Println("examples/withextra done")
}
