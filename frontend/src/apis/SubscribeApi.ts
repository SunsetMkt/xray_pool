import BaseApi from './BaseApi';
import type { ApiResponse } from '@/interfaces/common';
import type { SubscribeModel } from '@/interfaces/subscribe';

class SubscribeApi extends BaseApi {
  getList = (): Promise<ApiResponse<SubscribeModel>[]> => this.http('/v1/subscribe_list');
}

export default new SubscribeApi();
