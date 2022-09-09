export type ApiResponse<T> = [T, any] | [null, any];

export interface ApiResponseCommon {
  message: string;
}
