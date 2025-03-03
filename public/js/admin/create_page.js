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

async function uploadImage(event) {
    event.preventDefault();

    let formData = new FormData(event.target)

    try {
        const response = await fetch(`http://localhost:8080/image/temp`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`Ошибка: ${response.statusText}`);
        }

        const resp = await response.json();
        console.log(resp);

        // getImage(fileid.file_id)
        //     .then((base64Image) => {
        //         const img = document.createElement('img')
        //         img.setAttribute('id', fileid.file_id);
        //         img.width = IMAGE_GALLERY_WIDHT
        //         img.height = IMAGE_GALLERY_HEIGHT
        //         img.alt = 'image'
        //         img.src = base64Image

        //         document.getElementById('image-gallery').appendChild(img);

        //         createDeleteBtn(fileid.file_id)
        //     })
    } catch (error) {
        console.error('Ошибка при загрузке изображения:', error);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    optionCategory(document.getElementById("category"))
    optionBrands(document.getElementById("brand"))

    const createBtn = document.getElementById("create_btn")
    createBtn.addEventListener('click', create)

    // переключение форм с редактирование параметров и загрузкой изображения
    document.getElementById('upload_btn').addEventListener('click', function() {
        document.getElementById('item').style.display = 'none'
        document.getElementById('image_upload').style.display = 'block'
    })
    // переключение форм с загрузки изображения на редактирование параметров
    document.getElementById('back_btn').addEventListener('click', function() {
        document.getElementById('item').style.display = 'block';
        document.getElementById('image_upload').style.display = 'none';
    });

    // загрузка изображения
    document.getElementById('image-form').addEventListener('submit', (event) => {
        uploadImage(event)
    })
});