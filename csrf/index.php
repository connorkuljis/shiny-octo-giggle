<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>XSRF - Your trusted bank site</title>
</head>
<body style="font-family: sans-serif;">

<?php
  function getDatabaseConnection() {
    $servername = "localhost";
    $username = "admin";
    $password = "admin";
    $dbName = "testDB";

    $conn = new mysqli($servername, $username, $password, $dbName);

    if ($conn->connect_error) {
        die("Connection failed: " . $conn->connect_error);
    }

    return $conn;
  }

  mysqli_report(MYSQLI_REPORT_ERROR | MYSQLI_REPORT_STRICT);

    try {
      $conn = getDatabaseConnection();
      $query = "SELECT * FROM BankAccounts";
      echo "Query: " . $query;

      $result = $conn->query($query);

      if ($result->num_rows > 0) {
          echo '<br><br><table>';
          echo '<td style="width: 100px; height: 22px">' . "<b>Account id</b>" . '</td>';
          echo '<td style="width: 150px; height: 22px">' . "<b>User</b>" . '</td>';
          echo '<td style="width: 100px; height: 22px">' . "<b>Balance</b>" . '</td>';
          while($row = $result->fetch_assoc()) {
            echo '<tr>';
            echo '<td style="width: 100px; height: 18px">' . $row['account_id'] . '</td>';
            echo '<td style="width: 150px; height: 18px">' . $row['user_name'] . '</td>';
            echo '<td style="width: 100px; height: 18px">' . $row['account_balance'] . '</td>';
            echo '</tr>';
          }
          echo '</table>';
      } else {
          echo "<br><br>No results match your search:-(";
      }


      $conn->close($conn);
    } catch (Exception $e) {
      echo 'Error! ' + $e->getCode();
    }
  
?>

<h3> Account information page at your trusted bank</h3>
<form method="GET" action="<?php echo htmlspecialchars($_SERVER["PHP_SELF"]);?>">
   <br>
  <input type="submit" value="See all account details">
</form>
<br>
</body>
</html>
