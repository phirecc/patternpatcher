package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

type Rule struct {
	Desc        string
	Pattern     string
	Replacement string
	Dereference *Dereference
	Offset      int
}

type Dereference struct {
	NBytes      int `json:"nbytes"`
	OffsetAfter int `json:"offset_after"`
	Type        string
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
				t := i
				if rule.Dereference != nil {
					var x int
					for u := 0; u < rule.Dereference.NBytes; u++ {
						x += int(buffer[t+u+rule.Offset]) << (8 * u)
					}
					if rule.Dereference.Type == "rel" {
						if x&(1<<(rule.Dereference.NBytes*8-1)) != 0 {
							x ^= (1 << (rule.Dereference.NBytes * 8)) - 1
							x *= -1
							x -= 1
						}
						t += x + rule.Dereference.OffsetAfter
					} else if rule.Dereference.Type == "abs" {
						t = x + rule.Dereference.OffsetAfter
					} else {
						return fmt.Errorf("Unknown dereference type: %s", rule.Dereference.Type)
					}
				}
				fmt.Printf("Patching rule %d: \"%s\" (%s) at 0x%x\n", n, rule.Desc, rule.Pattern, t)
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
	return b, patchBuffer(b, rules)
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
	if err = ioutil.WriteFile(*outFile, buf, 0644); err != nil {
		log.Fatalln(err)
	}
}
