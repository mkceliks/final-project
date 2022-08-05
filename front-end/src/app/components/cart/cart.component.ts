import { Component, OnInit } from '@angular/core';
import { CartService } from 'src/app/services/cart.service';
import { Cart } from 'src/app/models/cart';

@Component({
  selector: 'app-cart',
  templateUrl: './cart.component.html',
  styleUrls: ['./cart.component.css']
})
export class CartComponent implements OnInit {

  cartItems : Cart[] = [];

  constructor(private cartService:CartService) { }

  ngOnInit(): void {
    this.getCartItems();
  }


  getCartItems(){
    this.cartService.getAllCartItems().subscribe(data => {
      this.cartItems = data;
    }
    );
  }

}
