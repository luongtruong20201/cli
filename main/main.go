package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IntSlice []int

func (i *IntSlice) String() string {
	return fmt.Sprint([]int(*i))
}

func (i *IntSlice) Set(value string) error {
	vals := strings.Split(value, ",")
	for _, v := range vals {
		val, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return err
		}
		*i = append(*i, int(val))
	}
	return nil
}

func main() {
	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)

	host := serverFlags.String("host", "localhost", "Server host")
	port := serverFlags.Int("port", 8080, "Server port")

	serverFlags.Parse(os.Args[1:])

	fmt.Println("Host:", *host)
	fmt.Println("Port:", *port)
}
