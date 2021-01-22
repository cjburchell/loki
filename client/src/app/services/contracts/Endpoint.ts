
export interface IEndpoint {
  name: string;
  path: string;
  method: string;
  string_body: string;
  content_type: string;
  response: number;
  reply_delay: number;
  headers: IHeader[];
}

export interface IHeader {
  key: string;
  value: string;
}

