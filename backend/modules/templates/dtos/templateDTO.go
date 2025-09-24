package dtos

import "elogika.vsb.cz/backend/models"

type TemplateDTO struct {
	ID            uint                 `json:"id"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	MixBlocks     bool                 `json:"mixBlocks"`
	MixEverything bool                 `json:"mixEverything"`
	Blocks        []TemplateBlockDTO   `json:"blocks"`
	Version       uint                 `json:"version"`
	CreatedBy     TemplateCreatedByDTO `json:"createdBy"`
}

func (TemplateDTO) From(d *models.Template) TemplateDTO {
	dto := TemplateDTO{
		ID:            d.ID,
		Title:         d.Title,
		Description:   d.Description,
		MixBlocks:     d.MixBlocks,
		MixEverything: d.MixEverything,
		Blocks:        make([]TemplateBlockDTO, len(d.Blocks)),
		Version:       d.Version,
		CreatedBy:     TemplateCreatedByDTO{}.From(d.CreatedBy),
	}

	for i, block := range d.Blocks {
		dto.Blocks[i] = TemplateBlockDTO{}.From(&block)
	}

	return dto
}
