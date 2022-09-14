import BaseApi from './BaseApi';

export interface MainProxySettings {
  PID: number;
  HttpPort: number;
  SocksPort: number;
  AllowLanConn: boolean;
  Sniffing: boolean;
  RelayUDP: boolean;
  DNSPort: number;
  DNSForeign: string;
  DNSDomestic: string;
  DNSDomesticBackup: string;
  BypassLANAndMainLand: boolean;
  RoutingStrategy: string;
  Mux: boolean;
}

export interface ApiResponseSettings {
  app_start_port: number;
  manual_lb_port: number;
  xray_port_range: string;
  xray_instance_count: number;
  xray_open_socks_and_http: boolean;
  one_node_test_time_out: number;
  batch_node_test_max_time_out: number;
  test_url: string;
  health_check_url: string;
  health_check_interval: number;
  test_url_thread: number;
  test_url_hard_way: boolean;
  test_url_failed_words: string[];
  test_url_failed_regex: string;
  test_url_status_code: string;
  glider_strategy: 'rr' | 'ha' | 'lha' | 'dh';
  main_proxy_settings: MainProxySettings;
}

class SettingsApi extends BaseApi {
  get = () => this.http('/v1/settings');

  getDefaultSettings = () => this.http('/v1/def_settings');

  update = (data: ApiResponseSettings) => this.http('/v1/settings', data, 'PUT');
}

export default new SettingsApi();
