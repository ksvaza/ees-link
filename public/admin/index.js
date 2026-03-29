document.getElementById('applicationForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = {
        teamName: document.getElementById('teamName').value,
        school: document.getElementById('school').value,
        members: parseInt(document.getElementById('members').value),
        supervisor: document.getElementById('supervisor').value
    };
    
    try {
        const response = await fetch('/api/applications', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });
        
        if (response.ok) {
            document.getElementById('successMessage').style.display = 'block';
            document.getElementById('applicationForm').reset();
            setTimeout(() => {
                document.getElementById('successMessage').style.display = 'none';
            }, 5000);
        } else {
            alert('Error submitting application');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('An error occurred while submitting your application');
    }
});
