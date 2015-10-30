const root = 'http://localhost:8000/';

const parse = function (xhr) {
  let body = xhr.response;
  let contentType = xhr.getResponseHeader('Content-Type');
  if (/json$/.test(contentType)) {
    return JSON.parse(body);
  } else {
    return body;
  }
};

export default function (url, options={}) {
  return new Promise((resolve, reject) => {
    let method = options.type || 'GET';
    url = root + url;
    let xhr = new window.XMLHttpRequest();
    xhr.onload = e => {
      if (xhr.status === 200) {
        resolve(parse(xhr));
      } else {
        reject(parse(xhr));
      }
    };
    xhr.open(method, url, true);
    xhr.send(options.data);
  });
};
