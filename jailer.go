package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type Pair struct {
	regex  *regexp.Regexp
	groups *[]string
}

func usage() {
	fmt.Printf("Usage: jailer key[,key]=regexp [key[,key]=regexp]\n\n")
	fmt.Printf("    Example: \n")
	fmt.Printf("    :$ echo 'INFO: [core_name] webapp=/solr path=/select' | jailer 'webapp,path=webapp=([^ ]+).*?path=([^ ]+)'\n")
	fmt.Printf("    { \"webapp\": \"/solr\", \"path\": \"/select\" }\n\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}

	regexen := make([]Pair, 0)

	for _, arg := range os.Args[1:] {
		segments := strings.Split(arg, "=")
		if len(segments) < 2 {
			log.Fatalln("Argument did not have an equals sign")
		}
		groups := strings.Split(segments[0], ",")
		re, err := regexp.Compile(strings.Join(segments[1:], "="))
		if err != nil {
			fmt.Printf("%s\n\n", err)
			usage()
		}
		pair := Pair{regex: re, groups: &groups}
		regexen = append(regexen, pair)
	}

	encoder := json.NewEncoder(os.Stdout)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err == nil || err == io.EOF {
			m := make(map[string]string)
			for _, pair := range regexen {
				matches := pair.regex.FindStringSubmatch(line)
				for i := 0; i < len(matches); i++ {
					key := (*pair.groups)[i]
					if key != "_" {
						m[key] = matches[i]
					}
				}
			}
			encoder.Encode(m)
		} else {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}
	}
}
