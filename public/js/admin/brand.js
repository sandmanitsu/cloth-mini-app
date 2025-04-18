export async function optionBrands(select) {
    try {
        const url = `http://localhost:8081/brand/get`;
        const response = await fetch(url)

        if (!response.ok) {
            throw new Error(`fetchBrands Ошибка HTTP: ${response.status}`)
        }

        const category = await response.json()
        // const select = document.getElementById("brand-search")

        let html = `<option value=""></option>`
        category.forEach(cat => {
            html += `<option value="${cat.brand_id}">${cat.brand_name}</option>`
        })

        select.insertAdjacentHTML('beforeend', html);
    } catch (error) {
        console.error('fetchBrands Ошибка', error.message)
    }
}

export async function fetchBrands() {
    try {
        const url = `http://localhost:8081/brand/get`;
        const response = await fetch(url)

        if (!response.ok) {
            throw new Error(`fetchBrands Ошибка HTTP: ${response.status}`)
        }

        return await response.json()
    } catch (error) {
        console.error('fetchBrands Ошибка', error.message)
    }
}