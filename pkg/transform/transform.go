package transform

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/BlueSageSolutions/courier/pkg/util"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

const (
	StageToken   string = "$(STAGE)"
	StageNull    string = "null"
	TenantID     string = "tenant-id"
	EventType    string = "data-type"
	TopicARN     string = "topic-arn"
	SubscribeURL string = "subscribe-url"
)

const (
	DestinationTypeSQS              string = "sqs"
	DestinationTypeLog              string = "log"
	DestinationTypeNull             string = "null"
	RuleTypeContains                string = "contains"
	RuleTypeNotContained            string = "not-contained"
	RuleTypeDefined                 string = "defined"
	RuleTypeUndefined               string = "undefined"
	RuleTypeIsTrue                  string = "is-true"
	RuleTypeIsFalse                 string = "is-false"
	TransformationTypeTemplate      string = "template"
	TransformationTypeTokenize      string = "tokenize"
	TransformationTypeJq            string = "jq"
	TransformationTypeRandom        string = "random"
	TransformationTypeTimestamp     string = "timestamp"
	TransformationTypeTranslate     string = "translate"
	TransformationTypeParse         string = "parse"
	TransformationTypeConcatenate   string = "concatenate"
	TransformationTypeStaticMapping string = "static-mapping"
	NameJqQuery                     string = "jq-query"
	NameLength                      string = "length"
	NameTranslate                   string = "translate"
	NameField                       string = "token-"
	NameDelimiter                   string = "delimiter"
	NameTemplate                    string = "template"
	NameRegularExpression           string = "regular-expression"
	NameJoinSeparator               string = "join-separator"
	NameStaticValue                 string = "static-value"
	NameTranslation                 string = "translation"
	NameTokenIndices                string = "token-indices"
	TransformerName                 string = "transformer-name"
)

type Transformation struct {
	Translation *Translations
	Ephemeral   bool        `yaml:"ephemeral"`
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"`
	Parameters  []Parameter `yaml:"parameters"`
}

type Parameter struct {
	Name   string   `yaml:"name"`
	Value  string   `yaml:"value"`
	Values []string `yaml:"values"`
}

type MatchingRule struct {
	Attribute  string   `yaml:"attribute"`
	Attributes []string `yaml:"attributes"`
	RuleType   string   `yaml:"rule-type"`
	Values     []string `yaml:"values"`
	Matcher    *regexp.Regexp
}

type Transformed map[string]string

type TransformerDefaults struct {
	Transformations []Transformation `yaml:"transformations"`
}

type Transformer struct {
	Name                string               `yaml:"transformer"`
	Destinations        []Destination        `yaml:"destinations"`
	MatchingRules       []MatchingRule       `yaml:"matching-rules"`
	Transformations     []Transformation     `yaml:"transformations"`
	TransformerDefaults *TransformerDefaults `yaml:"defaults"`
}

type TransformerList struct {
	Transformers []Transformer `yaml:"transformers"`
}

type TransformerContext map[string]string

// Buckets represents buckets configuration
type Translation struct {
	Name  string            `yaml:"translation"`
	Table map[string]string `yaml:"table,omitempty"`
}

type TranslationList []Translation
type Translations map[string]Translation

type Destination struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type TokenizedField struct {
	Name  string `yaml:"name"`
	Index int    `yaml:"index"`
}

type DataPayload struct {
	Input              json.RawMessage    `json:"input"`
	Outputs            Transformed        `json:"outputs"`
	TransformerContext TransformerContext `json:"context"`
}

type DataPayloads []*DataPayload

type SafeDataPayloads struct {
	sync.RWMutex
	DataPayloads DataPayloads
}

func NewTransformerContext() TransformerContext {
	ctx := make(TransformerContext)
	return ctx
}

func (payload *DataPayload) GetNamedValue(key string) string {
	value, ok := payload.Outputs[key]
	if !ok {
		return ""
	}
	return string(value)
}

func (payload *DataPayload) GetTransformerContextValue(key string) string {
	value, ok := payload.TransformerContext[key]
	if !ok {
		return ""
	}
	return value
}

func GenerateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}

func (safePayload *SafeDataPayloads) Append(message *DataPayload) {
	safePayload.Lock()
	safePayload.DataPayloads = append(safePayload.DataPayloads, message)
	safePayload.Unlock()
}

func (payload *DataPayload) Clone() *DataPayload {
	clone := &DataPayload{Input: payload.Input, Outputs: make(Transformed), TransformerContext: make(TransformerContext)}
	for key, value := range payload.Outputs {
		clone.Outputs[key] = value
	}
	for key, value := range payload.TransformerContext {
		clone.TransformerContext[key] = value
	}
	return clone
}

func NewDataPayload(transformerCtx TransformerContext, transformed Transformed, data json.RawMessage) *DataPayload {
	return &DataPayload{Input: data, Outputs: transformed, TransformerContext: transformerCtx}
}

func UnmarshalDataPayload(bytes []byte) (*DataPayload, error) {
	var newEnvelope DataPayload
	err := json.Unmarshal(bytes, &newEnvelope)
	if err != nil {
		util.GetLogger().Error("marshal", zap.Error(err))
		return nil, err
	}
	return &newEnvelope, nil
}

func NewTokenizedField(fieldIndex, fieldName string) (TokenizedField, error) {
	var tokenizedField TokenizedField
	indexAsString := strings.ReplaceAll(fieldIndex, NameField, "")
	index, err := strconv.Atoi(indexAsString)
	if err != nil {
		return tokenizedField, err
	}
	tokenizedField.Index = index
	tokenizedField.Name = fieldName
	return tokenizedField, nil
}

func GetAttributeAsString(data json.RawMessage, jqString string) (string, error) {
	result, err := GetStringValue(string(data), jqString)
	if err != nil {
		return "", err
	}
	if result == "null" {
		result = ""
	}
	result = strings.ReplaceAll(result, "\"", "")
	return result, nil
}

func GetAttributeAsBool(data json.RawMessage, jqString string) (bool, error) {
	result, err := GetStringValue(string(data), jqString)
	if err != nil {
		return false, err
	}
	return (result == "true"), nil
}

