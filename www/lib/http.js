const appjson = 'application/json';

// https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
// 'raw' indicates whether to parse the response or return the response itself.
function send(url, method, data, raw) {
  // Default options are marked with *
  return fetch(url, {
    method,
    // mode: 'cors', // no-cors, *cors, same-origin
    // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
    // credentials: 'same-origin', // include, *same-origin, omit
    headers: {
      'Content-Type': appjson
    },
    // redirect: 'follow', // manual, *follow, error
    // referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify(data) // body data type must match "Content-Type" header
  }).then((res) => {
    return raw ? res :
      (res.headers.get("Content-Type") === appjson) ?
      res.json() : res.text();
  });
}

export default {
  // promise json
  join(base, path) {
     return base + path; // fix: exception on dots
  },

  get(url, raw = false) {
    console.log("getting", url);
    return fetch(url).then((response) => {
      return raw? response: (!response.ok) ?
        Promise.reject({status: response.status, url}) :
        response.json().then(result => result);
    });
  },

  // raw indicates whether to parse the response as json formatted data.
  post(url, data = {}, raw = false) {
    return send(url, 'POST', data, raw);
  },

  // raw indicates whether to parse the response as json formatted data.
  put(url, data = {}, raw = false) {
    return send(url, 'PUT', data, raw);
  }
}
