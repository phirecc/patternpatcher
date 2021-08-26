# Patternpatcher
Patches binaries using IDA-Style signatures and rules configured in JSON

Supports:
- Regular bytepatch
- Dereferences
	- Absolute (e.g. r/m32 call)
	- Relative (e.g. rel32 call)


See `rules.json` for an example