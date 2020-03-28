package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"text/template"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

type Country struct {
	Name                string  `json:"country"`
	Cases               int     `json:"cases"`
	TodayCases          int     `json:"todayCases"`
	Deaths              int     `json:"deaths"`
	TodayDeaths         int     `json:"todayDeaths"`
	Recovere            int     `json:"recovered"`
	Active              int     `json:"active"`
	Critical            int     `json:"critical"`
	CasesPerOneMission  float32 `json:"casesPerOneMillion"`
	DeathsPerOneMillion float32 `json:"deathsPerOneMillion"`
}

func (country *Country) ToJSON() string {
	jsonCountry, err := json.Marshal(country)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(jsonCountry).String()
}

func (country *Country) String() string {
	const templateText = `
	Name               : {{.Name}}
	Cases              : {{.Cases}}
	TodayCases         : {{.TodayCases}}
	Deaths             : {{.Deaths}}
	TodayDeaths        : {{.TodayDeaths}}
	Recovere           : {{.Recovere}}
	Active             : {{.Active}}
	Critical           : {{.Critical}}
	CasesPerOneMission : {{.CasesPerOneMission}}
	DeathsPerOneMillion: {{.DeathsPerOneMillion}}
	`
	template, err := template.New("Country").Parse(templateText)
	if err != nil {
		panic(err)
	}
	var doc bytes.Buffer
	if err := template.Execute(&doc, country); err != nil {
		panic(err)
	}
	return doc.String()
}

func NewClient() (*Client, error) {
	baseURL, err := url.Parse("https://corona.lmao.ninja/")
	if err != nil {
		return nil, err
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}, nil
}

func (client *Client) GetCountries() ([]Country, error) {
	url := *client.BaseURL
	url.Path = path.Join(url.Path, "countries")
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/+json")
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var countries []Country
	if err = json.Unmarshal(body, &countries); err != nil {
		panic(err)
	}
	return countries, nil
}
