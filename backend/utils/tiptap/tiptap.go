package tiptap

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"time"

	"elogika.vsb.cz/backend/models"
	"elogika.vsb.cz/backend/modules/common"
	"elogika.vsb.cz/backend/utils"
	"gorm.io/gorm"
)

type TipTapParser struct {
	FileIDs []uint
}

func FindAndSaveRelations(dbRef *gorm.DB, userId uint, node *models.TipTapContent, obj any, prop string) *common.ErrorResponse {
	ttp := TipTapParser{}
	err := ttp.ParseContent(dbRef, userId, node)
	if err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to parse tiptap schema",
			Details: err.Error(),
		}
	}
	err = ttp.SaveRelations(dbRef, obj, prop)
	if err != nil {
		return &common.ErrorResponse{
			Code:    500,
			Message: "Failed to save file relations",
			Details: err.Error(),
		}
	}

	return nil
}

func (ttc *TipTapParser) SaveRelations(dbRef *gorm.DB, obj any, prop string) error {
	// Build the slice type *[]*models.File
	filePtrType := reflect.TypeOf(&models.File{})
	sliceType := reflect.SliceOf(filePtrType)
	slicePtrValue := reflect.New(sliceType)
	filesVal := slicePtrValue.Interface()

	// Load files
	if err := dbRef.
		Where("id IN ?", ttc.FileIDs).
		Find(filesVal).Error; err != nil {
		return fmt.Errorf("failed to load files: %w", err)
	}

	// Reflect on the target object
	v := reflect.ValueOf(obj)

	// Handle both pointer and non-pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("object must be a non-nil pointer")
		}
		v = v.Elem()
	} else {
		// If it's not a pointer, we can’t modify its fields
		return errors.New("object must be a pointer to struct, not struct value")
	}

	// If the inner object is still a pointer. Dereference it (special case for double pointered values)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return errors.New("object must be a non-nil pointer")
		}
		v = v.Elem()
	}

	// Get target field by name
	field := v.FieldByName(prop)
	if !field.IsValid() {
		return fmt.Errorf("property %q not found on %T", prop, obj)
	}
	if !field.CanSet() {
		return fmt.Errorf("cannot set property %q on %T", prop, obj)
	}

	// Assign loaded files to the property
	filesValElem := reflect.ValueOf(filesVal).Elem()
	field.Set(filesValElem)

	// Ensure ID field exists
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		return fmt.Errorf("%T does not have an 'ID' field", obj)
	}

	// If ID is zero, skip updating associations
	if idField.IsZero() {
		return nil
	}

	// Update association via GORM
	if err := dbRef.Model(obj).Association(prop).Replace(filesValElem.Interface()); err != nil {
		return fmt.Errorf("failed to update association %q: %w", prop, err)
	}

	return nil
}

func (ttc *TipTapParser) ParseContent(dbRef *gorm.DB, userId uint, node *models.TipTapContent) error {
	switch node.Type {
	case "custom-image":
		err := ttc.HandleImage(node)
		if err != nil {
			return err
		}

	case "codeBlock":
		err := ttc.HandleCodeBlock(dbRef, userId, node)
		if err != nil {
			return err
		}
	}

	// Recurse into content array if present
	for _, child := range node.Content {
		err := ttc.ParseContent(dbRef, userId, child)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ttc *TipTapParser) HandleImage(node *models.TipTapContent) error {
	if node.Attrs["mode"] == "storage" {
		ttc.FileIDs = append(ttc.FileIDs, uint(node.Attrs["id"].(float64)))
	}

	return nil
}

func (ttc *TipTapParser) HandleCodeBlock(dbRef *gorm.DB, userId uint, node *models.TipTapContent) error {
	if node.Attrs["language"] != "latex" {
		return nil
	}

	if len(node.Content) != 1 {
		return errors.New("code block must contain single text node")
	}

	if node.Content[0].Type != "text" {
		return errors.New("code block must contain single text node")
	}

	code := node.Content[0].Text

	// Compute hash
	h := sha256.Sum256([]byte(code))
	hash := hex.EncodeToString(h[:])

	if hash == node.Attrs["hash"] {

		existingId, ok := node.Attrs["id"].(float64)
		if !ok {
			return errors.New("image id not found")
		}

		ttc.FileIDs = append(ttc.FileIDs, uint(existingId))
		return nil
	}

	// If file hash does not match. Render new image
	fileName, size, err := utils.ConvertCodeToImage(dbRef, code)
	if err != nil {
		return err
	}

	node.Attrs["hash"] = hash
	node.Attrs["filename"] = fileName

	newFile := models.File{
		UserID:       userId,
		OriginalName: "GENERATED",
		StoredName:   fileName,
		MIMEType:     "image/svg",
		SizeBytes:    size,
		UploadedAt:   time.Now(),
	}

	if err := dbRef.Save(&newFile).Error; err != nil {
		return errors.New("failed to save file info to database")
	}

	node.Attrs["id"] = newFile.ID

	ttc.FileIDs = append(ttc.FileIDs, newFile.ID)
	return nil
}
