
var token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJhZG1pbiIsImV4cCI6MTcwMzc2Nzc4NCwiaXNzIjoiMzgzODQtU2VhcmNoRW5naW5lIn0.BaAAxKOjeHZpd_ciAHUtJfYlnsARr3BatDg4Bz0vCmE'
document.getElementById('connectButton').addEventListener('click', function() {
    var ns = document.getElementById('nsInput').value;
    var podName = document.getElementById('podNameInput').value;
    var container = document.getElementById('containerInput').value;

    var url = 'ws://127.0.0.1:8085/v1/logs?ns=' + encodeURIComponent(ns) + '&podName=' + encodeURIComponent(podName) + '&container=' + encodeURIComponent(container);
    var socket = new WebSocket(url);



    // 清空output元素的内容
    document.getElementById('output').innerHTML = '';

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



var selectedNamespace = '';
var selectedLabel = '';

document.getElementById('nsButton').addEventListener('click', function () {
    // Fetch namespaces from server
    fetch(`http://localhost:8082/v1/r/namespace/list`, {
        method: 'GET',
        headers: {
            'Authorization': token
        }
    }).then(response => response.json()).then(data => {
        var nsDiv = document.getElementById('nsDiv');
        nsDiv.innerHTML = '';
        data.data.forEach(ns => {
            var button = document.createElement('button');
            button.innerText = ns.namespace;
            button.addEventListener('click', function () {
                // Set the value of selectedNamespace to the clicked namespace
                selectedNamespace = ns.namespace;
                // Set the value of nsInput to the clicked namespace
                document.getElementById('nsInput').value = selectedNamespace;
                // Fetch deployments from server
                fetch(`http://localhost:8082/v1/r/deployments/list?namespace=${encodeURIComponent(selectedNamespace)}`, {
                    method: 'GET',
                    headers: {
                        'Authorization': token
                    }
                }).then(response => response.json()).then(data => {
                    var deploymentDiv = document.getElementById('deploymentDiv');
                    deploymentDiv.innerHTML = '';
                    data.data.forEach(deployment => {
                        var button = document.createElement('button');
                        button.innerText = deployment.name;
                        button.dataset.label = deployment.labels[0].split('=')[1];
                        deploymentDiv.appendChild(button);
                    });
                });
            });
            nsDiv.appendChild(button);
        });
    });
});

document.getElementById('deploymentDiv').addEventListener('click', function (event) {
    var deployment = event.target.innerText;
    selectedLabel = event.target.dataset.label;
    // Set the value of containerInput to the clicked deployment's label
    document.getElementById('containerInput').value = selectedLabel;
    // Fetch pods from server
    fetch(`http://127.0.0.1:8082/v1/r/pod/getPods?namespace=${encodeURIComponent(selectedNamespace)}&label=${encodeURIComponent(selectedLabel)}`, {
        method: 'GET',
        headers: {
            'Authorization': token
        }
    }).then(response => response.json()).then(data => {
        var podDiv = document.getElementById('podDiv');
        podDiv.innerHTML = '';
        data.data.forEach(pod => {
            var button = document.createElement('button');
            button.innerText = pod.name;
            podDiv.appendChild(button);
        });
    });
});

document.getElementById('podDiv').addEventListener('click', function (event) {
    var pod = event.target.innerText;
    // Set the value of podNameInput to the clicked pod
    document.getElementById('podNameInput').value = pod;
});