import axios, { AxiosHeaders, AxiosInstance, AxiosRequestConfig } from 'axios';

import { removeToken } from '@/lib/cookie.ts';
import { getBaseUrl } from '@/lib/service.ts';

type Response = {
  code: number;
  msg: string;
  data: any;
};

class Http {
  private instance: AxiosInstance;

  constructor() {
    const baseURL = getBaseUrl('http');

    this.instance = axios.create({
      baseURL,
      withCredentials: true,
      timeout: 60 * 1000
    });

    this.setInterceptors();
  }

  private setInterceptors() {
    this.instance.interceptors.request.use((config) => {
      if (config.headers) {
        config.headers.Accept = 'application/json';
      }

      return config;
    });

    this.instance.interceptors.response.use(
      (response) => {
        return response.data;
      },
      (error) => {
        console.log(error);
        const code = error.response?.status;
        if (code === 401) {
          removeToken();
          window.location.reload();
        }
        return Promise.reject(error);
      }
    );
  }

  public get(url: string, params?: any, headers?: AxiosHeaders): Promise<Response> {
    return this.instance.request({
      method: 'get',
      url,
      params,
      headers
    });
  }

  public post(url: string, data?: any): Promise<Response> {
    return this.instance.request({
      method: 'post',
      url,
      data
    });
  }

  public delete(url: string, data?: any): Promise<Response> {
    return this.instance.request({
      method: 'delete',
      url,
      data
    });
  }

  public request(config: AxiosRequestConfig): Promise<Response> {
    return this.instance.request(config);
  }
}

export const http = new Http();
