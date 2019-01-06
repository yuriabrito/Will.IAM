package api

import (
	"encoding/json"
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

func keepJSONFieldsSlBytes(
	isl interface{}, keep ...string,
) ([]byte, error) {
	m, err := keepJSONFieldsSl(isl, keep...)
	if err != nil {
		return nil, err
	}
	return json.Marshal(m)
}
