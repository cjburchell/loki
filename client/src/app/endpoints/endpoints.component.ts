import {Component, OnInit} from '@angular/core';
import {IEndpoint} from '../services/contracts/Endpoint';
import {IEndpointService} from '../services/endpoint.service';
import {Router} from '@angular/router';


@Component({
  selector: 'app-endpoints',
  templateUrl: './endpoints.component.html',
  styleUrls: ['./endpoints.component.scss']
})
export class EndpointsComponent implements OnInit {

  endpoints: IEndpoint[] | undefined;
  searchText = '';
  errorMsg = '';

  constructor(private endpointService: IEndpointService,
              private router: Router) { }

  async ngOnInit(): Promise<any> {
    try {
      this.endpoints = await this.endpointService.GetEndpoints();
    } catch (e) {
      this.errorMsg = `Error loading endpoints: ${e.status} (${e.statusText})`;
    }
  }

  public Show(endpoint: IEndpoint): void {
    this.router.navigate([`/endpoint/` + endpoint.name]).then(() => {});
  }

  public Create(): void {
    this.router.navigate([`/endpoint`]).then(() => {});
  }
}
