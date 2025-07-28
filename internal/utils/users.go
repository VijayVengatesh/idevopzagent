package utils

import (
	"github.com/shirou/gopsutil/v3/host"
)

// LoggedInUser holds user session information
type LoggedInUser struct {
	User     string
	Terminal string
	Host     string
	Started  int64 // Timestamp in seconds
}

// GetLoggedInUsers returns a list of currently logged-in users
func GetLoggedInUsers() ([]LoggedInUser, error) {
	users, err := host.Users()
	if err != nil {
		return nil, err
	}

	var result []LoggedInUser
	for _, u := range users {
		result = append(result, LoggedInUser{
			User:     u.User,
			Terminal: u.Terminal,
			Host:     u.Host,
			Started:  int64(u.Started),
		})
	}

	return result, nil
}
