export async function optionCategory(select) {
    try {
        const url = `http://localhost:8081/category/get`;
        const response = await fetch(url)

        if (!response.ok) {
            throw new Error(`fetchCategory Ошибка HTTP: ${response.status}`)
        }

        const category = await response.json()
        // const select = document.getElementById("category-search")

        let html = `<option value="">Все категории</option>`
        category.forEach(cat => {
            html += `<option value="${cat.category_id}">${cat.category_name}</option>`
        })

        select.insertAdjacentHTML('beforeend', html);
    } catch (error) {
        console.error('fetchCategory Ошибка', error.message)
    }
}

export async function fetchCategory() {
    try {
        const url = `http://localhost:8081/category/get`;
        const response = await fetch(url)

        if (!response.ok) {
            throw new Error(`fetchCategory Ошибка HTTP: ${response.status}`)
        }

        return await response.json()
    } catch (error) {
        console.error('fetchCategory Ошибка', error.message)
    }
}