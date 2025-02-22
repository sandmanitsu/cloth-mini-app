// Форматирует дату из ISO в формат ГГГГ-ММ-ДД - ЧЧ:ММ:СС
export function formatDate(dt) {
    if (!dt) {
        return ""
    }

    const date = new Date(dt);

    return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} - ${date.toLocaleTimeString()}`;
}