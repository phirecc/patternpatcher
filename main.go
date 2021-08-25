package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

type Rule struct {
	Pattern                string
	Replacement            string
	Dereference            bool
	DereferenceNBytes      int    `json:"dereference_nbytes"`
	DereferenceOffsetAfter int    `json:"dereference_offset_after"`
	DereferenceType        string `json:"dereference_type"`
	Offset                 int
}

func to_int(ch byte) byte {
	if ch >= 'A' && ch <= 'F' {
		return byte(ch - 'A' + 10)
	} else if ch >= 'a' && ch <= 'f' {
		return byte(ch - 'a' + 10)
	} else {
		return byte(ch - '0')
	}
}

func patchBuffer(buffer []byte, rules []Rule) error {
	for i := 0; i < len(buffer); i++ {
		for n, rule := range rules {
			var b int
			var k int
			for ; k < len(rule.Pattern); k++ {
				if rule.Pattern[k] == '?' {
					continue
				} else if rule.Pattern[k] == ' ' {
					b++
				} else {
					c := to_int(rule.Pattern[k])*16 + to_int(rule.Pattern[k+1])
					k++
					if c != buffer[i+b] {
						break
					}
				}
			}
			if k == len(rule.Pattern) {
				fmt.Printf("Patching rule %d (%s) at 0x%x\n", n, rule.Pattern, i)
				t := i
				if rule.Dereference {
					var x int
					for u := 0; u < rule.DereferenceNBytes; u++ {
						x += int(buffer[t+u+rule.Offset]) << (8 * u)
					}
					if rule.DereferenceType == "rel" {
						if x&(1<<31) != 0 {
							x ^= (1 << 32) - 1
							x *= -1
							x -= 1
						}
						t += x + rule.DereferenceOffsetAfter
					} else if rule.DereferenceType == "abs" {
						t = x + rule.DereferenceOffsetAfter
					}
				}
				for k = 0; k < len(rule.Replacement); k++ {
					if rule.Replacement[k] == ' ' {
						t++
					} else {
						c := to_int(rule.Replacement[k])*16 + to_int(rule.Replacement[k+1])
						k++
						buffer[t] = c
					}
				}
			}
		}
	}
	return nil
}

func patchFile(filename string, rules []Rule) ([]byte, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	patchBuffer(b, rules)
	return b, nil
}

func main() {
	rulesFile := flag.String("rules", "rules.json", "The file containing patching rules")
	targetFile := flag.String("target", "", "The file to patch")
	outFile := flag.String("out", "patched.out", "Output file")
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

	buf, err := patchFile(*targetFile, rules)
	if err != nil {
		log.Fatalln(err)
	}
	ioutil.WriteFile(*outFile, buf, 0644)
}
