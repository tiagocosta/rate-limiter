package limiter

// import (
// 	"context"
// 	"sync"
// 	"time"
// )

// type Limiter struct {
// 	burst           int
// 	window          time.Duration
// 	timer           *time.Timer
// 	expirationTime  time.Duration
// 	expirationTimer *time.Timer
// 	repository      RepositoryInterface

// 	sync.Mutex
// }

// func NewLimiter(burst int, expirationTime time.Duration, repository RepositoryInterface) *Limiter {
// 	limiter := &Limiter{
// 		burst:           burst,
// 		window:          time.Second,
// 		timer:           time.NewTimer(time.Second),
// 		expirationTime:  expirationTime,
// 		expirationTimer: time.NewTimer(expirationTime),
// 		repository:      repository,
// 	}

// 	limiter.startTimerVerification(context.Background())

// 	return limiter
// }

// func (limiter *Limiter) startTimerVerification(ctx context.Context) {
// 	go func() {
// 		for {
// 			select {
// 			case <-limiter.timer.C:
// 				if !limiter.isInExpirationTime() {

// 				}
// 				limiter.timer.Reset(limiter.window)
// 			case <-limiter.expirationTimer.C:
// 				limiter.Lock()
// 				limiter.expirationTimer = nil
// 				limiter.Unlock()

// 			}
// 		}
// 	}()
// }

// func (limiter *Limiter) Allow(consumed int) bool {
// 	if limiter.isInExpirationTime() {
// 		return false
// 	}

// 	if limiter.burst > consumed {
// 		return true
// 	} else {
// 		limiter.Lock()
// 		limiter.expirationTimer = time.NewTimer(limiter.expirationTime)
// 		limiter.Unlock()
// 	}

// 	return false
// }

// func (limiter *Limiter) isInExpirationTime() bool {
// 	return limiter.expirationTimer != nil
// }
