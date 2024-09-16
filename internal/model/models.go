package model

import "time"

type Bid struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	TenderId    string    `json:"tenderId,omitempty"`
	AuthorType  string    `json:"authorType"`
	AuthorId    string    `json:"authorId"`
	Version     int32     `json:"version"`
	CreatedAt   time.Time `json:"createdAt"`
}

type BidReview struct {
	CreatedAt   string `json:"createdAt"`
	Description string `json:"description"`
	Id          string `json:"id"`
}

type ErrorResponse struct {
	Reason string `json:"reason"`
}

type Tender struct {
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Status         string    `json:"status"`
	ServiceType    string    `json:"serviceType"`
	OrganizationId string    `json:"organizationId,omitempty"`
	Version        int32     `json:"version"`
	CreatedAt      time.Time `json:"createdAt"`
}

type GetUserBidsParams struct {
	Limit    int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset   int32  `form:"offset,omitempty" json:"offset,omitempty"`
	Username string `form:"username,omitempty" json:"username,omitempty"`
}

type CreateBidJSONBody struct {
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	TenderId    string `json:"tenderId"`
}

type EditBidJSONBody struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

type EditBidParams struct {
	Username string `form:"username" json:"username"`
}

type SubmitBidFeedbackParams struct {
	BidFeedback string `form:"bidFeedback" json:"bidFeedback"`
	Username    string `form:"username" json:"username"`
}

type RollbackBidParams struct {
	Username string `form:"username" json:"username"`
}

type GetBidStatusParams struct {
	Username string `form:"username" json:"username"`
}

type UpdateBidStatusParams struct {
	Status   string `form:"status" json:"status"`
	Username string `form:"username" json:"username"`
}

type SubmitBidDecisionParams struct {
	Decision string `form:"decision" json:"decision"`
	Username string `form:"username" json:"username"`
}

type GetBidsForTenderParams struct {
	Username string `form:"username" json:"username"`
	Limit    int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset   int32  `form:"offset,omitempty" json:"offset,omitempty"`
}

type GetBidReviewsParams struct {
	AuthorUsername    string `form:"authorUsername" json:"authorUsername"`
	RequesterUsername string `form:"requesterUsername" json:"requesterUsername"`
	Limit             int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset            int32  `form:"offset,omitempty" json:"offset,omitempty"`
}

type GetTendersParams struct {
	Limit       int32    `form:"limit,omitempty" json:"limit,omitempty"`
	Offset      int32    `form:"offset,omitempty" json:"offset,omitempty"`
	ServiceType []string `form:"service_type,omitempty" json:"service_type,omitempty"`
}

type GetUserTendersParams struct {
	Limit    int32  `form:"limit,omitempty" json:"limit,omitempty"`
	Offset   int32  `form:"offset,omitempty" json:"offset,omitempty"`
	Username string `form:"username,omitempty" json:"username,omitempty"`
}

type CreateTenderJSONBody struct {
	CreatorUsername string `json:"creatorUsername"`
	Description     string `json:"description"`
	Name            string `json:"name"`
	OrganizationId  string `json:"organizationId"`
	ServiceType     string `json:"serviceType"`
}

type EditTenderJSONBody struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	ServiceType string `json:"serviceType,omitempty"`
}

type EditTenderParams struct {
	Username string `form:"username" json:"username"`
}

type RollbackTenderParams struct {
	Username string `form:"username" json:"username"`
}

type GetTenderStatusParams struct {
	Username string `form:"username,omitempty" json:"username,omitempty"`
}

type UpdateTenderStatusParams struct {
	Status   string `form:"status" json:"status"`
	Username string `form:"username" json:"username"`
}
