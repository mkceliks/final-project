import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ProductComponent } from './components/product/product.component';
import { AddProductComponent } from './components/add-product/add-product.component';
import { CartComponent } from './components/cart/cart.component';
import { AddCustomerComponent } from './components/add-customer/add-customer.component';

const routes: Routes = [
  {path:"",component:ProductComponent},
  {path:"add-product",component:AddProductComponent},
  {path:"add-customer",component:AddCustomerComponent},
  {path:"cart",component:CartComponent},
  {path:"add-one-item/:productAddId",component:ProductComponent},
  {path:"delete-one-item/:productDeleteId",component:ProductComponent},
  {path:"delete-cart-row/:rowDelete",component:CartComponent},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
