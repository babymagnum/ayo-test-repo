package request

type AddPlayerRequest struct {
	Name         string  `json:"name" binding:"required,max=255"`
	HeightCm     float64 `json:"height_cm"`
	WeightKg     float64 `json:"weight_kg"`
	Position     string  `json:"position" binding:"required"`
	JerseyNumber int     `json:"jersey_number" binding:"required"`
}
