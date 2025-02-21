
async function fetchItems(limit = 20, offset = 0) {
    try {
        const url = `http://localhost:8080/item/get?limit=${limit}&offset=${offset}`;
        const response = await fetch(url)

        if (!response.ok) {
            throw new Error(`fetchItems Ошибка HTTP: ${response.status}`)
        }

        const items = await response.json()

        const container = document.getElementById("products-body")

        items.items.forEach((product) => {
            const itemCard = `
            <tr>
                <td>${product.id}</td>
                <td>${product.brand}</td>
                <td>${product.name}</td>
                <td>${product.description}</td>
                <td>${product.category_name}</td>
                <td>${product.price} руб.</td>
                <td><a href="${product.outer_link}" target="_blank">Товар в магазине</a></td>
                <td><button class="edit-button">Редактировать</button></td>
            </tr>`;

        container.insertAdjacentHTML('beforeend', itemCard);
        });
    } catch (error) {
        console.error('fetchItems Ошибка', error.message)
    }
}

document.addEventListener('DOMContentLoaded', () => {
    fetchItems()
});