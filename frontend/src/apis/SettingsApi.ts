import BaseApi from './BaseApi';

class SettingsApi extends BaseApi {
  get = () => this.http('/v1/settings');

  getDefaultSettings = () => this.http('/v1/def_settings');

  update = (data: any) => this.http('/v1/settings', data, 'PUT');
}

export default new SettingsApi();
