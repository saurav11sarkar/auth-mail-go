package user

// ProfileInput uses the authenticated identity, never a client-supplied body ID.
type ProfileInput struct{ UserID string }

type UpdateProfileInput struct {
	UserID string
	Name   *string
	Role   *string
	Status *string
}
