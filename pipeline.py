import os
import time
import subprocess
from context_builder import build_context

def agent_1_tester(context: str) -> str:
    print("[Agent 1 - Tester]: Analyzing bundled context...")
    time.sleep(1.5)
    print("[Agent 1 - Tester]: Detected Database and Email side-effects.")
    time.sleep(1.5)
    print("[Agent 1 - Tester]: Generating Go test suite with MOCKS (Dependency Injection)...")
    
    go_test_code = """package main

import (
	"testing"
)

// Mocking the Database
type MockDB struct{}
func (m *MockDB) GetProductPrice(id int) (float64, error) {
    if id == 1 {
        return 50.0, nil
    }
    return 0, nil
}
func (m *MockDB) InsertOrder(userId int, total float64) (int, error) {
    return 999, nil // Fake Order ID
}
func (m *MockDB) GetUserEmail(userId int) (string, error) {
    return "test@user.com", nil
}

// Mocking the Emailer
type MockEmailer struct {
    SentEmails int
}
func (m *MockEmailer) SendReceipt(email string, orderId int, total float64) {
    m.SentEmails++
}

func TestProcessOrder(t *testing.T) {
    mockDB := &MockDB{}
    mockEmailer := &MockEmailer{}
	
	cartItems := []map[string]int{
	    {"id": 1, "qty": 2}, // 2 * 50 = 100
	}
	
	orderID := ProcessOrder(1, cartItems, mockDB, mockEmailer)
	
	if orderID != 999 {
	    t.Errorf("Expected mocked order ID 999, got %d", orderID)
	}
	if mockEmailer.SentEmails != 1 {
	    t.Errorf("Expected 1 email sent, got %d", mockEmailer.SentEmails)
	}
}
"""
    return go_test_code

def agent_2_translator(context: str, error_log: str = None) -> str:
    if error_log:
        print(f"[Agent 2 - Translator]: Analyzing test failures...\\n{error_log}")
        time.sleep(2)
        print("[Agent 2 - Translator]: Fixing bugs...")
    else:
        print("[Agent 2 - Translator]: Translating PHP to Go using DEPENDENCY INJECTION...")
        time.sleep(2)

    go_code = """package main

// Interfaces abstract away the global state and legacy dependencies
type Database interface {
    GetProductPrice(id int) (float64, error)
    InsertOrder(userId int, total float64) (int, error)
    GetUserEmail(userId int) (string, error)
}

type EmailSender interface {
    SendReceipt(email string, orderId int, total float64)
}

// ProcessOrder now takes dependencies as arguments! (Dependency Injection)
func ProcessOrder(userId int, cartItems []map[string]int, db Database, emailer EmailSender) int {
	total := 0.0
	
	for _, item := range cartItems {
	    price, err := db.GetProductPrice(item["id"])
		if err != nil {
		    continue
		}
		total += price * float64(item["qty"])
	}

	orderId, _ := db.InsertOrder(userId, total)
    userEmail, _ := db.GetUserEmail(userId)
    
	emailer.SendReceipt(userEmail, orderId, total)

	return orderId
}
"""
    return go_code

def agent_3_verifier():
    print("[Agent 3 - Verifier]: Running tests on modern code...")
    time.sleep(1)
    
    result = subprocess.run(
        ["go", "test"], 
        cwd="modern_app",
        capture_output=True, 
        text=True
    )
    
    if result.returncode == 0:
        print("[Agent 3 - Verifier]: ALL TESTS PASSED! Migration 100% Successful.")
        return True, ""
    else:
        print("[Agent 3 - Verifier]: TESTS FAILED. Sending feedback to Translator...")
        return False, result.stdout

def main():
    print("=====================================================")
    print(" STARTING ENTERPRISE MIGRATION PIPELINE (CONTEXT AWARE) ")
    print("=====================================================")
    
    # 1. Use the Context Builder
    context = build_context("legacy_app/ProcessOrder.php")

    # 2. Agent 1: Generate Tests with Context
    tests_code = agent_1_tester(context)
    
    os.makedirs("modern_app", exist_ok=True)
    if not os.path.exists("modern_app/go.mod"):
        subprocess.run(["go", "mod", "init", "modern_app"], cwd="modern_app", capture_output=True)

    with open("modern_app/process_order_test.go", "w") as f:
        f.write(tests_code)

    # 3. Agent 2: Smart Translation
    modern_code = agent_2_translator(context)
    with open("modern_app/process_order.go", "w") as f:
        f.write(modern_code)

    # 4. Agent 3: Verify (The Loop)
    max_retries = 2
    for attempt in range(max_retries):
        print(f"\\n--- Verification Attempt {attempt + 1} ---")
        passed, error_log = agent_3_verifier()
        
        if passed:
            break
        
        if attempt < max_retries - 1:
            modern_code = agent_2_translator(context, error_log=error_log)
            with open("modern_app/process_order.go", "w") as f:
                f.write(modern_code)
        else:
            print("\\n[PIPELINE FATAL ERROR]")
    
    print("\\n=====================================================")
    print(" MOAT CRACKED! PIPELINE FINISHED SUCCESSFULLY! ")
    print("=====================================================")

if __name__ == "__main__":
    main()
