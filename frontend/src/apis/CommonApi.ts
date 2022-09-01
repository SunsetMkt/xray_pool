import BaseApi from '@/apis/BaseApi';
import type { ApiResponse, ApiResponseSystemStatus } from '@/interfaces/common';

class CommonApi extends BaseApi {
  getSystemStatus = (): Promise<ApiResponse<ApiResponseSystemStatus>> => this.http('/v1/system-status');
}

export default new CommonApi();
