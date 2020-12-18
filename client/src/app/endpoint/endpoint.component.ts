import { Component, OnInit } from '@angular/core';
import {IEndpointService} from '../services/endpoint.service';
import {ActivatedRoute} from '@angular/router';
import {IEndpoint} from '../services/contracts/Endpoint';

@Component({
  selector: 'app-endpoint',
  templateUrl: './endpoint.component.html',
  styleUrls: ['./endpoint.component.scss']
})
export class EndpointComponent implements OnInit {

  constructor(private activatedRoute: ActivatedRoute,  private endpointService: IEndpointService) { }

  public endpoint: IEndpoint | undefined;

  ngOnInit(): void {
    this.activatedRoute.params.subscribe(async params => {
      await this.endpointService.GetEndpoint(params.endpointId);
    });
  }

}
