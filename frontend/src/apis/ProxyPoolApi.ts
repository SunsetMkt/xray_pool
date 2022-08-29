import BaseApi from './BaseApi';
import type { ApiResponse } from '@/interfaces/common';
import type { ApiResponseProxyList } from '@/interfaces/proxy-pool';

class ProxyPoolApi extends BaseApi {
  getStatus = (): Promise<ApiResponse<ApiResponseProxyList>> => this.http('/v1/proxy_list');

  start = (data: any): Promise<ApiResponse<any>> => this.http('/v1/start_proxy_pool', data, 'POST');

  stop = (): Promise<ApiResponse<any>> => this.http('/v1/stop_proxy_pool', {}, 'POST');
}

export default new ProxyPoolApi();