func trimLeftDot(s string) string {
	if s[0] == '.' {
		return trimLeftChar(s)
	}
	return s
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func (transformation *Transformation) Transform(payload *DataPayload) (Transformed, error) {
	var transformed Transformed
	data, err := json.Marshal(payload)
	if err != nil {
		return transformed, err
	}
	switch transformation.Type {
	case TransformationTypeRandom:
		transformed, err = transformation.TransformFromRandom()
	case TransformationTypeTimestamp:
		transformed, err = transformation.TransformFromTimestamp()
	case TransformationTypeJq:
		transformed, err = transformation.TransformFromJq(data)
	case TransformationTypeConcatenate:
		transformed, err = transformation.TransformFromConcatenate(data)
	case TransformationTypeTranslate:
		transformed, err = transformation.TransformFromTranslate(data)
	case TransformationTypeParse:
		transformed, err = transformation.TransformFromParse(data)
	case TransformationTypeStaticMapping:
		transformed, err = transformation.TransformFromStaticMapping(data)
	case TransformationTypeTokenize:
		transformed, err = transformation.TransformFromTokenization(data)
	case TransformationTypeTemplate:
		transformed, err = transformation.TransformFromTemplate(data)

	}
	return transformed, err
}

func (transformation *Transformation) GetRegularExpressions() ([]string, error) {
	var regExes []string

	if len(transformation.Parameters) == 0 {
		return regExes, fmt.Errorf("regular expressions not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameRegularExpression {
			if param.Value != "" {
				regExes = append(regExes, param.Value)
			} else {
				regExes = param.Values
			}
			return regExes, nil
		}
	}
	return regExes, fmt.Errorf("regular expressions not found")
}

func (transformation *Transformation) GetJqString() (string, error) {
	if len(transformation.Parameters) == 0 {
		return "", fmt.Errorf("jq string not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameJqQuery {
			return param.Value, nil
		}
	}
	return "", fmt.Errorf("jq string not found")
}

func (transformation *Transformation) GetJqStrings() ([]string, error) {
	if len(transformation.Parameters) == 0 {
		return nil, fmt.Errorf("jq string not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameJqQuery {
			return param.Values, nil
		}
	}
	return nil, fmt.Errorf("jq string not found")
}

func (transformation *Transformation) GetLength() (int, error) {
	if len(transformation.Parameters) == 0 {
		return -1, fmt.Errorf("jq string not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameLength {
			length, err := strconv.Atoi(param.Value)
			if err != nil {
				return -1, err
			}
			return length, nil
		}
	}
	return -1, fmt.Errorf("jq string not found")
}

func (transformation *Transformation) GetStaticMapping() (string, error) {
	if len(transformation.Parameters) == 0 {
		return "", fmt.Errorf("static mappings not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameStaticValue {
			return param.Value, nil
		}
	}
	return "", fmt.Errorf("static mappings not found")
}

func (transformation *Transformation) GetTemplate() (string, error) {
	if len(transformation.Parameters) == 0 {
		return "", fmt.Errorf("template not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameTemplate {
			return param.Value, nil
		}
	}
	return "", fmt.Errorf("template not found")
}

func (transformation *Transformation) GetDelimiter() (string, error) {
	if len(transformation.Parameters) == 0 {
		return "", fmt.Errorf("delimiter not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameDelimiter {
			return param.Value, nil
		}
	}
	return "", fmt.Errorf("delimiter not found")
}

func (transformation *Transformation) GetJoinSeparator() (string, error) {
	if len(transformation.Parameters) == 0 {
		return "", fmt.Errorf("join separator not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameJoinSeparator {
			return param.Value, nil
		}
	}
	return "", fmt.Errorf("join separator not found")
}

func (transformation *Transformation) GetTranslation() (string, error) {
	if len(transformation.Parameters) == 0 {
		return "", fmt.Errorf("translation not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameTranslation {
			return param.Value, nil
		}
	}
	return "", fmt.Errorf("translation not found")
}

func (transformation *Transformation) GetTokenIndexesToTranslate() ([]int, error) {
	var indexes []int
	if len(transformation.Parameters) == 0 {
		return indexes, fmt.Errorf("token-indices not found")
	}
	for _, param := range transformation.Parameters {
		if param.Name == NameTokenIndices {
			if len(param.Values) > 0 {
				for _, value := range param.Values {
					index, err := strconv.Atoi(value)
					if err != nil {
						return indexes, fmt.Errorf("misconfigured token-indices index: %s", value)
					}
					indexes = append(indexes, index)
				}
			} else if param.Value != "" {
				index, err := strconv.Atoi(param.Value)
				if err != nil {
					return indexes, fmt.Errorf("misconfigured token-indices index: %s", param.Value)
				}
				indexes = append(indexes, index)
			}
			return indexes, nil
		}
	}
	return indexes, fmt.Errorf("token-indices not found")
}

func (transformation *Transformation) GetTokenizedFields(data json.RawMessage) ([]TokenizedField, error) {
	var fields []TokenizedField

	if len(transformation.Parameters) == 0 {
		return fields, fmt.Errorf("replacements not found")
	}
	for _, param := range transformation.Parameters {
		if strings.Contains(param.Name, NameField) {
			tokenizedField, err := NewTokenizedField(param.Name, param.Value)
			if err != nil {
				return nil, err
			}
			fields = append(fields, tokenizedField)
		}
	}
	return fields, nil

}

func (transformation *Transformation) GetFields(data json.RawMessage) ([]string, error) {
	var fields []string
	value := ""
	var err error

	if len(transformation.Parameters) == 0 {
		return fields, fmt.Errorf("parameters not found")
	}
	for _, param := range transformation.Parameters {
		if strings.Contains(param.Name, NameField) {
			if len(param.Values) == 0 {
				value, err = GetAttributeAsString(data, param.Value)
			} else {
				value, err = GetAttributeAsStringFirstFound(data, param.Values)
			}
			if err != nil {
				continue
			}
			fields = append(fields, value)
		}
	}
	return fields, nil
}

func (transformation *Transformation) GetTemplateValues(data json.RawMessage) (map[string]string, error) {
	templateValues := make(map[string]string)
	var err error
	value := ""
	if len(transformation.Parameters) == 0 {
		return templateValues, fmt.Errorf("template values not found")
	}

	for _, param := range transformation.Parameters {
		if !strings.Contains(param.Name, NameTemplate) {
			if len(param.Values) == 0 {
				value, err = GetAttributeAsString(data, param.Value)
			} else {
				value, err = GetAttributeAsStringFirstFound(data, param.Values)
			}
			if err != nil {
				continue
			}
			templateValues[param.Name] = value
		}
	}
	return templateValues, nil
}

func GetAttributeAsStringFirstFound(data json.RawMessage, parameterValues []string) (string, error) {
	var err error
	value := ""

	for _, parameterValue := range parameterValues {
		value, err = GetAttributeAsString(data, parameterValue)
		if err != nil {
			continue
		}
		if value != "" {
			break
		}
	}

	return value, nil
}

func (transformation *Transformation) TransformFromJqFirstFound(data json.RawMessage) (Transformed, error) {
	transformed := make(Transformed)
	jqStrings, err := transformation.GetJqStrings()
	if err != nil {
		return nil, err
	}
	value := ""
	for _, jqString := range jqStrings {
		value, err = GetAttributeAsString(data, jqString)
		if err != nil {
			continue
		}
		if value != "" {
			break
		}
	}

	transformed[transformation.Name] = value
	return transformed, nil
}

func (transformation *Transformation) TransformFromTimestamp() (Transformed, error) {
	transformed := make(Transformed)
	transformed[transformation.Name] = strconv.FormatInt(time.Now().Unix(), 10)
	return transformed, nil
}

func (transformation *Transformation) TransformFromRandom() (Transformed, error) {
	transformed := make(Transformed)
	length, err := transformation.GetLength()
	if err != nil {
		return transformed, err
	}
	rando, err := GenerateRandomString(length)
	if err != nil {
		return transformed, err
	}
	template, err := transformation.GetTemplate()
	if err == nil {
		templatizedKey := fmt.Sprintf("{{%s}}", transformation.Name)
		template = strings.ReplaceAll(template, templatizedKey, rando)
		transformed[transformation.Name] = template
	} else {
		transformed[transformation.Name] = rando
	}

	return transformed, nil
}

func (transformation *Transformation) TransformFromJq(data json.RawMessage) (Transformed, error) {
	transformed := make(Transformed)

	jqStrings, err := transformation.GetJqStrings()
	if err != nil {
		return nil, err
	}
	if len(jqStrings) > 0 {
		return transformation.TransformFromJqFirstFound(data)
	}

	jqString, err := transformation.GetJqString()
	if err != nil {
		return nil, err
	}
	value, err := GetAttributeAsString(data, jqString)
	// jq is tricky. in this context, a jq error should be handled as a miss
	if err != nil {
		return transformed, nil
	}

	// split because we *arbitrarily* assume all data that is comma delimeted is an array
	// in the case of EC2 autoscaling cloudwatch events this is true
	// this may need to change
	transformed[transformation.Name] = value
	return transformed, nil
}

func (transformation *Transformation) TransformFromConcatenate(data json.RawMessage) (Transformed, error) {
	transformed := make(Transformed)

	delimiter, err := transformation.GetDelimiter()
	if err != nil {
		return nil, err
	}

	fields, err := transformation.GetFields(data)
	if err != nil {
		return nil, err
	}
	// test for empty
	_transformed := strings.Join(fields, "")
	if _transformed != "" {
		_transformed = strings.Join(fields, delimiter)
	}

	transformed[transformation.Name] = _transformed

	return transformed, nil
}

func (transformation *Transformation) GetTokenToTranslateBySplit(value, joinSeparator string, indexes []int) (string, error) {
	var tokensToJoin []string
	delimiter, err := transformation.GetDelimiter()
	if err != nil {
		return "", err
	}

	pieces := strings.Split(value, delimiter)

	for _, index := range indexes {
		if len(pieces) >= index {
			tokensToJoin = append(tokensToJoin, pieces[index-1])
		}
	}
	return strings.ToLower(strings.Join(tokensToJoin, joinSeparator)), nil
}

func (transformation *Transformation) GetTokenToTranslateByRegEx(value, joinSeparator string, indexes []int) (string, error) {
	var tokensToJoin []string
	regExes, _ := transformation.GetRegularExpressions()

	for _, regExText := range regExes {
		regEx := regexp.MustCompile(regExText)
		matches := regEx.FindStringSubmatch(value)
		if len(matches) > 0 {
			for _, index := range indexes {
				if len(matches) >= index-1 {
					tokensToJoin = append(tokensToJoin, matches[index])
				}
			}
			return strings.ToLower(strings.Join(tokensToJoin, joinSeparator)), nil
		}
	}
	return "", fmt.Errorf("no match for regular expression")
}

func (transformation *Transformation) GetTokenToTranslate(value string) (string, error) {
	joinSeparator, err := transformation.GetJoinSeparator()
	if err != nil {
		joinSeparator = ""
	}

	indexes, err := transformation.GetTokenIndexesToTranslate()
	if err != nil {
		return "", err
	}
	_, err = transformation.GetRegularExpressions()
	if err != nil {
		return transformation.GetTokenToTranslateBySplit(value, joinSeparator, indexes)
	}

	return transformation.GetTokenToTranslateByRegEx(value, joinSeparator, indexes)
}

func (transformation *Transformation) TransformFromParse(data json.RawMessage) (Transformed, error) {
	transformed, err := transformation.TransformFromJq(data)
	if err != nil {
		return nil, err
	}

	token, err := transformation.GetTokenToTranslate(string(transformed[transformation.Name]))
	if err != nil {
		return nil, err
	}
	if token != "" {
		transformed[transformation.Name] = token
		return transformed, nil
	}

	return nil, fmt.Errorf("no replacement")
}

func (transformation *Transformation) TransformFromTranslate(data json.RawMessage) (Transformed, error) {
	transformed, err := transformation.TransformFromJq(data)
	if err != nil {
		return nil, err
	}

	translation, err := transformation.GetTranslation()
	if err != nil {
		return nil, err
	}
	if (transformation.Translation) == nil {
		return nil, fmt.Errorf("no translation table")
	}
	table, ok := (*transformation.Translation)[translation]
	if !ok {
		return nil, fmt.Errorf("invalid translation table: %s", translation)
	}
	replacement, ok := table.Table[string(transformed[transformation.Name])]
	if ok {
		transformed[transformation.Name] = replacement
		return transformed, nil
	}
	transformed[transformation.Name] = ""
	return transformed, nil
}

func (transformation *Transformation) TransformFromTemplate(data json.RawMessage) (Transformed, error) {
	transformed := make(Transformed)
	template, err := transformation.GetTemplate()
	if err != nil {
		return nil, err
	}
	templateValues, err := transformation.GetTemplateValues(data)
	if err != nil {
		return nil, err
	}
	for key, value := range templateValues {
		templatizedKey := fmt.Sprintf("{{%s}}", key)
		template = strings.ReplaceAll(template, templatizedKey, value)
	}
	transformed[transformation.Name] = template
	return transformed, nil

}

func (transformation *Transformation) TransformFromTokenization(data json.RawMessage) (Transformed, error) {
	transformed := make(Transformed)
	fromTransformed, err := transformation.TransformFromJq(data)
	if err != nil {
		return nil, err
	}
	// get delimiter
	delimiter, err := transformation.GetDelimiter()
	if err != nil {
		return nil, err
	}
	// get tokens
	tokenizedFields, err := transformation.GetTokenizedFields(data)
	if err != nil {
		return nil, err
	}
	// parse
	stringToParse := string(fromTransformed[transformation.Name])
	tokens := strings.Split(stringToParse, delimiter)
	// range tokens, add _transformed attr
	for _, tokenizedField := range tokenizedFields {
		if len(tokens) > tokenizedField.Index {
			transformed[tokenizedField.Name] = tokens[tokenizedField.Index]
		}
	}
	return transformed, nil
}

func (transformation *Transformation) TransformFromStaticMapping(data json.RawMessage) (Transformed, error) {
	transformed := make(Transformed)
	staticMapping, err := transformation.GetStaticMapping()
	if err != nil {
		return transformed, err
	}

	transformed[transformation.Name] = staticMapping
	return transformed, nil
}

func (rule *MatchingRule) IsTrue(data json.RawMessage) bool {
	isTrue, _ := GetAttributeAsBool(data, rule.Attribute)
	return isTrue
}

func (rule *MatchingRule) Contains(data json.RawMessage) bool {
	var attributeValue string
	var err error
	if rule.Attributes == nil {
		attributeValue, err = GetAttributeAsString(data, rule.Attribute)
	} else {
		attributeValue, err = GetAttributeAsStringFirstFound(data, rule.Attributes)
	}

	if err != nil {
		util.GetLogger().Error("error querying data json", zap.Error(err))
	}
	if rule.Matcher != nil {
		return rule.Matcher.MatchString(attributeValue)
	}
	return false
}

func (rule *MatchingRule) Defined(data json.RawMessage) bool {
	return strings.Contains(string(data), rule.Attribute)
}

func (rule *MatchingRule) setMatcher(matcher *regexp.Regexp) {
	rule.Matcher = matcher
}

func (transformerList *TransformerList) SetTransformerDefaults(defaults *TransformerDefaults) {
	for index := range transformerList.Transformers {
		transformerList.Transformers[index].TransformerDefaults = defaults
	}
}

func (transformer *Transformer) SetTranslations(translation *Translations) {
	for index := range transformer.Transformations {
		transformer.Transformations[index].Translation = translation
	}
}

func (transformer *Transformer) TransformDefaultAttributes(payload *DataPayload) (*DataPayload, error) {
	if transformer.TransformerDefaults != nil {
		for _, transformation := range transformer.TransformerDefaults.Transformations {
			transformed, err := transformation.Transform(payload)
			if err != nil {
				return payload, err
			}

			for key, value := range transformed {
				if len(value) > 0 {
					payload.Outputs[key] = value
				}
			}
		}
	}

	return payload, nil
}

func (transformer *Transformer) Transform(payload *DataPayload) (*DataPayload, error) {
	payload.Outputs[TransformerName] = transformer.Name
	for _, transformation := range transformer.Transformations {
		transformed, err := transformation.Transform(payload)
		if err != nil {
			return payload, err
		}
		for key, value := range transformed {
			if len(value) > 0 {
				payload.Outputs[key] = value
			}
		}
	}

	return payload, nil
}
func (transformerList *TransformerList) GetTransformerByName(name string) *Transformer {
	for _, transformer := range transformerList.Transformers {
		if transformer.Name == name {
			return &transformer
		}
	}
	return nil
}

func (transformer *Transformer) Matches(payload *DataPayload) bool {
	data, err := json.Marshal(payload)
	if err != nil {
		return false
	}
	evaluated := false
	matched := true
	for _, rule := range transformer.MatchingRules {
		switch rule.RuleType {
		case RuleTypeDefined:
			matched = rule.Defined(data) && matched
			evaluated = true
		case RuleTypeUndefined:
			matched = !rule.Defined(data) && matched
			evaluated = true
		case RuleTypeContains:
			matched = rule.Contains(data) && matched
			evaluated = true
		case RuleTypeNotContained:
			matched = !rule.Contains(data) && matched
			evaluated = true
		case RuleTypeIsFalse:
			matched = !rule.IsTrue(data) && matched
			evaluated = true
		case RuleTypeIsTrue:
			matched = rule.IsTrue(data) && matched
			evaluated = true
		}
	}

	return evaluated && matched
}

func (destination *Destination) Send(message *DataPayload) (string, error) {
	messageID := time.Now().Format("2006-01-02 15:04:05.000000000")
	var err error
	switch destination.Type {
	case DestinationTypeLog:
		util.GetLogger().Info("%s\n", zap.Any("outputs", message.Outputs))
	}
	return messageID, err
}

func (transformer *Transformer) RemoveEphemeralTransformations(transformed *Transformed) {
	for _, transformation := range transformer.Transformations {
		if transformation.Ephemeral {
			delete(*transformed, transformation.Name)
		}
	}
}

func (transformer *Transformer) Transformer(payload *DataPayload) *DataPayload {
	var err error
	if transformer.Matches(payload) {
		payload, err = transformer.TransformDefaultAttributes(payload)
		if err != nil {
			util.GetLogger().Error(fmt.Sprintf("transformer %s default transformation error: ", transformer.Name), zap.Error(err))
			return nil
		}
		payload, err = transformer.Transform(payload)
		if err != nil {
			util.GetLogger().Error(fmt.Sprintf("transformer %s transformation error: ", transformer.Name), zap.Error(err))
			return nil
		}
		transformer.RemoveEphemeralTransformations(&payload.Outputs)

		for _, destination := range transformer.Destinations {
			_, err := destination.Send(payload)
			if err != nil {
				util.GetLogger().Error("error sending", zap.Error(err))
			}
		}
		return payload
	}
	return nil
}

func (transformers *TransformerList) Compile(translations *Translations) {
	for transformerIndex, transformer := range transformers.Transformers {
		transformers.Transformers[transformerIndex].SetTranslations(translations)
		for ruleIndex, rule := range transformer.MatchingRules {
			if rule.RuleType == RuleTypeContains || rule.RuleType == RuleTypeNotContained {
				transformers.Transformers[transformerIndex].MatchingRules[ruleIndex].setMatcher(regexp.MustCompile(strings.Join(rule.Values, "|")))
			}
		}
	}
}

func (transformers *TransformerList) Stage(stage string) {
	for transformerIndex, transformer := range transformers.Transformers {
		for destinationIndex := range transformer.Destinations {
			transformers.Transformers[transformerIndex].Destinations[destinationIndex].Name = strings.ReplaceAll(transformers.Transformers[transformerIndex].Destinations[destinationIndex].Name, StageToken, stage)
			if stage == StageNull {
				transformers.Transformers[transformerIndex].Destinations[destinationIndex].Type = StageNull
			}
		}
	}
}

func (transformers *TransformerList) Transformer(payload *DataPayload) []*DataPayload {
	var transformations SafeDataPayloads
	g, _ := errgroup.WithContext(context.Background())

	for index := range transformers.Transformers {
		transformer := transformers.Transformers[index]
		g.Go(func() error {
			transformation := transformer.Transformer(payload.Clone())
			if transformation != nil {
				transformations.Append(transformation)
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil
	}

	return transformations.DataPayloads
}

func IsYamlFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fileInfo.IsDir() {
		return false
	}
	return filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml"
}

func LoadTransformer(file, stage string, translations *Translations, defaults *TransformerDefaults) (*TransformerList, error) {
	var list []Transformer
	yfile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yfile, &list)

	if err != nil {
		return nil, err
	}
	transformerList := &TransformerList{Transformers: list}
	transformerList.Compile(translations)
	transformerList.Stage(stage)
	transformerList.SetTransformerDefaults(defaults)
	return transformerList, nil
}

func LoadTransformers(path, stage string, translations *Translations, defaults *TransformerDefaults) (*TransformerList, error) {
	transformerList := &TransformerList{}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if IsYamlFile(path) {
			loaded, err := LoadTransformer(path, stage, translations, defaults)
			if err != nil {
				util.GetLogger().Error(fmt.Sprintf("initialization: transformer file in error (%s)", path), zap.Error(err))
			}
			transformerList.Transformers = append(transformerList.Transformers, loaded.Transformers...)
		}
		return nil
	})

	return transformerList, err
}

