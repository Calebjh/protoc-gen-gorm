package plugin

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	gorm "github.com/infobloxopen/protoc-gen-gorm/options"
)

// retrieves the GormMessageOptions from a message
func getMessageOptions(message *generator.Descriptor) *gorm.GormMessageOptions {
	if message.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(message.Options, gorm.E_Opts)
	if err != nil {
		return nil
	}
	opts, ok := v.(*gorm.GormMessageOptions)
	if !ok {
		return nil
	}
	return opts
}

func getFieldOptions(field *descriptor.FieldDescriptorProto) *gorm.GormFieldOptions {
	if field.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(field.Options, gorm.E_Field)
	if err != nil {
		return nil
	}
	opts, ok := v.(*gorm.GormFieldOptions)
	if !ok {
		return nil
	}
	return opts
}

func addGORMTag(tagString, key, value string) string {
	if tagString == "" {
		if value == "" {
			return fmt.Sprintf("`gorm:\"%s\"`", key)
		}
		return fmt.Sprintf("`gorm:\"%s:%s\"`", key, value)
	}
	lcTagString := strings.ToLower(tagString)
	if strings.Contains(lcTagString, "gorm:") { // gorm tag already there
		index := strings.Index(tagString, "gorm:")
		if value == "" {
			return fmt.Sprintf("%s%s;%s", tagString[:index+6],
				key, tagString[index+6:])
		}
		return fmt.Sprintf("%s%s:%s;%s", tagString[:index+6],
			key, value, tagString[index+6:])
	}
	// no gorm tag yet
	if value == "" {
		return fmt.Sprintf("`gorm:\"%s\" %s", key, tagString[1:])
	}
	return fmt.Sprintf("`gorm:\"%s:%s\" %s", key, value, tagString[1:])

}
func getGORMTag(tagString, key string) (string, bool) {
	if !strings.Contains(tagString, "gorm:") {
		return "", false
	}
	start := strings.Index(tagString, "gorm:\"")
	end := strings.Index(tagString[start:], `"`)
	for _, tag := range strings.Split(tagString[start:end], `,`) {
		parts := strings.Split(tag, `:`)
		if parts[0] == key {
			if len(parts) == 2 {
				return parts[1], true
			}
			return "", true
		}
	}
	return "", false
}

// If key already exists in tag return tagString unchanged and true,
// otherwise inserts tag:key pair and returns false
func safeAddGORMTag(tagString, key, value string) (string, bool) {
	if _, ok := getGORMTag(tagString, key); ok {
		return tagString, true
	}
	return addGORMTag(tagString, key, value), false
}
