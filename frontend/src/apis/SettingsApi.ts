import BaseApi from './BaseApi';

class SettingsApi extends BaseApi {
  get = () => this.http('/v1/settings');
}

export default new SettingsApi();
