package entity

type (
	Movie struct {
		ID        string  `json:"id"`
		NameUz    string  `json:"name_uz"`
		NameRu    string  `json:"name_ru"`
		NameEn    string  `json:"name_en"`
		Embedding string  `json:"-"`
		Distance  float32 `json:"distance"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}

	MovieSingleRequest struct {
		ID     string `json:"id"`
		NameUz string `json:"name_uz"`
		NameRu string `json:"name_ru"`
		NameEn string `json:"name_en"`
	}

	MovieList struct {
		Items []Movie `json:"movie"`
		Count int     `json:"count"`
	}
)
