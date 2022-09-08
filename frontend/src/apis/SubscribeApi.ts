import BaseApi from './BaseApi';
import type { ApiResponse, ApiResponseCommon } from '@/interfaces/common';
import type { SubscribeModel, SubscribeNodeModel } from '@/interfaces/subscribe';

export interface ApiRequestAddSubscribe {
  name: string;
  url: string;
}

export interface ApiRequestUpdateSubscribe {
  index: number;
  name: string;
  url: string;
  using: boolean;
}

class SubscribeApi extends BaseApi {
  getList = (): Promise<ApiResponse<SubscribeModel>> => this.http('/v1/subscribe_list');

  add = (data: ApiRequestAddSubscribe): Promise<ApiResponse<ApiResponseCommon>> =>
    this.http('/v1/add_subscribe', data, 'POST');

  remove = (index: number): Promise<ApiResponse<ApiResponseCommon>> =>
    this.http('/v1/del_subscribe', { index: `${index}` }, 'POST');

  update = (data: ApiRequestUpdateSubscribe): Promise<ApiResponse<ApiResponseCommon>> =>
    this.http('/v1/update_subscribe', data, 'POST');

  getNodeList = (): Promise<ApiResponse<SubscribeNodeModel>> => this.http('/v1/node_list');
}

export default new SubscribeApi();
