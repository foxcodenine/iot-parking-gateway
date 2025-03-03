// Access Level Utility

// Returns a descriptive string based on the provided access level
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
    getAccessLevelDescriptionUtil,
};
