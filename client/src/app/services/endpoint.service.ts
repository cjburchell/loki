import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {IEndpoint} from './contracts/Endpoint';

export abstract class IEndpointService{
  abstract GetEndpoints(): Promise<IEndpoint[]>;
  abstract GetEndpoint(id: string): Promise<IEndpoint>;
  abstract CreateEndpoint(endpoint: IEndpoint): Promise<any>;
  abstract UpdateEndpoint(endpoint: IEndpoint): Promise<any>;
  abstract DeleteEndpoint(id: string): Promise<any>;
}

@Injectable({
  providedIn: 'root'
})
export class EndpointService implements IEndpointService {
  constructor(private http: HttpClient) { }

  GetEndpoints(): Promise<IEndpoint[]> {
    return this.http.get<IEndpoint[]>(`/@mock/endpoint`).toPromise();
  }

  CreateEndpoint(endpoint: IEndpoint): Promise<any> {
    return this.http.post( `/@mock/endpoint`, endpoint).toPromise();
  }

  DeleteEndpoint(id: string): Promise<any> {
    return this.http.delete(`/@mock/endpoint/${id}`).toPromise();
  }

  GetEndpoint(id: string): Promise<IEndpoint> {
    return this.http.get<IEndpoint>(`/@mock/endpoint/${id}`).toPromise();
  }

  UpdateEndpoint(endpoint: IEndpoint): Promise<any> {
    return this.http.put( `/@mock/endpoint/${endpoint.name}`, endpoint).toPromise();
  }
}
