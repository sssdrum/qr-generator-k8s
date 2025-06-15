# Operators

An operator extends Kubernetes API to manage complex applications using domain-specific knowledge

# Operator SDK

Provides tools to:

- build
- test
- package

Kubernetes operators.
It uses the **controller runtime**

## Components

- **CRD: (Custom Resource Definition)** = API schema
- **Controller:** == Logic that reconciles desired vs actual state
- **Reconciler** == Function that gets called when resources change
- **Manger**: == Runs controller and handles leader election

## Reconciliation loop

```go
// This pattern shows up in every interview
func (r *MyAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // 1. Fetch the custom resource
    // 2. Check current state of managed resources
    // 3. Create/Update/Delete resources to match desired state
    // 4. Update status
    // 5. Return result (requeue if needed)
}
```

## Misc

### Controller pattern:

- Watch for desired state
- observe actual state
- take action to reconcile -
- repeat
