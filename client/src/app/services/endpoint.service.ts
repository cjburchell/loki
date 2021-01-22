import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {IEndpoint} from './contracts/Endpoint';
import {ISettings} from './contracts/Settings';

export abstract class IEndpointService{
  abstract GetEndpoints(): Promise<IEndpoint[]>;
  abstract GetEndpoint(id: string): Promise<IEndpoint>;
  abstract CreateEndpoint(endpoint: IEndpoint): Promise<any>;
  abstract UpdateEndpoint(endpoint: IEndpoint): Promise<any>;
  abstract DeleteEndpoint(id: string): Promise<any>;

  abstract GetSettings(): Promise<ISettings>;
  abstract UpdateSettings(settings: ISettings): Promise<any>;
}

@Injectable({
  providedIn: 'root'
})
export class EndpointService implements IEndpointService {
  constructor(private http: HttpClient) { }

  GetEndpoints(): Promise<IEndpoint[]> {
    return this.http.get<IEndpoint[]>(`/@mock/endpoint`).toPromise();
  }

  async CreateEndpoint(endpoint: IEndpoint): Promise<any> {
    return this.http.post<any>(`/@mock/endpoint`, endpoint).toPromise();
  }

  async DeleteEndpoint(id: string): Promise<any> {
    return this.http.delete<any>(`/@mock/endpoint/${id}`).toPromise();
  }

  GetEndpoint(id: string): Promise<IEndpoint> {
    return this.http.get<IEndpoint>(`/@mock/endpoint/${id}`).toPromise();
  }

  async UpdateEndpoint(endpoint: IEndpoint): Promise<any> {
    await this.http.put<any>(`/@mock/endpoint/${endpoint.name}`, endpoint).toPromise();
  }

  GetSettings(): Promise<ISettings> {
    return this.http.get<ISettings>(`/@mock/settings`).toPromise();
  }

  UpdateSettings(settings: ISettings): Promise<any> {
    return this.http.put<any>(`/@mock/settings`, settings).toPromise();
  }
}
