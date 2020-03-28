# ShopifyTools

## Description

ShopifyTools is a HTTP server application written in Go. The purpose of the application is to provide API endpoints for useful features that Shopify doesn't include out of the box.

## Installation

To allow a Shopify client to use these API endpoints they must be added to the config file. The config file can be found on the Server path (home/ubuntu/shopifytools) and is named **conf.json**. New Shopify clients can be easily added by adding a record to the JSON contents of the **conf.json** file where the shop url (example.myshopify.com) is the object key and the API password is the value.

## Routes

### POST /validate-discount
**Returns `bool`**

Discount data in the form of JSON is sent to the route in the request body and the server will respond with a boolean stating if the code sent in the request body is a valid code or not.

This route accepts the following paramters:

| Parameter                 | Description                                                                                              | Example                                                 |
| ------------------------- | -------------------------------------------------------------------------------------------------------- | -----------------------------------------------------   |
| shop  *required*          | URL of the shopify store. (Must match the shop url added to the **conf.json** file) *string*             | example.myshopify.com                                   |
| discount_code  *required* | The discount code to be validated. *string*                                                              | SUMMER20OFF                                             |
| variant_ids               | Array of variant ID's to be tested as discount codes may only apply to specific variants. *[]int*        | `[15378478301229, 14642142576331]`                      |
| product_ids               | Array of product ID's to be tested as discount codes may only apply to specific products. *[]int*        | `[1836365643821, 4506936016941]`                        |
| collection_ids            | Array of collection ID's to be tested as discount codes may only apply to specific collectionss. *[]int* | `[81355931693, 81354653741, 80957046829, 153633652781]` |


A liquid snippet has been created already, which can be found [here](./theme/snippets/cart-discount-code.liquid), for easy integration. Add that file to your snippets folder of your shopify theme and then add the snippet to your cart.liquid template by adding `{% render 'cart-discount-code', cart_items: cart.items %}`. This snippet will display a discount code form along with Javascript that makes a POST HTTP request to this endpoint with the neccessary data when the form is submitted. 

## Testing

To build Go HTTP application to test locally you must use the following command while in the root directory of this Repository:
```
go build -o ShopifyTools -v ./src
```

Run the program easily by refering to the compiled binary. While in the directory of the Binary
```
./ShopifyTools
```

You can also build & run the Go HTTP application by chaining to two commands above like so:
```
go build -o ShopifyTools -v ./src && ./ShopifyTools
```

The Go HTTP application will run on Port 8080. You can then make requests using a request client such as cURL or Postman to http://localhost:8080.

#### ngrok

You could also use a tool called [ngrok](https://ngrok.com/download) which can create a publicly accessible url to a local server with an SSL certificate. This will enable you to test requests on your Shopify theme. 
 
Use ngrok by running the following command in the directory that ngrok is located while your ShopifyTools is running locally.

```
./ngrok http 8080
```

## Deployment

If any code changes are ready to be deployed the Go application must be compiled to run on Linux with a AMD64 architecure which can be acheived by running the following command:
```
env GOOS=linux GOARCH=amd64 go build -o ShopifyTools -v ./src
```

This will compile the contents of the **src** directory to a Binary named ShopifyTools.

The output binary file can then be uploaded to a Linux server where you can run it as a systemd service