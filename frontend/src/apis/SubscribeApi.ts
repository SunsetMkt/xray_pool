import BaseApi from './BaseApi';
import type { ApiResponse } from '@/interfaces/common';
import type { SubscribeModel, SubscribeNodeModel } from '@/interfaces/subscribe';

class SubscribeApi extends BaseApi {
  getList = (): Promise<ApiResponse<SubscribeModel>> => this.http('/v1/subscribe_list');

  add = (data: any): Promise<ApiResponse<any>> => this.http('/v1/add_subscribe', data, 'POST');

  remove = (index: number): Promise<ApiResponse<any>> => this.http('/v1/del_subscribe', { index }, 'POST');

  update = (data: any): Promise<ApiResponse<any>> => this.http('/v1/update_subscribe', data, 'POST');

  getNodeList = (): Promise<ApiResponse<SubscribeNodeModel>> => this.http('/v1/node_list');
}

export default new SubscribeApi();
