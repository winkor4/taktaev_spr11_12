package storage

import "github.com/winkor4/taktaev_spr11_12/internal/model"

// Описание результата функции checkTextData
type checkTextDataResult struct {
	resultData []model.TextDataList           // Cлайс данных для дальнейшей обработки
	response   []model.UploadTextDataResponse // Слайс данных для ответа на запрос
}
