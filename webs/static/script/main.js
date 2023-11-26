function ajax(url) {
  const p = new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open('GET', url, true)
    xhr.onreadystatechange = () => {
      if (xhr.readyState === 4) {
        if ((xhr.status >= 200 && xhr.status < 300) || xhr.status === 304) {
          resolve(
            JSON.parse(xhr.response)
          )
        } else {
          reject(new Error('Response error'))
        }
      }
    }
    xhr.send(null)
  })
  return p
}

function getAllJson(jsons) {
  var res = ""

  for (k1 in jsons) {
    res += "<h5>" + jsons[k1].index;
    res += " " + jsons[k1].name;
    res += " MAC: " + jsons[k1].mac + "</h5>";

    if ('nips' in jsons[k1]) {
      for (k2 in jsons[k1].nips) {
        res += "<ul>"
        res += "<li>ip:  " + jsons[k1].nips[k2].ip + "</li>"
        if ('arps' in jsons[k1].nips[k2]) {
          res += "<ul>"
          for (k3 in jsons[k1].nips[k2].arps) {
            res += "<li>" + jsons[k1].nips[k2].arps[k3].ip + "  [" + jsons[k1].nips[k2].arps[k3].mac + "]</li>"
          }
          res += "</ul>"
        }
        res += "</ul>"
      }
    }
  }
  return res
};

function SendMac() {
  
}