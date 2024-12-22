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

function timeElapsed(utcDateString) {
    const parsedDate = new Date(utcDateString); // Parse the UTC date string
    const now = new Date(); // Get the current date and time
    const diffMs = now - parsedDate; // Calculate the difference in milliseconds

    if (diffMs < 0) {
        return "In the future";
    }

    const seconds = Math.floor(diffMs / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);
    const weeks = Math.floor(days / 7);
    const years = Math.floor(days / 365);

    // Handle special case for devices with no data for over 10 years
    if (years > 10) {
        return "Device has not sent any data";
    }

    const remainingWeeks = weeks % 52;
    const remainingDays = days % 7;
    const remainingHours = hours % 24;

    if (seconds < 60) {
        return "A few seconds ago";
    } else if (minutes < 60) {
        return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
    } else if (hours < 24) {
        return `${hours} hour${hours > 1 ? 's' : ''}, ${Math.floor(minutes % 60)} minute${minutes % 60 > 1 ? 's' : ''} ago`;
    } else if (days < 7) {
        return `${days} day${days > 1 ? 's' : ''}, ${remainingHours} hour${remainingHours > 1 ? 's' : ''} ago`;
    } else if (weeks < 52) {
        return `${weeks} week${weeks > 1 ? 's' : ''}, ${remainingDays} day${remainingDays > 1 ? 's' : ''} ago`;
    } else {
        return `${years} year${years > 1 ? 's' : ''}, ${remainingWeeks} week${remainingWeeks > 1 ? 's' : ''} ago`;
    }
}

function formatToLocalDateTime(utcDateString) {
    const parsedDate = new Date(utcDateString);

    if (isNaN(parsedDate)) {
        return null; // Invalid date, return null
    }

    const now = new Date();
    const diffMs = now - parsedDate;

    // Check if the elapsed time is more than 10 years
    const tenYearsInMs = 10 * 365 * 24 * 60 * 60 * 1000; // Approximation for 10 years
    if (diffMs > tenYearsInMs) {
        return null; // More than 10 years ago, return null
    }

    // Custom format for the date
    const options = { month: "long" };
    const month = new Intl.DateTimeFormat("en-US", options).format(parsedDate);
    const day = parsedDate.getDate();
    const year = parsedDate.getFullYear();
    const hours = parsedDate.getHours() % 12 || 12; // Convert to 12-hour format
    const minutes = parsedDate.getMinutes().toString().padStart(2, "0");
    const ampm = parsedDate.getHours() >= 12 ? "PM" : "AM";

    return `${month} ${day}, ${year} ${hours}:${minutes} ${ampm}`;
}


export {
    formatDateUtil,
    getAccessLevelDescriptionUtil,
    timeElapsed,
    formatToLocalDateTime,
};
