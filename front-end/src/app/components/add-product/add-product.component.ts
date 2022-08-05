import { Component, OnInit } from '@angular/core';
import {
  FormGroup,
  FormBuilder,
  Validators,
  FormControl,
} from '@angular/forms';
import { ProductService } from 'src/app/services/product.service';

@Component({
  selector: 'app-add-product',
  templateUrl: './add-product.component.html',
  styleUrls: ['./add-product.component.css']
})
export class AddProductComponent implements OnInit {

  productAddForm : FormGroup;

  constructor(private formBuilder: FormBuilder,private productService:ProductService) { }

  ngOnInit(): void {
    this.createProductAddForm();
  }

  createProductAddForm(){
    this.productAddForm = this.formBuilder.group({
      Name: ['', Validators.required],
      Description: ['', Validators.required],
      Price: ['', Validators.required],
      Discount: ['', Validators.required],
      Tax: ['', Validators.required],
    });
  }

  addProduct(){
    if (this.productAddForm.valid) {
      let productModel = Object.assign({}, this.productAddForm.value);
      this.productService.addProduct(productModel).subscribe(data => {
        console.log(data);
      });
    } else {
      console.log('Form is not valid');
    }
}
}
