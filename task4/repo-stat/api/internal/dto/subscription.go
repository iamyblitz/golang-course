package dto

type SubscriptionRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type SubscriptionResponse struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type SubscriptionsResponse struct {
	Subscriptions []SubscriptionResponse `json:"subscriptions"`
}
