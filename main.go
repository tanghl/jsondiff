package main

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"strings"
	"io/ioutil"
	"strconv"
	"time"
	"regexp"
)

func main() {
	start := time.Now().Nanosecond()
	key1 := make(map[string]interface{})
	b1, _ := ioutil.ReadFile("p.json")
	json1, _ := simplejson.NewJson(b1)
	js1, _ := json1.MarshalJSON()
	fmt.Println(string(js1))
	ParseJsonKeys(&key1, "", json1)

	key1C := make(map[string]interface{})
	for k, v := range all {
		key1C[k] = v
		delete(all, k)
	}

	key2 := make(map[string]interface{})
	b2, _ := ioutil.ReadFile("p1.json")
	json2, _ := simplejson.NewJson(b2)

	ParseJsonKeys(&key2, "", json2)

	key2C := make(map[string]interface{})
	for k, v := range all {
		key2C[k] = v
		delete(all, k)
	}

	diff2, _ := CompareJsonP(key1C, key2C, json2)
	FillData(diff2, json2)
	FillData(key1C, json2)
	js2, _ := json2.MarshalJSON()
	fmt.Println(string(js2))

	end := time.Now().Nanosecond()
	fmt.Println((end - start))
}

var all = make(map[string]interface{})

func ParseJsonKeys(keys *map[string]interface{}, pk string, json *simplejson.Json) (nk map[string]interface{}, err error) {
	nk = make(map[string]interface{})
	if strings.EqualFold(pk, "") {
		if node_map, err := json.Map(); err == nil {
			for k, v := range node_map {
				nk[k] = v
				ret, _ := ParseJsonKeys(&nk, k, json.Get(k))
				for kk, vv := range ret {
					all[kk] = vv
				}
			}
		}
	} else {
		if m, ok := json.Map(); ok == nil {
			for k, v := range m {
				delete(*keys, pk)
				nk[pk + "." + k] = v
			}
			for k, _ := range nk {
				var lastKey string
				if strings.Contains(k, ".") {
					ks := strings.Split(k, ".")
					lastKey = ks[len(ks) - 1]
				}
				ret, _ := ParseJsonKeys(&nk, k, json.Get(lastKey))
				for kk, vv := range ret {
					all[kk] = vv
				}
			}
		} else if a, ok := json.Array(); ok == nil {
			for k, v := range a {
				delete(*keys, pk)
				nk[pk + "." + strconv.Itoa(k)] = v
			}
			for k, _ := range nk {
				var lastKey int
				if strings.Contains(k, ".") {
					ks := strings.Split(k, ".")
					lastKey, _ = strconv.Atoi(ks[len(ks) - 1])
				}
				ret, _ := ParseJsonKeys(&nk, k, json.GetIndex(lastKey))
				for kk, vv := range ret {
					all[kk] = vv
				}
			}
		}
	}
	return nk, nil
}

func CompareJsonP(a, b map[string]interface{}, json *simplejson.Json) (diffb map[string]interface{}, err error) {
	diffa := make(map[string]interface{})
	diffb = make(map[string]interface{})
	for k, v := range a {
		if b[k] == nil {
			delete(a, k)
			diffa[k] = v
		}
	}
	for k, _ := range b {
		nks := strings.Split(k, ".")
		nk := nks[len(nks) - 1]
		for kk, vv := range diffa {
			nkks := strings.Split(kk, ".")
			nkk := nkks[len(nkks) - 1]
			if strings.EqualFold(nk, nkk) {
				diffb[k] = vv
			}
		}
	}
	return diffb, nil
}

func FillData(a map[string]interface{}, json *simplejson.Json) (error) {
	for k, v := range a {
		ks := strings.Split(k, ".")
		var temp = json
		for in, item := range ks {
			if b, err := regexp.Match("[0-9]", []byte(item)); err == nil && b == true {
				index, _ := strconv.Atoi(item)
				temp = temp.GetIndex(index)
			} else if in < len(ks) - 1 {
				temp = temp.Get(item)
			} else if in == len(ks) - 1 {
				temp.Set(item, v)
			}
		}
	}
	return nil
}