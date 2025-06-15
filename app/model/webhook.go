package model

type WebhookPayload struct {
	Detail struct {
		Overrides struct {
			ContainerOverrides []struct {
				Environment []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"environment"`
			} `json:"containerOverrides"`
		} `json:"overrides"`
		DesiredStatus string `json:"desiredStatus"`
		StoppedReason string `json:"stoppedReason"`
	} `json:"detail"`
}
