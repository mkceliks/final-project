import { Component, OnInit } from '@angular/core';
import { CartService } from 'src/app/services/cart.service';
import { Cart } from 'src/app/models/cart';
import { Product } from 'src/app/models/product';
import { ProductService } from 'src/app/services/product.service';
import { ActivatedRoute } from '@angular/router';
import { Router } from '@angular/router';
import { Customer } from 'src/app/models/customer';
import { CustomerService } from 'src/app/services/customer.service';
import {
  FormGroup,
  FormBuilder,
  Validators,
  FormControl,
} from '@angular/forms';

@Component({
  selector: 'app-cart',
  templateUrl: './cart.component.html',
  styleUrls: ['./cart.component.css']
})
export class CartComponent implements OnInit {

  orderForm : FormGroup;
  totalPrice:number;
  customers : Customer[] = [];
  cartItems : Cart[] = [];
  products : Product[] = [];

  constructor(private cartService:CartService,private productService:ProductService,private customerService:CustomerService,private activatedRoot:ActivatedRoute,private router:Router,private formBuilder: FormBuilder,) { }

  ngOnInit(): void {

  

    this.getCartItems();
    this.getProducts();
    this.getCustomers();
    this.createOrderCompleteForm();

    this.activatedRoot.params.subscribe(params => {

      if(params["rowDelete"]){
        this.deleteCartRow(params["rowDelete"]);
      }else if(params["productDeleteId"]){
        this.deleteCartItem(params["productDeleteId"]);
      }else if(params["productAddId"]){
        this.addToCart(params["productAddId"]);
      }
    })
  }

  createOrderCompleteForm(){
    this.orderForm = this.formBuilder.group({
      CustomerID: ['', Validators.required],
      TotalPrice: ['', Validators.required],
    });
  }


  getCustomers(){
    this.customerService.getAllCustomers().subscribe(data => {
      this.customers = data;
    });
  }

  getCartItems(){
    this.cartService.getAllCartItems().subscribe(data => {
      this.cartItems = data;
      this.totalPrice = 0;
      this.cartItems.forEach(element => {
        this.totalPrice += element.Price;
      });
      console.log(this.totalPrice);
    }
    );

  

  }

  getProducts(){
    this.productService.getAllProducts().subscribe(data => {
      this.products = data;
    }
    );
  }

  getProductById(productId:number){
    this.productService.getProductById(productId).subscribe(data => {
      return data;
    });
  }

  addToCart(productAddId:number){
    this.cartService.addCartItem(productAddId).subscribe((data) => {
      this.router.navigate(['cart']);
    });
    alert("Product added to cart");
    
}

  deleteCartItem(productDeleteId:number){
    this.cartService.deleteCartItem(productDeleteId).subscribe((data) => {
        });
        alert("Product deleted from cart");
        this.router.navigate(['cart']);
  }

  deleteCartRow(cartId:number){
    this.cartService.deleteCartRow(cartId).subscribe(data => {
    });
    alert("Product deleted from cart");
    this.router.navigate(['cart']);
  }

  orderComplete(){
    if (this.orderForm.valid) {
      let orderModel = Object.assign({}, this.orderForm.value);
      orderModel.CustomerID = Number(orderModel.CustomerID);
      orderModel.TotalPrice = Number(orderModel.TotalPrice);
      console.log(orderModel);
      this.cartService.orderComplete(orderModel).subscribe(data => {
        alert("Order completed successfully");
      });
      alert("Purchase successfully completed");
    } else {
      console.log('Form is not valid');
    }
  }

}
