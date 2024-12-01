function formatDateUtil(dateString) {
    // Parse the input date string into a Date object
    const date = new Date(dateString);

    // Format the date to the desired format: e.g., "Nov 30, 2024"
    // 'en-US' locale can be changed to any locale as needed
    return date.toLocaleDateString('en-US', {
        month: 'short', // abbreviated month name
        day: 'numeric', // numeric day
        year: 'numeric' // four digit year
    });
}

function getAccessLevelDescriptionUtil(accessLevel) {
    switch (accessLevel) {
        case 0:
            return "Root";
        case 1:
            return "Admin";
        case 2:
            return "Editor";
        case 3:
            return "Viewer";
        default:
            return "Unknown Access Level"; // Handles any integer not 0-3
    }
}

export {
    formatDateUtil,
    getAccessLevelDescriptionUtil,
};
