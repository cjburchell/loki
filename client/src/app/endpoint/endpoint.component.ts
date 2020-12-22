import { Component, OnInit } from '@angular/core';
import {IEndpointService} from '../services/endpoint.service';
import {ActivatedRoute, Router} from '@angular/router';
import {IEndpoint} from '../services/contracts/Endpoint';

@Component({
  selector: 'app-endpoint',
  templateUrl: './endpoint.component.html',
  styleUrls: ['./endpoint.component.scss']
})
export class EndpointComponent implements OnInit {

  constructor(private activatedRoute: ActivatedRoute,
              private endpointService: IEndpointService,
              private router: Router) { }

  public endpoint: IEndpoint | undefined;

  ngOnInit(): void {
    this.activatedRoute.params.subscribe(async params => {
      this.endpoint = await this.endpointService.GetEndpoint(params.endpointId);
    });
  }

  async Save(): Promise<any> {
    if (this.endpoint) {
      await this.endpointService.UpdateEndpoint(this.endpoint);
      await this.router.navigate([`/endpoints`]);
    }
  }

  async Delete(): Promise<any> {
    if (this.endpoint) {
      await this.endpointService.DeleteEndpoint(this.endpoint.name);
      await this.router.navigate([`/endpoints`]);
    }
  }
}
