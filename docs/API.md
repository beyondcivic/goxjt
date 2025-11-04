# API Documentation

This document provides comprehensive API documentation for the goxjt library.

## Package Overview

The `github.com/beyondcivic/goxjt/pkg/goxjt` package provides functionality for converting XML documents into JSON objects based on user-defined schemas using XPath expressions.

## Core Functions

### XML to JSON Mapping

#### `MapXMLToJSON(xmlData, schemaData []byte) ([]byte, error)`

Maps an XML document to a JSON object based on a user-defined schema. This is the main entry point for the library.

**Parameters:**

- `xmlData`: The source XML document as bytes
- `schemaData`: The mapping schema definition as JSON bytes

**Returns:**

- `[]byte`: The resulting JSON document with proper indentation
- `error`: Any error that occurred during parsing or mapping

**Example:**

```go
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

jsonResult, err := goxjt.MapXMLToJSON(xmlData, schemaData)
if err != nil {
    log.Fatalf("Error: %v", err)
}
fmt.Printf("Result: %s\n", jsonResult)
```

### Internal Processing Functions

#### `processProperty(prop SchemaProperty, parentContext *xmlquery.Node) (interface{}, error)`

Core recursive function that processes a single schema property against an XML context node. This function is used internally by `MapXMLToJSON`.

#### `castValue(xmlVal string, targetType string) (interface{}, error)`

Converts an XML string value to the specified Go type. Supports conversion to string, int, float, and bool types.

**Parameters:**

- `xmlVal`: The string value extracted from XML
- `targetType`: Target data type ("string", "int", "float", "bool")

**Returns:**

- `interface{}`: The converted value
- `error`: Any error that occurred during conversion

## Data Structures

### `SchemaProperty`

Defines the structure of a single property in the user-supplied schema. It is recursive to allow for nested objects and arrays.

**Fields:**

- `Type string`: Specifies the target data type in the output JSON. Supported values: "object", "array", "string", "int", "float", "bool"
- `XPath string`: XPath expression used to find corresponding data in the XML document
- `Properties map[string]SchemaProperty`: Sub-properties for "object" type (map key is the desired property name in output JSON)
- `Items *SchemaProperty`: Schema definition for each element in an "array" type

**Example:**

```go
schema := goxjt.SchemaProperty{
    Type: "object",
    Properties: map[string]goxjt.SchemaProperty{
        "title": {
            Type:  "string",
            XPath: "//title",
        },
        "chapters": {
            Type:  "array",
            XPath: "//chapter",
            Items: &goxjt.SchemaProperty{
                Type: "object",
                Properties: map[string]goxjt.SchemaProperty{
                    "name": {
                        Type:  "string",
                        XPath: "@name",
                    },
                    "pages": {
                        Type:  "int",
                        XPath: "pages",
                    },
                },
            },
        },
    },
}
```

## Schema Definition Guide

### Data Types

The schema supports the following data types:

| Type     | Description                     | Example XPath Result | JSON Output |
| -------- | ------------------------------- | -------------------- | ----------- |
| `object` | Container for nested properties | N/A                  | `{}`        |
| `array`  | Collection of items             | Multiple nodes       | `[]`        |
| `string` | Text data                       | `"Hello World"`      | `"Hello"`   |
| `int`    | Integer numbers                 | `"42"`               | `42`        |
| `float`  | Floating-point numbers          | `"3.14"`             | `3.14`      |
| `bool`   | Boolean values                  | `"true"`             | `true`      |

### XPath Usage Patterns

#### For Objects

- **With XPath**: Sets the context node for sub-properties
- **Without XPath**: Uses parent context for sub-properties

#### For Arrays

- **Must have XPath**: Returns a set of nodes to iterate over
- **Uses QueryAll**: Processes each matching node with the Items schema

#### For Primitives

- **With XPath**: Finds the specific node containing the value
- **Without XPath**: Uses current context node (valid only inside arrays)

### Schema Validation Rules

1. **Root must be object**: The top-level schema type must be "object"
2. **Array requires XPath**: Arrays must specify an XPath to select nodes
3. **Array requires Items**: Arrays must define an Items schema
4. **Object requires Properties**: Objects must define Properties (except when used as array items with implicit structure)
5. **XPath expressions**: Must be valid XPath 1.0 expressions

## XPath Reference

### Common XPath Patterns

