# 📜 xsd-codegen

**xsd-codegen** is a CLI tool and Go library for generating dummy XML output from XSD (XML Schema Definition) files. It parses complex schemas including support for `<import>` and `<include>` directives, producing idiomatic Go types and mockable XML data for testing and prototyping.

---

## 🚀 Features

- 🔗 Resolves `<include>` and `<import>` directives recursively
- 📝 Outputs XML with dummy values, respecting constraints like `enumeration`, `pattern`, `minInclusive`, `maxLength`, etc.
- 🧪 Includes integration and unit tests for all modules
- 💡 Designed with modular architecture (library + CLI)

---

## 🛠️ Installation

```bash
git clone https://github.com/patrick-ivann/xsd-codegen.git
cd xsd-codegen
go build ./cmd/xsd-codegen
```

## 📦 Usage
```bash
./xsd-codegen -xsd complete.xsd -out example.xml
```
Parses the given .xsd

Resolves all nested imports/includes

Generates sample XML output

Generated outputs can be customized or piped into files depending on your CLI extension.