package tpsmartapi

import (
	"context"
	"time"

	tpshp "github.com/ppacher/tplink-smart-home-protocol"
)

// EMeter provides access to the engery meter built into some TP-Link smart home
// devices
type EMeter interface {
	// GetRealtime queries the current realtime energy consumption
	GetRealtime(context.Context) <-chan *RealtimeInfo

	// GetDayStats returns statistics per day for the given month and year
	GetDayStats(context.Context, time.Month, int) <-chan *DailyStats

	// GetMonthStats returns statistics per month for the given year
	GetMonthStats(context.Context, int) <-chan *MonthlyStats

	// EraseStats erases the statistics memory of the energy meter
	EraseStats(context.Context) <-chan RPCError
}

type emeter struct {
	ns     Namespaces
	client tpshp.Client
}

// NewEmeter returns a new energy meter client
func NewEmeter(client tpshp.Client, ns Namespaces) EMeter {
	return &emeter{
		ns:     ns,
		client: client,
	}
}

func (em *emeter) GetRealtime(ctx context.Context) <-chan *RealtimeInfo {
	res := make(chan *RealtimeInfo, 1)
	var realtime RealtimeInfo

	req := tpshp.NewRequest().AddCommand(em.ns["emeter"], "get_realtime", struct{}{}, &realtime)

	Call(ctx, em.client, req, &realtime.NetErr, func() { res <- &realtime })

	return res
}

func (em *emeter) GetDayStats(ctx context.Context, month time.Month, year int) <-chan *DailyStats {
	res := make(chan *DailyStats, 1)
	var result DailyStats
	req := tpshp.NewRequest().AddCommand(em.ns["emeter"], "get_daystat", map[string]interface{}{
		"month": int(month),
		"year":  year,
	}, &result)

	Call(ctx, em.client, req, &result.NetErr, func() { res <- &result })
	return res
}

func (em *emeter) GetMonthStats(ctx context.Context, year int) <-chan *MonthlyStats {
	res := make(chan *MonthlyStats, 1)
	var result MonthlyStats
	req := tpshp.NewRequest().AddCommand(em.ns["emeter"], "get_monthstat", map[string]interface{}{
		"year": year,
	}, &result)

	Call(ctx, em.client, req, &result.NetErr, func() { res <- &result })
	return res
}

func (em *emeter) EraseStats(ctx context.Context) <-chan RPCError {
	res := make(chan RPCError, 1)
	var result ErrorHandler

	req := tpshp.NewRequest().AddCommand(em.ns["emeter"], "erase_emeter_stat", struct{}{}, &result)
	Call(ctx, em.client, req, &result.NetErr, func() { res <- &result })
	return res
}
