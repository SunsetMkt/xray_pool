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
  xray_port_range: string;
  xray_instance_count: number;
  xray_open_socks_and_http: boolean;
  one_node_test_time_out: number;
  batch_node_test_max_time_out: number;
  test_url: string;
  test_url_thread: number;
  main_proxy_settings: MainProxySettings;
}
