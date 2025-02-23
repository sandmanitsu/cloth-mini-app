
import { optionBrands } from './brand.js';
import { optionCategory } from './category.js';
import { formatDate } from './date.js';

async function performSearch() {
    const params = new URLSearchParams();

    // Собираем значения из формы поиска
    if (document.getElementById('id-search').value) {
        params.append('id', document.getElementById('id-search').value);
    }
    if (document.getElementById('brand-search').value) {
        params.append('brand_id', document.getElementById('brand-search').value);
    }
    if (document.getElementById('name-search').value) {
        params.append('name', document.getElementById('name-search').value);
    }
    if (document.getElementById('category-search').value && document.getElementById('category-search').value !== '') {
        params.append('category_id', document.getElementById('category-search').value);
    }
    if (document.getElementById('gender-search').value && document.getElementById('gender-search').value !== '') {
        params.append('sex', document.getElementById('gender-search').value);
    }
    if (document.getElementById('min-price-search').value) {
        params.append('min_price', document.getElementById('min-price-search').value);
    }
    if (document.getElementById('max-price-search').value) {
        params.append('max_price', document.getElementById('max-price-search').value);
    }
    if (document.getElementById('discount-search').value) {
        params.append('discount', document.getElementById('discount-search').value);
    }

    fetchItems(20, 0, '&'+params.toString())
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
        console.error('fetchItems Ошибка', error.message, error)
    }
}

// Отрисовывает список товаров
function renderItems(items) {
    const container = document.getElementById("products-body")

    container.innerHTML = ''

    items.items.forEach((product) => {
        let sex = ''
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
        <tr class="product">
            <td>${product.id}</td>
            <td>${product.brand_name}</td>
            <td>${product.name}</td>
            <td>${product.category_name}</td>
            <td>${sex}</td>
            <td>${product.price} руб.</td>
            <td>${product.discount} %</td>
            <td><a href="${product.outer_link}" target="_blank">Товар в магазине</a></td>
            <td>${formatDate(product.created_at)}</td>
            <td>${formatDate(product.updated_at)}</td>
            <td><a class="button button-primary" href="update/${product.id}">Редактировать</a></td>
        </tr>`;

    container.insertAdjacentHTML('beforeend', itemCard);
    });
}

function fetchNextItems() {
    let currOffset = document.getElementById("offset-counter")

    if (document.getElementsByClassName("product").length !== 20) {
        return
    }

    let offset = parseInt(currOffset.innerHTML)+20
    fetchItems(20, offset, "")

    currOffset.innerHTML = offset
}

function fetchPrevItems() {
    let currOffset = document.getElementById("offset-counter")

    if (parseInt(currOffset.innerHTML) == 0) {
        return
    }

    let offset = parseInt(currOffset.innerHTML)-20
    fetchItems(20, offset, "")

    currOffset.innerHTML = offset
}

document.addEventListener('DOMContentLoaded', () => {
    fetchItems()
    optionCategory(document.getElementById("category-search"))
    optionBrands(document.getElementById("brand-search"))

    const searchBtn = document.getElementById("search_btn")
    searchBtn.addEventListener('click', performSearch)

    const nextBtn = document.getElementById("next_btn")
    nextBtn.addEventListener('click', fetchNextItems)

    const prevBtn = document.getElementById("prev_btn")
    prevBtn.addEventListener('click', fetchPrevItems)
});