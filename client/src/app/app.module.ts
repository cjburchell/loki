import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { EndpointComponent } from './endpoint/endpoint.component';
import { EndpointsComponent } from './endpoints/endpoints.component';
import { EndpointItemComponent } from './endpoints/endpoint-item/endpoint-item.component';
import {EndpointService, IEndpointService} from './services/endpoint.service';
import {environment} from '../environments/environment';
import {MockDataService} from './services/mockdata.service';

@NgModule({
  declarations: [
    AppComponent,
    EndpointComponent,
    EndpointsComponent,
    EndpointItemComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [
    { provide: IEndpointService, useClass: !environment.mockData ? EndpointService : MockDataService }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
