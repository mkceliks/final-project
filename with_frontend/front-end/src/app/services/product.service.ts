import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Product } from '../models/product';

@Injectable({
  providedIn: 'root'
})
export class ProductService {

  apiUrl = 'http://localhost:8080/'

  constructor(private httpClient: HttpClient) { }

  getAllProducts(): Observable<Product[]> {

    let newPath = this.apiUrl + 'products';

    return this.httpClient.get<Product[]>(newPath);

  }

  getProductById(id: number): Observable<Product> {

    let newPath = this.apiUrl + 'getProductById/' + id;

    return this.httpClient.get<Product>(newPath);


  }

  addProduct(product: Product): Observable<Product> {
      
      let newPath = this.apiUrl + 'addProduct';
  
      return this.httpClient.post<Product>(newPath, product);
  
    }






}
