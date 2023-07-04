package rangerversion

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

/* Requests the ranger version from the ranger server, parses the returned
 * json and returns a Struct object, or an error.
 */
func RequestRangerVersion(ranger_url string) (*Config, error) {
	var http_url = "https://" + ranger_url + "/rangerversion"
	req, err := http.Get(http_url)
	if err != nil {
		return nil, err
	}
	byte_object, err := ioutil.ReadAll(req.Body)
	if err != nil || byte_object == nil {
		return nil, err
	}
	var jsonObject map[string]interface{}
	err = json.Unmarshal(byte_object, &jsonObject)
	if err != nil {
		return nil, err
	}
	config_object := new(Config)
	config_object.IsPrime, _ = strconv.ParseBool(jsonObject["RangerPrime"].(string))
	config_object.RangerVersion = jsonObject["Version"].(string)
	config_object.GitCommit = jsonObject["GitCommit"].(string)
	return config_object, nil
}
