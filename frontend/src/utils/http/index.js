import { createRequest, registerInterceptor } from './http-client';

// 扫描interceptors下的目录，并且注册到http-client中
const modulesFiles = import.meta.globEager('./interceptors/**.js');

Object.keys(modulesFiles).forEach((key) => {
  const module = modulesFiles[key].default;
  registerInterceptor(module);
});

export { createRequest };
