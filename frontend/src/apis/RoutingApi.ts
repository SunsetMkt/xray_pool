import BaseApi from './BaseApi';
import type { ApiResponse, ApiResponseCommon } from '@/interfaces/common';
import type { RoutingResponseModel } from '@/interfaces/Routing';

export interface ApiRequestAddRouting {
  routing_type: string;
  rules: string[];
}

export interface ApiRequestRemoveRouting {
  routing_type: string;
  index_list: number[];
}

class RoutingApi extends BaseApi {
  getList = (): Promise<ApiResponse<RoutingResponseModel>> => this.http('/v1/routing_list');

  add = (data: ApiRequestAddRouting): Promise<ApiResponse<ApiResponseCommon>> =>
    this.http('/v1/routing_add', data, 'POST');

  remove = (data: ApiRequestRemoveRouting): Promise<ApiResponse<ApiResponseCommon>> =>
    this.http('/v1/routing_delete', data, 'POST');
}

export default new RoutingApi();
