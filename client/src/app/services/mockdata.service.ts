import {Injectable} from '@angular/core';
import {IEndpointService} from './endpoint.service';
import {IEndpoint} from './contracts/Endpoint';
import {ISettings} from './contracts/Settings';

export const mockEndpoints: IEndpoint[] = [
  {
    name: 'test text',
    path: '/test1/test',
    method: 'GET',
    string_body: 'this is a test',
    content_type: 'text/plain',
    response: 200,
    reply_delay: 0,
    headers: [ {
      key: 'key',
      value: 'value'
    }],
  },
  {
    name: 'test json',
    path: '/test2/test',
    method: 'GET',
    string_body: '{\'test\':true}',
    content_type: 'application/json',
    response: 200,
    reply_delay: 0,
    headers: [],
  },
  {
    name: 'test html',
    path: '/test3/test',
    method: 'GET',
    string_body: '<h1>this is a test</h1>',
    content_type: 'text/html',
    response: 200,
    reply_delay: 0,
    headers: [],
  }
];

export const settings: ISettings = {
  default_reply: 404,
  partial_mock_server_address: '',
};

@Injectable({
  providedIn: 'root'
})
export class MockDataService implements IEndpointService{

  constructor() { }

  async CreateEndpoint(endpoint: IEndpoint): Promise<void> {
    mockEndpoints.push(endpoint);
    console.log(`Create Endpoint: ${JSON.stringify(endpoint)}`);
    await new Promise<void>(resolve => resolve());
  }

  async DeleteEndpoint(id: string): Promise<void> {
    const endpoint = await this.GetEndpoint(id);
    const index = mockEndpoints.indexOf(endpoint, 0);
    if (index > -1) {
      mockEndpoints.splice(index, 1);
    }

    console.log(`Delete: ${id}`);

    return new Promise<any>(resolve => resolve('ok'));
  }

  GetEndpoint(id: string): Promise<IEndpoint> {
    return new Promise<IEndpoint>(resolve => resolve(mockEndpoints.find(item => item.name === id)));
  }

  GetEndpoints(): Promise<IEndpoint[]> {
    return new Promise<IEndpoint[]>(resolve =>  resolve(mockEndpoints));
  }

  async UpdateEndpoint(endpoint: IEndpoint): Promise<void> {

    console.log(`Update Endpoint: ${JSON.stringify(endpoint)}`);
    await new Promise<void>(resolve => resolve());
  }

  GetSettings(): Promise<ISettings> {
    return new Promise<ISettings>(resolve =>  resolve(settings));
  }

  // tslint:disable-next-line:no-shadowed-variable
  UpdateSettings(settings: ISettings): Promise<any> {
    console.log(`Update Settings: ${JSON.stringify(settings)}`);
    throw {statusText: 'error'};
    return new Promise<any>(resolve =>  resolve());
  }
}
