# 📜 xsd-codegen

**xsd-codegen** is a CLI tool and Go library for generating Go code and dummy XML output from XSD (XML Schema Definition) files. It parses complex schemas including support for `<import>` and `<include>` directives, producing idiomatic Go types and mockable XML data for testing and prototyping.

---

## 🚀 Features

- ✅ Parses XSD files and builds internal schema models
- 🔗 Resolves `<include>` and `<import>` directives recursively
- 🧬 Generates Go structs from complex XSD types
- 📝 Outputs XML with dummy values, respecting constraints like `enumeration`, `pattern`, `minInclusive`, `maxLength`, etc.
- 🧪 Includes integration and unit tests for all modules
- 💡 Designed with modular architecture (library + CLI)

---

## 🛠️ Installation

```bash
git clone https://github.com/yourname/xsd-codegen.git
cd xsd-codegen
go build ./cmd/xsd-codegen
```

## 📦 Usage
```bash
./xsd-codegen path/to/schema.xsd
```
Parses the given .xsd

Resolves all nested imports/includes

Generates Go structs and sample XML output

Generated outputs can be customized or piped into files depending on your CLI extension.