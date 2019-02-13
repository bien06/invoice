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

**Step 3:** Go back to your terminal and go to the **node** folder. Run the following commands:  
```./startFabric.sh```  
```npm install```  
```node enrollAdmin.js```  
```node registerUser.js```  
```node app.js```  

**Step 4:** Open **Postman** or **Insomnia** and test the application.  
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
