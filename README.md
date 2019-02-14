# Invoice Hyperledger Application  

## Development Environment:  
**Go version:** go version go1.11.5 linux/amd64  
**Operating System**: Ubuntu 18.04 LTS  
  
## Setup Network
**Step 1:** From **blockchain-training-labs** directory, create a new folder named **invoice**. If not cloned yet, get it from https://github.com/hyperledger/fabric-samples.git  
**Step 2:** Open terminal and clone the repository using the command  
```git clone https://github.com/bien06/invoice.git```  

Copy the contents of **node** folder to **invoice** folder. Your **invoice** folder should contain the following files:  
- invoice
    - app.js
    - enrollAdmin.js
    - package.json
    - registerUser.js
    - startFabric.sh  
<br>

**Step 3:** Create another invoice folder inside chaincode folder. Copy the **go** folder to **invoice** folder. Your go file should be located in the following path:  
- chaincode  
  - invoice
    - go
      - invoice.go

**Step 4:** Go back to your terminal and go to the **node** folder. Run the following commands:  
```./startFabric.sh```  
```npm install```  
```node enrollAdmin.js```  
```node registerUser.js```  
```node app.js```  

**Step 5:** Open **Postman** or **Insomnia** and test the application.  
Create the following requests:  
  
**Request Name:** Raise Invoice  
**Method:** POST  
**URL:** localhost:3000/invoice  
**Arguments**: username  
  
**Request Name:** Display All Invoice  
**Method:** GET  
**URL:** localhost:3000  
**Arguments**: 
- username  
- invoiceid  
- invoicenumber  
- billedto  
- invoicedate  
- invoiceamount  
- itemdescription  
- gr  
- ispaid  
- paidamount  
- isrepaid  
- repaymentamount  
  
**Request Name:** Goods Received  
**Method:** PUT  
**URL:** localhost:3000/invoice  
**Arguments**:  
- invoiceid  
- gr  
- username  
  
**Request Name:** Paid to Supplier  
**Method:** PUT  
**URL:** localhost:3000/invoice  
**Arguments**:  
- invoiceid  
- ispaid  
- username 
  
**Request Name:** Repaid to Bank  
**Method:** PUT  
**URL:** localhost:3000/invoice  
**Arguments**:  
- invoiceid  
- isrepaid  
- username  
  
**Request Name:** Audit History  
**Method:** GET  
**URL:** localhost:3000/  
**Arguments**: username  
<br/>
Make sure you use form URL encoded for the arguments
