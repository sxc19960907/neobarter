package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/neobarter/server/internal/service"
)

// Scheduler 运行后台周期任务。
type Scheduler struct {
	tradeSvc *service.TradeService
}

func New(tradeSvc *service.TradeService) *Scheduler {
	return &Scheduler{tradeSvc: tradeSvc}
}

// Start 启动所有后台任务，随 ctx 取消而停止。
func (s *Scheduler) Start(ctx context.Context) {
	go s.runTradeExpiry(ctx)
}

// runTradeExpiry 每分钟扫描一次超时的 pending 交易并置为 expired。
func (s *Scheduler) runTradeExpiry(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// 启动时先跑一次，处理停机期间累积的过期交易
	s.expireOnce()

	for {
		select {
		case <-ctx.Done():
			log.Println("Trade expiry scheduler stopped")
			return
		case <-ticker.C:
			s.expireOnce()
		}
	}
}

func (s *Scheduler) expireOnce() {
	n, err := s.tradeSvc.ExpireStale()
	if err != nil {
		log.Printf("Trade expiry scan error: %v", err)
		return
	}
	if n > 0 {
		log.Printf("Trade expiry: %d 笔超时交易已置为 expired", n)
	}
}
