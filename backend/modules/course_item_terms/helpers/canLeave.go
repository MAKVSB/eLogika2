package helpers

import (
	"errors"
	"time"

	"elogika.vsb.cz/backend/models"
)

func CanLeave(isJoined bool, term models.Term) error {
	if !isJoined {
		return errors.New("user not signed in")
	}

	if time.Now().Truncate(time.Minute).Before(term.SignOutFrom.Truncate(time.Minute)) {
		return errors.New("cannot sign out before allowed time")
	}

	if time.Now().Truncate(time.Minute).After(term.SignOutTo.Truncate(time.Minute)) {
		return errors.New("cannot sign out after allowed time")
	}

	return nil
}
