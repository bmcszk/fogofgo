package world

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type WorldService struct {
	serverAddress string
	client        *http.Client
}

type WorldRequest struct {
	X, Y          int
	Width, Height int
}

type WorldResponse struct {
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Rows   [][]Tile `json:"rows"`
}

type Tile struct {
	Point           Point  `json:"point"`
	Value           string `json:"value"`
	LandType        string `json:"landType"`
	FrontStyleClass string `json:"frontStyleClass"`
	BackStyleClass  string `json:"backStyleClass"`
	GroundLevel     int    `json:"groundLevel"`
	WaterLevel      *int   `json:"waterLevel"`
	PostGlacial     bool   `json:"postGlacial"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewWorldService() WorldService {
	return WorldService{
		serverAddress: "http://localhost:8080",
		client:        http.DefaultClient,
	}
}

func (m WorldService) DeadSimpleLoad() *WorldResponse {
	resp, err := m.client.Get(m.serverAddress + "/api/map?x=1&y=2&width=100&height=100")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response *WorldResponse
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		log.Fatal(err)
	}
	return response

}

func (m WorldService) Load(request WorldRequest) (*WorldResponse, error) {
	u, err := url.Parse(m.serverAddress + "/api/map")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("x", fmt.Sprintf("%d", request.X))
	q.Add("y", fmt.Sprintf("%d", request.Y))
	q.Add("width", fmt.Sprintf("%d", request.Width))
	q.Add("height", fmt.Sprintf("%d", request.Height))
	u.RawQuery = q.Encode()
	resp, err := m.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *WorldResponse
	err = json.Unmarshal(responseData, &response)
	if err != nil {
		return nil, err
	}
	return response, nil

}
