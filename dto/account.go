package dto

type ChangeUsernameRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
}

type ChangeEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	RecoveryKey string `json:"recoveryKey"`
	NewPassword string `json:"newPassword"`
}
