package models

// struct for rating object
type Rating struct {
	Rater  string  `json:"rater"`
	Rating float32 `json:"rating"`
}
