import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { EndpointComponent } from './endpoint/endpoint.component';
import { EndpointsComponent } from './endpoints/endpoints.component';
import {EndpointService, IEndpointService} from './services/endpoint.service';
import {environment} from '../environments/environment';
import {MockDataService} from './services/mockdata.service';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { NavComponent } from './nav/nav.component';
import { SearchComponent } from './common/search/search.component';
import { FilterPipe } from './pipes/filter.pipe';
import {FormsModule} from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { SettingsComponent } from './settings/settings.component';

@NgModule({
  declarations: [
    AppComponent,
    EndpointComponent,
    EndpointsComponent,
    NavComponent,
    SearchComponent,
    FilterPipe,
    SettingsComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    NgbModule,
    FormsModule,
    FontAwesomeModule
  ],
  providers: [
    { provide: IEndpointService, useClass: !environment.mockData ? EndpointService : MockDataService }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
