// Package handlers contains all http handlers for server
package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ole-larsen/green-api/internal/httpclient"
)

type StatusResponse struct {
	Status string `json:"status"`
}

// Status godoc
// @Tags Info
// @Summary server status
// @ID serverStatus
// @Accept  json
// @Produce json
// @Success 200 {object} StatusResponse
// @Router /status [get].
func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	_, err := rw.Write([]byte(`{"status":"ok"}`))
	if err != nil {
		InternalServerErrorRequest(rw, r)
		return
	}
}

func BadRequest(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(rw, fmt.Sprintf("%d", http.StatusBadRequest)+" bad request", http.StatusBadRequest)
}

func NotFoundRequest(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(rw, fmt.Sprintf("%d", http.StatusNotFound)+" page not found", http.StatusNotFound)
}

func NotAllowedRequest(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(rw, fmt.Sprintf("%d", http.StatusMethodNotAllowed)+" method not allowed", http.StatusMethodNotAllowed)
}

func ForbiddenRequest(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(rw, fmt.Sprintf("%d", http.StatusForbidden)+" forbidden", http.StatusForbidden)
}

func InternalServerErrorRequest(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.Error(rw, fmt.Sprintf("%d", http.StatusInternalServerError)+" internal server error", http.StatusInternalServerError)
}

