package main

import (
	"cli"
	"fmt"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "greet"
	app.Usage = "Chào hỏi người dùng"
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "formal",
			Usage: "Sử dụng cách chào trang trọng",
		},
		cli.StringFlag{
			Name:  "name",
			Value: "World",
			Usage: "Tên người cần chào",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:      "hello",
			ShortName: "h",
			Usage:     "Chào hỏi thân thiện",
			Action: func(c *cli.Context) {
				name := c.GlobalString("name")
				if c.GlobalBool("formal") {
					fmt.Printf("Xin chào, %s!\n", name)
				} else {
					fmt.Printf("Chào %s!\n", name)
				}
			},
		},
		{
			Name:      "goodbye",
			ShortName: "g",
			Usage:     "Chào tạm biệt",
			Action: func(c *cli.Context) {
				name := c.GlobalString("name")
				fmt.Printf("Tạm biệt %s!\n", name)
			},
		},
	}

	app.Action = func(c *cli.Context) {
		fmt.Println("Sử dụng 'greet help' để xem các lệnh có sẵn")
	}

	app.Run(os.Args)
}
