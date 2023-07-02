document.getElementById('generate').addEventListener('click', () => {
	const file = document.getElementById('csv').files[0]

	const formData = new FormData();
	formData.append('file', file);

	fetch('', {
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
	alink.download = path.substring(5);
	alink.href = path;
	alink.click();
}

const handleErrors = (errors) => {
	console.log(errors)
	document.getElementById('message').innerHTML = errors 
}


