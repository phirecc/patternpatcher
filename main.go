package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

type Rule struct {
	Pattern string
	Replacement string
	Dereference bool
	Offset int
}

func to_int(ch byte) byte {
	if ch >= 'A' && ch <= 'F' {
		return byte(ch - 'A' + 10)
	} else {
		return byte(ch - '0')
	}
}

func patchBuffer(buffer []byte, rules []Rule) error {
	for i := 0; i < len(buffer); i++ {
		for _, rule := range rules {
			var b int
			var k int
			// for k, c := range rule.Pattern {
			for ; k < len(rule.Pattern); k++ {
				if rule.Pattern[k] == '?' {
					continue
				} else if rule.Pattern[k] == ' ' {
					b++
				} else {
					c := to_int(rule.Pattern[k])*16 + to_int(rule.Pattern[k+1])
					k++
					if c != buffer[i+b] {
						break;
					}
				}
			}
			if k == len(rule.Pattern) {
				log.Println("FOUND PATTERN", rule.Pattern, "AT", i)
			}
		}
	}
	return nil
}

func patchFile(filename string, rules []Rule) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	patchBuffer(b, rules)
	return nil
}

func main() {
	rulesFile := flag.String("rules", "rules.json", "The file containing patching rules")
	targetFile := flag.String("target", "", "The file to patch")
	flag.Parse()
	var rules []Rule
	b, err := ioutil.ReadFile(*rulesFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(b, &rules)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(rules[0].Pattern)

	patchFile(*targetFile, rules)
}