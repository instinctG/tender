package server

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/instinctG/tender/internal/model"
	"log"
	"net/http"
	"strconv"
)

type TenderService interface {
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
	GetTenders(params model.GetTendersParams) ([]*model.Tender, error)
	GetUserTenders(params model.GetUserTendersParams) ([]*model.Tender, error)
	CreateTender(params model.CreateTenderJSONBody) (*model.Tender, error)
	EditTender(tenderId string, par model.EditTenderParams, params model.EditTenderJSONBody) (*model.Tender, error)
	RollbackTender(tenderId string, version int32, params model.RollbackTenderParams) (*model.Tender, error)
	GetTenderStatus(tenderId string, params model.GetTenderStatusParams) (string, error)
	UpdateTenderStatus(tenderId string, params model.UpdateTenderStatusParams) (*model.Tender, error)
}

func (h *Handler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var params model.CreateTenderJSONBody

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		jsonRespond(w, http.StatusBadRequest, "error in decoding a body")
		return
	}

	tender, err := h.Service.CreateTender(params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot create a tender")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(tender); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a tender")
		return
	}

}

func (h *Handler) GetTenders(w http.ResponseWriter, r *http.Request) {
	var params model.GetTendersParams
	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	params.Limit = int32(limit)
	params.Offset = int32(offset)
	params.ServiceType = queryParams["service_type"]

	tenders, err := h.Service.GetTenders(params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot get tenders from service")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tenders); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a list of tenders")
		return
	}
}

// todo : доделать обработку ошибок на 400,401,500
func (h *Handler) GetUserTenders(w http.ResponseWriter, r *http.Request) {
	var params model.GetUserTendersParams
	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	params.Limit = int32(limit)
	params.Offset = int32(offset)
	params.Username = queryParams.Get("username")

	myTenders, err := h.Service.GetUserTenders(params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot get user tenders from service")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(myTenders); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a list of user tenders")
		return
	}
}

func (h *Handler) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	var params model.GetTenderStatusParams
	vars := mux.Vars(r)
	tenderId := vars["tenderId"]
	queryParams := r.URL.Query()
	params.Username = queryParams.Get("username")
	if !IsValidUUID(tenderId) {
		jsonRespond(w, http.StatusBadRequest, "tender id is invalid")
		return
	}

	status, err := h.Service.GetTenderStatus(tenderId, params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "tender status not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(status); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a tender status")
		return
	}
}

func (h *Handler) UpdateTenderStatus(w http.ResponseWriter, r *http.Request) {
	var params model.UpdateTenderStatusParams
	vars := mux.Vars(r)
	tenderId := vars["tenderId"]
	if !IsValidUUID(tenderId) {
		jsonRespond(w, http.StatusBadRequest, "tender id is invalid")
	}
	queryParams := r.URL.Query()
	params.Status, params.Username = queryParams.Get("status"), queryParams.Get("username")

	tender, err := h.Service.UpdateTenderStatus(tenderId, params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "tender can not be updated")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tender); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode an updated tender")
		return
	}
}

func (h *Handler) EditTender(w http.ResponseWriter, r *http.Request) {
	var params model.EditTenderJSONBody
	var par model.EditTenderParams
	var zeroStruct model.EditTenderJSONBody
	vars := mux.Vars(r)
	tenderId := vars["tenderId"]
	queryParams := r.URL.Query()
	par.Username = queryParams.Get("username")

	if !IsValidUUID(tenderId) {
		jsonRespond(w, http.StatusBadRequest, "tender id is invalid")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		if params != zeroStruct {
			jsonRespond(w, http.StatusInternalServerError, "json body cannot be decoded")
			return
		}
	}
	editedTender, err := h.Service.EditTender(tenderId, par, params)
	if err != nil && err.Error() != "no updates provided" {
		jsonRespond(w, http.StatusInternalServerError, "tender can not be edited")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(editedTender); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode an edited tender status")
		return
	}
}

func (h *Handler) CreateBid(w http.ResponseWriter, r *http.Request) {
	var params model.CreateBidJSONBody
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		jsonRespond(w, http.StatusBadRequest, "error in decoding a bid body")
		return
	}

	bid, err := h.Service.CreateBid(params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot create a bid")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(bid); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a bid")
		return
	}
}

func (h *Handler) GetUserBids(w http.ResponseWriter, r *http.Request) {
	var params model.GetUserBidsParams
	queryParams := r.URL.Query()

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	params.Limit = int32(limit)
	params.Offset = int32(offset)
	params.Username = queryParams.Get("username")

	myBids, err := h.Service.GetUserBids(params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot get user bids from service")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(myBids); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a list of user bids")
		return
	}
}

func (h *Handler) GetBidsForTender(w http.ResponseWriter, r *http.Request) {
	var params model.GetBidsForTenderParams
	queryParams := r.URL.Query()

	vars := mux.Vars(r)
	tenderId := vars["tenderId"]

	limit, err := strconv.Atoi(queryParams.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 5
	}
	if limit > 50 {
		limit = 50
	}

	offset, err := strconv.Atoi(queryParams.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	params.Limit = int32(limit)
	params.Offset = int32(offset)
	params.Username = queryParams.Get("username")

	bids, err := h.Service.GetBidsForTender(tenderId, params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot get bids for tender from service")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bids); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a list of bids for tender")
		return
	}
}

func (h *Handler) GetBidStatus(w http.ResponseWriter, r *http.Request) {
	var params model.GetBidStatusParams
	vars := mux.Vars(r)
	bidId := vars["bidId"]
	queryParams := r.URL.Query()
	params.Username = queryParams.Get("username")

	if !IsValidUUID(bidId) {
		jsonRespond(w, http.StatusBadRequest, "bid id is invalid")
		return
	}

	status, err := h.Service.GetBidStatus(bidId, params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "bid status not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(status); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode a bid status")
		return
	}
}

func (h *Handler) UpdateBidStatus(w http.ResponseWriter, r *http.Request) {
	var params model.UpdateBidStatusParams
	vars := mux.Vars(r)
	bidId := vars["bidId"]
	if !IsValidUUID(bidId) {
		jsonRespond(w, http.StatusBadRequest, "bid id is invalid")
	}

	queryParams := r.URL.Query()
	params.Status, params.Username = queryParams.Get("status"), queryParams.Get("username")

	bid, err := h.Service.UpdateBidStatus(bidId, params)
	if err != nil {
		jsonRespond(w, http.StatusInternalServerError, "bid status can not be updated")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(bid); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode an updated bid status")
		return
	}
}

func (h *Handler) EditBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body model.EditBidJSONBody
	var params model.EditBidParams
	var zeroStruct model.EditBidJSONBody
	vars := mux.Vars(r)
	bidId := vars["bidId"]
	queryParams := r.URL.Query()
	params.Username = queryParams.Get("username")

	if !IsValidUUID(bidId) {
		jsonRespond(w, http.StatusBadRequest, "bid id is invalid")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		if body != zeroStruct {
			jsonRespond(w, http.StatusInternalServerError, "json body cannot be decoded")
			return
		}
	}
	editedBid, err := h.Service.EditBid(bidId, params, body)
	if err != nil && err.Error() != "no updates provided" {
		jsonRespond(w, http.StatusInternalServerError, "bid can not be edited")
		return
	}

	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(editedBid); err != nil {
		jsonRespond(w, http.StatusInternalServerError, "cannot encode an edited bid status")
		return
	}
}

func jsonRespond(w http.ResponseWriter, statusCode int, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(model.ErrorResponse{Reason: data}); err != nil {
		log.Println(err)
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
