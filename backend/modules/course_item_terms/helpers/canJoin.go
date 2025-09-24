package helpers

import (
	"errors"
	"time"

	"elogika.vsb.cz/backend/models"
)

func CanJoin(isJoined bool, term models.Term, alreadyJoinedCount int) error {
	if isJoined {
		return errors.New("user already signed in")
	}

	if time.Now().Truncate(time.Minute).Before(term.SignInFrom.Truncate(time.Minute)) {
		return errors.New("cannot sign in before allowed time")
	}

	if time.Now().Truncate(time.Minute).After(term.SignInTo.Truncate(time.Minute)) {
		return errors.New("cannot sign in after allowed time")
	}

	if alreadyJoinedCount >= int(term.StudentsMax) {
		return errors.New("capacity already full")
	}

	// TODO zkontrolovat, že už nevyčerpal možné pokusy
	// TODO zkontrolovat, že nepřekročil choice max
	// TODO mám validovat že není ještě na termínu který ještě neproběhl ?

	return nil
}
