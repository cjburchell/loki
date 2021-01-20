import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {EndpointComponent} from './endpoint/endpoint.component';
import {EndpointsComponent} from './endpoints/endpoints.component';
import {SettingsComponent} from './settings/settings.component';

const routes: Routes = [
  {path: 'settings', component: SettingsComponent },
  {path: 'endpoints', component: EndpointsComponent },
  {path: 'endpoint/:endpointId', component: EndpointComponent },
  {path: 'endpoint', component: EndpointComponent },
  {path: '', pathMatch: 'full', redirectTo: 'endpoints'}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
