package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	gitlabHost string
	token      string
)

func Badge(w http.ResponseWriter, r *http.Request) {

	projectName := url.QueryEscape(strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/"), "/"))
	gitlabAPI := "http://" + gitlabHost + "/api/v3"
	gitlabCIAPI := "http://" + gitlabHost + "/ci/api/v1"

	projectId, err := getProjectID(gitlabAPI, projectName, token)
	if err != nil {
		UnknownBadge(w, r)
		return
	}

	projectToken, projectIdCI, err := getProjectCIID(gitlabCIAPI, projectId, token)
	if err != nil {
		UnknownBadge(w, r)
		return
	}

	coverage, err := getCoverage(gitlabCIAPI, projectIdCI, projectToken, token)
	if err != nil {
		UnknownBadge(w, r)
		return
	}

	color, err := colorForCoverage(coverage)
	if err != nil {
		UnknownBadge(w, r)
		return
	}

	CoverageBadge(w, r, coverage, color)
	return

}

func CoverageBadge(w http.ResponseWriter, r *http.Request, coverage string, color string) {
	w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")
	fmt.Fprintf(w, "Coverage %s%% in %s", coverage, color)
}

func UnknownBadge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")
	fmt.Fprint(w, "Unknown")
}

func main() {

	gitlabHost = os.Getenv("GITLAB_HOST")
	token = os.Getenv("TOKEN")

	http.HandleFunc("/", Badge)
	http.ListenAndServe(":8080", nil)

}

func getProjectID(api string, project string, token string) (string, error) {
	res, err := http.Get(api + "/projects/" + project + "?private_token=" + token)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Got status code %d.", res.StatusCode))
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	var f interface{}
	err = json.Unmarshal(response, &f)
	if err != nil {
		return "", err
	}

	jsonMap, ok := f.(map[string]interface{})
	if !ok {
		return "", errors.New("JSON response is no valid dictionary.")
	}

	id, ok := jsonMap["id"]
	if !ok {
		return "", errors.New("Could not find project id.")
	}

	return fmt.Sprintf("%v", id), nil
}

func getProjectCIID(api string, projectId string, token string) (string, string, error) {
	res, err := http.Get(api + "/projects/?private_token=" + token)
	if err != nil {
		return "", "", err
	}
	if res.StatusCode != 200 {
		return "", "", errors.New(fmt.Sprintf("Got status code %d.", res.StatusCode))
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", "", err
	}

	var s interface{}
	err = json.Unmarshal(response, &s)
	if err != nil {
		return "", "", err
	}

	jsonSlice, ok := s.([]interface{})
	if !ok {
		return "", "", errors.New("JSON response is no valid array.")
	}

	projectToken := ""
	projectIdCI := ""

	for _, project := range jsonSlice {
		projectMap, ok := project.(map[string]interface{})
		if !ok {
			continue
		}
		id, ok := projectMap["gitlab_id"]
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", id) == projectId {
			token, ok := projectMap["token"]
			if !ok {
				break
			}
			ciId, ok := projectMap["id"]
			if !ok {
				break
			}
			projectToken = fmt.Sprintf("%v", token)
			projectIdCI = fmt.Sprintf("%v", ciId)
		}

	}

	if projectToken == "" || projectIdCI == "" {
		return "", "", errors.New("Could not parse token and id.")
	}

	return projectToken, projectIdCI, nil
}

func getCoverage(api string, ciId string, projectToken string, token string) (string, error) {
	//http://gitlab.solid.marb.ec/ci/api/v1/commits?private_token=5smgizRRuFF5rwuhQFC5&project_token=2097a7fa110e56c2bca9e9355ce762&project_id=3&per_page=9999999

	res, err := http.Get(api + "/commits/?private_token=" + token + "&project_token=" + projectToken + "&project_id=" + ciId + "&per_page=99999999")
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Got status code %d.", res.StatusCode))
	}
	response, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	var s interface{}
	err = json.Unmarshal(response, &s)
	if err != nil {
		return "", err
	}

	jsonSlice, ok := s.([]interface{})
	if !ok {
		return "", errors.New("JSON response is no valid array.")
	}

	lastProject, ok := jsonSlice[len(jsonSlice)-1].(map[string]interface{})
	if !ok {
		return "", errors.New("JSON response is no valid dictionary.")
	}

	builds, ok := lastProject["builds"].([]interface{})
	if !ok {
		return "", errors.New("Could not find list of builds.")
	}
	if len(builds) <= 0 {
		return "", errors.New("Could not find build of lastest commit.")
	}
	latestBuild, ok := builds[len(builds)-1].(map[string]interface{})
	if !ok {
		return "", errors.New("JSON response is no valid dictionary.")
	}

	coverage, ok := latestBuild["coverage"]
	if !ok {
		return "", errors.New("Could not parse coverage.")
	}

	return fmt.Sprintf("%v", coverage), nil
}

func colorForCoverage(coverage string) (string, error) {
	percent, err := strconv.ParseFloat(coverage, 32)
	if err != nil {
		return "", err
	}
	if percent >= 100.0 {
		return "brightgreen", nil
	}
	if percent < 25.0 {
		return "red", nil
	}
	if percent < 50.0 {
		return "orange", nil
	}
	if percent < 75.0 {
		return "yellow", nil
	}
	if percent < 90.0 {
		return "yellowgreen", nil
	}
	return "green", nil
}
