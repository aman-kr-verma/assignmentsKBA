# Name of Project - `Electronics Supply Chain`

## Steps to install/setup the application
    
### 1. Go to electronics_network directory
```
cd electronics_network
```

### 2. Run startNetwork.sh
```
./startNetwork.sh
```

This script setups the necessary components, like ca-admins, peers, etc, packages, installs, approves and commit the chaincode to each of the individual peers. 

Since the network is up and chaincode is installed, we can invoke the chaincode and run the transactions.

### 3. Go to Electronics-App
```
cd ../Electronics-App
```

### 4. Run the `main.go` file
```
go run .
```

### 5. Commands to run transactions on chaincode

Open Postman to hit the apis for performing transactions on chaincode.

For eg, to create an instance of electronic item, hit the following api with appropriate body.
```
http://localhost:8080/api/electronicItem
```

### 6. Stop the network

```
cd ../electronics_network/
```
```
./stopNetwork.sh
```




