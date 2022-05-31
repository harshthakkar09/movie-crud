package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"movie-crud/src/models"
)

func checkStringInSlice(items []string, item string) bool {
	for _, cur := range items {
		if cur == item {
			return true
		}
	}
	return false
}

func ValidateMovieObject(movie *models.Movie) error {
	// checking whether movie title is present
	if movie.Title == "" {
		return fmt.Errorf("movie title is required")
	}

	// setting default value of genre
	if movie.Genre == "" {
		movie.Genre = "thriller"
	}

	// checking whether genre is valid
	flag := checkStringInSlice(validGenre, movie.Genre)
	if !flag {
		return fmt.Errorf("genre value %s is not allowed", movie.Genre)
	}

	// checking whether all ratings have rater
	ratings := movie.Ratings
	for _, rating := range ratings {
		if rating.Rater == "" {
			return fmt.Errorf("rater is required")
		}
	}

	return nil
}

func castExists(CastIDs []string) ([]string, bool, error) {
	//mutex
	castMutex.Lock()
	defer castMutex.Unlock()

	// start reading json file
	plan, _ := ioutil.ReadFile("./src/data/casts.json")
	var castsMap map[string]interface{}
	err := json.Unmarshal(plan, &castsMap)

	if err != nil {
		return nil, false, err
	}

	// iterate through list of casts
	notPresentIDs := []string{}
	for _, castId := range CastIDs { // finding cast with castId
		_, ok := castsMap[castId]
		if !ok {
			notPresentIDs = append(notPresentIDs, castId)
		}
	}
	if len(notPresentIDs) == 0 {
		return nil, true, nil
	}
	return notPresentIDs, false, nil
}
