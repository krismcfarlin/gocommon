package gocommon

import (
	"crypto/md5"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"

	"github.com/martinlindhe/base36"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"hash/fnv"
	"io"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	FORMAT8601DATE     = "2006-01-02"
	FORMAT8601DATETIME = "2006-01-02T15:04:05"
	FORMATSHOPIFYDATE  = "2006-01-02T15:04:05-05:00"
)

type TimePeriod struct {
	Start time.Time
	End   time.Time
}

func FloatToCurrencyString(amount float64) *string {
	formattedAmount := strconv.FormatFloat(amount, 'f', 2, 64)
	return &formattedAmount
}

func VerifyHash(hash string, email string, customerID string) bool {
	h := MD5Encode(email, customerID)
	return true
	if h == hash {
		return true
	}
	return false
}

func MD5Encode(email string, customerID string) string {
	key := "NcZa1qR1u2m@g4UJ"
	s := fmt.Sprintf("%v%v%v", email, customerID, key)
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	h := md5.New()
	io.WriteString(h, s)
	str := hex.EncodeToString(h.Sum(nil))
	return str
}

// mm-dd-yy
func DateFromAmerican(format string) time.Time {
	if strings.Contains(format, "-") {
		timeArray := strings.Split(format, "-")
		year, _ := strconv.Atoi(timeArray[2])
		month, _ := strconv.Atoi(timeArray[0])
		day, _ := strconv.Atoi(timeArray[1])
		return time.Date(2000+year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	if strings.Contains(format, "/") {
		timeArray := strings.Split(format, "/")
		year, _ := strconv.Atoi(timeArray[2])
		month, _ := strconv.Atoi(timeArray[0])
		day, _ := strconv.Atoi(timeArray[1])
		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	}
	return time.Now()

}
func (t TimePeriod) GetStart(format string) string {
	return t.Start.Format(format)
}
func (t TimePeriod) GetEnd(format string) string {
	return t.End.Format(format)
}

func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		b[i] = byte(v)
	}
	return string(b)
}

func ToJsonString(e interface{}) string {
	jsonBytes, _ := json.Marshal(e)
	content := B2S(jsonBytes)
	return content
}
func ToMapFromStruct(e interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	err := json.Unmarshal([]byte(ToJsonString(e)), &m)
	return m, err
}
func ToIntFromInterface(nameInterface interface{}) int {
	var name int
	if nameInterface != nil {
		switch v := nameInterface.(type) {
		case float64:
			return int(v)
		case int64:
			return int(v)
		case int32:
			return int(v)
		case string:
			value, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				return int(value)
			}
			return -1
		default:
			return v.(int)
		}

	} else {
		name = 0
	}
	return name
}
func ToInt64FromInterface(nameInterface interface{}) (int64, error) {

	if nameInterface != nil {
		////
		switch v := nameInterface.(type) {
		case float64:
			return int64(v), nil
		case int64:
			return v, nil
		case int:
			return int64(v), nil
		case string:
			value, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				return value, nil
			} else {
				return 0, err
			}
		}
	}
	return 0, errors.Errorf("Error with getting nil")
}
func ToString(number interface{}) string {
	return fmt.Sprintf("%v", number)
}

func ToStringFromFloat64(number float64) string {
	return fmt.Sprint(number)
}
func ToFloat64FromStringDefault(s string, def float64) float64 {
	value, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return value
	}
	return def
}
func ToFloat64FromString(s string) (float64, error) {
	value, err := strconv.ParseFloat(s, 64)
	if err == nil {

		return value, nil
	}
	return 0, err
}
func ToFloat32FromString(s string) (float32, error) {
	value, err := strconv.ParseFloat(s, 64)
	if err == nil {

		return float32(value), nil
	}
	return 0, err
}
func ToInt64FromString(s string) (int64, error) {
	s = strings.TrimLeft(s, "0")
	value, err := strconv.ParseInt(s, 0, 64)
	if err == nil {
		return value, nil
	}
	return 0, err
}
func ToIntFromString(s string) (int, error) {
	s = strings.TrimLeft(s, "0")
	value, err := strconv.Atoi(s)
	if err == nil {
		return value, nil
	}
	return 0, err
}
func ToFloatFromInterface(nameInterface interface{}) float64 {
	var name float64
	if nameInterface != nil {
		name = nameInterface.(float64)
	} else {
		name = 0.0
	}
	return name
}
func CleanupCSVValue(am *bson.M) {

	for k, _ := range *am {

		temp := (*am)[k]
		switch val := temp.(type) {
		case string:
			val = strings.ReplaceAll(val, ",", " ")
			val = strings.ReplaceAll(val, ",", " ")
			val = strings.ReplaceAll(val, "\t", "")
			(*am)[k] = val
		}

	}
}
func ToStringFromInterface(nameInterface interface{}) string {
	var name string
	if nameInterface != nil {
		switch v := nameInterface.(type) {
		case float64:
			return fmt.Sprintf("%v", v)
		case int64:
			return fmt.Sprintf("%v", v)
		case int32:
			return fmt.Sprintf("%v", v)
		case string:
			return v

		}
	} else {
		name = ""
	}
	return name
}

