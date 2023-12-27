document.getElementById('connectButton').addEventListener('click', function() {
    var ns = document.getElementById('nsInput').value;
    var podName = document.getElementById('podNameInput').value;
    var container = document.getElementById('containerInput').value;

    var url = 'ws://127.0.0.1:8085/v1/logs?ns=' + encodeURIComponent(ns) + '&podName=' + encodeURIComponent(podName) + '&container=' + encodeURIComponent(container);
    var socket = new WebSocket(url);


    // 清空output元素的内容
    document.getElementById('output').innerHTML = '';

    socket.onopen = function() {
        socket.send('eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJhZG1pbiIsImV4cCI6MTcwMzY3MTA3MywiaXNzIjoiMzgzODQtU2VhcmNoRW5naW5lIn0.NVoEb8q7yHWCHss1BbZZ-8snfXoJu8VAroaNvJ6CfnU');
    };

    socket.onmessage = function(event) {
        var output = document.getElementById('output');
        output.innerHTML += '<p>' + event.data + '</p>';
    };

    socket.onerror = function(error) {
        console.log('WebSocket Error: ' + error);
        // 清空output元素的内容
        document.getElementById('output').innerHTML = '';
    };

    socket.onclose = function(event) {
        console.log('WebSocket connection closed: ' + event.code);
    };
});

document.getElementById('output').addEventListener('DOMNodeInserted', function() {
    this.scrollTop = this.scrollHeight;
});