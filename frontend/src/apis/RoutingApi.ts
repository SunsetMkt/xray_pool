import BaseApi from './BaseApi';
import type { ApiResponse } from '@/interfaces/common';
import type { RoutingResponseModel } from '@/interfaces/Routing';

class RoutingApi extends BaseApi {
  getList = (): Promise<ApiResponse<RoutingResponseModel>> => this.http('/v1/routing_list');

  add = (data: any): Promise<ApiResponse<any>> => this.http('/v1/routing_add', data, 'POST');

  remove = (data: any): Promise<ApiResponse<any>> => this.http('/v1/routing_delete', data, 'POST');
}

export default new RoutingApi();
