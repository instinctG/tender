package service

import (
	"fmt"
	"github.com/instinctG/tender/internal/model"
)

type Store interface {
	GetUserBids(params model.GetUserBidsParams) ([]*model.Bid, error)
	CreateBid(params model.CreateBidJSONBody) (*model.Bid, error)
	EditBid(bidId string, params model.EditBidParams, body model.EditBidJSONBody) (*model.Bid, error)
	SubmitBidFeedback(bidId string, params model.SubmitBidFeedbackParams) *model.Bid
	RollbackBid(bidId string, version int32, params model.RollbackBidParams) *model.Bid
	GetBidStatus(bidId string, params model.GetBidStatusParams) (string, error)
	UpdateBidStatus(bidId string, params model.UpdateBidStatusParams) (*model.Bid, error)
	SubmitBidDecision(bidId string, params model.SubmitBidDecisionParams) (*model.Bid, error)
	GetBidsForTender(tenderId string, params model.GetBidsForTenderParams) ([]*model.Bid, error)
	GetBidReviews(tenderId string, params model.GetBidReviewsParams) []*model.BidReview
	GetTenders(params model.GetTendersParams) ([]*model.Tender, error) //todo : доделать
	GetUserTenders(params model.GetUserTendersParams) ([]*model.Tender, error)
	CreateTender(params model.CreateTenderJSONBody) (*model.Tender, error)
	EditTender(tenderId string, par model.EditTenderParams, params model.EditTenderJSONBody) (*model.Tender, error)
	RollbackTender(tenderId string, version int32, params model.RollbackTenderParams) (*model.Tender, error)
	GetTenderStatus(tenderId string, params model.GetTenderStatusParams) (string, error)
	UpdateTenderStatus(tenderId string, params model.UpdateTenderStatusParams) (*model.Tender, error)
}

type Service struct {
	Store Store
}

// NewService создает новый экземпляр Service.
func NewService(store Store) *Service {
	return &Service{Store: store}
}

func (s *Service) GetUserBids(params model.GetUserBidsParams) ([]*model.Bid, error) {
	bids, err := s.Store.GetUserBids(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bids, nil
}

func (s *Service) CreateBid(params model.CreateBidJSONBody) (*model.Bid, error) {
	bid, err := s.Store.CreateBid(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bid, nil
}

func (s *Service) EditBid(bidId string, params model.EditBidParams, body model.EditBidJSONBody) (*model.Bid, error) {
	editedBid, err := s.Store.EditBid(bidId, params, body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return editedBid, nil
}

func (s *Service) SubmitBidFeedback(bidId string, params model.SubmitBidFeedbackParams) *model.Bid {
	return s.Store.SubmitBidFeedback(bidId, params)
}

func (s *Service) RollbackBid(bidId string, version int32, params model.RollbackBidParams) *model.Bid {
	return s.Store.RollbackBid(bidId, version, params)
}

func (s *Service) GetBidStatus(bidId string, params model.GetBidStatusParams) (string, error) {
	status, err := s.Store.GetBidStatus(bidId, params)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return status, nil
}

func (s *Service) UpdateBidStatus(bidId string, params model.UpdateBidStatusParams) (*model.Bid, error) {
	bid, err := s.Store.UpdateBidStatus(bidId, params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bid, nil
}

func (s *Service) SubmitBidDecision(bidId string, params model.SubmitBidDecisionParams) (*model.Bid, error) {
	return s.Store.SubmitBidDecision(bidId, params)
}

func (s *Service) GetBidsForTender(tenderId string, params model.GetBidsForTenderParams) ([]*model.Bid, error) {
	bids, err := s.Store.GetBidsForTender(tenderId, params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return bids, nil
}

func (s *Service) GetBidReviews(tenderId string, params model.GetBidReviewsParams) []*model.BidReview {
	return s.Store.GetBidReviews(tenderId, params)
}

func (s *Service) GetTenders(params model.GetTendersParams) ([]*model.Tender, error) {
	tenders, err := s.Store.GetTenders(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return tenders, nil
}

func (s *Service) GetUserTenders(params model.GetUserTendersParams) ([]*model.Tender, error) {
	tenders, err := s.Store.GetUserTenders(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return tenders, nil
}

func (s *Service) CreateTender(params model.CreateTenderJSONBody) (*model.Tender, error) {
	tender, err := s.Store.CreateTender(params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return tender, nil
}

func (s *Service) EditTender(tenderId string, par model.EditTenderParams, params model.EditTenderJSONBody) (*model.Tender, error) {
	editedTender, err := s.Store.EditTender(tenderId, par, params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return editedTender, nil
}

func (s *Service) RollbackTender(tenderId string, version int32, params model.RollbackTenderParams) (*model.Tender, error) {
	return s.Store.RollbackTender(tenderId, version, params)
}

func (s *Service) GetTenderStatus(tenderId string, params model.GetTenderStatusParams) (string, error) {
	status, err := s.Store.GetTenderStatus(tenderId, params)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return status, nil
}

func (s *Service) UpdateTenderStatus(tenderId string, params model.UpdateTenderStatusParams) (*model.Tender, error) {
	tender, err := s.Store.UpdateTenderStatus(tenderId, params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return tender, nil
}
