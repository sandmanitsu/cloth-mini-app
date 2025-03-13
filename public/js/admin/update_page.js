import { formatDate } from "./date.js";
import { fetchBrands } from "./brand.js";
import { fetchCategory } from "./category.js";
import { getImage } from './image.js';

const IMAGE_GALLERY_WIDHT = 323
const IMAGE_GALLERY_HEIGHT = 430

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
    let genderOptions = `<option value="${item.sex}">${sex}</option>
    <option value="1">Мужской</option>
    <option value="2">Женский</option>
    <option value="3">Унисекс</option>`

    // заголовок с id
    document.getElementById('item-id-text').innerHTML = `ID: ${item.id} | Создан: ${formatDate(item.created_at)} | Обновлен: ${formatDate(item.updated_at)}`
    document.getElementById('item-id').value = item.id

    // опции для брендов
    document.getElementById('brand').innerHTML = brandOptions

    // опции для категории
    document.getElementById('category').innerHTML = categoryOptions

    // опции для выбора пола
    document.getElementById('gender').innerHTML = genderOptions

    // value для поля Название
    document.getElementById('item-name').value = item.name

    // value для поля Цена
    document.getElementById('price').value = item.price

    // value для скидки
    document.getElementById('discount').value = item.discount

    // description
    document.getElementById('description').value = item.description

    // подставляет изображения вместо мокового, если такое есть
    let imageId = ''
    if (item?.image_id && Array.isArray(item.image_id) && item.image_id.length > 0) {
        imageId = item.image_id[0]
    }

    getImage(imageId)
        .then((base64Image) => {
            if (base64Image == '') {
                return
            }

            document.getElementById('image-main').src = base64Image
    })

    // вставляем все изображения в галлерею
    if (item.image_id.length > 0) {
        item.image_id.forEach(image => {
            getImage(image)
                .then((base64Image) => {
                    if (base64Image == '') {
                        return
                    }

                    const img = document.createElement('img')
                    img.setAttribute('id', image);
                    img.width = IMAGE_GALLERY_WIDHT
                    img.height = IMAGE_GALLERY_HEIGHT
                    img.alt = 'image'
                    img.src = base64Image

                    document.getElementById('image-gallery').appendChild(img);

                    createDeleteBtn(image)
                })
        })
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

        try {
            const response = await fetch(`http://localhost:8080/image/delete?image_id=${imageId}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                throw new Error(`Ошибка при удалении изображения: ${response.status}`);
            }

            const imageElement = document.getElementById(imageId);
            if (imageElement) {
                imageElement.remove();
            }

            event.target.remove();
        } catch (error) {
            console.error('Ошибка при удалении изображения:', error);
        }
    });

    // Вставка кнопки в контейнер
    const div = document.createElement('div');
    div.classList.add('three', 'columns');
    div.appendChild(button);

    container.appendChild(div);
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

async function uploadImage(event) {
    event.preventDefault();

    let formData = new FormData(event.target)
    let id = document.getElementById('item-id').value

    try {
        const response = await fetch(`http://localhost:8080/image/create?itemId=${id}`, {
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error(`Ошибка: ${response.statusText}`);
        }

        const fileid = await response.json();

        getImage(fileid.file_id)
            .then((base64Image) => {
                const img = document.createElement('img')
                img.setAttribute('id', fileid.file_id);
                img.width = IMAGE_GALLERY_WIDHT
                img.height = IMAGE_GALLERY_HEIGHT
                img.alt = 'image'
                img.src = base64Image

                document.getElementById('image-gallery').appendChild(img);

                createDeleteBtn(fileid.file_id)
            })
    } catch (error) {
        console.error('Ошибка при загрузке изображения:', error);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    fetchItem()

    const updateBtn = document.getElementById("update_btn")
    updateBtn.addEventListener('click', update)

    // переключение форм с редактирование параметров и загрузкой изображения
    document.getElementById('update_image_btn').addEventListener('click', function() {
        document.getElementById('item').style.display = 'none'
        document.getElementById('image_update').style.display = 'block'
    })
    // переключение форм с загрузки изображения на редактирование параметров
    document.getElementById('back_btn').addEventListener('click', function() {
        document.getElementById('item').style.display = 'block';
        document.getElementById('image_update').style.display = 'none';
    });

    // загрузка изображения
    document.getElementById('image-form').addEventListener('submit', (event) => {
        uploadImage(event)
    })
});