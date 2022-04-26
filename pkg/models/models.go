// Members: Miguel Avila, Federico Rosado

package models

import (
	"errors"
	"time"
)

var ErrRecordNotFound = errors.New("BlogsL no matching record found")

// A struct that hold Blog
type Blog struct {
	ID           int
	FirstName    string
	LastName     string
	Email        string
	Subject      string
	Message      string
	Date_Created time.Time
}
