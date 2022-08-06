import { Component, OnInit } from '@angular/core';
import { Product } from '../../models/product';
import { ProductService } from '../../services/product.service';
import { CartService } from 'src/app/services/cart.service';
import { ActivatedRoute } from '@angular/router';
@Component({
  selector: 'app-product',
  templateUrl: './product.component.html',
  styleUrls: ['./product.component.css']
})
export class ProductComponent implements OnInit {

  products : Product[] = [];

  constructor(private productService:ProductService,private cartService:CartService,private activatedRoot:ActivatedRoute) { }

  ngOnInit(): void {
    this.getProducts();

    this.activatedRoot.params.subscribe(params => {

      if(params["productAddId"]){
        this.addToCart(params["productAddId"]);
      }else if(params["productDeleteId"]){
        this.deleteCartItem(params["productDeleteId"]);
      }
      
    })
    
  }

  getProducts(){
    this.productService.getAllProducts().subscribe(data => {
      this.products = data;
    })
  }

  addToCart(productAddId:number){
    this.cartService.addCartItem(productAddId).subscribe((data) => {
    alert("Product added to cart");
    });
}

  deleteCartItem(productDeleteId:number){
    this.cartService.deleteCartItem(productDeleteId).subscribe((data) => {
      alert("Product deleted from cart");
        });
  }
}
