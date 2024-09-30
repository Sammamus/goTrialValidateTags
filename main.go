package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"
)

type EnvVariables struct {
	token string
	owner string
	repo  string
}

type RepoRelease []struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []any     `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
}

func collect_inputs() <-chan [2]string {
	out := make(chan [2]string)

	go func() {
		for _, e := range os.Environ() {

			pair := strings.SplitN(e, "=", 2)

			matched, err := regexp.Match(`INPUT_.*`, []byte(pair[0]))
			if err != nil {
				log.Fatal(err)
			}

			if matched {
				out_array := [2]string{pair[0], pair[1]}
				out <- out_array
			}
		}
		close(out)
	}()
	return out
}

func get_variables() EnvVariables {
	sort_map := make(map[string]string)

	for entry := range collect_inputs() {
		sort_map[entry[0]] = entry[1]
	}

	if sort_map["INPUT_TOKEN"] == "" {
		log.Fatal("Missing Token")
	}
	if sort_map["INPUT_OWNER"] == "" {
		log.Fatal("Missing Owner")
	}
	if sort_map["INPUT_REPO"] == "" {
		log.Fatal("Missing Repo")
	}

	return EnvVariables{
		sort_map["INPUT_TOKEN"],
		sort_map["INPUT_OWNER"],
		sort_map["INPUT_REPO"],
	}
}

func get_releases(inputs EnvVariables) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", inputs.owner, inputs.repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", inputs.token))
	req.Header.Add("X-Github-Api-Version", "2022-11-28")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var jsonResponse RepoRelease
	//json.Unmarshal([]byte(body), &jsonResponse)
	//fmt.Println(jsonResponse)
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		log.Fatal(err)
	}

	fmt.Println(reflect.TypeOf(jsonResponse))
	for _, i := range jsonResponse {
		fmt.Println(i)
	}
}

func main() {
	get_releases(get_variables())

}
