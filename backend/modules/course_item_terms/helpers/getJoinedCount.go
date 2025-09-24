package helpers

import (
	"math"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"gorm.io/gorm"
)

func GetJoinedLocking(transaction *gorm.DB, termID uint, perClassLimit bool, studentClasses []uint, locking bool) (int, *common.ErrorResponse) {
	if locking {
		transaction = transaction.
			Table(models.UserTerm{}.TableName() + " WITH (XLOCK, TABLOCK)")
	} else {
		transaction = transaction.
			Table(models.UserTerm{}.TableName())
	}

	if perClassLimit {
		transaction = transaction.
			Where("user_terms.deleted_at is NULL").
			Where("user_terms.term_id = ?", termID).
			Select("class_id, count(class_students.id) as \"cnt\"").
			Joins("JOIN class_students on class_students.user_id = user_terms.user_id AND class_students.deleted_at is NULL AND class_students.class_id in ?", studentClasses).
			Group("class_id").
			Order("\"cnt\" ASC")
	} else {
		transaction = transaction.
			Where("deleted_at is NULL").
			Where("term_id = ?", termID).
			Select("count(id) as \"cnt\"")
	}

	var joinedRecords struct {
		Cnt int
	}
	if err := transaction.
		Find(&joinedRecords).Error; err != nil {
		return math.MaxInt, &common.ErrorResponse{
			Code:    500,
			Message: "Failed to get list of joined students",
		}
	}

	return joinedRecords.Cnt, nil
}
