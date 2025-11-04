/*
Package goxjt provides tools to convert an XML document into a JSON object
based on a user-defined schema.
*/
package goxjt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/antchfx/xmlquery"
)

// SchemaProperty defines the structure of a single property in the user-supplied schema.
// It is recursive to allow for nested objects and arrays.
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

// MapXMLToJSON is the main entry point. It parses an XML document and a schema,
// then builds and returns a JSON byte slice based on the schema's mapping.
func MapXMLToJSON(xmlData, schemaData []byte) ([]byte, error) {
	// 1. Parse the user-defined schema
	var schema SchemaProperty
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema JSON: %w", err)
	}

	// The root of the schema must be an object
	if schema.Type != "object" {
		return nil, fmt.Errorf("root of schema must be of type 'object', got %q", schema.Type)
	}

	// 2. Parse the source XML document
	doc, err := xmlquery.Parse(bytes.NewReader(xmlData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse source XML: %w", err)
	}

	// 3. Start the recursive mapping process
	// The root context is the document itself.
	result, err := processProperty(schema, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to map XML: %w", err)
	}

	// 4. Serialize the final result to JSON
	return json.MarshalIndent(result, "", "  ")
}

// processProperty is the core recursive function that walks the schema.
// It processes a single SchemaProperty against a given XML context node.
func processProperty(prop SchemaProperty, parentContext *xmlquery.Node) (interface{}, error) {
	switch prop.Type {
	case "object":
		// This object may define a new context node with its own XPath,
		// or it may be a container for properties from the parent context.
		currentContext := parentContext
		if prop.XPath != "" {
			node, err := xmlquery.Query(parentContext, prop.XPath)
			if err != nil {
				return nil, fmt.Errorf("xpath error for object %q: %w", prop.XPath, err)
			}
			if node == nil {
				return nil, fmt.Errorf("xpath for object %q returned no node", prop.XPath)
			}
			currentContext = node
		}

		if prop.Properties == nil {
			return nil, fmt.Errorf("schema 'object' type (xpath: %q) has no 'properties' field", prop.XPath)
		}

		// Build the resulting map
		resultMap := make(map[string]interface{})
		for key, subProp := range prop.Properties {
			val, err := processProperty(subProp, currentContext)
			if err != nil {
				// Add context to the error
				return nil, fmt.Errorf("error processing property %q: %w", key, err)
			}
			resultMap[key] = val
		}
		return resultMap, nil

	case "array":
		if prop.XPath == "" {
			return nil, fmt.Errorf("schema 'array' type missing 'xpath'")
		}
		if prop.Items == nil {
			return nil, fmt.Errorf("schema 'array' type missing 'items' definition")
		}

		// "array" type *must* use QueryAll to get a list of nodes.
		nodes, err := xmlquery.QueryAll(parentContext, prop.XPath)
		if err != nil {
			return nil, fmt.Errorf("xpath error for array %q: %w", prop.XPath, err)
		}

		resultSlice := make([]interface{}, 0, len(nodes))
		for _, node := range nodes {
			// Process the 'items' schema for each node found by the array XPath.
			// The 'node' itself becomes the new parent context for this item.
			itemVal, err := processProperty(*prop.Items, node)
			if err != nil {
				return nil, fmt.Errorf("error processing array item (%s): %w", prop.XPath, err)
			}
			resultSlice = append(resultSlice, itemVal)
		}
		return resultSlice, nil

	case "string", "int", "float", "bool":
		// This is a primitive type.
		var xmlVal string

		if prop.XPath == "" {
			// This is valid *only* if we are inside an "array" of primitives.
			// In this case, 'parentContext' is the node itself.
			xmlVal = parentContext.InnerText()
		} else {
			// This is a standard property. Find the node.
			node, err := xmlquery.Query(parentContext, prop.XPath)
			if err != nil {
				return nil, fmt.Errorf("xpath error for primitive %q: %w", prop.XPath, err)
			}
			if node == nil {
				return nil, fmt.Errorf("xpath for primitive %q returned no node", prop.XPath)
			}
			xmlVal = node.InnerText()
		}

		// Cast the string value to the target type
		return castValue(xmlVal, prop.Type)

	default:
		return nil, fmt.Errorf("unsupported schema type: %q", prop.Type)
	}
}

// castValue converts an XML string value to the specified Go type.
func castValue(xmlVal string, targetType string) (interface{}, error) {
	switch targetType {
	case "string":
		return xmlVal, nil
	case "int":
		i, err := strconv.ParseInt(xmlVal, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to cast %q to int: %w", xmlVal, err)
		}
		return i, nil
	case "float":
		f, err := strconv.ParseFloat(xmlVal, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to cast %q to float64: %w", xmlVal, err)
		}
		return f, nil
	case "bool":
		b, err := strconv.ParseBool(xmlVal)
		if err != nil {
			return nil, fmt.Errorf("failed to cast %q to bool: %w", xmlVal, err)
		}
		return b, nil
	default:
		// This should not be reachable if processProperty is correct.
		return nil, fmt.Errorf("unsupported cast type: %q", targetType)
	}
}
