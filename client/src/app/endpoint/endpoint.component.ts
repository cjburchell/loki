import { Component, OnInit } from '@angular/core';
import {IEndpointService} from '../services/endpoint.service';
import {ActivatedRoute, Router} from '@angular/router';
import {IEndpoint, IHeader} from '../services/contracts/Endpoint';

@Component({
  selector: 'app-endpoint',
  templateUrl: './endpoint.component.html',
  styleUrls: ['./endpoint.component.scss']
})
export class EndpointComponent implements OnInit {
  public errorMsg = '';
  constructor(private activatedRoute: ActivatedRoute,
              private endpointService: IEndpointService,
              private router: Router) { }

  public endpoint: IEndpoint | undefined;

  public  create = false;

  ngOnInit(): void {
    this.activatedRoute.params.subscribe(async params => {
      if (params.endpointId === undefined){
        this.endpoint = {content_type: '', headers: [], method: 'GET', name: '', path: '', reply_delay: 0, response: 200, string_body: ''};
        this.create = true;
      }
      else {
        try {
          this.endpoint = await this.endpointService.GetEndpoint(params.endpointId);
          this.create = false;
        }catch (e) {
          this.errorMsg = `Unable to display endpoint ${params.endpointId}: ${e.status} (${e.statusText})`;
        }
      }
    });
  }

  async Save(): Promise<any> {
    if (this.endpoint && !this.create) {
      try {
          await this.endpointService.UpdateEndpoint(this.endpoint);
          this.router.navigate([`/endpoints`]).then(() => {
        });
      } catch (e) {
        this.errorMsg = `Unable to save endpoint: ${e.status} (${e.statusText})`;
      }
    }
  }

  async Create(): Promise<any> {
    if (this.endpoint && this.create) {
      try {
          await this.endpointService.CreateEndpoint(this.endpoint);
          this.router.navigate([`/endpoints`]).then(() => {
        });
      } catch (e) {
        this.errorMsg = `Unable to create endpoint: ${e.status} (${e.statusText})`;
      }
    }
  }

  async Delete(): Promise<any> {
    if (this.endpoint && !this.create) {
      try {
        await this.endpointService.DeleteEndpoint(this.endpoint.name);
        this.router.navigate([`/endpoints`]).then( () => {});
      }catch (e) {
        this.errorMsg = `Unable to delete endpoint: ${e.status} (${e.statusText})`;
      }
    }
  }

  public AddHeaderItem(): void {
    if (this.endpoint) {
      this.endpoint.headers.push({key: '', value: ''});
    }
  }

  public DeleteHeaderItem(header: IHeader): void {
    if (this.endpoint && this.endpoint.headers) {
      const index = this.endpoint.headers.indexOf(header, 0);
      if ( index === 0 && this.endpoint.headers.length === 1){
        this.endpoint.headers = [];
      }
      else if (index > -1) {
        this.endpoint.headers.splice(index, 1);
      }
    }
  }
}
