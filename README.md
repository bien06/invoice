# Invoice Hyperledger Application

## Development Environment:
**Go version:** go version go1.11.5 linux/amd64
**Operating System**: Ubuntu 18.04 LTS

## Setup Network
**Step 1:** Find a good location and create a new folder named **invoice**. 
**Step 2:** Open terminal and clone the repository using the command 
```git clone https://github.com/bien06/invoice.git```

Your **invoice** folder should contain the following files:
- invoice
  - node
    - app.js
    - enrollAdmin.js
    - package.json
    - registerUser.js
    - startFabric.sh
  - go
    - invoice.go

**Step 3:** Go back to your terminal and go to the **node** folder. Run the following commands:__
```./startFabric.sh```__
```npm install```__
```node enrollAdmin.js```__
```node registerUser.js```__
```node app.js```__

**Step 4:** Open **Postman** or **Insomnia** and test the application.
Create the following requests:

**Request Name:** Raise Invoice
**Method:** POST
**URL:** localhost:3000/invoice

**Request Name:** Display All Invoice
**Method:** GET
**URL:** localhost:3000

**Request Name:** Goods Received
**Method:** PUT
**URL:** localhost:3000/invoice

**Request Name:** Paid to Supplier
**Method:** PUT
**URL:** localhost:3000/invoice

**Request Name:** Repaid to Bank
**Method:** PUT
**URL:** localhost:3000/invoice
 

