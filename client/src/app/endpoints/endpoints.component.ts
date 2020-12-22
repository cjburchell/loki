import { Component, OnInit } from '@angular/core';
import {IEndpoint} from '../services/contracts/Endpoint';
import {IEndpointService} from '../services/endpoint.service';


@Component({
  selector: 'app-endpoints',
  templateUrl: './endpoints.component.html',
  styleUrls: ['./endpoints.component.scss']
})
export class EndpointsComponent implements OnInit {

  endpoints: IEndpoint[] | undefined;
  searchText = '';

  constructor(private endpointService: IEndpointService) { }

  async ngOnInit(): Promise<any> {
    this.endpoints = await this.endpointService.GetEndpoints();
  }
}
