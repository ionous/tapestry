<!-- 
  view state used for go specific user interactions. 
  ( for example: opening a file )
  wails normally handles these by exposing a go to javascript api;
  tapestry handles these by posting the requested action, then polling using the returned key until the user ( or server ) has completed the desired action. 
 -->
<template>
  <div class="mk-user-action">    
  </div>
</template>
 <script>

import http from '/lib/http.js'
import routes from './routes.js'
import endpoint from './endpoints.js';

// promise a file name
function request(action) {
  // appcfg comes through vite conifg.
  let poll = function(resp) {
    if (!resp.ok) {
      console.warn(resp.status, resp.statusText);
    } else {
      return resp.json().then((res) => {
        console.log("polled", JSON.stringify(res, 0, 2));
        if (res.value!== undefined) { 
          return res.value;

        } else if (res.token!== undefined) {
          const where = endpoint.action(res.token);
          return http.get(where, true).then(poll);
        } else {
          throw new Error("unknown result", res);
        }
      });
    }
  };
  return http.post(appcfg.mosaic + '/actions/' + action, action, true).
    then(poll);
}

export default {
  mounted() {
    const action = routes.getAction();
    switch (action) {
      case 'new':
      case 'open':
        request(action).then((at) => {
          if (at) {
            routes.openFile(at);
          } else {
            routes.goHome();
          }
        }).catch((e) => {
          routes.goHome();
        });

      break;
      default: 
        console.log("unknown action", action);
        routes.goHome();
      break;
    };
  },
}
</script>


