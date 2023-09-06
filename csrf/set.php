<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>XSRF - Withdraw funds from trusted bank</title>
</head>
<body style="font-family: sans-serif;">

<?php
  function getDatabaseConnection() {
    // All the code in this function is just to set up an
  $servername = "localhost";
  $username = "admin";
  $password = "admin";
  $dbName = "testDB";

    $conn = new mysqli($servername, $username, $password, $dbName);

    if ($conn->connect_error) {
        die("Connection failed: " . $conn->connect_error);
    }

    // Comment out if the table does not exist.
    $drop_query = "DROP TABLE BankAccounts";
    if ($conn->query($drop_query) !== TRUE) {
        echo "Error dropping table: " . $conn->error;
    }

    $create_query = <<<_CREATE_
        CREATE TABLE IF NOT EXISTS BankAccounts (
          account_id INT(10) NOT NULL PRIMARY KEY AUTO_INCREMENT,
          user_name VARCHAR(100) NOT NULL,
          account_balance INT(100) NOT NULL
        )
_CREATE_;

    $result = $conn->query($create_query);
    if($result === FALSE) {
        die("Create failed: " . mysql_error());
    }

  $insert_queries = array();
  $insert_queries[] = <<<_INSERT_
      INSERT INTO BankAccounts
      VALUES
        (1234, 'Redowan', 10000)
_INSERT_;
  $insert_queries[] = <<<_INSERT_
      INSERT INTO BankAccounts
      VALUES
        (1111, 'Andrew', 10000)
_INSERT_;
  $insert_queries[] = <<<_INSERT_
      INSERT INTO BankAccounts
      VALUES
        (2222, 'Jordan', 10000)
_INSERT_;
  $insert_queries[] = <<<_INSERT_
      INSERT INTO BankAccounts
      VALUES
        (3333, 'Peter', 10000)
_INSERT_;

    foreach ($insert_queries as $query) {
      if ($conn->query($query)) {
          //echo "Inserted row: $query\n";
      } else {
          echo "Error inserting row: " . mysqli_error($conn) . "\n";
      }
    }
    return $conn;
  }

  mysqli_report(MYSQLI_REPORT_ERROR | MYSQLI_REPORT_STRICT);

  try {
    $conn = getDatabaseConnection();

    //echo "Transferring funds <br>";
    $from_id = $conn->real_escape_string($_GET['from_id']);
    $to_id = $conn->real_escape_string($_GET['to_id']);
    $amount = $conn->real_escape_string($_GET['amount']);
    

    if (!empty($from_id) && !empty($to_id) && !empty($amount)) {
      $subtract_query =
        "UPDATE BankAccounts SET account_balance = account_balance - $amount where account_id = $from_id";
      echo $subtract_query;
      $conn->query($subtract_query);

      $add_query =
        "UPDATE BankAccounts SET account_balance = account_balance + $amount where account_id = $to_id";
      echo $add_query;
      $conn->query($add_query);
      echo "From: $from_id, To: $to_id, Amount: $amount <br>";

      $conn->close();
    }
  } catch (Exception $e) {
    echo 'Error! ' + $e->getCode();
  }
?>


<h3> Account information page at your trusted bank</h3>
<form method="GET" action="<?php echo htmlspecialchars($_SERVER["PHP_SELF"]);?>">
  <p>Transfer Money</p>
  <br>
  <input type="text" name="from_id" value="1234">
  <br>
  <input type="text" name="to_id" placeholder="Receiver Account ID">
  <br>
  <input type="number" name="amount" placeholder="amount">
  <br>
  <input type="submit" value="send">
</form>
<br>
</body>
</html>
