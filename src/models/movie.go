package models

// struct for movie object
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn,omitempty"`
	Title    string    `json:"title"`
	Director *Director `json:"director,omitempty"`
	CastIDs  []string  `json:"casts,omitempty"`
	Ratings  []Rating  `json:"ratings,omitempty"`
	Genre    string    `json:"genre"`
}
