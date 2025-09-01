package main

import (
	"fmt"
	"phantom-types/bank"
)

func main() {
	fmt.Println("=== True Phantom Types in Go ===")

	// Create accounts in different states
	activeAcc := &bank.ActiveAccount{ID: "ACT-123", Balance: 100.0}
	closedAcc := &bank.ClosedAccount{ID: "CLS-456", Balance: 200.0}
	pendingAcc := &bank.PendingAccount{ID: "PND-789", Balance: 50.0}

	fmt.Printf("Initial balances:\n")
	fmt.Printf("  Active: $%.2f\n", activeAcc.GetBalance())
	fmt.Printf("  Closed: $%.2f\n", closedAcc.GetBalance())
	fmt.Printf("  Pending: $%.2f\n", pendingAcc.GetBalance())

	// Withdraw and Close active account (composite operation)
	nowClosedAcct := bank.WithdrawAndClose(activeAcc, 30.0)
	fmt.Printf("After withdrawing $30 and closing active account: $%.2f\n", nowClosedAcct.GetBalance())

	// Invalid operations (uncommenting these lines will cause compile-time errors)
	// bank.Withdraw(closedAcc, 50.0) // Cannot withdraw from closed account
	// bank.Deposit(pendingAcc, 20.0) // Cannot deposit to pending account
	// bank.Close(closedAcc)           // Cannot close an already closed account
	// bank.Withdraw(pendingAcc, 10.0) // Cannot withdraw from pending account

	// Operations that work
	fmt.Printf("\nValid operations:\n")
	bank.Deposit(activeAcc, 50.0)
	fmt.Printf("  Deposited $50 to active account: $%.2f\n", activeAcc.GetBalance())

	if bank.Withdraw(activeAcc, 25.0) {
		fmt.Printf("  Withdrew $25 from active account: $%.2f\n", activeAcc.GetBalance())
	}

	// State transitions
	fmt.Printf("\nState transitions:\n")
	newActiveAcc := bank.ActivatePending(pendingAcc)
	fmt.Printf("  Activated pending account: $%.2f\n", newActiveAcc.GetBalance())

	newClosedAcc := bank.Close(activeAcc)
	fmt.Printf("  Closed active account: $%.2f\n", newClosedAcc.GetBalance())
}