// HTMLHandler godoc
// @Tags Html
// @Summary показ 
// @Description показ
// @ID html
// @Produce text/html; charset=utf-8
// @Success 200 {string} string "value"
// @Failure 404 {string} string "404 page not found"
// @Failure 500 {string} string "500 internal server error"
// @Router / [get].
func HTMLHandler(done chan struct{}) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		getSettingsUrl := httpclient.ApiURL + "/waInstance{{idInstance}}/{{method}}/{{apiTokenInstance}}"
		getStateInstanceUrl := httpclient.ApiURL + "/waInstance{{idInstance}}/{{method}}/{{apiTokenInstance}}"
		postMessageUrl := httpclient.ApiURL + "/waInstance{{idInstance}}/{{method}}/{{apiTokenInstance}}"
		postFileUrl := httpclient.ApiURL + "/waInstance{{idInstance}}/{{method}}/{{apiTokenInstance}}"

		template := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Web Interface</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        .container {
            display: flex;
            gap: 20px;
        }
        .form-section {
            flex: 1;
            max-width: 300px;
        }
        .output-section {
            flex: 2;
            border: 1px solid #ccc;
            padding: 10px;
            border-radius: 5px;
            background-color: #f9f9f9;
            word-wrap: break-word;
        }
        input[type="text"] {
            width: 100%;
            margin-bottom: 10px;
            padding: 8px;
            box-sizing: border-box;
        }
        button {
            width: 100%;
            padding: 10px;
            margin-bottom: 10px;
            background-color: #007BFF;
            color: white;
            border: none;
            border-radius: 5px;
            cursor: pointer;
        }
        button:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <h1>Web Interface</h1>
    <div class="container">
        <div class="form-section">
            <input type="text" id="idInstance" placeholder="idInstance">
            <input type="text" id="apiTokenInstance" placeholder="ApiTokenInstance">
            <button id="getSettings">getSettings</button>
            <button id="getStateInstance">getStateInstance</button>
            
            <input type="text" id="chatId" placeholder="Chat ID">
            <input type="text" id="chatMessage" placeholder="Message">
            <button id="sendMessage">sendMessage</button>
            
            <input type="text" id="fileChatId" placeholder="Chat ID">
            <input type="text" id="fileUrl" placeholder="File URL">
            <button id="sendFileByUrl">sendFileByUrl</button>
        </div>
        <div class="output-section" id="output">
            <pre>{ "idMessage": "" }</pre>
        </div>
    </div>
    <script>
		async function getData(url) {
			try {
				const response = await fetch(url);
				if (!response.ok) {
					throw new Error('Response status: ' + response.status);
				}
				return await response.json()
			} catch (e) {
				throw e;
			}
		}
		async function sendData(url, body) {
			try {
				const response = await fetch(url, {
					method: 'POST',
  					body: JSON.stringify(body),
				});
				if (!response.ok) {
					throw new Error('Response status: ' + response.status);
				}
				return await response.json()
			} catch (e) {
				throw e;
			}
		}
		function getRandomInt(min, max) {
			min = Math.ceil(min);
			max = Math.floor(max);
			return Math.floor(Math.random() * (max - min + 1)) + min;
		}
		function isNumber(val) {
			return /^\d+$/.test(val);
		} 
		function formatOutput(message, chatId = '') {
			const idMessage = "GeneratedMessageId" + getRandomInt(1, 1000000);
			chatId = chatId ? chatId : "chatId" + getRandomInt(1, 1000000);
			return JSON.stringify({
				idMessage,
				chatId,
				message
			}, null, 4);
		}
		function checkErrors(id, token, chatId = null, chatMessage = null, fileUrl = null) {
			let message = '';
			if (!id) {
			    message = 'Please fill in idInstance!';
			} else  if (!token) {
			    message = 'Please fill in apiTokenInstance!';
			}

			if (!isNumber(id)) {
    			message = 'idInstance is incorrect!';
			}
			if (chatId !== null && chatId === '') {
    			message = 'Please fill chatId!';
			}
			if (chatMessage !== null && chatMessage === '') {
    			message = 'Please fill message!';
			}
			
			if (fileUrl !== null && fileUrl === '') {
    			message = 'Please fill fileUrl!';
			}
			return message;
		}
		document.getElementById('getSettings').addEventListener('click', (e) => {
            
			e.preventDefault();

			let idInstance = document.getElementById('idInstance').value,
				apiTokenInstance = document.getElementById('apiTokenInstance').value,
			    message = checkErrors(idInstance, apiTokenInstance);

			if (message !== '') {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(message) + "</pre>";
				return;
			}
			
			const getSettingsUrl = '` + getSettingsUrl + `'
				.replace("{{idInstance}}", idInstance)
				.replace("{{method}}", "getSettings")
				.replace("{{apiTokenInstance}}", apiTokenInstance);

					
			getData(getSettingsUrl).then(response => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(response) + "</pre>";
			}).catch(e => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(e.message) + "</pre>";
			});
        });
        
		document.getElementById('getStateInstance').addEventListener('click', (e) => {
			e.preventDefault();

			let idInstance = document.getElementById('idInstance').value,
				apiTokenInstance = document.getElementById('apiTokenInstance').value,
			    message = checkErrors(idInstance, apiTokenInstance);

			if (message !== '') {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(message) + "</pre>";
				return;
			}
			
			const getStateInstance = '` + getStateInstanceUrl + `'
				.replace("{{idInstance}}", idInstance)
				.replace("{{method}}", "getStateInstance")
				.replace("{{apiTokenInstance}}", apiTokenInstance);

			getData(getStateInstance).then(response => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(response) + "</pre>";
			}).catch(e => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(e.message) + "</pre>";
			});
        });
		document.getElementById('sendMessage').addEventListener('click', (e) => {
            
			e.preventDefault();

			let idInstance = document.getElementById('idInstance').value,
				apiTokenInstance = document.getElementById('apiTokenInstance').value,
				chatId = document.getElementById('chatId').value,
				chatMessage = document.getElementById('chatMessage').value,
			    message = checkErrors(idInstance, apiTokenInstance, chatId, chatMessage);

			if (message !== '') {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(message) + "</pre>";
				return;
			}
			
			const sendMessageUrl = '` + postMessageUrl + `'
				.replace("{{idInstance}}", idInstance)
				.replace("{{method}}", "sendMessage")
				.replace("{{apiTokenInstance}}", apiTokenInstance);

			const body = {
				chatId,
				message: chatMessage
			}
		
			sendData(sendMessageUrl, body).then(response => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(response) + "</pre>";
			}).catch(e => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(e.message) + "</pre>";
			});

        });
        document.getElementById('sendFileByUrl').addEventListener('click', (e) => {
            
			e.preventDefault();

			let idInstance = document.getElementById('idInstance').value,
				apiTokenInstance = document.getElementById('apiTokenInstance').value,
				chatId = document.getElementById('fileChatId').value,
				urlFile = document.getElementById('fileUrl').value,
			    message = checkErrors(idInstance, apiTokenInstance, chatId, null, urlFile);

			if (message !== '') {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(message) + "</pre>";
				return;
			}
			
			const sendFileUrl = '` + postFileUrl + `'
				.replace("{{idInstance}}", idInstance)
				.replace("{{method}}", "sendFileByUrl")
				.replace("{{apiTokenInstance}}", apiTokenInstance);

			const parts = urlFile.split("/")	
			const filename = parts[parts.length-1];

			const body = {
				chatId,
				urlFile: urlFile,
    			fileName: filename,
    			caption: filename
			}
			sendData(sendFileUrl, body).then(response => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(response) + "</pre>";
			}).catch(e => {
				document.getElementById('output').innerHTML = "<pre>" + formatOutput(e.message) + "</pre>";
			});

        });
    </script>
</body>
</html>
`

		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(http.StatusOK)

		_, err := io.WriteString(rw, template)
		if err != nil {
			InternalServerErrorRequest(rw, r)
			return
		}
	}
}
