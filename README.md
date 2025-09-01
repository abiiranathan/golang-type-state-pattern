This is an excellent and comprehensive explanation of the type-state pattern implementation. You've clearly articulated the philosophy, design choices, and technical details. The emphasis on **compile-time guarantees through unexported methods and anchored interfaces** is precisely correct for achieving robust capability-based security in Go.

Let's distill this into a more concise, actionable README suitable for the repository, keeping all the critical insights.

---

# `bank` Package: Type-State Pattern with Compile-Time Guarantees

This package implements a bank account using the **type-state pattern** (also known as the behavior pattern). Account state (`active`, `pending`, `closed`) is encoded as a compile-time type parameter on the `Account` struct. The public API consists of top-level functions (e.g., `Deposit`, `Withdraw`) that are only available for accounts in the correct state, enforced by the Go compiler.

## 🧠 Philosophy: Capabilities as Types

An `*Account[active]` is a **capability token**. Holding this type grants you the right to call `Deposit`, `Withdraw`, or `Close`. If you only have an `*Account[closed]`, these operations are statically unavailable. This design moves runtime state checks to compile-time type safety, eliminating entire classes of bugs.

**Crucially, this is not achieved through exported methods on `*Account[...]`.** Instead, we use a combination of generic interfaces, unexported methods, and type anchoring to ensure:
1.  **True Compile-Time Guarantees:** Invalid operations fail to compile.
2.  **Encapsulation:** External packages cannot fake the required capabilities.
3.  **Minimal API Surface:** A simple, consistent set of functions for all operations.

## 🏗️ Implementation Overview

### 1. State Markers and Core Type
State is defined by empty structs used as type parameters.

```go
// Phantom type markers for state
type active struct{}
type pending struct{}
type closed struct{}

// Core generic type
type Account[S accountState] struct {
	ID      string
	Balance float64
}
```

### 2. Unexported Methods and Anchored Interfaces
The real logic is defined as **unexported methods** on the concrete types (e.g., `*Account[active]`). Public functions are constrained by interfaces that **anchor** the type parameter to an `*Account[T]` and require an unexported method.

**Interface Definition (Anchoring):**
```go
// CanDeposit is anchored to *Account[active] and requires the unexported deposit() method.
type CanDeposit[T active] interface {
	*Account[T] // Anchor: T must be a type parameter that eventually resolves to 'active'
	deposit(amount float64) // Unexported method, only implementable in this package
}
```

**Method Implementation (Package-private):**
```go
// deposit is only defined on *Account[active] and is unexported.
func (a *Account[active]) deposit(amount float64) {
	a.Balance += amount
}
```

### 3. Public Wrapper Functions
The public API consists of generic functions constrained by the anchored interfaces.

```go
// Deposit is a top-level function. It only compiles if `acc` is an *Account[active].
func Deposit[T active, A CanDeposit[T]](acc A, amount float64) {
	// Can call the unexported method because `A` satisfies CanDeposit.
	acc.deposit(amount)
}
```

## 🚀 Usage Example

```go
package main

// Import your package here.
import "typestate/bank"

func main() {
    // Create a new account in a 'pending' state.
    p := &bank.PendingAccount{ID: "acct-123", Balance: 0}

    // Activate it. Statically allowed because p is *Account[pending].
    a := bank.Activate(p) // returns *Account[active]

    // Perform operations on the active account.
    bank.Deposit(a, 100.0) // Compiles
    ok := bank.Withdraw(a, 50.0)

    // Close the account.
    c := bank.Close(a) // returns *Account[closed]

    // This line would FAIL TO COMPILE.
    // bank.Deposit(c, 10.0) // Error: *Account[closed] does not satisfy constraint CanDeposit.
}
```

## ⚠️ Why Exported Methods Are Not Used

Using exported methods (e.g., `func (a *Account[active]) Deposit(...)`) as the public API is tempting but **inadequate** for our goals:
1.  **Loose Constraints:** A generic function constrained only by `interface { Deposit(float64) }` could be satisfied by any external type, breaking the guarantee that the value is an `*Account[active]`.
2.  **No Anchoring:** It doesn't force the type parameter `T` to be `active`; it only requires a method named `Deposit`.
3.  **Poor Encapsulation:** It exposes internal implementation details as the public API.

Our pattern (unexported methods + anchored interfaces + wrapper functions) solves all these problems, providing rigorous compile-time safety.

## 🔧 Testing

*   **Internal Tests:** Tests for unexported logic must reside in the
*   same package.
  
## ✅ FAQ

**Q: Can code outside this package create a type that satisfies `CanDeposit`?**
**A: No.** The `CanDeposit` interface requires an embedded `*Account[active]` *and* an unexported method `deposit()`. Since external packages cannot define this unexported method, they cannot satisfy the interface.

**Q: Does this pattern have any runtime overhead?**
**A: No.** All safety is enforced at compile time. The wrapper functions are simple calls to the underlying unexported methods.

**Q: How are state transitions handled?**
**A: Transitions are explicit and type-safe.** Functions like `Activate` and `Close` return a new account value with a new state type (`*Account[active]`, `*Account[closed]`).

---

## License

MIT