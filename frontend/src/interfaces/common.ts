export type ApiResponse<T> = [T, any] | [null, any];

export interface ApiResponseSystemStatus {
  is_setup: boolean;
}
