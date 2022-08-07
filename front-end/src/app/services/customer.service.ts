import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Customer } from '../models/customer';

@Injectable({
  providedIn: 'root'
})
export class CustomerService {

  apiUrl = 'http://localhost:8080/'

  constructor(private httpClient: HttpClient) { }


  addCustomer(username:Customer): Observable<Customer> {
      
    let newPath = this.apiUrl + 'addCustomer';

    return this.httpClient.post<Customer>(newPath, username);

  }

  getAllCustomers(): Observable<Customer[]> {
        
      let newPath = this.apiUrl + 'customers';
  
      return this.httpClient.get<Customer[]>(newPath);
  
    }


}