func LoadTranslationFile(file string) (TranslationList, error) {
	var translationList TranslationList
	yfile, err := os.ReadFile(file)
	if err != nil {
		return translationList, err
	}
	lower := strings.ToLower(string(yfile))
	err = yaml.Unmarshal([]byte(lower), &translationList)

	if err != nil {
		return translationList, err
	}
	return translationList, nil
}

func LoadTranslations(path string) (*Translations, error) {
	translations := make(Translations)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if IsYamlFile(path) {
			translationList, err := LoadTranslationFile(path)
			if err != nil {
				util.GetLogger().Error(fmt.Sprintf("initialization: translation file in error (%s)", path), zap.Error(err))
			}
			for index, translation := range translationList {
				translations[translation.Name] = translationList[index]
			}
		}
		return nil
	})

	return &translations, err
}

func LoadDefaults(path string) (*TransformerList, error) {
	var config []Transformer
	yfile, err := os.ReadFile(path)
	util.GetLogger().Info(fmt.Sprintf("initialization: LoadDefaults(%s)", path))
	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("initialization: failed to read %s", path), zap.Error(err))
		return nil, err
	}
	err = yaml.Unmarshal(yfile, &config)

	if err != nil {
		util.GetLogger().Error(fmt.Sprintf("initialization: failed to read %s", path), zap.Error(err))
		return nil, err
	}
	transformerList := &TransformerList{Transformers: config}
	return transformerList, nil
}

