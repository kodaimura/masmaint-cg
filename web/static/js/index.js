document.getElementById('generate').addEventListener('click', () => {
	document.getElementById('message').innerHTML = '';

	const file = document.getElementById('csv').files[0];
	const lang = document.getElementById('lang').value;
	const rdbms = document.getElementById('rdbms').value;

	if (file === undefined) {
		renderMessage("csvファイルが選択されていません。", false);
		return;
	}

	const formData = new FormData();
	formData.append('file', file);
	formData.append('lang', lang);
	formData.append('rdbms', rdbms);

	fetch('/csv', {
		method: 'POST',
		body: formData
	})
	.then(response => {
		return response.json()
		.then(data => {
			if (response.ok) {
				download(data.path)
			} else {
				handleErrors(data.errors)
			}
		});
    })
	.catch(console.error);
});

const download = (path) => {
	let alink = document.createElement('a');
	alink.download = path.substring(2);
	alink.href = path;
	alink.click();
	document.getElementById('csv').value = ''
	renderMessage(`${path.substring(2).replace('/', '_')} がダウンロードされました。`, true);
}

const handleErrors = (errors) => {
	for (err of errors) {
		renderMessage(err, false);
	}
}

const renderMessage = (msg, isSuccess) => {
	let message = document.createElement('div');
	message.textContent = msg;
	message.className = `alert alert-${isSuccess? 'success' : 'danger'} alert-custom my-1`;
	document.getElementById('message').appendChild(message);
}