package bank

// Define phantom type markers
type active struct{}
type closed struct{}
type pending struct{}

type state interface {
	active | pending | closed
}

type pendingOrClosed interface {
	pending | closed
}

type Account[State state] struct {
	ID      string
	Balance float64

	// we could add flags or metadata here
}

// ==================================================
// private method forces type-specific implementation
// ==================================================

// CanDeposit enforces that only certain states can deposit
type CanDeposit[T active] interface {
	*Account[T]
	deposit(amount float64)
}

// CanWithdraw enforces that only certain states can withdraw
type CanWithdraw[T active] interface {
	*Account[T]
	withdraw(amount float64) bool
}

// CanClose enforces that only certain states can close
type CanClose[T active] interface {
	*Account[T]
	close() *Account[closed]
}

// Intreface composition
type CanWithdrawAndClose[T active] interface {
	CanWithdraw[T]
	CanClose[T]
}

type CanActivate[T pendingOrClosed] interface {
	*Account[T]
	activate() *Account[active]
}

// Only ActiveAccount can deposit.
func (a *Account[ActiveState]) deposit(amount float64) {
	a.Balance += amount
}

// Only ActiveAccount can withdraw.
func (a *Account[ActiveState]) withdraw(amount float64) bool {
	if a.Balance >= amount {
		a.Balance -= amount
		return true
	}
	return false
}

// Only ActiveAccount can close.
func (a *Account[ActiveState]) close() *Account[closed] {
	return &Account[closed]{
		ID:      a.ID,
		Balance: a.Balance,
	}
}

// Public API functions with type constraints.
// Deposit, Withdraw, and Close can only be called on accounts in the correct state.
func Deposit[T active, A CanDeposit[T]](acc A, amount float64) {
	acc.deposit(amount)
}

// Withdraw returns true if successful, false if insufficient funds.
func Withdraw[T active, A CanWithdraw[T]](acc A, amount float64) bool {
	return acc.withdraw(amount)
}

func WithdrawAndClose[T active, A CanWithdrawAndClose[T]](acc A, amount float64) *Account[closed] {
	acc.withdraw(amount)
	return acc.close()
}

// Close transitions an active account to a closed account.
func Close[T active, A CanClose[T]](acc A) *Account[closed] {
	return acc.close()
}

// Universal operations (available on all account types)
func (a *Account[AccountState]) GetBalance() float64 {
	return a.Balance
}

// GetID returns the account ID.
func (a *Account[AccountState]) GetID() string {
	return a.ID
}

func (a *Account[PendingState]) activate() *Account[active] {
	return &Account[active]{
		ID:      a.ID,
		Balance: a.Balance,
	}
}

// Activate transitions a pending or closed account to an active account.
func Activate[T pendingOrClosed, A CanActivate[T]](acc A) *Account[active] {
	return acc.activate()
}

// Type specializations for Activate for Pending account.
// If we don't specialize, then the caller MUST specify the type parameter explicitly.
// But we may want to keep the state private, so we provide these helpers.
func ActivatePending(acc *Account[pending]) *Account[active] {
	return acc.activate()
}

// Type specializations for Activate for Closed account.
func ActivateClosed(acc *Account[closed]) *Account[active] {
	return acc.activate()
}

// Helpers
type ActiveAccount = Account[active]
type PendingAccount = Account[pending]
type ClosedAccount = Account[closed]
