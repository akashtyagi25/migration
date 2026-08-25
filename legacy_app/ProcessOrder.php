<?php
include 'db_connect.php';
include 'Emailer.php';

// Bad Practice: Using global connection inside a function
function processOrder($userId, $cartItems) {
    global $conn;
    
    $total = 0;
    foreach ($cartItems as $item) {
        // SQL Injection vulnerability!
        $sql = "SELECT price FROM products WHERE id = " . $item['id'];
        $result = mysql_query($sql, $conn);
        $row = mysql_fetch_assoc($result);
        $total += $row['price'] * $item['qty'];
    }

    // Insert order
    $insertSql = "INSERT INTO orders (user_id, total) VALUES ($userId, $total)";
    mysql_query($insertSql, $conn);
    $orderId = mysql_insert_id($conn);

    // Get user email
    $userSql = "SELECT email FROM users WHERE id = $userId";
    $userResult = mysql_query($userSql, $conn);
    $userRow = mysql_fetch_assoc($userResult);

    // Side effect: Send email
    Emailer::sendReceipt($userRow['email'], $orderId, $total);

    return $orderId;
}
?>
