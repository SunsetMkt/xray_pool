export interface RoutingItem {
  routing_type: string;
  rules: string[];
}

export interface RoutingResponseModel {
  block_list: RoutingItem;
  direct_list: RoutingItem;
  proxy_list: RoutingItem;
}

export type RoutingType = 'Block' | 'Direct' | 'Proxy';
