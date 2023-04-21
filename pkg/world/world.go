package world

import (
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
)

type WorldService struct {
	serverAddress string
	client        *http.Client
}

type WorldRequest struct {
	MinX, MinY, MaxX, MaxY int
}

type WorldResponse struct {
	Tiles []Tile `json:"map"`
	MinX  int    `json:"minX"`
	MinY  int    `json:"minY"`
	MaxX  int    `json:"maxX"`
	MaxY  int    `json:"maxY"`
}

type Tile struct {
	Point           image.Point `json:"point"`
	Value           string      `json:"value"`
	LandType        string      `json:"landType"`
	FrontStyleClass string      `json:"frontStyleClass"`
	BackStyleClass  string      `json:"backStyleClass"`
	GroundLevel     int         `json:"groundLevel"`
	WaterLevel      *int        `json:"waterLevel"`
	PostGlacial     bool        `json:"postGlacial"`
}

func NewWorldService() WorldService {
	return WorldService{
		serverAddress: "http://localhost:8080",
		client:        http.DefaultClient,
	}
}

func (m WorldService) Load(request WorldRequest) (*WorldResponse, error) {
	u, err := url.Parse(m.serverAddress + "/api/map/rect")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("minX", fmt.Sprintf("%d", request.MinX))
	q.Add("minY", fmt.Sprintf("%d", request.MinY))
	q.Add("maxX", fmt.Sprintf("%d", request.MaxX))
	q.Add("maxY", fmt.Sprintf("%d", request.MaxY))
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
