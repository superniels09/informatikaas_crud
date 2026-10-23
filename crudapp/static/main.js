function proxyURL(urlString) {
  const url = new URL(urlString, window.location.origin);

  const targetPath = url.pathname;
  url.pathname = '/crudapp/apiproxy/';
  url.searchParams.set('__target_path', targetPath);

  const isRelative = !urlString.startsWith('http://') && !urlString.startsWith('https://');
  return isRelative
    ? url.pathname + url.search
    : url.toString();
}

async function postData(url, data) {
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data), // Convert JS object to JSON string
  });

	if (!response.ok) {
		let errorMessage = ""
		try {
			const responseJSON = await response.json()
			errorMessage = responseJSON.error
		} catch {
			throw new Error(`HTTP error! status: ${response.status}`);
		}
		throw new Error(`HTTP error! status: ${response.status}, message: ${errorMessage}`);
  }

  return await response.json(); // Parse the response JSON
}

function sendMessage(message) {
	return postData(proxyURL("/messages/post"), {
		message: message
	})
}

async function listMessages() {
	const response = await fetch(proxyURL("/messages/list"))
	const resp = await response.json()
	return resp.messages
}

function main() {
	const messageField = document.getElementById("message-field")
	const messageSendButton = document.getElementById("message-send-button")

	messageSendButton.addEventListener("click", async () => {
		/** @type {string} */
		const message = messageField.value
		if (message.trim() === "") {
			alert("no sending empty messages >:(")
			return
		}
		const resp = await sendMessage(message)
		console.log(resp)

		alert("Sent successfully")
		messageField.value = ""
	})
}
document.addEventListener('DOMContentLoaded', main)
