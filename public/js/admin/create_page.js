import { optionBrands } from './brand.js';
import { optionCategory } from './category.js';

/**
 * @typedef  {Object} updateData
 * @property {number} brand_id - Идентификатор бренда
 * @property {string} brandName - Название бренда
 * @property {string} category_id - Идентификатор категории товара
 * @property {number} sex - Пол (1 - муж, 2 - жен, 3 - уни)
 * @property {number} price - Цена товара
 * @property {number} discount - Процент скидки
 * @property {string} description - Описание товара
 * @property {string} outer_link - ссылка на товар
 */
async function create() {
    /**
    * @type {createData}
    */
    let createData = {
        brand_id: parseInt(document.getElementById('brand').value),
        name: document.getElementById('item-name').value,
        category_id: parseInt(document.getElementById('category').value),
        sex: parseInt(document.getElementById('gender').value),
        price: parseInt(document.getElementById('price').value),
        discount: parseInt(document.getElementById('discount').value),
        description: document.getElementById('description').value,
        outer_link: document.getElementById('outer-link').value
    }

    console.log(createData);

    try {
        const response = await fetch(`http://localhost:8080/item/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(createData)
        });

        if (response.ok) {
            window.location.replace('/admin/')
        } else {
            alert("Не удалось создать товар.");
            throw new Error(`Ошибка HTTP: ${response.status}`);
        }
    } catch (error) {
        console.error('Ошибка при отправке данных: ', error.message)
    }
}

document.addEventListener('DOMContentLoaded', () => {
    optionCategory(document.getElementById("category"))
    optionBrands(document.getElementById("brand"))

    const createBtn = document.getElementById("create_btn")
    createBtn.addEventListener('click', create)
});