| Pattern          | Description                     | Example              |
| ---------------- | ------------------------------- | -------------------- |
| `//element`      | Find all elements anywhere      | `//book`             |
| `element`        | Direct child element            | `title`              |
| `@attribute`     | Attribute value                 | `@id`                |
| `/root/element`  | Absolute path from root         | `/books/book`        |
| `element[1]`     | First occurrence                | `book[1]`            |
| `element[@attr]` | Element with specific attribute | `book[@id='1']`      |
| `text()`         | Text content                    | `title/text()`       |
| `.`              | Current context node            | `.` (in array items) |

### XPath Context Rules

- **Object XPath**: Changes context for all sub-properties
- **Array XPath**: Each matching node becomes context for Items processing
- **Primitive XPath**: Extracts text content from the matched node
- **Empty XPath in primitives**: Valid only inside array Items, uses current array element as context

## Error Handling

The library provides detailed error messages for common issues:

### Schema Errors

- **Invalid JSON**: Malformed schema JSON syntax
- **Invalid root type**: Root schema must be type "object"
- **Missing properties**: Object types must define properties
- **Missing items**: Array types must define items schema
- **Missing XPath**: Array types must specify XPath

### XML Errors

- **Parse errors**: Invalid XML syntax or structure
- **XPath errors**: Invalid XPath expressions
- **Node not found**: XPath returns no matching nodes
- **Context errors**: Issues with XML context navigation

### Type Conversion Errors

- **Invalid int**: String cannot be parsed as integer
- **Invalid float**: String cannot be parsed as float
- **Invalid bool**: String cannot be parsed as boolean

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
    case strings.Contains(err.Error(), "failed to cast"):
        log.Printf("Type conversion error: %v", err)
    default:
        log.Printf("Mapping error: %v", err)
    }
    return
}
```

## Best Practices

### Schema Design

1. **Start simple**: Begin with basic object structures before adding arrays
2. **Test XPath expressions**: Validate XPath patterns against your XML structure
3. **Use meaningful names**: Choose descriptive property names for JSON output
4. **Consider context**: Design object hierarchy to minimize XPath complexity

### Performance Optimization

1. **Efficient XPath**: Use specific paths rather than broad searches when possible
2. **Minimize array nesting**: Deep array structures can impact performance
3. **Context reuse**: Leverage object XPath to set context for multiple sub-properties

### Error Prevention

1. **Validate inputs**: Check XML and schema validity before processing
2. **Handle missing data**: Consider optional fields and default values
3. **Test edge cases**: Verify behavior with empty elements and missing attributes

## Command Line Interface

The `goxjt` command-line tool provides access to the library functionality:

### Basic Usage

```bash
# Map XML to JSON using schema file
goxjt map input.xml schema.json

# Save output to file
goxjt map input.xml schema.json -o output.json

# Display version information
goxjt version
```

### Example Files

**books.xml:**

```xml
<?xml version="1.0" encoding="UTF-8"?>
<library>
    <book id="1" genre="programming">
        <title>Go in Action</title>
        <authors>
            <author>William Kennedy</author>
            <author>Brian Ketelsen</author>
        </authors>
        <publication year="2015"/>
        <price currency="USD">39.99</price>
        <available>true</available>
    </book>
    <book id="2" genre="web">
        <title>Web Development with Go</title>
        <authors>
            <author>Jon Calhoun</author>
        </authors>
        <publication year="2019"/>
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
              "id": {
                "type": "string",
                "xpath": "@id"
              },
              "genre": {
                "type": "string",
                "xpath": "@genre"
              },
              "title": {
                "type": "string",
                "xpath": "title"
              },
              "authors": {
                "type": "array",
                "xpath": "authors/author",
                "items": {
                  "type": "string"
                }
              },
              "publicationYear": {
                "type": "int",
                "xpath": "publication/@year"
              },
              "price": {
                "type": "float",
                "xpath": "price"
              },
              "currency": {
                "type": "string",
                "xpath": "price/@currency"
              },
              "available": {
                "type": "bool",
                "xpath": "available"
              }
            }
          }
        }
      }
    }
  }
}
```

## Version Compatibility

This library is compatible with:

- Go 1.24+
- XPath 1.0 specification
- JSON-LD 1.1 compatible output
- XML 1.0 documents

## Dependencies

- `github.com/antchfx/xmlquery`: XML parsing and XPath evaluation
- `github.com/spf13/cobra`: Command-line interface framework
- `github.com/spf13/viper`: Configuration management

## Related Documentation

- [Command Line Reference](cmd/goxjt.md)
- [Version Information](cmd/goxjt_version.md)
- [Go Documentation](../godoc/goxjt.md)
