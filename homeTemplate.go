package main

const HomeTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <title></title>
    <meta charset="UTF-8">
    <meta http-equiv="Cache-control" content="no-cache">
    <meta http-equiv="Expires" content="-1">
    <meta name="viewport" content="width=device-width, initial-scale=1">
  </head>
  <style>
    body {
      margin: 0;
      padding: 0;
      overflow: hidden;
    }
    iframe {
      margin: 0;
      width: 100vw;
      height: 100vh;
    }
  
  </style>
  <body>
    <iframe id="view" src="http://{{.Host}}/watcher{{.Prefix}}" frameborder="0"></iframe>
  </body>
  <script>
    
    let iframe = document.querySelector("#view");
    let view = iframe.contentWindow || iframe.contentDocument.document || iframe.contentDocument;

    
		let socket = new WebSocket("ws://{{.Host}}/ws")
		
    socket.onmessage = (message) => {
      console.log(message)
      iframe.src = iframe.src
    }

    socket.onerror = (err) => {
      console.log(err)
    }

    socket.onopen = (q) => {
      console.log(q)
      
    }
  </script>
</html>
`
