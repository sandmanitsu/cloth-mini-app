
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

    let brandOptions = `<option value=""></option>`
    brands.forEach(cat => {
        brandOptions += `<option value="${cat.brand_id}">${cat.brand_name}</option>`
    })

    let categoryOptions = `<option value="">Все категории</option>`
    category.forEach(cat => {
        categoryOptions += `<option value="${cat.category_id}">${cat.category_name}</option>`
    })

    let html = `
            <div class="container">
                <h6>ID: ${item.id} | Создан: ${formatDate(item.created_at)} | Обновлен: ${formatDate(item.updated_at)}</h6>
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
                        <label for="brand-name">Название:</label>
                        <input class="u-full-width"  type="text" id="brand-name" placeholder="${item.name}" />
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
                            <option value="${item.sex}"></option>
                            <option value="1">Мужской</option>
                            <option value="2">Женский</option>
                            <option value="3">Унисекс</option>
                        </select>
                    </div>

                    <div class="three columns">
                        <label for="price">Цена:</label>
                        <input class="u-full-width"  type="number" id="price" placeholder="${item.price}" />
                    </div>

                    <div class="two columns">
                        <label for="discount">Скидка:</label>
                        <input class="u-full-width"  type="number" id="discount" placeholder="${item.discount}" />
                    </div>
                </div>

                <div class="eight columns">
                    <label for="description">Описание:</label>
                    <input class="u-full-width"  type="text" id="description" placeholder="${item.description}" />
                </div>
            </div>`

            container.insertAdjacentHTML('beforeend', html)
}

document.addEventListener('DOMContentLoaded', () => {
    fetchItem()
});