<?php
// Horrible legacy DB connection using deprecated mysql_connect
$conn = mysql_connect("localhost", "root", "");
if (!$conn) {
    die("Could not connect: " . mysql_error());
}
mysql_select_db("old_ecommerce", $conn);
?>
