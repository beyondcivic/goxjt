# goxjt

```go
import "github.com/beyondcivic/goxjt/pkg/goxjt"
```

Package goxjt provides tools to convert an XML document into a JSON object based on a user\-defined schema.

## Index

- [func MapXMLToJSON\(xmlData, schemaData \[\]byte\) \(\[\]byte, error\)](<#MapXMLToJSON>)
- [type SchemaProperty](<#SchemaProperty>)


<a name="MapXMLToJSON"></a>
## func [MapXMLToJSON](<https://github.com:beyondcivic/goxjt/blob/main/pkg/goxjt/goxjt.go#L41>)

```go
func MapXMLToJSON(xmlData, schemaData []byte) ([]byte, error)
```

MapXMLToJSON is the main entry point. It parses an XML document and a schema, then builds and returns a JSON byte slice based on the schema's mapping.

<a name="SchemaProperty"></a>
## type [SchemaProperty](<https://github.com:beyondcivic/goxjt/blob/main/pkg/goxjt/goxjt.go#L18-L37>)

SchemaProperty defines the structure of a single property in the user\-supplied schema. It is recursive to allow for nested objects and arrays.

```go
type SchemaProperty struct {
    // Type specifies the target data type in the output JSON.
    // Supported: "object", "array", "string", "int", "float", "bool".
    Type string `json:"type"`

    // XPath is the expression used to find the corresponding data in the XML.
    // - For "object" (with XPath), it sets the context node for sub-properties.
    // - For "array", it must return a set of nodes to iterate over.
    // - For primitives ("string", "int", etc.), it must return a single node.
    // - For primitives inside an "array", XPath should be empty (or ".")
    //   as the "array" XPath already selected the node.
    XPath string `json:"xpath,omitempty"`

    // Properties defines the sub-properties for a "object" type.
    // The map key is the desired property name in the output JSON.
    Properties map[string]SchemaProperty `json:"properties,omitempty"`

    // Items defines the schema for each element in an "array" type.
    Items *SchemaProperty `json:"items,omitempty"`
}
```