package domain

type AdvisoryLockId int

const (
	TempImageAdvisoryLockId AdvisoryLockId = 10
	OutboxAdvisoryLockId    AdvisoryLockId = 20
)
