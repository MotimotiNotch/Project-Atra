# Atra Usage Guide

This guide explains how to use the Atra compiler to analyze your own "Semantic Universe."

## 1. Setup Go Environment

Ensure you have [Go](https://go.dev/) installed on your system.

1.  **Check Installation**: Open your terminal and run:
    ```bash
    go version
    ```
    You should see `go version go1.x.x ...`.

## 2. Running Atra

1.  **Navigate to the Atra directory**:
    ```bash
    cd atra
    ```
2.  **Execute the compiler**:
    You can run the compiler directly using `go run`:
    ```bash
    go run main.go -file ../experiments/slime_morphology/mindset.atra
    ```
    This will parse the provided `.atra` file and output the calculated gravity scores.

## 3. Describing Your Own "Universe"

You can create your own `.atra` files to define semantic attractors and phase structures.

Example syntax:
```atra
Phase(Identity) {
  Phase(Core) {
    Event(Introspection) {
      viscosity: +0.80
      crystallinity: +0.20
    }
  }
}
```

## 4. Development & Testing

- **Run Tests**:
    ```bash
    go test ./...
    ```
- **Custom Mapping**: You can modify the mapping between keywords and phases in `main.go` to suit your specific experiment.

---
"In Atra, we don't just set prompts; we build the gravity that shapes the model's soul."
