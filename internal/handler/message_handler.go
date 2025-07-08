package handler

import (
	"context"
	"encoding/json"
	"errors"
	"go-sender/internal/service"
	"go-sender/internal/util"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var ErrAlreadyRunning = errors.New("sender is already running")

type MsgController struct {
	service service.IMessageService
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex
	running bool
}

func NewMsgController(service service.IMessageService, ctx context.Context) *MsgController {
	return &MsgController{service: service, ctx: ctx}
}

// @Summary Start sender
// @Tags    control
// @Success 200
// @Failure 400  {string}  string  "already running"
// @Router  /start [post]
func (c *MsgController) Start(w http.ResponseWriter, r *http.Request) {
	if err := c.StartService(); err != nil {
		if errors.Is(err, ErrAlreadyRunning) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Fatalf("failed to auto‐start message service: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Stop sender
// @Tags    control
// @Success 200
// @Failure 400  {string}  string  "not running"
// @Router  /stop [post]
func (c *MsgController) Stop(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		http.Error(w, "sender is not running", http.StatusBadRequest)
		return
	}

	c.service.Stop()
	c.running = false

	w.WriteHeader(http.StatusOK)
}

// @Summary      List sent messages
// @Description  Get message content, phone, sent-status.
// @Tags         messages
// @Produce      json
// @Success      200  {array}  message.Dto
// @Router       /message [get]
func (c *MsgController) GetSentMessages(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	paging := &util.Paging{
		Limit: limit,
		Page:  page,
	}

	messages, err := c.service.GetMessages(true, paging)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *MsgController) StartService() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return ErrAlreadyRunning
	}
	c.service.Initialize()
	if err := c.service.Start(); err != nil {
		return err
	}
	c.running = true
	return nil
}
