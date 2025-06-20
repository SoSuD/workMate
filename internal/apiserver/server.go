package apiserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nwtgck/go-fakelish"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"time"
	"workMate/internal/service/faketask"
	"workMate/internal/store"
)

type Server struct {
	router *gin.Engine
	db     store.Store
	logger *zap.Logger
	config Config
}

func NewServer(db store.Store, config Config) *Server {
	logger, _ := zap.NewDevelopment()
	level, err := zapcore.ParseLevel(config.Server.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	zap.IncreaseLevel(level)

	return &Server{
		router: gin.Default(),
		db:     db,
		logger: logger,
	}
}

func (s *Server) configureRouter() {
	s.router.Use(gin.Recovery())
	task := s.router.Group("/task")
	task.GET("/:id", s.getTask)
	task.DELETE("/:id", s.deleteTask)
	task.POST("", s.createTask)
}

type createTaskRequest struct {
	Name string `json:"name"`
}

type createTaskResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}

func (s *Server) createTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Name = fakelish.GenerateFakeWord(7, 9)
	}

	t, err := s.db.Task().Create(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrCreateTask})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)

	go func() {
		defer cancel()

		resCh := make(chan faketask.Result, 1)
		go faketask.LongTask(resCh)
		s.logger.Debug("task created", zap.String("id", t.ID.String()))
		select {
		case <-ctx.Done():
			fmt.Println("ctx done")

		case x := <-resCh:
			if x.Err != nil {
				s.logger.Error(fmt.Sprintf("task finished with error: %s", x.Err))
				if err := s.db.Task().Finish(t.ID, "FAILED", ""); err != nil {
					s.logger.Error("failed to finish task", zap.Error(err))
				}
				return
			} else {
				s.logger.Debug(fmt.Sprintf("task finished: %s", t.ID))
				if err := s.db.Task().Finish(t.ID, "DONE", x.Value); err != nil {
					s.logger.Error("failed to finish task", zap.Error(err))
				}
				return
			}

		}
	}()
	var res createTaskResponse
	res.ID = t.ID
	res.Name = t.Name
	res.CreatedAt = t.CreatedAt
	res.Status = t.Status
	c.JSON(http.StatusOK, res)
}

type getTaskResponse struct {
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"created_at"`
	ProcessDuration float64   `json:"process_duration"`
	FinishedAt      time.Time `json:"finished_at,omitzero"`
	Status          string    `json:"status"`
}

func (s *Server) getTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrParseUUID})
		return
	}

	task, err := s.db.Task().Get(id)
	if err != nil {
		s.logger.Error("failed to fetch task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGetTask})
		return
	}
	var res getTaskResponse
	res.Name = task.Name
	res.CreatedAt = task.CreatedAt
	res.ProcessDuration = time.Since(task.CreatedAt).Round(time.Millisecond).Seconds()
	res.Status = task.Status
	res.FinishedAt = task.FinishedAt
	c.JSON(http.StatusOK, res)
	return

}

func (s *Server) deleteTask(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrParseUUID})
		return
	}
	if err := s.db.Task().Delete(id); err != nil {
		s.logger.Error("failed to delete task", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrFailedDelete})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	return
}
