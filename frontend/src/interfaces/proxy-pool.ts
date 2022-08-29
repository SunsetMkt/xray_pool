export interface ProxyItem {
  name: string;
  proto_mode: string;
  socks_port: number;
  http_port: number;
}

export interface ApiResponseProxyList {
  status: 'starting' | 'running' | 'stopped';
  lib_port: number;
  open_result_list: ProxyItem[];
}
