package main
import (
	"fmt"
	)

type TransactionContextInterface interface {
	// GetStub should provide a way to access the stub set by Init/Invoke
	GetStub() string
}

type TransactionContext struct {
	stub string
}

func (ctx *TransactionContext) GetStub() string {
	return ctx.stub
}


// SetStub stores the passed stub in the transaction context
func InitLedger(ctx TransactionContextInterface) error {
	fmt.Println(ctx.GetStub())
	return nil
}


func main() {
	ctx := new(TransactionContext)
	ctx.stub = "1"
	InitLedger(ctx)
}