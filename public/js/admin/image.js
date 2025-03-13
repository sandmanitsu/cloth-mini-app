
export async function getImage(imageId) {
    if (!imageId) {
        return ''
    }

    try {
        let response = await fetch(`http://localhost:8080/image/get/${imageId}`)

        if (!response.ok) {
            throw new Error(`Ошибка: ${response.status}`);
        }

        const blob = await response.blob();

        const base64Data = await new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.readAsDataURL(blob);
            reader.onload = () => resolve(reader.result);
            reader.onerror = (error) => reject(error);
        });

        return base64Data;
    } catch (error) {
        console.error('Ошибка получения изображения:', error);
    }
}