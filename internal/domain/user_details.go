package domain

import (
	"errors"
	"strconv"
	"time"

	"github.com/omniful/go_commons/log"
)

type UserDetails struct {
	UserName        string
	UserID          string
	UserEmail       string
	TenantName      string
	TenantID        string
	CreateByUserID  uint64
	UpdatedAt       time.Time
	UpdatedByUserID uint64
}

func (u *UserDetails) SetCreated() (userDetails *UserDetails) {
	userID, err := u.getUserIDUint64()
	if err != nil {
		return
	}

	u.CreateByUserID = userID
	return u
}

func (u *UserDetails) SetUpdated() (userDetails *UserDetails) {
	userID, err := u.getUserIDUint64()
	if err != nil {
		return
	}

	u.UpdatedByUserID = userID
	u.UpdatedAt = time.Now()
	return u
}

func (u *UserDetails) getUserIDUint64() (uint64, error) {
	if u == nil {
		return 0, errors.New("user details is empty")
	}

	if len(u.UserID) == 0 {
		log.Errorf("user id is empty")
		return 0, errors.New("user id is empty")
	}

	userID, err := strconv.ParseUint(u.UserID, 10, 64)
	if err != nil {
		log.Errorf("unable to parse user id :: [%s]", u.UserID)
		return 0, err
	}

	return userID, nil
}
