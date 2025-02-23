
import { formatDate } from "./date.js";
import { fetchBrands } from "./brand.js";
import { fetchCategory } from "./category.js";

async function fetchItem() {
    let url = window.location.href
    let id = url.match(/\d+$/)[0]

    try {
        const url = `http://localhost:8080/item/get/${id}`;
        const response = await fetch(url)
        const brands = await fetchBrands()
        const category = await fetchCategory()

        if (!response.ok) {
            throw new Error(`fetchItems Ошибка HTTP: ${response.status}`)
        }

        const item = await response.json()

        renderItem(item, brands, category)
    } catch (error) {
        console.error('fetchItems Ошибка', error.message, error)
    }
}

function renderItem(item, brands, category) {
    const container = document.getElementById("item")

    let brandOptions = `<option value="${item.brand_id}">${item.brand_name}</option>`
    brands.forEach(cat => {
        brandOptions += `<option value="${cat.brand_id}">${cat.brand_name}</option>`
    })

    let categoryOptions = `<option value="${item.category_id}">${item.category_name}</option>`
    category.forEach(cat => {
        categoryOptions += `<option value="${cat.category_id}">${cat.category_name}</option>`
    })

    let sex = ''
    switch (item.sex) {
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

    let html = `
            <div class="container">
                <input type="hidden" id="item-id" value="${item.id}">
                <h6 item_id="${item.id}">ID: ${item.id} | Создан: ${formatDate(item.created_at)} | Обновлен: ${formatDate(item.updated_at)}</h6>
            </div>

            <div class="container">
                <div class="four columns">
                    <div class="row">
                        <img src="../static/img/cardigan_mock.jpg" width="100%" height="auto" alt="mock image">
                        <button class="u-full-width" id="update_btn">Обновить изображение</button>
                    </div>
                </div>

                <div class="eight columns">
                    <div class="five columns">
                        <label for="brand">Бренд - ${item.brand_name}</label>
                        <select class="u-full-width"  id="brand">
                            ${brandOptions}
                        </select>
                    </div>

                    <div class="seven columns">
                        <label for="item-name">Название:</label>
                        <input class="u-full-width"  type="text" id="item-name" placeholder="Название" value="${item.name}" />
                    </div>
                </div>

                <div class="eight columns">
                    <div class="five columns">
                        <label for="category">Категория - ${item.category_name}</label>
                        <select class="u-full-width"  id="category">
                            ${categoryOptions}
                        </select>
                    </div>

                    <div class="two columns">
                        <label for="gender">Пол:</label>
                        <select class="u-full-width"  id="gender">
                            <option value="${item.sex}">${sex}</option>
                            <option value="1">Мужской</option>
                            <option value="2">Женский</option>
                            <option value="3">Унисекс</option>
                        </select>
                    </div>

                    <div class="three columns">
                        <label for="price">Цена:</label>
                        <input class="u-full-width"  type="number" id="price" placeholder="Цена" value="${item.price}"/>
                    </div>

                    <div class="two columns">
                        <label for="discount">Скидка:</label>
                        <input class="u-full-width"  type="number" id="discount" placeholder="Скидка" value="${item.discount}" />
                    </div>
                </div>

                <div class="eight columns">
                    <label for="description">Описание:</label>
                    <textarea id="description" cols="50" rows="5" placeholder="">${item.description}</textarea>
                </div>
            </div>`

            container.insertAdjacentHTML('beforeend', html)
}

/**
 * @typedef {Object} updateData
 * @property {number} brand_id - Идентификатор бренда
 * @property {string} brandName - Название бренда
 * @property {string} category_id - Идентификатор категории товара
 * @property {number} gender - Пол (1 - муж, 2 - жен, 3 - уни)
 * @property {number} price - Цена товара
 * @property {number} discount - Процент скидки
 * @property {string} description - Описание товара
 */
async function update() {
    const id = document.getElementById("item-id").value

    /**
    * @type {updateData}
    */
    let updateData = {
        brand_id: parseInt(document.getElementById('brand').value),
        name: document.getElementById('item-name').value,
        category_id: parseInt(document.getElementById('category').value),
        sex: parseInt(document.getElementById('gender').value),
        price: parseInt(document.getElementById('price').value),
        discount: parseInt(document.getElementById('discount').value),
        description: document.getElementById('description').value
    }

    console.log(updateData, id);

    try {
        const response = await fetch(`http://localhost:8080/item/update/${id}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updateData)
        });

        if (response.ok) {
            window.location.replace('/admin/')
        } else {
            alert("Не удалось обновить данные.");
            throw new Error(`Ошибка HTTP: ${response.status}`);
        }
    } catch (error) {
        console.error('Ошибка при отправке данных: ', error.message)
    }
}

document.addEventListener('DOMContentLoaded', () => {
    fetchItem()

    const updateBtn = document.getElementById("update_btn")
    updateBtn.addEventListener('click', update)
});