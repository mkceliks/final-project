# Final Project of Property Finder Go Bootcamp

This project is a final project for the property finder company.

With this project you can;

- Show all the products  -->  ```http://localhost:8080/products  ( GET )```
- Show all the cart items  -->  ```http://localhost:8080/cart  ( GET )```
- Show all the customers  --> ```http://localhost:8080/customers  ( GET )```
- Show the product by using productid  --> ```http://localhost:8080/getProductById/{id}  ( GET )```
- Add product to the system  --> ```http://localhost:8080/addProduct  ( POST )```
- Add customer to the system  --> ```http://localhost:8080/addCustomer  ( POST )```
- Add products into the cart with quantity  --> ```http://localhost:8080/addToCart/{id}  ( POST )```
- Add one item into the cart  --> ```http://localhost:8080/addOneItemToCart/{id}  ( POST )```
- Delete an item from cart  --> ```http://localhost:8080/deleteOneItemFromCart/{id}  ( DELETE )```
- Delete the row from cart  --> ```http://localhost:8080/deleteAllCart/{id}  ( DELETE )```

In this project i used;

- Go ( gorilla/mux, rs/cors, database/sql, lib/pq )
- PostgreSQL database
- Angular ( Bootstrap, fontawesome )


# 1.0 Go

## 1.1 Install Dependencies ##
- ```go get -u github.com/gorilla/mux```
- ```go get github.com/rs/cors```
- ```go get github.com/lib/pq```

## 1.2 Run

- ```go run main.go```

## 1.3 Files
### 1.3.1 main.go File

This file includes the APIs.

### 1.3.2 service/service.go File

This file contains the functions that APIs run in the background.

### 1.3.3 service/business.go File

This file includes the some business logics of the system. ( Using this logics in the service.go file. )

### 1.3.4 models Folder

This folder includes the cart,customer,order_details,product struct models to modelize the tables of the database at the db connection.

### 1.3.5 db.go File

This file starts the connection with PostgreSQL and generates a global variable to use this connection in other files. And includes the constants needs for connection.

# 2.0 PostgreSQL

 I chose to use PostgreSQL as database in this project because I know Property Finder company uses PostgreSQL.
 
### 2.1 Tables;
 
- cart
- customer
- products
- order_details

 **2.1.1 Table of cart**
 
 | id     | product_id      | quantity   | total_price | total_discount |
| ------------- | ------------- | --------    |--------|--------|
| `PRIMARY KEY`        | `REFERENCE TO products.id`         | INTEGER   |NUMERIC(8,2)|NUMERIC(8,2)|


**2.1.2 Table of products**

 | id     | name      | description   | price | tax |
| ------------- | ------------- | --------    |--------|--------|
| `PRIMARY KEY`        | STRING         | STRING   |NUMERIC(8,2)|NUMERIC(8,2)|

**2.1.3 Table of order_details**

 | id     | customer_id      | total_price   |
| ------------- | ------------- | --------    
| `PRIMARY KEY`        | `REFERENCE TO customer.id`          | NUMERIC(8,2)   |

**2.1.4 Table of customer**

 | id     | username      |
| ------------- | ------------- |    
| `PRIMARY KEY`        | STRING    |

# 3.0 Angular

This is a front-end side of the project that integrated APIs with services,models and components.

## 3.1 Install Dependencies ##
- ```npm install -g @angular/cli```
- ```ng add @ng-bootstrap/ng-bootstrap```
- ```npm i @fortawesome/fontawesome-free```

## 3.2 Run

- ```ng serve --open```
