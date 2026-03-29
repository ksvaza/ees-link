document.getElementById('applicationForm').addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = {
        teamName: document.getElementById('teamName').value.trim(),
        school: document.getElementById('school').value.trim(),
        members: Number(document.getElementById('members').value),
        supervisor: document.getElementById('supervisor').value.trim(),
    };

    if (!formData.teamName || !formData.school || !formData.members || !formData.supervisor) {
        alert('Lūdzu aizpildiet visus obligātos laukus.');
        return;
    }

    try {
        const response = await fetch('/api/applications', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            throw new Error('Servera kļūda');
        }

        document.getElementById('successMessage').style.display = 'block';
        document.getElementById('applicationForm').reset();
        setTimeout(() => {
            document.getElementById('successMessage').style.display = 'none';
        }, 6000);
    } catch (error) {
        console.error('Kļūda:', error);
        alert('Notika kļūda pieteikuma nosūtīšanas laikā. Lūdzu mēģiniet vēlreiz.');
    }
});
