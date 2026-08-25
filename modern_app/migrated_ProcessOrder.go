package main
import "fmt"
type Database interface { GetProductPrice(id int) (float64, error) }
type EmailSender interface { SendReceipt(email string, orderId int, total float64) }
func ProcessOrder(userId int, cartItems []map[string]int, db Database, emailer EmailSender) int {
	fmt.Println("Processing Order via Simulation...")
	return 999
}
func main() {}
