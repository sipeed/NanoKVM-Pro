import { http } from '@/lib/http.ts';

// get image list
export function getImages() {
  return http.get('/api/storage/image');
}

// get mounted image
export function getMountedImage() {
  return http.get('/api/storage/image/mounted');
}

// mount/unmount image
export function mountImage(file?: string, cdrom?: boolean, readOnly?: boolean) {
  const data = {
    file: file ? file : '',
    cdrom: cdrom,
    readOnly: readOnly
  };
  return http.post('/api/storage/image/mount', data);
}

// upload image
export function uploadImage(file: FormData, fid: string = '') {
  return http.request({
    url: '/api/storage/image/upload?' + `fid=${fid}`,
    method: 'post',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data: file
  });
}

// delete image
export function deleteImage(file: string) {
  const data = {
    file
  };
  return http.post('/api/storage/image/delete', data);
}

// download image
export function downloadImage(file?: string) {
  const data = {
    file: file ? file : ''
  };
  return http.post('/api/storage/download/image', data);
}

// get image download status
export function statusImage() {
  return http.get('/api/storage/download/image/status');
}

// check if download is supported
export function imageEnabled() {
  return http.get('/api/storage/download/image/enabled');
}
