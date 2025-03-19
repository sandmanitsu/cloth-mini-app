package lock

import (
	domain "cloth-mini-app/internal/domain/lock"
	"context"
)

type LockRepository interface {
	AdvisoryLock(ctx context.Context, id domain.AdvisoryLockId) error
	AdvisoryUnlock(ctx context.Context, id domain.AdvisoryLockId) error
}

type LockService struct {
	lockRepo LockRepository
}

func NewLockService(lr LockRepository) *LockService {
	return &LockService{
		lockRepo: lr,
	}
}

func (l *LockService) AdvisoryLock(ctx context.Context, id domain.AdvisoryLockId) error {
	return l.lockRepo.AdvisoryLock(ctx, id)
}

func (l *LockService) AdvisoryUnlock(ctx context.Context, id domain.AdvisoryLockId) error {
	return l.lockRepo.AdvisoryUnlock(ctx, id)
}
