import { Component, OnInit } from '@angular/core';
import {
  FormGroup,
  FormBuilder,
  Validators,
  FormControl,
} from '@angular/forms';
import { CustomerService } from 'src/app/services/customer.service';




@Component({
  selector: 'app-add-customer',
  templateUrl: './add-customer.component.html',
  styleUrls: ['./add-customer.component.css']
})
export class AddCustomerComponent implements OnInit {

  customerAddForm : FormGroup;

  constructor(private formBuilder: FormBuilder,private customerService:CustomerService) { }

  ngOnInit(): void {
    this.createCustomerAddForm();
  }

  createCustomerAddForm(){
    this.customerAddForm = this.formBuilder.group({
      UserName: ['', Validators.required]
    });
  }

  addCustomer(){
    if (this.customerAddForm.valid) {
      let customerModel = Object.assign({}, this.customerAddForm.value);
      this.customerService.addCustomer(customerModel).subscribe(data => {
        
      });
      alert("Customer added successfully");
    } else {
      console.log('Form is not valid');
    }
}

}
