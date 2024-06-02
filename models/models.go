package models

import "encoding/json"

type AssetType string

const (
    Chart    AssetType = "chart"
    Insight  AssetType = "insight"
    Audience AssetType = "audience"
)

type Asset struct {
    ID          string    `json:"id"`
    Type        AssetType `json:"type"`
    Description string    `json:"description"`
    Data        string    `json:"data"`
}

type UserFavorites struct {
    UserID string `json:"userId"`
    Assets []Asset `json:"assets"`
}

type User struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"-"`
    CreatedAt string `json:"created_at"`
}

type Credentials struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// MarshalJSON customizes JSON encoding to exclude the Password field.
func (u User) MarshalJSON() ([]byte, error) {
    type Alias User
    return json.Marshal(&struct {
        Password string `json:"password,omitempty"`
        *Alias
    }{
        Password: "",
        Alias:    (*Alias)(&u),
    })
}

// UnmarshalJSON customizes JSON decoding to include the Password field.
func (u *User) UnmarshalJSON(data []byte) error {
    type Alias User
    aux := &struct {
        Password string `json:"password"`
        *Alias
    }{
        Alias: (*Alias)(u),
    }
    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }
    u.Password = aux.Password
    return nil
}

