package consts

import "elogika.vsb.cz/backend/models"

var DefaultContent = &models.TipTapContent{
	Type: "doc",
	Content: []*models.TipTapContent{
		{
			Type: "paragraph",
			Content: []*models.TipTapContent{
				{
					Type: "text",
					Text: "asdasdasdasdasd",
				},
			},
		},
	},
}
