<?php
class Emailer {
    public static function sendReceipt($email, $orderId, $total) {
        // Pretend this sends a real email via SMTP
        $subject = "Order #$orderId Receipt";
        $body = "Thank you for your order! Total: $$total";
        mail($email, $subject, $body);
        
        // Log it to a file
        file_put_contents("email_logs.txt", "Sent to $email\n", FILE_APPEND);
    }
}
?>
