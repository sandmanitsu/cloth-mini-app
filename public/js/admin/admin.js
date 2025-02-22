
async function performSearch() {
    const params = new URLSearchParams();

    // Собираем значения из формы поиска
    if (document.getElementById('id-search').value) {
        params.append('id', document.getElementById('id-search').value);
    }
    if (document.getElementById('brand-search').value) {
        params.append('brand', document.getElementById('brand-search').value);
    }
    if (document.getElementById('name-search').value) {
        params.append('name', document.getElementById('name-search').value);
    }
    if (document.getElementById('category-search').value && document.getElementById('category-search').value !== '') {
        params.append('category', document.getElementById('category-search').value);
    }
    if (document.getElementById('gender-search').value && document.getElementById('gender-search').value !== '') {
        params.append('sex', document.getElementById('gender-search').value);
    }
    if (document.getElementById('price-search').value) {
        params.append('price', document.getElementById('price-search').value);
    }
    if (document.getElementById('discount-search').value) {
        params.append('discount', document.getElementById('discount-search').value);
    }

    fetchItems(20, 0, '&'+params.toString())
}

// Форматирует дату из ISO в формат ГГГГ-ММ-ДД - ЧЧ:ММ:СС
function formatDate(dt) {
    if (!dt) {
        return ""
    }

    const date = new Date(dt);

    return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} - ${date.toLocaleTimeString()}`;
}

// Получение списка товаров
async function fetchItems(limit = 20, offset = 0, queryParams = "") {
    try {
        const url = `http://localhost:8080/item/get?limit=${limit}&offset=${offset}${queryParams}`;
        const response = await fetch(url)

        if (!response.ok) {
            throw new Error(`fetchItems Ошибка HTTP: ${response.status}`)
        }

        const items = await response.json()
        renderItems(items)
    } catch (error) {
        console.error('fetchItems Ошибка', error.message)
    }
}

// Отрисовывает список товаров
function renderItems(items) {
    const container = document.getElementById("products-body")

    container.innerHTML = ''

    items.items.forEach((product) => {
        switch (product.sex) {
            case 1:
                sex = 'муж';
                break;
            case 1:
                sex = 'жен';
                break;
            default:
                sex = 'уни';
                break;
        }

        const itemCard = `
        <tr>
            <td>${product.id}</td>
            <td>${product.brand}</td>
            <td>${product.name}</td>
            <td>${product.category_name}</td>
            <td>${sex}</td>
            <td>${product.price} руб.</td>
            <td>${product.discount} %</td>
            <td><a href="${product.outer_link}" target="_blank">Товар в магазине</a></td>
            <td>${formatDate(product.created_at)}</td>
            <td>${formatDate(product.updated_at)}</td>
            <td><button class="edit-button">Редактировать</button></td>
        </tr>`;

    container.insertAdjacentHTML('beforeend', itemCard);
    });
}

async function fetchCategory() {
    try {
        const url = `http://localhost:8080/category/get`;
        const response = await fetch(url)

        if (!response.ok) {
            throw new Error(`fetchCategory Ошибка HTTP: ${response.status}`)
        }

        const category = await response.json()
        const select = document.getElementById("category-search")

        let html = `<option value="">Все категории</option>`
        category.forEach(cat => {
            html += `<option value="${cat.category_id}">${cat.category_name}</option>`
        })

        select.insertAdjacentHTML('beforeend', html);
    } catch (error) {
        console.error('fetchCategory Ошибка', error.message)
    }
}

document.addEventListener('DOMContentLoaded', () => {
    fetchItems()
    fetchCategory()
});