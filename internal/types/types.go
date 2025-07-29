package types

type Student struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string ` json:"email" validate:"required"`
	Age   int    ` json:"age" validate:"required"`
}

type UpdateStudentRequest struct {
	Id    int64   `json:"id"`    // usually required to identify
	Name  *string `json:"name"`  // optional
	Email *string `json:"email"` // optional
	Age   *int    `json:"age"`   // optional
}
