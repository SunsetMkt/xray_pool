import BaseApi from './BaseApi';
import type { ApiResponseCommon } from '@/interfaces/common';

class SettingsApi extends BaseApi {
  get = () => this.http('/v1/settings');

  getDefaultSettings = () => this.http('/v1/def_settings');

  update = (data: ApiResponseCommon) => this.http('/v1/settings', data, 'PUT');
}

export default new SettingsApi();
