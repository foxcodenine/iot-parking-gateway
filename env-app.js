const express = require('express');
const fs = require('fs');

const app = express();
const port = 9090;



function corsMiddleware(req, res, next) {
    // Allow any origin to access the resource
    res.setHeader('Access-Control-Allow-Origin', '*');

    res.setHeader('Access-Control-Allow-Methods', 'GET, OPTIONS');
    res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
    res.setHeader('Access-Control-Allow-Credentials', 'true');
    res.setHeader('Access-Control-Expose-Headers', ''); // Specify if there are any specific headers to expose

    // Handle preflight requests
    if (req.method === 'OPTIONS') {
        res.status(204).send();
    } else {
        next();
    }
}

// Use CORS middleware to allow cross-origin requests
app.use(corsMiddleware);

// Manually read the .env file and parse it into an object
const readEnvVariables = () => {
    const envVariables = {};
    try {
        const data = fs.readFileSync('.env', 'utf8');
        data.split('\n').forEach(line => {
            const [key, value] = line.split('=');
            if (key && value) {
                envVariables[key.trim()] = value.trim();
            }
        });
    } catch (err) {
        console.error('Failed to read .env file:', err);
    }
    return envVariables;
};

// Define a route for /env that returns environment variables
app.get('/env', (req, res) => {
    res.json(readEnvVariables());
});

// Catch-all route for handling 404 errors
app.use((req, res) => {
    res.status(404).send('Not Found');
});

// Start the server
app.listen(port, () => {
    console.log(`Server running on http://localhost:${port}/env`);
});
