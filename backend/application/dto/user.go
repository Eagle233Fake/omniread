package dto

type UpdateProfileReq struct {
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Bio       string `json:"bio"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdatePreferencesReq struct {
	FontFamily string `json:"font_family"`
	FontSize   int    `json:"font_size"`
}
