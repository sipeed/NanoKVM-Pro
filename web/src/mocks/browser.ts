import { http, HttpResponse } from 'msw';
import { setupWorker } from 'msw/browser';

export const handlers = [
  http.post('http://localhost/api/auth/login', () => {
    return HttpResponse.json({
      code: 0,
      data: {
        token: 'mocked_token'
      }
    });
  })
];
export const worker = setupWorker(...handlers);