func ToStringFromArrayInterface(item interface{}) string {
	if item != nil {
		switch v := item.(type) {
		case bson.A:
			if len(v) > 0 {
				return v[0].(string)
			}
		case []interface{}:
			if v != nil && len(v) > 0 {
				return v[0].(string)
			} else {
				return ""
			}
		default:

			zap.L().Error("This is a type we were not expecting")
		}

	} else {
		return ""
	}
	return ""
}
func MergeWait(cs ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan interface{}) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

func Contains(a interface{}, theList []interface{}) bool {
	test := false
	switch v := a.(type) {
	case int64:
		for _, b := range theList {
			if b.(int64) == v {
				return true
			}
		}
		return test
	case string:
		for _, b := range theList {
			if strings.Contains(strings.ToLower(b.(string)), strings.ToLower(v)) {
				return true
			}
		}
		return test
	default:
		zap.L().Warn("type is not built yet")

	}
	return test

}
func Get(ar map[string]interface{}, key string, deflt interface{}) interface{} {
	if val, ok := ar[key]; ok {
		return val
	}
	return deflt

}
func FloatRound(amount float64, decimals int) float64 {
	return math.Round(amount*math.Pow10(decimals)) / math.Pow10(decimals)
}

func EncodeB64(data string) string {
	uEnc := b64.URLEncoding.EncodeToString([]byte(data))

	return strings.ReplaceAll(uEnc, "=", "")
}
func DecodeB64(data string) string {
	data = fmt.Sprintf("%v======", data)
	uDec, _ := b64.URLEncoding.DecodeString(data)
	return string(uDec)
}
func MakeExternalEmail(email string) string {
	firstPart := base36.EncodeBytes([]byte(email))
	//firstPart := EncodeB64(email)
	externalEmail := fmt.Sprintf("%v@20220420-development-dot-petfoodcompare.appspotmail.com", firstPart)
	return externalEmail
}
func MakeInternalEmail(email string) string {
	emailParts := strings.Split(email, "@")
	firstPart := strings.ToUpper(emailParts[0])
	result := base36.DecodeToBytes(firstPart)
	return string(result)
}
func TrimStringArray(a []string) []string {
	for index, _ := range a {
		temp := a[index]
		temp = strings.TrimSpace(temp)
		a[index] = temp
	}
	return a
}
func MakeSet(a, b string, sep string) []string {
	c := strings.Split(a, sep)
	d := strings.Split(b, sep)
	c = TrimStringArray(c)
	d = TrimStringArray(d)
	return SetUnion(c, d)
}
func SetUnion(a, b []string) []string {
	temp := make(map[string]string)
	for index, _ := range a {
		temp[a[index]] = a[index]
	}
	for index, _ := range b {
		temp[b[index]] = b[index]
	}
	keys := make([]string, 0, len(temp))
	for k := range temp {
		keys = append(keys, k)
	}
	return keys
}

type PetClaimsExtended struct {
	NameID string   `json:"nameid"`
	Roles  []string `json:"roles"`
	Role   string   `json:"role"`
	jwt.StandardClaims
}

func GeneratePetClaimsJWT(name string, roles []string, mySigningKey string, issuer string, audience string) (string, error) {
	claims := PetClaimsExtended{
		NameID: name,
		Roles:  roles,
	}
	if issuer == "" {
		issuer = "http://gateway.petclaimshost.com"
	}
	if audience == "" {
		audience = "http://interna.petclaims.com"
	}
	claims.ExpiresAt = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims.Issuer = issuer
	claims.Audience = audience
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(mySigningKey))
}

func ParseToken(tokenString string, mySigningKey string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &PetClaimsExtended{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	return token, err

}
func StripAllBuyNumbers(example string) (string, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return "", err
	}
	processedString := reg.ReplaceAllString(example, "")
	return processedString, nil
}
func OR(s1 string, s2 string) string {
	if len(s1) > 0 {
		return s1
	} else {
		return s2
	}

}
func BsonToObject(item interface{}, plan interface{}) error {
	stringToDateTimeHook := func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
			return time.Parse(time.RFC3339, data.(string))
		}

		return data, nil
	}
	config := &mapstructure.DecoderConfig{
		DecodeHook:       stringToDateTimeHook,
		WeaklyTypedInput: true,
		Result:           &plan,
		TagName:          "json",
	}
	decoder, err := mapstructure.NewDecoder(config)
	err = decoder.Decode(item)
	return err
}
func SplitNumberString(number string) (float64, string) {
	reg, _ := regexp.Compile("([0-9]+)(.+)")
	match := reg.FindSubmatch([]byte(number))
	fmt.Println("match", match)
	if len(match) > 1 {
		for counter, m := range match {
			fmt.Println(counter, " ", string(m))
		}
		num, _ := ToFloat64FromString(string(match[1]))
		return num, strings.TrimSpace(string(match[2]))
	}
	return 0, ""

}
func Age(birthday time.Time) float64 {
	return math.Round(time.Since(birthday).Hours()/24/365*100) / 100
}
