package api

import (
	"encoding/json"
	"reflect"
)

func keepJSONFieldsSl(
	isl interface{}, keep ...string,
) ([]map[string]interface{}, error) {
	bts, err := json.Marshal(isl)
	if err != nil {
		return nil, err
	}
	var mSl []map[string]interface{}
	if err := json.Unmarshal(bts, &mSl); err != nil {
		return nil, err
	}
	kmSl := make([]map[string]interface{}, len(mSl))
	for i := range mSl {
		kmSl[i] = map[string]interface{}{}
		for _, f := range keep {
			kmSl[i][f] = mSl[i][f]
		}
	}
	return kmSl, nil
}

func keepJSONFieldsOne(
	i interface{}, keep ...string,
) (map[string]interface{}, error) {
	bts, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(bts, &m); err != nil {
		return nil, err
	}
	km := map[string]interface{}{}
	for _, f := range keep {
		km[f] = m[f]
	}
	return km, nil
}

func keepJSONFields(i interface{}, keep ...string) (interface{}, error) {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		return keepJSONFieldsSl(i, keep...)
	default:
		return keepJSONFieldsOne(i, keep...)
	}
}

func keepJSONFieldsBytes(i interface{}, keep ...string) ([]byte, error) {
	ri, err := keepJSONFields(i, keep...)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ri)
}
