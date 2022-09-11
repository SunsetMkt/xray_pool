import BaseApi from '@/apis/BaseApi';
import type { ApiResponse } from '@/interfaces/common';

export interface ApiResponseSystemStatus {
  is_setup: boolean;
  os: string;
  arch: string;
  glider_downloaded: boolean;
  xray_downloaded: boolean;
}

class CommonApi extends BaseApi {
  getSystemStatus = (): Promise<ApiResponse<ApiResponseSystemStatus>> => this.http('/v1/system-status');
}

export default new CommonApi();
