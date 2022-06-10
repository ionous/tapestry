// https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch
// readJson indicates whether to parse the response as json formatted data.
async function send(url, method, data, readJson) {
  // Default options are marked with *
  const response = await fetch(url, {
    method,
    // mode: 'cors', // no-cors, *cors, same-origin
    // cache: 'no-cache', // *default, no-cache, reload, force-cache, only-if-cached
    // credentials: 'same-origin', // include, *same-origin, omit
    headers: {
      'Content-Type': 'application/json'
    },
    // redirect: 'follow', // manual, *follow, error
    // referrerPolicy: 'no-referrer', // no-referrer, *no-referrer-when-downgrade, origin, origin-when-cross-origin, same-origin, strict-origin, strict-origin-when-cross-origin, unsafe-url
    body: JSON.stringify(data) // body data type must match "Content-Type" header
  });
  return readJson? response.json(): true;
}

export default {
  // promise json
  join(base, path){
     const url= base+path; // fix to exception on dots
     return url;
  },

  get(url) {
    console.log("getting", url);
    return fetch(url).then((response) => {
      return (!response.ok) ?
        Promise.reject({status: response.status, url}) :
        response.json().then(result => result);
    }).catch((error) => {
      console.log('error:', error)
    });
  },

  // readJson indicates whether to parse the response as json formatted data.
  async post(url, data = {}, readJson=false) {
    return send(url, 'POST', data, readJson);
  },

  // readJson indicates whether to parse the response as json formatted data.
  async put(url, data = {}, readJson=false) {
    return send(url, 'PUT', data, readJson);
  }
}
