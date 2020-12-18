import { Injectable } from '@angular/core';
import {IEndpointService} from './endpoint.service';
import {IEndpoint} from './contracts/Endpoint';

const mockEndpoints: IEndpoint[] = [
  {
    name: 'test0',
    path: '/test1/test',
    method: 'get',
    string_body: 'test',
    content_type: 'text',
    response: 200,
    reply_delay: 0
  },
  {
    name: 'test1',
    path: '/test1/test',
    method: 'get',
    string_body: 'test',
    content_type: 'text',
    response: 200,
    reply_delay: 0
  },
  {
    name: 'test2',
    path: '/test1/test',
    method: 'get',
    string_body: 'test',
    content_type: 'text',
    response: 200,
    reply_delay: 0
  }
];

@Injectable({
  providedIn: 'root'
})
export class MockDataService implements IEndpointService{

  constructor() { }

  CreateEndpoint(endpoint: IEndpoint): Promise<any> {
    mockEndpoints.push(endpoint);
    return new Promise<string>(resolve => resolve('test'));
  }

  async DeleteEndpoint(id: string): Promise<any> {
    const endpoint = await this.GetEndpoint(id);
    const index = mockEndpoints.indexOf(endpoint, 0);
    if (index > -1) {
      mockEndpoints.splice(index, 1);
    }

    return new Promise<any>(resolve => resolve('ok'));
  }

  GetEndpoint(id: string): Promise<IEndpoint> {
    return new Promise<IEndpoint>(resolve => resolve(mockEndpoints[0]));
  }

  GetEndpoints(): Promise<IEndpoint[]> {
    return new Promise<IEndpoint[]>(resolve =>  resolve(mockEndpoints));
  }

  UpdateEndpoint(endpoint: IEndpoint): Promise<any> {
    return new Promise<string>(resolve => resolve('test'));
  }
}
