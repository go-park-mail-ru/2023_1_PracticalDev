package profile

// API requests
type fullUpdateRequest struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}

type partialUpdateRequest struct {
	Username     *string `json:"username"`
	Name         *string `json:"name"`
	ProfileImage *string `json:"profile_image"`
	WebsiteUrl   *string `json:"website_url"`
}

// API responses
type fullUpdateResponse struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}

type partialUpdateResponse struct {
	Username     string `json:"username"`
	Name         string `json:"name"`
	ProfileImage string `json:"profile_image"`
	WebsiteUrl   string `json:"website_url"`
}