func (transformation *Transformation) MungeJqStrings(jqStrings []string) []string {
	mungedJqStrings := make([]string, 0)
	for _, rawJqString := range jqStrings {
		trimmed := trimLeftDot(rawJqString)
		mungedJqStrings = append(mungedJqStrings, fmt.Sprintf(".%s", trimmed))
		mungedJqStrings = append(mungedJqStrings, fmt.Sprintf(".input.%s", trimmed))
		mungedJqStrings = append(mungedJqStrings, fmt.Sprintf(".outputs.\"%s\".[]", trimmed))
	}
	return mungedJqStrings
}

func (transformation *Transformation) FindAndMungeJqStrings() {
	rawJqStrings := make([]string, 0)
	for index, param := range transformation.Parameters {
		if param.Name == NameJqQuery {
			if len(param.Value) > 0 {
				rawJqStrings = append(rawJqStrings, param.Value)
			}
			if len(param.Values) > 0 {
				rawJqStrings = append(rawJqStrings, param.Values...)
			}
			transformation.Parameters[index].Value = ""
			transformation.Parameters[index].Values = transformation.MungeJqStrings(rawJqStrings)
		}
	}
}

func (transformation *Transformation) MungeAndTransform(payload *DataPayload) (Transformed, error) {
	switch transformation.Type {
	case TransformationTypeJq:
		transformation.FindAndMungeJqStrings()
	case TransformationTypeTranslate:
		transformation.FindAndMungeJqStrings()
	case TransformationTypeTokenize:
		transformation.FindAndMungeJqStrings()
	case TransformationTypeParse:
		transformation.FindAndMungeJqStrings()
	case TransformationTypeConcatenate:
		for index, param := range transformation.Parameters {
			if param.Name != NameDelimiter {
				rawJqStrings := make([]string, 0)
				if len(param.Value) > 0 {
					rawJqStrings = append(rawJqStrings, param.Value)
				}
				if len(param.Values) > 0 {
					rawJqStrings = append(rawJqStrings, param.Values...)
				}
				transformation.Parameters[index].Value = ""
				transformation.Parameters[index].Values = transformation.MungeJqStrings(rawJqStrings)
			}
		}
	case TransformationTypeTemplate:
		for index, param := range transformation.Parameters {
			if param.Name != NameTemplate {
				rawJqStrings := make([]string, 0)
				if len(param.Value) > 0 {
					rawJqStrings = append(rawJqStrings, param.Value)
				}
				if len(param.Values) > 0 {
					rawJqStrings = append(rawJqStrings, param.Values...)
				}
				transformation.Parameters[index].Value = ""
				transformation.Parameters[index].Values = transformation.MungeJqStrings(rawJqStrings)
			}
		}
	}
	return transformation.Transform(payload)
}
