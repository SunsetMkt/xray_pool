import BaseApi from './BaseApi';
import type { ApiResponseSettings } from '@/interfaces/settings';

class SettingsApi extends BaseApi {
  get = () => this.http('/v1/settings');

  getDefaultSettings = () => this.http('/v1/def_settings');

  update = (data: ApiResponseSettings) => this.http('/v1/settings', data, 'PUT');
}

export default new SettingsApi();
