package adapt

import (
	"reflect"
)

// parseStructTags parse structural tag, forming a map
// consisting of sets of tag name and its value.
func parseStructTag(tag reflect.StructTag) tagsList {
	tagsList := make(tagsList)

	// Robustly read only supported tags via tag.Get
	if v := tag.Get(RST_MIN); v != "" {
		tagsList[tagName(RST_MIN)] = tagValue(v)
	}
	if v := tag.Get(RST_MAX); v != "" {
		tagsList[tagName(RST_MAX)] = tagValue(v)
	}
	if v := tag.Get(RST_REGEX); v != "" {
		tagsList[tagName(RST_REGEX)] = tagValue(v)
	}
	if v := tag.Get(RST_DEFAULT); v != "" {
		tagsList[tagName(RST_DEFAULT)] = tagValue(v)
	}
	if v := tag.Get(RST_CHOICE); v != "" {
		tagsList[tagName(RST_CHOICE)] = tagValue(v)
	}
	if v := tag.Get(RST_FORBIDDEN); v != "" {
		tagsList[tagName(RST_FORBIDDEN)] = tagValue(v)
	}

	return tagsList
}
