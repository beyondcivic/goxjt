# goxjt

[![Version](https://img.shields.io/badge/version-v0.1.0-blue)](https://github.com/beyondcivic/goxjt/releases/tag/v0.1.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/beyondcivic/goxjt.svg)](https://pkg.go.dev/github.com/beyondcivic/goxjt)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go implementation for converting XML documents to JSON using user-defined schemas with XPath expressions. Provides precise control over XML-to-JSON transformation through declarative mapping configurations.

## Overview

goxjt (Go XML JSON Transform) simplifies the conversion of complex XML structures into clean, well-structured JSON objects. Instead of direct XML-to-JSON conversion that often results in verbose or awkward JSON structures, goxjt uses schema-based mapping to produce the exact JSON format you need.

This project provides both a command-line interface and a Go library for XML-to-JSON transformation with:

- **Schema-driven mapping** with precise control over output structure
- **XPath-based field selection** for flexible data extraction
- **Type conversion support** (string, int, float, bool)
- **Nested object and array handling** with recursive processing
- **Comprehensive error reporting** with detailed validation messages

## Key Features

- ✅ **Schema-Based Mapping**: Define exact JSON structure using declarative schemas
- ✅ **XPath Integration**: Use powerful XPath expressions for precise data selection
- ✅ **Type Safety**: Automatic type conversion (string, int, float, bool)
- ✅ **Nested Structures**: Support for complex nested objects and arrays
- ✅ **Context Management**: Efficient XML context handling for performance
- ✅ **CLI & Library**: Both command-line tool and Go library interfaces
- ✅ **Error Handling**: Comprehensive validation and detailed error reporting
- ✅ **Cross-platform**: Works on Linux, macOS, and Windows

## Getting Started

### Prerequisites

- Go 1.24 or later
- Nix 2.25.4 or later (optional but recommended)
- PowerShell v7.5.1 or later (for building)

### Installation

#### Option 1: Install from Source

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/goxjt.git
cd goxjt
```

2. Build the application:

```bash
go build -o goxjt .
```

#### Option 2: Using Nix (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/goxjt.git
cd goxjt
```

2. Prepare the environment using Nix flakes:

```bash
nix develop
```

3. Build the application:

```bash
./build.ps1
```

#### Option 3: Go Install

```bash
go install github.com/beyondcivic/goxjt@latest
```

## Quick Start

### Command Line Interface

The `goxjt` tool provides commands for XML-to-JSON transformation:

```bash
# Convert XML to JSON using a schema
goxjt map input.xml schema.json

# Save output to a file
goxjt map input.xml schema.json -o output.json

# Show version information
goxjt version
```

### Go Library Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/beyondcivic/goxjt/pkg/goxjt"
)

func main() {
	// Example XML document
	xmlData := []byte(`
	<books>
		<book id="1">
			<title>Go Programming</title>
			<author>John Doe</author>
			<price>29.99</price>
			<available>true</available>
		</book>
		<book id="2">
			<title>XML Processing</title>
			<author>Jane Smith</author>
			<price>34.99</price>
			<available>false</available>
		</book>
	</books>`)

	// Schema defining the desired JSON structure
	schemaData := []byte(`{
		"type": "object",
		"properties": {
			"books": {
				"type": "array",
				"xpath": "//book",
				"items": {
					"type": "object",
					"properties": {
						"id": {
							"type": "string",
							"xpath": "@id"
						},
						"title": {
							"type": "string",
							"xpath": "title"
						},
						"author": {
							"type": "string",
							"xpath": "author"
						},
						"price": {
							"type": "float",
							"xpath": "price"
						},
						"available": {
							"type": "bool",
							"xpath": "available"
						}
					}
				}
			}
		}
	}`)

	// Convert XML to JSON
	jsonResult, err := goxjt.MapXMLToJSON(xmlData, schemaData)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Result: %s\n", jsonResult)
}
```

## Detailed Command Reference

### `map` - Convert XML to JSON

Transform an XML document to JSON using a user-defined schema with XPath expressions.

```bash
goxjt map [XML_FILE] [SCHEMA_FILE] [OPTIONS]
```

**Arguments:**

- `XML_FILE`: Path to the source XML file
- `SCHEMA_FILE`: Path to the JSON schema file that defines the mapping

**Options:**

- `-o, --output`: Output file path (if not specified, prints to stdout)

**Examples:**

```bash
# Basic conversion with stdout output
goxjt map data.xml mapping.json

# Save to specific output file
goxjt map data.xml mapping.json -o result.json

# Process complex nested structures
goxjt map inventory.xml complex-schema.json -o structured-data.json
```

### `version` - Show Version Information

Display version, build information, and system details.

```bash
goxjt version
```

**Example output:**

```
goxjt version dev
  Built with gc on 2024-11-04T10:30:00Z
  Git ref: abc123def456
  Go version go1.24.9, GOOS linux, GOARCH amd64
```

## Schema Definition Guide

### Schema Structure

The schema is a JSON document that defines how XML elements map to JSON properties. It supports nested objects, arrays, and primitive types.

#### Basic Schema Format

```json
{
  "type": "object",
  "properties": {
    "propertyName": {
      "type": "string|int|float|bool|object|array",
      "xpath": "XPath expression"
    }
  }
}
```

### Data Types

| Type     | Description                     | XPath Result          | JSON Output       |
| -------- | ------------------------------- | --------------------- | ----------------- |
| `object` | Container for nested properties | Context node          | `{}`              |
| `array`  | Collection of items             | Multiple nodes        | `[]`              |
| `string` | Text data                       | `"Hello World"`       | `"Hello World"`   |
| `int`    | Integer numbers                 | `"42"`                | `42`              |
| `float`  | Floating-point numbers          | `"3.14"`              | `3.14`            |
| `bool`   | Boolean values                  | `"true"` or `"false"` | `true` or `false` |

### XPath Usage Patterns

#### For Objects

Objects can define a new context for their properties:

```json
{
  "type": "object",
  "xpath": "//person",
  "properties": {
    "name": {
      "type": "string",
      "xpath": "name"
    }
  }
}
```

#### For Arrays

Arrays must specify an XPath that returns multiple nodes:

```json
{
  "type": "array",
  "xpath": "//book",
  "items": {
    "type": "object",
    "properties": {
      "title": {
        "type": "string",
        "xpath": "title"
      }
    }
  }
}
```

#### For Primitives

Primitive types extract text content or attribute values:

```json
{
  "title": {
    "type": "string",
    "xpath": "title"
  },
  "id": {
    "type": "string",
    "xpath": "@id"
  }
}
```

### Common XPath Patterns

| Pattern          | Description                     | Example         |
| ---------------- | ------------------------------- | --------------- |
| `//element`      | Find all elements anywhere      | `//book`        |
| `element`        | Direct child element            | `title`         |
| `@attribute`     | Attribute value                 | `@id`           |
| `/root/element`  | Absolute path from root         | `/books/book`   |
| `element[1]`     | First occurrence                | `book[1]`       |
| `element[@attr]` | Element with specific attribute | `book[@id='1']` |
| `text()`         | Text content explicitly         | `title/text()`  |

## Examples

### Example 1: Simple Book Catalog

**books.xml:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<library>
    <book id="1" genre="programming">
        <title>Go in Action</title>
        <author>William Kennedy</author>
        <price currency="USD">39.99</price>
        <available>true</available>
    </book>
    <book id="2" genre="web">
        <title>Web Development with Go</title>
        <author>Jon Calhoun</author>
        <price currency="USD">29.99</price>
        <available>false</available>
    </book>
</library>
```

**schema.json:**

```json
{
  "type": "object",
  "properties": {
    "library": {
      "type": "object",
      "xpath": "/library",
      "properties": {
        "books": {
          "type": "array",
          "xpath": "book",
          "items": {
            "type": "object",
            "properties": {
              "id": { "type": "string", "xpath": "@id" },
              "genre": { "type": "string", "xpath": "@genre" },
              "title": { "type": "string", "xpath": "title" },
              "author": { "type": "string", "xpath": "author" },
              "price": { "type": "float", "xpath": "price" },
              "currency": { "type": "string", "xpath": "price/@currency" },
              "available": { "type": "bool", "xpath": "available" }
            }
          }
        }
      }
    }
  }
}
```

**Command:**

```bash
goxjt map books.xml schema.json -o output.json
```

**Output (output.json):**

```json
{
  "library": {
    "books": [
      {
        "id": "1",
        "genre": "programming",
        "title": "Go in Action",
        "author": "William Kennedy",
        "price": 39.99,
        "currency": "USD",
        "available": true
      },
      {
        "id": "2",
        "genre": "web",
        "title": "Web Development with Go",
        "author": "Jon Calhoun",
        "price": 29.99,
        "currency": "USD",
        "available": false
      }
    ]
  }
}
```

### Example 2: Complex Nested Structure

**inventory.xml:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<inventory>
    <categories>
        <category name="Electronics">
            <product id="E001">
                <name>Laptop</name>
                <specs>
                    <cpu>Intel i7</cpu>
                    <ram>16GB</ram>
                    <storage>512GB SSD</storage>
                </specs>
                <price>999.99</price>
                <stock>5</stock>
            </product>
            <product id="E002">
                <name>Smartphone</name>
                <specs>
                    <cpu>A15 Bionic</cpu>
                    <ram>6GB</ram>
                    <storage>256GB</storage>
                </specs>
                <price>799.99</price>
                <stock>12</stock>
            </product>
        </category>
    </categories>
</inventory>
```

**complex-schema.json:**

```json
{
  "type": "object",
  "properties": {
    "inventory": {
      "type": "object",
      "xpath": "/inventory",
      "properties": {
        "categories": {
          "type": "array",
          "xpath": "categories/category",
          "items": {
            "type": "object",
            "properties": {
              "name": { "type": "string", "xpath": "@name" },
              "products": {
                "type": "array",
                "xpath": "product",
                "items": {
                  "type": "object",
                  "properties": {
                    "id": { "type": "string", "xpath": "@id" },
                    "name": { "type": "string", "xpath": "name" },
                    "specifications": {
                      "type": "object",
                      "xpath": "specs",
                      "properties": {
                        "processor": { "type": "string", "xpath": "cpu" },
                        "memory": { "type": "string", "xpath": "ram" },
                        "storage": { "type": "string", "xpath": "storage" }
                      }
                    },
                    "price": { "type": "float", "xpath": "price" },
                    "inStock": { "type": "int", "xpath": "stock" }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}
```

### Example 3: Array of Primitives

**tags.xml:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<article>
    <title>Getting Started with Go</title>
    <tags>
        <tag>programming</tag>
        <tag>golang</tag>
        <tag>tutorial</tag>
        <tag>beginner</tag>
    </tags>
</article>
```

**tags-schema.json:**

```json
{
  "type": "object",
  "properties": {
    "article": {
      "type": "object",
      "xpath": "/article",
      "properties": {
        "title": { "type": "string", "xpath": "title" },
        "tags": {
          "type": "array",
          "xpath": "tags/tag",
          "items": {
            "type": "string"
          }
        }
      }
    }
  }
}
```

## API Reference

### Core Functions

#### `MapXMLToJSON(xmlData, schemaData []byte) ([]byte, error)`

Converts an XML document to JSON using a user-defined schema.

**Parameters:**

- `xmlData []byte`: The source XML document
- `schemaData []byte`: The mapping schema as JSON

**Returns:**

- `[]byte`: The resulting JSON document with proper indentation
- `error`: Any error that occurred during processing

**Example:**

```go
jsonResult, err := goxjt.MapXMLToJSON(xmlData, schemaData)
if err != nil {
    log.Fatalf("Conversion failed: %v", err)
}
fmt.Printf("Result: %s\n", jsonResult)
```

### Data Structures

#### `SchemaProperty`

Defines the structure of a schema property with support for nested objects and arrays.

**Fields:**

- `Type string`: Target data type ("object", "array", "string", "int", "float", "bool")
- `XPath string`: XPath expression for data selection
- `Properties map[string]SchemaProperty`: Sub-properties for object types
- `Items *SchemaProperty`: Schema definition for array elements

### Internal Functions

#### `processProperty(prop SchemaProperty, parentContext *xmlquery.Node) (interface{}, error)`

Recursively processes schema properties against XML context nodes.

#### `castValue(xmlVal string, targetType string) (interface{}, error)`

Converts XML string values to specified Go types with validation.

## Error Handling

The library provides comprehensive error reporting for various failure scenarios:

### Schema Validation Errors

- **Invalid JSON syntax**: Malformed schema JSON
- **Invalid root type**: Root schema must be "object"
- **Missing required fields**: Objects must have properties, arrays must have items
- **Missing XPath**: Arrays require XPath expressions

### XML Processing Errors

- **XML parse errors**: Invalid XML syntax or structure
- **XPath evaluation errors**: Invalid XPath expressions
- **Node selection errors**: XPath returns no matching nodes
- **Context navigation errors**: Issues traversing XML structure

### Type Conversion Errors

- **Invalid numeric conversion**: String cannot be parsed as int/float
- **Invalid boolean conversion**: String cannot be parsed as boolean
- **Type mismatch**: Unexpected data format for target type

**Example Error Handling:**

```go
jsonResult, err := goxjt.MapXMLToJSON(xmlData, schemaData)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "failed to parse schema JSON"):
        log.Printf("Schema JSON is malformed: %v", err)
    case strings.Contains(err.Error(), "failed to parse source XML"):
        log.Printf("XML document is invalid: %v", err)
    case strings.Contains(err.Error(), "xpath error"):
        log.Printf("XPath expression error: %v", err)
    default:
        log.Printf("Processing error: %v", err)
    }
    return
}
```

## Best Practices

### Schema Design

1. **Start with simple structures**: Begin with basic objects before adding complexity
2. **Test XPath expressions**: Validate XPath patterns against your XML structure
3. **Use descriptive property names**: Choose meaningful names for JSON output
4. **Leverage context efficiently**: Use object XPath to set context for multiple properties
5. **Handle optional data**: Consider how to handle missing XML elements

### Performance Optimization

1. **Efficient XPath expressions**: Use specific paths rather than broad searches
2. **Minimize deep nesting**: Complex nested structures impact performance
3. **Context reuse**: Set context at appropriate levels to avoid redundant queries
4. **Batch processing**: Process multiple similar documents with the same schema

### Error Prevention

1. **Validate inputs**: Check XML and schema validity before processing
2. **Handle edge cases**: Test with empty elements, missing attributes, and malformed data
3. **Use type-safe schemas**: Ensure XPath results match expected data types
4. **Test thoroughly**: Validate transformations with representative data sets

## Development

### Project Structure

```
goxjt/
├── cmd/goxjt/          # Command-line interface
├── pkg/goxjt/          # Core library
├── pkg/version/        # Version information
├── docs/               # Documentation
├── bin/                # Built binaries
├── main.go             # Application entry point
├── go.mod              # Go module definition
├── build.ps1           # Build script
└── flake.nix           # Nix development environment
```

### Adding New Features

To extend goxjt with new functionality:

1. **Add new data types**: Extend the `castValue` function for new type conversions
2. **Enhance XPath support**: Add specialized XPath handling in `processProperty`
3. **Improve error handling**: Add specific error types for better debugging
4. **Add validation**: Implement schema validation before processing

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and add tests
4. Ensure all tests pass: `go test ./...`
5. Update documentation as needed
6. Commit your changes: `git commit -am 'Add new feature'`
7. Push to the branch: `git push origin feature/new-feature`
8. Submit a pull request

### Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Build Environment

### Using Nix (Recommended)

Set up the development environment with Nix:

```bash
nix develop
```

This provides all necessary dependencies and tools for development.

### Manual Build

Build the application manually:

```bash
go build -o goxjt .
```

For optimized builds, use the provided build script:

```bash
./build.ps1
```

## Configuration

The application supports configuration through environment variables:

- `GOXJT_*`: Environment variables with GOXJT prefix are automatically loaded

## Use Cases

### Data Integration

- **API Response Transformation**: Convert XML API responses to JSON for modern applications
- **Legacy System Integration**: Transform XML data from legacy systems to JSON for new services
- **Data Pipeline Processing**: Convert XML data sources in ETL pipelines

### Document Processing

- **Configuration File Conversion**: Transform XML configuration to JSON format
- **Content Management**: Convert XML content to JSON for web applications
- **Data Migration**: Migrate XML databases to JSON-based systems

### Enterprise Integration

- **Message Format Translation**: Convert XML messages to JSON in enterprise service buses
- **Protocol Adaptation**: Bridge XML-based and JSON-based microservices
- **Data Format Standardization**: Standardize diverse XML formats to consistent JSON schemas

## Related Projects

- [xmlquery](https://github.com/antchfx/xmlquery) - XPath query library for Go
- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go
- [jq](https://stedolan.github.io/jq/) - JSON processor (complementary tool)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
