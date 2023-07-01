/* 初期設定 */
window.addEventListener('DOMContentLoaded', (event) => {
	setUp();
});

/* リロードボタン押下 */
document.getElementById('reload').addEventListener('click', (event) => {
	clearMessage();
	document.getElementById('records').innerHTML = '';
	setUp();
})

/* 保存モーダル確定ボタン押下 */
document.getElementById('ModalSaveAllOk').addEventListener('click', (event) => {
	clearMessage();
	doPutAll();
	doPost();
})

/* 削除モーダル確定ボタン押下 */
document.getElementById('ModalDeleteAllOk').addEventListener('click', (event) => {
	clearMessage();
	doDeleteAll();
})

/* チェックボックスの選択一覧取得 */
const getDeleteTarget = () => {
	let dels = document.getElementsByName('del');
	let ret = [];

	for (let x of dels) {
		if (x.checked) {
			ret.push(x.value);
		}
	}
	return ret
}

const renderMessage = (msg, count, isSuccess) => {
	if (count !== 0) {
		let message = document.createElement('div');
		message.textContent = `${count}件の${msg}に${isSuccess? '成功' : '失敗'}しました。`
		message.className = `alert alert-${isSuccess? 'success' : 'danger'} alert-custom my-1`;
		document.getElementById('message').appendChild(message);
	}
}

const clearMessage = () => {
	document.getElementById('message').innerHTML = '';
}

const nullToEmpty = (s) => {
	return (s == null)? '' : s;
}

/* チェンジアクション */
const changeAction = (event) => {
	let target = event.target;
	let target_bk = target.nextElementSibling;

	if (target.value !== target_bk.value) {
		target.classList.add('changed');
	} else {
		target.classList.remove('changed');
	}
}

/* <tbody></tbody>内のレコードにチェンジアクション追加 */
const addChangedAction = (columnName) => {
	let elems = document.getElementsByName(columnName);
	for (const elem of elems) {
		elem.addEventListener('change', changeAction);
	}
}

/* <tbody></tbody>レンダリング */
const renderTbody = (data) => {
	let tbody= '';
	for (const elem of data) {
		tbody += createTr(elem);
	}
	tbody += createTrNew();

	document.getElementById('records').innerHTML = tbody;
}