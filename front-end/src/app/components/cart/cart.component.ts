import { Component, OnInit } from '@angular/core';
import { CartService } from 'src/app/services/cart.service';
import { Cart } from 'src/app/models/cart';
import { Product } from 'src/app/models/product';
import { ProductService } from 'src/app/services/product.service';
import { find } from 'rxjs';

@Component({
  selector: 'app-cart',
  templateUrl: './cart.component.html',
  styleUrls: ['./cart.component.css']
})
export class CartComponent implements OnInit {

  cartItems : Cart[] = [];
  products : Product[] = [];

  constructor(private cartService:CartService,private productService:ProductService) { }

  ngOnInit(): void {
    this.getCartItems();
    this.getProducts();
  }


  getCartItems(){
    this.cartService.getAllCartItems().subscribe(data => {
      this.cartItems = data;
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

}
