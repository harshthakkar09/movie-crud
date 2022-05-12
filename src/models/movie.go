package models

// struct for movie object
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
	CastIDs  []string  `json:"casts"`
	Ratings  []Rating  `json:"ratings"`
}
