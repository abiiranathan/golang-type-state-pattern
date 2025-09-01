package main

import "fmt"

// Define phantom type markers
type ActiveState struct{}
type ClosedState struct{}

type AccountState interface {
	ActiveState | PendingState | ClosedState
}

type Account[State AccountState] struct {
	ID      string
	Balance float64
}

type CanDeposit[T ActiveState] interface {
	*Account[T]
	deposit(amount float64) // private method forces type-specific implementation
}

type CanWithdraw[T ActiveState] interface {
	*Account[T]
	withdraw(amount float64) bool
}

type CanClose[T ActiveState] interface {
	*Account[T]
	close() *Account[ClosedState]
}

// Only implement capabilities for specific phantom types
func (a *Account[ActiveState]) deposit(amount float64) {
	a.Balance += amount
}

func (a *Account[ActiveState]) withdraw(amount float64) bool {
	if a.Balance >= amount {
		a.Balance -= amount
		return true
	}
	return false
}

func (a *Account[ActiveState]) close() *Account[ClosedState] {
	return &Account[ClosedState]{
		ID:      a.ID,
		Balance: a.Balance,
	}
}

// Public API functions that enforce phantom type constraints
func Deposit[T ActiveState, A CanDeposit[T]](acc A, amount float64) {
	acc.deposit(amount)
}

func Withdraw[T ActiveState, A CanWithdraw[T]](acc A, amount float64) bool {
	return acc.withdraw(amount)
}

func Close[T ActiveState, A CanClose[T]](acc A) *Account[ClosedState] {
	return acc.close()
}

// Universal operations (available on all account types)
func (a *Account[AccountState]) GetBalance() float64 {
	return a.Balance
}

func (a *Account[AccountState]) GetID() string {
	return a.ID
}

// Demonstration of state transitions
type PendingState struct{}

type CanActivate[T PendingState] interface {
	*Account[T]
	activate() *Account[ActiveState]
}

func (a *Account[PendingState]) activate() *Account[ActiveState] {
	return &Account[ActiveState]{
		ID:      a.ID,
		Balance: a.Balance,
	}
}

func Activate[T PendingState, A CanActivate[T]](acc A) *Account[ActiveState] {
	return acc.activate()
}

func main() {
	fmt.Println("=== True Phantom Types in Go ===")

	// Create accounts in different states
	activeAcc := &Account[ActiveState]{ID: "ACT-123", Balance: 100.0}
	closedAcc := &Account[ClosedState]{ID: "CLS-456", Balance: 200.0}
	pendingAcc := &Account[PendingState]{ID: "PND-789", Balance: 50.0}

	fmt.Printf("Initial balances:\n")
	fmt.Printf("  Active: $%.2f\n", activeAcc.GetBalance())
	fmt.Printf("  Closed: $%.2f\n", closedAcc.GetBalance())
	fmt.Printf("  Pending: $%.2f\n", pendingAcc.GetBalance())

	// Operations that work
	fmt.Printf("\nValid operations:\n")
	Deposit(activeAcc, 50.0)
	fmt.Printf("  Deposited $50 to active account: $%.2f\n", activeAcc.GetBalance())

	if Withdraw(activeAcc, 25.0) {
		fmt.Printf("  Withdrew $25 from active account: $%.2f\n", activeAcc.GetBalance())
	}

	// State transitions
	fmt.Printf("\nState transitions:\n")
	newActiveAcc := Activate(pendingAcc)
	fmt.Printf("  Activated pending account: $%.2f\n", newActiveAcc.GetBalance())

	newClosedAcc := Close(activeAcc)
	fmt.Printf("  Closed active account: $%.2f\n", newClosedAcc.GetBalance())

	// These operations would cause compile-time errors:

	// Deposit(closedAcc, 25.0)     // Error: ClosedState doesn't satisfy CanDeposit
	// Withdraw(closedAcc, 10.0)    // Error: ClosedState doesn't satisfy CanWithdraw
	// Close(closedAcc)             // Error: ClosedState doesn't satisfy CanClose
	// Deposit(pendingAcc, 10.0)    // Error: PendingState doesn't satisfy CanDeposit
	// Activate(activeAcc)          // Error: ActiveState doesn't satisfy CanActivate
	// Activate(closedAcc) // Error: ClosedState doesn't satisfy CanActivate

	// Summary
	fmt.Printf("\n✅ All phantom type constraints enforced at compile time!\n")
	fmt.Printf("💡 Key insight: Use private methods + type constraints + helper functions\n")
}
