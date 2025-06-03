import { optionBrands } from './brand.js';
import { optionCategory } from './category.js';
import { getImage } from './image.js';

const IMAGE_GALLERY_WIDHT = 323
const IMAGE_GALLERY_HEIGHT = 430
const DEFAULT_IMAGE_ADDR = "http://localhost:8081/admin/static/img/no_image.jpg"

// галлерея изображение {file_id => base64image}
let imagesIds = [];

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
 * @property {array} temp_images - галлерия изображений
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
        outer_link: document.getElementById('outer-link').value,
        temp_images: Object.keys(imagesIds),
    }

    try {
        const response = await fetch(`http://localhost:8081/item/create`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(createData)
        });

        if (response.ok) {
            // window.location.replace('/admin/')
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
    formData.append('uuid', crypto.randomUUID())

    try {
        const response = await fetch(`http://localhost:8081/image/temp`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`Ошибка: ${response.statusText}`);
        }

        const resp = await response.json();
        console.log(resp);

        getImage(resp.file_id)
            .then((base64Image) => {
                imagesIds[resp.file_id] = base64Image
                if (document.getElementById('image-main').src == DEFAULT_IMAGE_ADDR) {
                    document.getElementById('image-main').src = base64Image
                }

                const img = document.createElement('img')
                img.setAttribute('id', resp.file_id);
                img.width = IMAGE_GALLERY_WIDHT
                img.height = IMAGE_GALLERY_HEIGHT
                img.alt = 'image'
                img.src = base64Image

                document.getElementById('image-gallery').appendChild(img);

                createDeleteBtn(resp.file_id)
            })
    } catch (error) {
        console.error('Ошибка при загрузке изображения:', error);
    }
}

function createDeleteBtn(image) {
    const container = document.getElementById('delete-btns')

    // создаем кнопку
    const button = document.createElement('button');
    button.classList.add('u-full-width');
    button.setAttribute('id', 'delete_btn');
    button.setAttribute('image_id', image);
    button.textContent = 'Удалить';

    // Привязываем слушатель событий к кнопке
    button.addEventListener('click', async (event) => {
        const imageId = event.target.getAttribute('image_id');
        const imageElement = document.getElementById(imageId);
        if (imageElement) {
            imageElement.remove();
        }

        event.target.remove();

        // меняем/удаляем мейн изображение
        if (imagesIds[image]) {
            delete imagesIds[image]

            const entries = Object.entries(imagesIds)

            if (entries.length > 0) {
                const [file_id, base64image] = entries[0];
                document.getElementById('image-main').src = base64image
            } else {
                document.getElementById('image-main').src = DEFAULT_IMAGE_ADDR
            }
        }
    });

    // Вставка кнопки в контейнер
    const div = document.createElement('div');
    div.classList.add('three', 'columns');
    div.appendChild(button);

    container.appendChild(div);
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