export interface SubscribeItem {
  name: string;
  url: string;
  using?: boolean;
}

export interface SubscribeModel {
  subscribe_list: SubscribeItem[];
}

export interface SubscribeNodeItem {
  name: string;
  proto_model: string;
}

export interface SubscribeNodeModel {
  node_info_list: SubscribeNodeItem[];
}
