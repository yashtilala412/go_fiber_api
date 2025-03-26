package models

type App struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Rating      string `json:"rating"`
	Reviews     string `json:"reviews"`
	Size        string `json:"size"`
	Installs    string `json:"installs"`
	Type        string `json:"type"`
	Price       string `json:"price"`
	ContentRate string `json:"content_rating"`
}
