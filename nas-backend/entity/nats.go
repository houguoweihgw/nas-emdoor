package entity

type SCRequest struct {
	ID      int
	UserID  int
	Picture []byte
}

type SCResponse struct {
	ID     int
	UserID int
	TAGS   []string
}

type FCResponse struct {
	ID     int        `json:"id"`
	Flag   bool       `json:"flag"`
	UserID int        `json:"user_id"`
	Faces  []FaceInfo `json:"faces"`
}

type FaceInfo struct {
	X          int       `json:"x"`
	Y          int       `json:"y"`
	Width      int       `json:"w"`
	Height     int       `json:"h"`
	Embeddings []float64 `json:"embeddings"`
}
