import BaseApi from './BaseApi';
import type { ApiResponse, ApiResponseCommon } from '@/interfaces/common';
import type { ApiResponseProxyList } from '@/interfaces/proxy-pool';

interface ApiRequestStartProxyPool {
  target_site_url: string;
}

class ProxyPoolApi extends BaseApi {
  getStatus = (): Promise<ApiResponse<ApiResponseProxyList>> => this.http('/v1/proxy_list');

  start = (data: ApiRequestStartProxyPool): Promise<ApiResponse<ApiResponseCommon>> =>
    this.http('/v1/start_proxy_pool', data, 'POST');

  stop = (): Promise<ApiResponse<ApiResponseCommon>> => this.http('/v1/stop_proxy_pool', {}, 'POST');

  updateNodeList = (): Promise<ApiResponse<ApiResponseCommon>> => this.http('/v1/update_nodes', {}, 'POST');
}

export default new ProxyPoolApi();
