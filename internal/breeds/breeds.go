package breeds

import (
	"encoding/json"
	"log"
	"net/http"
)

type Breed struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

func FetchBreeds() ([]Breed, error) {
    url := "https://api.thecatapi.com/v1/breeds"
    req, _ := http.NewRequest("GET", url, nil)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error fetching breeds: %v", err)
        return nil, err
    }
    defer resp.Body.Close()

    var breeds []Breed
    if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
        log.Printf("Error decoding breed response: %v", err)
        return nil, err
    }

    return breeds, nil
}
