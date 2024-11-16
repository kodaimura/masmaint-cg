document.getElementById('generate').addEventListener('click', () => {
	document.getElementById('message').innerHTML = '';

	const ddl = document.getElementById('ddl').files[0];
	const lang = document.getElementById('lang').value;
	const rdbms = document.getElementById('rdbms').value;

	if (ddl === undefined) {
		renderMessage("DDLファイルが選択されていません。", false);
		return;
	}

	const formData = new FormData();
	formData.append('ddl', ddl);
	formData.append('lang', lang);
	formData.append('rdbms', rdbms);

	fetch('/generate', {
		method: 'POST',
		body: formData
	})
	.then(response => {
		return response.json()
		.then(data => {
			if (response.ok) {
				download(data.zip)
			} else {
				handleErrors(data.errors)
			}
		});
    })
	.catch(console.error);
});

const download = (zip) => {
	let alink = document.createElement('a');
	alink.download = zip;
	alink.href = `output/${zip}`;
	alink.click();
	document.getElementById('ddl').value = ''
	renderMessage(`${path} がダウンロードされました。`, true);
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