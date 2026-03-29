async function loadNews() {
    const newsContainer = document.getElementById('newsItems');
    newsContainer.innerHTML = '<p>Notiek ielāde...</p>';

    try {
        const response = await fetch('/api/news');
        if (!response.ok) {
            throw new Error('Neizdevās ielādēt ziņas');
        }
        const data = await response.json();
        if (!Array.isArray(data) || data.length === 0) {
            newsContainer.innerHTML = '<p>Nav pieejamu ziņu.</p>';
            return;
        }

        newsContainer.innerHTML = data.slice(0, 4).map(item => {
            const title = item.title || item.name || 'Bez nosaukuma';
            const date = item.date || item.publishedAt || 'Nav datuma';
            const preview = item.summary || item.description || 'Ziņa bez apraksta.';

            return `<article class="news-item"><h4>${title}</h4><small>${date}</small><p>${preview}</p></article>`;
        }).join('');
    } catch (error) {
        newsContainer.innerHTML = '<p class="error-message">Neizdevās ielādēt ziņas no servera.</p>';
        console.error(error);
    }
}

window.addEventListener('DOMContentLoaded', () => {
    loadNews();
});
