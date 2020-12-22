import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {IEndpoint} from './contracts/Endpoint';
import {ISettings} from './contracts/Settings';

export abstract class IEndpointService{
  abstract GetEndpoints(): Promise<IEndpoint[]>;
  abstract GetEndpoint(id: string): Promise<IEndpoint>;
  abstract CreateEndpoint(endpoint: IEndpoint): Promise<void>;
  abstract UpdateEndpoint(endpoint: IEndpoint): Promise<void>;
  abstract DeleteEndpoint(id: string): Promise<void>;

  abstract GetSettings(): Promise<ISettings>;
  abstract UpdateSettings(settings: ISettings): Promise<void>;
}

@Injectable({
  providedIn: 'root'
})
export class EndpointService implements IEndpointService {
  constructor(private http: HttpClient) { }

  GetEndpoints(): Promise<IEndpoint[]> {
    return this.http.get<IEndpoint[]>(`/@mock/endpoint`).toPromise();
  }

  async CreateEndpoint(endpoint: IEndpoint): Promise<void> {
    await this.http.post(`/@mock/endpoint`, endpoint).toPromise();
  }

  async DeleteEndpoint(id: string): Promise<void> {
    await this.http.delete(`/@mock/endpoint/${id}`).toPromise();
  }

  GetEndpoint(id: string): Promise<IEndpoint> {
    return this.http.get<IEndpoint>(`/@mock/endpoint/${id}`).toPromise();
  }

  async UpdateEndpoint(endpoint: IEndpoint): Promise<void> {
    await this.http.put(`/@mock/endpoint/${endpoint.name}`, endpoint).toPromise();
  }

  GetSettings(): Promise<ISettings> {
    return this.http.get<ISettings>(`/@mock/settings`).toPromise();
  }

  async UpdateSettings(settings: ISettings): Promise<void> {
    await this.http.put(`/@mock/settings`, settings).toPromise();
  }
}
