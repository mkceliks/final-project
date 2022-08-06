import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Cart } from '../models/cart';
@Injectable({
  providedIn: 'root'
})
export class CartService {

  apiUrl = 'http://localhost:8080/'

  constructor(private httpClient: HttpClient) { }


  getAllCartItems(): Observable<Cart[]> {
      
      let newPath = this.apiUrl + 'cart';
  
      return this.httpClient.get<Cart[]>(newPath);
  
    }

    addCartItem(id:number): Observable<Cart> {
        
        let newPath = this.apiUrl + 'addOneItemToCart/' + id;
    
        return this.httpClient.post<Cart>(newPath,id);
    
      }

      deleteCartItem(id:number): Observable<Cart> {

        let newPath = this.apiUrl + 'deleteOneItemFromCart/' + id;

        return this.httpClient.delete<Cart>(newPath);

      }

}
