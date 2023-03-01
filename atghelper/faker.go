/*
 * Copyright 2022 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package atghelper

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	util "github.com/typa01/go-utils"

	"github.com/bxcodec/faker/v3"

	"github.com/bxcodec/faker/v3/support/slice"
)

var (
	mu = &sync.Mutex{}
	// Sets nil if the value type is struct or map and the size of it equals to zero.
	shouldSetNil = false
	// Sets random integer generation to zero for slice and maps
	testRandZero = false
	// Sets the default number of string when it is created randomly.
	randomStringLen = 25
	// Sets the boundary for random value generation. Boundaries can not exceed integer(4 byte...)
	nBoundary = numberBoundary{start: 0, end: 100}
	// Sets the random size for slices and maps.
	randomSize = 100
	// Sets the single fake data generator to generate unique values
	generateUniqueValues = false
	// Sets whether interface{}s should be ignored.
	ignoreInterface = false
	// Unique values are kept in memory so the generator retries if the value already exists
	uniqueValues = map[string][]interface{}{}
	// Lang is selected language for random string generator
	lang = LangENG
	// How much tries for generating random string
	maxGenerateStringRetries = 1000000
)

type numberBoundary struct {
	start int
	end   int
}

type langRuneBoundary struct {
	start   rune
	end     rune
	exclude []rune
}

// Language rune boundaries here
var (
	// LangENG is for english language
	LangENG = langRuneBoundary{65, 122, []rune{91, 92, 93, 94, 95, 96}}
	// LangCHI is for chinese language
	LangCHI = langRuneBoundary{19968, 40869, nil}
	// LangRUS is for russian language
	LangRUS = langRuneBoundary{1025, 1105, nil}
)

// Supported tags
const (
	letterIdxBits         = 6                    // 6 bits to represent a letter index
	letterIdxMask         = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax          = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	maxRetry              = 10000                // max number of retry for unique values
	tagName               = "faker"
	keep                  = "keep"
	unique                = "unique"
	ID                    = "uuid_digit"
	HyphenatedID          = "uuid_hyphenated"
	EmailTag              = "email"
	MacAddressTag         = "mac_address"
	DomainNameTag         = "domain_name"
	UserNameTag           = "username"
	URLTag                = "url"
	IPV4Tag               = "ipv4"
	IPV6Tag               = "ipv6"
	PASSWORD              = "password"
	JWT                   = "jwt"
	LATITUDE              = "lat"
	LONGITUDE             = "long"
	CreditCardNumber      = "cc_number"
	CreditCardType        = "cc_type"
	PhoneNumber           = "phone_number"
	TollFreeNumber        = "toll_free_number"
	E164PhoneNumberTag    = "e_164_phone_number"
	TitleMaleTag          = "title_male"
	TitleFemaleTag        = "title_female"
	FirstNameTag          = "first_name"
	FirstNameMaleTag      = "first_name_male"
	FirstNameFemaleTag    = "first_name_female"
	LastNameTag           = "last_name"
	NAME                  = "name"
	GENDER                = "gender"
	UnixTimeTag           = "unix_time"
	DATE                  = "date"
	TIME                  = "time"
	MonthNameTag          = "month_name"
	YEAR                  = "year"
	DayOfWeekTag          = "day_of_week"
	DayOfMonthTag         = "day_of_month"
	TIMESTAMP             = "timestamp"
	CENTURY               = "century"
	TIMEZONE              = "timezone"
	TimePeriodTag         = "time_period"
	WORD                  = "word"
	SENTENCE              = "sentence"
	PARAGRAPH             = "paragraph"
	CurrencyTag           = "currency"
	AmountTag             = "amount"
	AmountWithCurrencyTag = "amount_with_currency"
	SKIP                  = "-"
	Length                = "len"
	SliceLength           = "slice_len"
	Language              = "lang"
	BoundaryStart         = "boundary_start"
	BoundaryEnd           = "boundary_end"
	Equals                = "="
	comma                 = ","
	colon                 = ":"
	ONEOF                 = "oneof"
	// period                = "."
	// hyphen = "-"
)

var defaultTag = map[string]string{
	EmailTag:              EmailTag,
	MacAddressTag:         MacAddressTag,
	DomainNameTag:         DomainNameTag,
	URLTag:                URLTag,
	UserNameTag:           UserNameTag,
	IPV4Tag:               IPV4Tag,
	IPV6Tag:               IPV6Tag,
	PASSWORD:              PASSWORD,
	JWT:                   JWT,
	CreditCardType:        CreditCardType,
	CreditCardNumber:      CreditCardNumber,
	LATITUDE:              LATITUDE,
	LONGITUDE:             LONGITUDE,
	PhoneNumber:           PhoneNumber,
	TollFreeNumber:        TollFreeNumber,
	E164PhoneNumberTag:    E164PhoneNumberTag,
	TitleMaleTag:          TitleMaleTag,
	TitleFemaleTag:        TitleFemaleTag,
	FirstNameTag:          FirstNameTag,
	FirstNameMaleTag:      FirstNameMaleTag,
	FirstNameFemaleTag:    FirstNameFemaleTag,
	LastNameTag:           LastNameTag,
	NAME:                  NAME,
	GENDER:                GENDER,
	UnixTimeTag:           UnixTimeTag,
	DATE:                  DATE,
	MonthNameTag:          MonthNameTag,
	DayOfWeekTag:          DayOfWeekTag,
	TIMESTAMP:             TIMESTAMP,
	CENTURY:               CENTURY,
	TIMEZONE:              TIMEZONE,
	WORD:                  WORD,
	SENTENCE:              SENTENCE,
	PARAGRAPH:             PARAGRAPH,
	CurrencyTag:           CurrencyTag,
	AmountTag:             AmountTag,
	AmountWithCurrencyTag: AmountWithCurrencyTag,
	ID:                    ID,
	HyphenatedID:          HyphenatedID,
}

// TaggedFunction used as the standard layout function for tag providers in struct.
// This type also can be used for custom provider.
type TaggedFunction func(v reflect.Value) (interface{}, error)

var MapperTag = mapperTag

var mapperTag = map[string]TaggedFunction{
	EmailTag:              faker.GetNetworker().Email,
	MacAddressTag:         faker.GetNetworker().MacAddress,
	DomainNameTag:         faker.GetNetworker().DomainName,
	URLTag:                faker.GetNetworker().URL,
	UserNameTag:           faker.GetNetworker().UserName,
	IPV4Tag:               faker.GetNetworker().IPv4,
	IPV6Tag:               faker.GetNetworker().IPv6,
	PASSWORD:              faker.GetNetworker().Password,
	JWT:                   faker.GetNetworker().Jwt,
	CreditCardType:        faker.GetPayment().CreditCardType,
	CreditCardNumber:      faker.GetPayment().CreditCardNumber,
	LATITUDE:              faker.GetAddress().Latitude,
	LONGITUDE:             faker.GetAddress().Longitude,
	PhoneNumber:           faker.GetPhoner().PhoneNumber,
	TollFreeNumber:        faker.GetPhoner().TollFreePhoneNumber,
	E164PhoneNumberTag:    faker.GetPhoner().E164PhoneNumber,
	TitleMaleTag:          faker.GetPerson().TitleMale,
	TitleFemaleTag:        faker.GetPerson().TitleFeMale,
	FirstNameTag:          faker.GetPerson().FirstName,
	FirstNameMaleTag:      faker.GetPerson().FirstNameMale,
	FirstNameFemaleTag:    faker.GetPerson().FirstNameFemale,
	LastNameTag:           faker.GetPerson().LastName,
	NAME:                  faker.GetPerson().Name,
	GENDER:                faker.GetPerson().Gender,
	UnixTimeTag:           faker.GetDateTimer().UnixTime,
	DATE:                  faker.GetDateTimer().Date,
	TIME:                  faker.GetDateTimer().Time,
	MonthNameTag:          faker.GetDateTimer().MonthName,
	YEAR:                  faker.GetDateTimer().Year,
	DayOfWeekTag:          faker.GetDateTimer().DayOfWeek,
	DayOfMonthTag:         faker.GetDateTimer().DayOfMonth,
	TIMESTAMP:             faker.GetDateTimer().Timestamp,
	CENTURY:               faker.GetDateTimer().Century,
	TIMEZONE:              faker.GetDateTimer().TimeZone,
	TimePeriodTag:         faker.GetDateTimer().TimePeriod,
	WORD:                  faker.GetLorem().Word,
	SENTENCE:              faker.GetLorem().Sentence,
	PARAGRAPH:             faker.GetLorem().Paragraph,
	CurrencyTag:           faker.GetPrice().Currency,
	AmountTag:             faker.GetPrice().Amount,
	AmountWithCurrencyTag: faker.GetPrice().AmountWithCurrency,
	ID:                    faker.GetIdentifier().Digit,
	HyphenatedID:          faker.GetIdentifier().Hyphenated,
}

// Generic Error Messages for tags
// 		ErrUnsupportedKindPtr: Error when get fake from ptr
// 		ErrUnsupportedKind: Error on passing unsupported kind
// 		ErrValueNotPtr: Error when value is not pointer
// 		ErrTagNotSupported: Error when tag is not supported
// 		ErrTagAlreadyExists: Error when tag exists and call AddProvider
// 		ErrTagDoesNotExist: Error when tag does not exist and call RemoveProvider
// 		ErrMoreArguments: Error on passing more arguments
// 		ErrNotSupportedPointer: Error when passing unsupported pointer
var (
	ErrUnsupportedKindPtr  = "Unsupported kind: %s Change Without using * (pointer) in Field of %s"
	ErrUnsupportedKind     = "Unsupported kind: %s"
	ErrValueNotPtr         = "Not a pointer value"
	ErrTagNotSupported     = "Tag unsupported: %s"
	ErrTagAlreadyExists    = "Tag exists"
	ErrTagDoesNotExist     = "Tag does not exist"
	ErrMoreArguments       = "Passed more arguments than is possible : (%d)"
	ErrNotSupportedPointer = "Use sample:=new(%s)\n faker.FakeData(sample) instead"
	ErrSmallerThanZero     = "Size:%d is smaller than zero."
	ErrSmallerThanOne      = "Size:%d is smaller than one."
	ErrUniqueFailure       = "Failed to generate a unique value for field \"%s\""

	ErrStartValueBiggerThanEnd = "Start value can not be bigger than end value."
	ErrWrongFormattedTag       = "Tag \"%s\" is not written properly"
	ErrUnknownType             = "Unknown Type"
	ErrNotSupportedTypeForTag  = "Type is not supported by tag."
	ErrUnsupportedTagArguments = "Tag arguments are not compatible with field type."
	ErrDuplicateSeparator      = "Duplicate separator for tag arguments."
	ErrNotEnoughTagArguments   = "Not enough arguments for tag."
	ErrUnsupportedNumberType   = "Unsupported Number type."
)

// Compiled regexp
var (
	findLangReg     *regexp.Regexp
	findLenReg      *regexp.Regexp
	findSliceLenReg *regexp.Regexp
)

func init() {
	rand.Seed(time.Now().UnixNano())
	findLangReg, _ = regexp.Compile("lang=[a-z]{3}")
	findLenReg, _ = regexp.Compile(`len=\d+`)
	findSliceLenReg, _ = regexp.Compile(`slice_len=\d+`)
}

// It is only for ptr.
func GetValue(t reflect.Type, level int) (reflect.Value, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("[GetValue] has error: %v\n", err)
		}
	}()
	if level >= 2 {
		return reflect.Zero(t), nil
	}
	if t == nil {
		return reflect.Value{}, fmt.Errorf("interface{} not allowed")
	}

	k := t.Kind()
	a := reflect.New(t)
	switch k {
	case reflect.Ptr:
		v := reflect.New(t.Elem())
		var val reflect.Value
		var err error
		if a != reflect.Zero(t).Interface() {
			val, err = GetValue(t.Elem(), level+1)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("[GetValue] has GetValue on t.Elem() err: %v", err)
			}
		} else {
			val, err = GetValue(v.Elem().Type(), level+1)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("[GetValue] has GetValue on v.Elem().Type() err: %v", err)
			}
		}
		v.Elem().Set(val.Convert(t.Elem()))
		return v, nil
	case reflect.Struct:
		switch t.String() {
		case "time.Time":
			ft := time.Now().Add(time.Duration(rand.Int63()))
			return reflect.ValueOf(ft), nil
		default:
			originalDataVal := reflect.New(t)
			v := reflect.New(t).Elem()
			if v.NumField() >= 10 {
				return reflect.Value{}, fmt.Errorf("too many struct field")
			}
			for i := 0; i < v.NumField(); i++ {
				if !v.Field(i).CanSet() {
					continue // to avoid panic to set on unexported field in struct
				}
				tags := decodeTags(t, i)
				switch {
				case tags.fieldType == "":
					val, err := GetValue(v.Field(i).Type(), level)
					if err != nil {
						return reflect.Value{}, fmt.Errorf("[GetValue] has GetValue on v.Field(i).Type() err: %v", err)
					}
					val = val.Convert(v.Field(i).Type())
					v.Field(i).Set(val)
				case tags.fieldType == SKIP:
					item := originalDataVal.Field(i).Interface()
					if v.CanSet() && item != nil {
						v.Field(i).Set(reflect.ValueOf(item))
					}
				default:
					err := setDataWithTag(v.Field(i).Addr(), tags.fieldType)
					if err != nil {
						return reflect.Value{}, fmt.Errorf("[GetValue] has setDataWithTag err: %v", err)
					}
				}
			}
			return v, nil
		}

	case reflect.String:
		res, err := randomString(randomStringLen, &lang)
		return reflect.ValueOf(res), err
	case reflect.Slice:
		len := randomSliceAndMapSize()
		if shouldSetNil && len == 0 {
			return reflect.Zero(t), nil
		}
		v := reflect.MakeSlice(t, len, len)
		for i := 0; i < v.Len(); i++ {
			val, err := GetValue(v.Index(i).Type(), level)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("[GetValue] has GetValue in reflect.Slice case err: %v", err)
			}
			val = val.Convert(v.Index(i).Type())
			v.Index(i).Set(val)
		}
		return v, nil
	case reflect.Array:
		v := reflect.New(t).Elem()
		for i := 0; i < v.Len(); i++ {
			val, err := GetValue(v.Index(i).Type(), level+1)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("[GetValue] has GetValue on v.Index(i).Type() in reflect.Array case err: %v", err)
			}
			val = val.Convert(v.Index(i).Type())
			v.Index(i).Set(val)
		}
		return v, nil
	case reflect.Int:
		return reflect.ValueOf(randomInteger()), nil
	case reflect.Int8:
		return reflect.ValueOf(int8(randomInteger())), nil
	case reflect.Int16:
		return reflect.ValueOf(int16(randomInteger())), nil
	case reflect.Int32:
		return reflect.ValueOf(int32(randomInteger())), nil
	case reflect.Int64:
		return reflect.ValueOf(int64(randomInteger())), nil
	case reflect.Float32:
		return reflect.ValueOf(rand.Float32()), nil
	case reflect.Float64:
		return reflect.ValueOf(rand.Float64()), nil
	case reflect.Bool:
		val := rand.Intn(2) > 0
		return reflect.ValueOf(val), nil

	case reflect.Uint:
		return reflect.ValueOf(uint(randomInteger())), nil

	case reflect.Uint8:
		return reflect.ValueOf(uint8(randomInteger())), nil

	case reflect.Uint16:
		return reflect.ValueOf(uint16(randomInteger())), nil

	case reflect.Uint32:
		return reflect.ValueOf(uint32(randomInteger())), nil

	case reflect.Uint64:
		return reflect.ValueOf(uint64(randomInteger())), nil

	case reflect.Map:
		len := randomSliceAndMapSize()
		if shouldSetNil && len == 0 {
			return reflect.Zero(t), nil
		}
		v := reflect.MakeMap(t)
		for i := 0; i < len; i++ {
			keyInstance := reflect.New(t.Key()).Elem().Type()
			if t.Elem().PkgPath() != "" && !IsExportByName(keyInstance.Name()) {
				var vNil interface{} = nil
				return reflect.ValueOf(vNil), fmt.Errorf("[GetValue] has GetValue on keyInstance unexport field Name: %v", keyInstance.Name())
			}
			key, err := GetValue(keyInstance, level)

			if err != nil {
				return reflect.Value{}, fmt.Errorf("[GetValue] has GetValue on keyInstance err: %v", err)
			}
			valueInstance := reflect.New(t.Elem()).Elem().Type()
			if t.Elem().PkgPath() != "" && !IsExportByName(valueInstance.Name()) {
				var vNil interface{} = nil
				return reflect.ValueOf(vNil), fmt.Errorf("[GetValue] has GetValue on valueInstance unexport field Name: %v", valueInstance.Name())
			}
			val, err := GetValue(valueInstance, level)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("[GetValue] has GetValue on valueInstance err: %v", err)
			}
			v.SetMapIndex(key, val)
		}
		return v, nil
	default:
		err := fmt.Errorf("no support for kind %+v", t)
		return reflect.Value{}, err
	}

}

func isZero(field reflect.Value) (bool, error) {
	if field.Kind() == reflect.Map {
		return field.Len() == 0, nil
	}

	for _, kind := range []reflect.Kind{reflect.Struct, reflect.Slice, reflect.Array} {
		if kind == field.Kind() {
			return false, fmt.Errorf("keep not allowed on struct")
		}
	}
	return reflect.Zero(field.Type()).Interface() == field.Interface(), nil
}

func decodeTags(typ reflect.Type, i int) structTag {
	tags := strings.Split(typ.Field(i).Tag.Get(tagName), ",")

	keepOriginal := false
	uni := false
	res := make([]string, 0)
	for _, tag := range tags {
		if tag == keep {
			keepOriginal = true
			continue
		} else if tag == unique {
			uni = true
			continue
		}
		res = append(res, tag)
	}

	return structTag{
		fieldType:    strings.Join(res, ","),
		unique:       uni,
		keepOriginal: keepOriginal,
	}
}

type structTag struct {
	fieldType    string
	unique       bool
	keepOriginal bool
}

func setDataWithTag(v reflect.Value, tag string) error {

	if v.Kind() != reflect.Ptr {
		return errors.New(ErrValueNotPtr)
	}
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Ptr:
		if _, exist := mapperTag[tag]; !exist {
			return fmt.Errorf(ErrTagNotSupported, tag)
		}
		if _, def := defaultTag[tag]; !def {
			res, err := mapperTag[tag](v)
			if err != nil {
				return fmt.Errorf("[setDataWithTag] has mapperTag with v err: %v", err)
			}
			v.Set(reflect.ValueOf(res))
			return nil
		}

		t := v.Type()
		newv := reflect.New(t.Elem())
		res, err := mapperTag[tag](newv.Elem())
		if err != nil {
			return fmt.Errorf("[setDataWithTag] has mapperTag with newv.Elem() err: %v", err)
		}
		rval := reflect.ValueOf(res)
		newv.Elem().Set(rval)
		v.Set(newv)
		return nil
	case reflect.String:
		return userDefinedString(v, tag)
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return userDefinedNumber(v, tag)
	case reflect.Slice, reflect.Array:
		return userDefinedArray(v, tag)
	case reflect.Map:
		return userDefinedMap(v, tag)
	default:
		if _, exist := mapperTag[tag]; !exist {
			return fmt.Errorf(ErrTagNotSupported, tag)
		}
		res, err := mapperTag[tag](v)
		if err != nil {
			return fmt.Errorf("[setDataWithTag] has mapperTag in default case err: %v", err)
		}
		v.Set(reflect.ValueOf(res))
	}
	return nil
}

func userDefinedMap(v reflect.Value, tag string) error {
	if tagFunc, ok := mapperTag[tag]; ok {
		res, err := tagFunc(v)
		if err != nil {
			return fmt.Errorf("[userDefinedMap] has tagFunc err: %v", err)
		}

		v.Set(reflect.ValueOf(res))
		return nil
	}

	len := randomSliceAndMapSize()
	if shouldSetNil && len == 0 {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	definedMap := reflect.MakeMap(v.Type())
	for i := 0; i < len; i++ {
		key, err := getValueWithTag(v.Type().Key(), tag)
		if err != nil {
			return fmt.Errorf("[userDefinedMap] has getValueWithTag with v.Type().Key() err: %v", err)
		}
		val, err := getValueWithTag(v.Type().Elem(), tag)
		if err != nil {
			return fmt.Errorf("[userDefinedMap] has getValueWithTag with v.Type().Elem() err: %v", err)
		}
		definedMap.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
	}
	v.Set(definedMap)
	return nil
}

func getValueWithTag(t reflect.Type, tag string) (interface{}, error) {
	switch t.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16, reflect.Uint, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64:
		res, err := extractNumberFromTag(tag, t)
		if err != nil {
			return nil, fmt.Errorf("[getValueWithTag] has extractNumberFromTag err: %v", err)
		}
		return res, nil
	case reflect.String:
		res, err := extractStringFromTag(tag)
		if err != nil {
			return nil, fmt.Errorf("[getValueWithTag] has extractStringFromTag err: %v", err)
		}
		return res, nil
	default:
		return 0, errors.New(ErrUnknownType)
	}
}

func userDefinedArray(v reflect.Value, tag string) error {
	_, tagExists := mapperTag[tag]
	if tagExists {
		res, err := mapperTag[tag](v)
		if err != nil {
			return fmt.Errorf("[userDefinedArray] has mapperTag err: %v", err)
		}
		v.Set(reflect.ValueOf(res))
		return nil
	}
	sliceLen, err := extractSliceLengthFromTag(tag)
	if err != nil {
		return fmt.Errorf("[userDefinedArray] has extractSliceLengthFromTag err: %v", err)
	}
	if shouldSetNil && sliceLen == 0 {
		v.Set(reflect.Zero(v.Type()))
		return nil
	}
	array := reflect.MakeSlice(v.Type(), sliceLen, sliceLen)
	for i := 0; i < sliceLen; i++ {
		res, err := getValueWithTag(v.Type().Elem(), tag)
		if err != nil {
			return fmt.Errorf("[userDefinedArray] has getValueWithTag err: %v", err)
		}
		array.Index(i).Set(reflect.ValueOf(res))
	}
	v.Set(array)
	return nil
}

func userDefinedString(v reflect.Value, tag string) error {
	var res interface{}
	var err error

	if tagFunc, ok := mapperTag[tag]; ok {
		res, err = tagFunc(v)
		if err != nil {
			return fmt.Errorf("[userDefinedString] has tagFunc err: %v", err)
		}
	} else {
		res, err = extractStringFromTag(tag)
		if err != nil {
			return fmt.Errorf("[userDefinedString] has extractStringFromTag err: %v", err)
		}
	}
	if res == nil {
		return fmt.Errorf(ErrTagNotSupported, tag)
	}
	val, _ := res.(string)
	v.SetString(val)
	return nil
}

func userDefinedNumber(v reflect.Value, tag string) error {
	var res interface{}
	var err error

	if tagFunc, ok := mapperTag[tag]; ok {
		res, err = tagFunc(v)
		if err != nil {
			return fmt.Errorf("[userDefinedNumber] has tagFunc err: %v", err)
		}
	} else {
		res, err = extractNumberFromTag(tag, v.Type())
		if err != nil {
			return fmt.Errorf("[userDefinedNumber] has extractNumberFromTag err: %v", err)
		}
	}
	if res == nil {
		return fmt.Errorf(ErrTagNotSupported, tag)
	}

	v.Set(reflect.ValueOf(res))
	return nil
}

// extractSliceLengthFromTag checks if the sliceLength tag 'slice_len' is set, if so, returns its value, else return a random length
func extractSliceLengthFromTag(tag string) (int, error) {
	if strings.Contains(tag, SliceLength) {
		lenParts := strings.SplitN(findSliceLenReg.FindString(tag), Equals, -1)
		if len(lenParts) != 2 {
			return 0, fmt.Errorf(ErrWrongFormattedTag, tag)
		}
		sliceLen, err := strconv.Atoi(lenParts[1])
		if err != nil {
			return 0, fmt.Errorf("the given sliceLength has to be numeric, tag: %s", tag)
		}
		if sliceLen < 0 {
			return 0, fmt.Errorf("slice length can not be negative, tag: %s", tag)
		}
		return sliceLen, nil
	}

	return randomSliceAndMapSize(), nil // Returns random slice length if the sliceLength tag isn't set
}

func extractStringFromTag(tag string) (interface{}, error) {
	var err error
	strlen := randomStringLen
	strlng := &lang
	isOneOfTag := strings.Contains(tag, ONEOF)
	if !strings.Contains(tag, Length) && !strings.Contains(tag, Language) && !isOneOfTag {
		return nil, fmt.Errorf(ErrTagNotSupported, tag)
	}
	if strings.Contains(tag, Length) {
		lenParts := strings.SplitN(findLenReg.FindString(tag), Equals, -1)
		if len(lenParts) != 2 {
			return nil, fmt.Errorf(ErrWrongFormattedTag, tag)
		}
		strlen, _ = strconv.Atoi(lenParts[1])
	}
	if strings.Contains(tag, Language) {
		strlng, err = extractLangFromTag(tag)
		if err != nil {
			return nil, fmt.Errorf(ErrWrongFormattedTag, tag)
		}
	}
	if isOneOfTag {
		items := strings.Split(tag, colon)
		argsList := items[1:]
		if len(argsList) != 1 {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		if strings.Contains(argsList[0], ",,") {
			return nil, fmt.Errorf(ErrDuplicateSeparator)
		}
		args := strings.Split(argsList[0], comma)
		if len(args) < 2 {
			return nil, fmt.Errorf(ErrNotEnoughTagArguments)
		}
		toRet := args[rand.Intn(len(args))]
		return strings.TrimSpace(toRet), nil
	}
	res, err := randomString(strlen, strlng)
	return res, err
}

func extractLangFromTag(tag string) (*langRuneBoundary, error) {
	text := findLangReg.FindString(tag)
	texts := strings.SplitN(text, Equals, -1)
	if len(texts) != 2 {
		return nil, fmt.Errorf(ErrWrongFormattedTag, text)
	}
	switch strings.ToLower(texts[1]) {
	case "eng":
		return &LangENG, nil
	case "rus":
		return &LangRUS, nil
	case "chi":
		return &LangCHI, nil
	default:
		return &LangENG, nil
	}
}

func extractNumberFromTag(tag string, t reflect.Type) (interface{}, error) {
	hasOneOf := strings.Contains(tag, ONEOF)
	hasBoundaryStart := strings.Contains(tag, BoundaryStart)
	hasBoundaryEnd := strings.Contains(tag, BoundaryEnd)
	usingOneOfTag := hasOneOf && (!hasBoundaryStart && !hasBoundaryEnd)
	usingBoundariesTags := !hasOneOf && (hasBoundaryStart && hasBoundaryEnd)
	if !usingOneOfTag && !usingBoundariesTags {
		return nil, fmt.Errorf(ErrTagNotSupported, tag)
	}

	// handling oneof tag
	if usingOneOfTag {
		argsList := strings.Split(tag, colon)[1:]
		if len(argsList) != 1 {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		if strings.Contains(argsList[0], ",,") {
			return nil, fmt.Errorf(ErrDuplicateSeparator)
		}
		args := strings.Split(argsList[0], comma)
		if len(args) < 2 {
			return nil, fmt.Errorf(ErrNotEnoughTagArguments)
		}
		switch t.Kind() {
		case reflect.Float64:
			{
				toRet, err := extractFloat64FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractFloat64FromTagArgs err: %v", err)
				}
				return toRet.(float64), nil
			}
		case reflect.Float32:
			{
				toRet, err := extractFloat32FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractFloat32FromTagArgs err: %v", err)
				}
				return toRet.(float32), nil
			}
		case reflect.Int64:
			{
				toRet, err := extractInt64FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractInt64FromTagArgs err: %v", err)
				}
				return toRet.(int64), nil
			}
		case reflect.Int32:
			{
				toRet, err := extractInt32FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractInt32FromTagArgs err: %v", err)
				}
				return toRet.(int32), nil
			}
		case reflect.Int16:
			{
				toRet, err := extractInt16FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractInt16FromTagArgs err: %v", err)
				}
				return toRet.(int16), nil
			}
		case reflect.Int8:
			{
				toRet, err := extractInt8FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractInt8FromTagArgs err: %v", err)
				}
				return toRet.(int8), nil
			}
		case reflect.Int:
			{
				toRet, err := extractIntFromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractIntFromTagArgs err: %v", err)
				}
				return toRet.(int), nil
			}
		case reflect.Uint64:
			{
				toRet, err := extractUint64FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractUint64FromTagArgs err: %v", err)
				}
				return toRet.(uint64), nil
			}
		case reflect.Uint32:
			{
				toRet, err := extractUint32FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractUint32FromTagArgs err: %v", err)
				}
				return toRet.(uint32), nil
			}
		case reflect.Uint16:
			{
				toRet, err := extractUint16FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractUint16FromTagArgs err: %v", err)
				}
				return toRet.(uint16), nil
			}
		case reflect.Uint8:
			{
				toRet, err := extractUint8FromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractUint8FromTagArgs err: %v", err)
				}
				return toRet.(uint8), nil
			}
		case reflect.Uint:
			{
				toRet, err := extractUintFromTagArgs(args)
				if err != nil {
					return nil, fmt.Errorf("[extractNumberFromTag] has extractUintFromTagArgs err: %v", err)
				}
				return toRet.(uint), nil
			}
		default:
			{
				return nil, fmt.Errorf(ErrUnsupportedNumberType)
			}
		}
	}

	// handling boundary tags
	valuesStr := strings.SplitN(tag, comma, -1)
	if len(valuesStr) != 2 {
		return nil, fmt.Errorf(ErrWrongFormattedTag, tag)
	}
	startBoundary, err := extractNumberFromText(valuesStr[0])
	if err != nil {
		return nil, fmt.Errorf("[extractNumberFromTag] has extractNumberFromText(valuesStr[0]) err: %v", err)
	}
	endBoundary, err := extractNumberFromText(valuesStr[1])
	if err != nil {
		return nil, fmt.Errorf("[extractNumberFromTag] has extractNumberFromText(valuesStr[1]) err: %v", err)
	}
	boundary := numberBoundary{start: startBoundary, end: endBoundary}
	switch t.Kind() {
	case reflect.Uint:
		return uint(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint8:
		return uint8(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint16:
		return uint16(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint32:
		return uint32(randomIntegerWithBoundary(boundary)), nil
	case reflect.Uint64:
		return uint64(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int:
		return randomIntegerWithBoundary(boundary), nil
	case reflect.Int8:
		return int8(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int16:
		return int16(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int32:
		return int32(randomIntegerWithBoundary(boundary)), nil
	case reflect.Int64:
		return int64(randomIntegerWithBoundary(boundary)), nil
	default:
		return nil, errors.New(ErrNotSupportedTypeForTag)
	}
}

func extractNumberFromText(text string) (int, error) {
	text = strings.TrimSpace(text)
	texts := strings.SplitN(text, Equals, -1)
	if len(texts) != 2 {
		return 0, fmt.Errorf(ErrWrongFormattedTag, text)
	}
	return strconv.Atoi(texts[1])
}

// TODO: add the random string logic
func randomString(n int, lang *langRuneBoundary) (string, error) {
	mLen := rand.Intn(28)
	builder := util.NewStringBuilder()
	for i := 0; i < mLen; i++ {
		builder.Append(RandStringBytes(1))
	}
	return builder.ToString(), nil
}

// randomIntegerWithBoundary returns a random integer between input start and end boundary. [start, end)
func randomIntegerWithBoundary(boundary numberBoundary) int {
	span := boundary.end - boundary.start
	if span <= 0 {
		return boundary.start
	}
	return rand.Intn(span) + boundary.start
}

// randomInteger returns a random integer between start and end boundary. [start, end)
// TODO: add random number here
func randomInteger() int {
	return 0
}

// randomSliceAndMapSize returns a random integer between [0,randomSliceAndMapSize). If the testRandZero is set, returns 0
// Written for test purposes for shouldSetNil
func randomSliceAndMapSize() int {
	if testRandZero {
		return 0
	}
	return 1
}

func randomElementFromSliceString(s []string) string {
	return s[rand.Int()%len(s)]
}
func randomStringNumber(n int) string {
	return ""
}

// RandomInt Get three parameters , only first mandatory and the rest are optional
// 		If only set one parameter :  This means the minimum number of digits and the total number
// 		If only set two parameters : First this is min digit and second max digit and the total number the difference between them
// 		If only three parameters: the third argument set Max count Digit
func RandomInt(parameters ...int) (p []int, err error) {
	switch len(parameters) {
	case 1:
		minCount := parameters[0]
		p = rand.Perm(minCount)
		for i := range p {
			p[i] += minCount
		}
	case 2:
		minDigit, maxDigit := parameters[0], parameters[1]
		p = rand.Perm(maxDigit - minDigit + 1)

		for i := range p {
			p[i] += minDigit
		}
	default:
		err = fmt.Errorf(ErrMoreArguments, len(parameters))
	}
	return p, err
}

func generateUnique(dataType string, fn func() interface{}) (interface{}, error) {
	for i := 0; i < maxRetry; i++ {
		value := fn()
		if !slice.ContainsValue(uniqueValues[dataType], value) { // Retry if unique value already found
			uniqueValues[dataType] = append(uniqueValues[dataType], value)
			return value, nil
		}
	}
	return reflect.Value{}, fmt.Errorf(ErrUniqueFailure, dataType)
}

func extractFloat64FromTagArgs(args []string) (interface{}, error) {
	bytes := 64
	var floatValues []float64
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseFloat(k, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, j)
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractFloat32FromTagArgs(args []string) (interface{}, error) {
	bytes := 32
	var floatValues []float32
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseFloat(k, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, float32(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt64FromTagArgs(args []string) (interface{}, error) {
	bytes := 64
	var floatValues []int64
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, j)
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt32FromTagArgs(args []string) (interface{}, error) {
	bytes := 32
	var floatValues []int32
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int32(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt16FromTagArgs(args []string) (interface{}, error) {
	bytes := 16
	var floatValues []int16
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int16(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractInt8FromTagArgs(args []string) (interface{}, error) {
	bytes := 8
	var floatValues []int8
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int8(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractIntFromTagArgs(args []string) (interface{}, error) {
	bytes := 0
	var floatValues []int
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseInt(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, int(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint64FromTagArgs(args []string) (interface{}, error) {
	bytes := 64
	var floatValues []uint64
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, j)
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint32FromTagArgs(args []string) (interface{}, error) {
	bytes := 32
	var floatValues []uint32
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint32(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint16FromTagArgs(args []string) (interface{}, error) {
	bytes := 16
	var floatValues []uint16
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint16(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUint8FromTagArgs(args []string) (interface{}, error) {
	bytes := 8
	var floatValues []uint8
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint8(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}

func extractUintFromTagArgs(args []string) (interface{}, error) {
	bytes := 0
	var floatValues []uint
	for _, i := range args {
		k := strings.TrimSpace(i)
		j, err := strconv.ParseUint(k, 10, bytes)
		if err != nil {
			return nil, fmt.Errorf(ErrUnsupportedTagArguments)
		}
		floatValues = append(floatValues, uint(j))
	}
	toRet := floatValues[rand.Intn(len(floatValues))]
	return toRet, nil
}
func IsExportByName(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(ch)
}
