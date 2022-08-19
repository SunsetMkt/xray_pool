const handleError = (error) => {
  // eslint-disable-next-line
  console.error('interceptor catch the error!\n', error);
  const errorMessageText = error.data?.message || error.message || '网络错误';

  const rtData = {
    error,
    message: errorMessageText,
  };

  return Promise.reject(rtData);
};

export default {
  onRequestRejected: (error) => handleError(error),
  onResponseFullFilled: (response) => {
    const { data } = response;
    // 正常返回但是code是错误码的情况也需要异常处理
    if (data?.code && data?.code > 300) {
      return handleError(response);
    }
    return response;
  },
  onResponseRejected: (error) => handleError(error?.response || error),
};
