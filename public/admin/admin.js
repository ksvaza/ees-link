const applicationsContainer = document.getElementById('applicationsContainer');
const loadButton = document.getElementById('loadApplications');

function renderApplications(applications) {
    if (!Array.isArray(applications) || applications.length === 0) {
        applicationsContainer.innerHTML = '<p>Nav pieejamu pieteikumu.</p>';
        return;
    }

    const rows = applications.map(app => {
        const status = app.status || 'nav statusa';
        const id = app.id || app._id || '-';

        return `<tr>
            <td>${id}</td>
            <td>${app.teamName || '—'}</td>
            <td>${app.school || '—'}</td>
            <td>${app.members || '—'}</td>
            <td>${app.supervisor || '—'}</td>
            <td>${status}</td>
            <td><button class="btn btn-secondary" data-id="${id}" data-current-status="${status}">Apstiprināt</button></td>
        </tr>`;
    }).join('');

    applicationsContainer.innerHTML = `
        <table>
            <thead>
                <tr><th>ID</th><th>Komanda</th><th>Skola</th><th>Dalībnieki</th><th>Mentors</th><th>Statuss</th><th>Darbība</th></tr>
            </thead>
            <tbody>${rows}</tbody>
        </table>
        <p class="small-muted">Nospied "Apstiprināt" lai pārietu statusu.</p>
    `;

    applicationsContainer.querySelectorAll('button[data-id]').forEach(button => {
        button.addEventListener('click', async () => {
            const id = button.dataset.id;
            const currentStatus = button.dataset.currentStatus || 'pending';
            const newStatus = currentStatus === 'accepted' ? 'in_progress' : 'accepted';

            try {
                const response = await fetch('/api/applications/' + encodeURIComponent(id), {
                    method: 'PATCH',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ status: newStatus })
                });

                if (!response.ok) {
                    throw new Error('Neizdevās atjaunināt statusu');
                }

                await loadApplications();
            } catch (error) {
                alert('Kļūda statusa atjaunināšanā.');
                console.error(error);
            }
        });
    });
}

async function loadApplications() {
    applicationsContainer.innerHTML = '<p>Notiek ielāde...</p>';

    try {
        const response = await fetch('/api/applications');
        if (!response.ok) throw new Error('Neizdevās ielādēt');
        const data = await response.json();
        renderApplications(data);
    } catch (error) {
        applicationsContainer.innerHTML = '<p class="error-message">Neizdevās ielādēt pieteikumus.</p>';
        console.error(error);
    }
}

loadButton.addEventListener('click', loadApplications);
