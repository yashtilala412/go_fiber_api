package models

// App represents an app in the Google Play Store.
type App struct {
	Name          string `json:"name"`
	Category      string `json:"category"`
	Rating        string `json:"rating"`
	Reviews       string `json:"reviews"`
	Size          string `json:"size"`
	Installs      string `json:"installs"`
	Type          string `json:"type"`
	Price         string `json:"price"`
	ContentRating string `json:"content_rating"`
	Genres        string `json:"genres"`
	LastUpdated   string `json:"last_updated"`
	CurrentVer    string `json:"current_ver"`
	AndroidVer    string `json:"android_ver"`
}
