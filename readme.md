# Patternpatcher
Tool that patches binaries using a set of rules with IDA-Style signatures

```
$ patternpatcher
Usage of patternpatcher:
  -out string
    	Output file (default "patched.out")
  -rules string
    	The file containing patching rules (default "rules.json")
  -target string
    	The file to patch
```

Supports:
- Regular bytepatch
- Dereferences
	- Absolute (e.g. r/m32 call)
	- Relative (e.g. rel32 call)


See `rules.json` for an example config
