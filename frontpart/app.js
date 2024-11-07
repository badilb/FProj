// Обработчик для формы добавления записей
document.getElementById('record-form').addEventListener('submit', function(event) {
    event.preventDefault();

    const recordInput = document.getElementById('record-input');
    const recordText = recordInput.value;

    // Добавляем запись на сервер
    addRecordToServer(recordText).then((newRecord) => {
        // Добавляем запись в таблицу
        addRecordToTable(newRecord);
        recordInput.value = ''; // Очищаем поле ввода
    });
});

// Функция для добавления записи на сервер
async function addRecordToServer(text) {
    const response = await fetch('http://localhost:9070/records', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ text }),
    });
    if (!response.ok) {
        console.error('Ошибка при добавлении записи:', response.statusText);
        return;
    }
    return await response.json();
}

// Функция для добавления записи в таблицу
function addRecordToTable(record) {
    const tableBody = document.getElementById('record-table').getElementsByTagName('tbody')[0];
    const newRow = tableBody.insertRow();

    const textCell = newRow.insertCell(0);
    const actionsCell = newRow.insertCell(1);

    textCell.textContent = record.text;
    newRow.setAttribute('data-id', record.id); // Сохраняем ID записи для дальнейшего использования

    // Создаем кнопки "Редактировать" и "Удалить"
    const editButton = document.createElement('button');
    editButton.textContent = 'Редактировать';
    editButton.onclick = function() {
        editRecord(newRow, textCell);
    };

    const deleteButton = document.createElement('button');
    deleteButton.textContent = 'Удалить';
    deleteButton.onclick = function() {
        deleteRecordFromServer(record.id).then(() => {
            deleteRecord(newRow);
        });
    };

    actionsCell.appendChild(editButton);
    actionsCell.appendChild(deleteButton);
}

// Функция для редактирования записи
async function editRecord(row, textCell) {
    const newText = prompt('Введите новый текст:', textCell.textContent);
    if (newText) {
        const id = row.getAttribute('data-id');
        await updateRecordOnServer(id, newText);
        textCell.textContent = newText;
    }
}

// Функция для обновления записи на сервере
async function updateRecordOnServer(id, text) {
    const response = await fetch(`http://localhost:9070/records/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ text }),
    });
    if (!response.ok) {
        console.error('Ошибка при обновлении записи:', response.statusText);
    }
}

// Функция для удаления записи на сервере
async function deleteRecordFromServer(id) {
    const response = await fetch(`http://localhost:9070/records/${id}`, {
        method: 'DELETE',
    });
    if (!response.ok) {
        console.error('Ошибка при удалении записи:', response.statusText);
    }
}

// Функция для удаления записи из таблицы
function deleteRecord(row) {
    const tableBody = document.getElementById('record-table').getElementsByTagName('tbody')[0];
    tableBody.deleteRow(row.rowIndex - 1); // Удаляем строку из таблицы
}
