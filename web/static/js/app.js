// Theme toggle
const themeToggle = document.querySelector('.theme-toggle');
const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');

function setTheme(isDark) {
    document.documentElement.classList.toggle('dark', isDark);
    localStorage.setItem('theme', isDark ? 'dark' : 'light');
}

// Initialize theme
const savedTheme = localStorage.getItem('theme');
if (savedTheme) {
    setTheme(savedTheme === 'dark');
} else {
    setTheme(prefersDark.matches);
}

themeToggle.addEventListener('click', () => {
    setTheme(!document.documentElement.classList.contains('dark'));
});

// API Interaction
const API_BASE = '/v1';
let currentToken = null;

async function makeRequest(endpoint, data = null) {
    const headers = {
        'Content-Type': 'application/json'
    };

    if (currentToken) {
        headers['Authorization'] = `Bearer ${currentToken}`;
    }

    console.log('Making request:', {
        endpoint,
        headers,
        method: data ? 'POST' : 'GET'
    });

    try {
        const response = await fetch(API_BASE + endpoint, {
            method: data ? 'POST' : 'GET',
            headers,
            body: data ? JSON.stringify(data) : null
        });

        const responseText = await response.text();
        console.log('Raw response:', responseText);

        if (!response.ok) {
            throw new Error(responseText || `HTTP error! status: ${response.status}`);
        }

        const result = JSON.parse(responseText);
        console.log('Parsed response:', result);
        return result;
    } catch (error) {
        console.error('Request failed:', error);
        throw error;
    }
}

// UI Elements
const generateTokenBtn = document.getElementById('generate-token');
const validateAddressBtn = document.getElementById('validate-address');
const resolveEnsBtn = document.getElementById('resolve-ens');
const checkContractBtn = document.getElementById('check-contract');

const addressInput = document.getElementById('eth-address');
const ensInput = document.getElementById('ens-name');
const contractInput = document.getElementById('contract-address');

const tokenDisplay = document.getElementById('jwt-token');

// Copy functionality
document.querySelectorAll('.copy-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        const text = document.getElementById(btn.dataset.clipboard).textContent;
        navigator.clipboard.writeText(text).then(() => {
            const originalHTML = btn.innerHTML;
            btn.innerHTML = '<svg viewBox="0 0 24 24"><path d="M20 6L9 17l-5-5"/></svg>';
            setTimeout(() => btn.innerHTML = originalHTML, 2000);
        });
    });
});

// Disable all buttons except generate token initially
[validateAddressBtn, resolveEnsBtn, checkContractBtn].forEach(btn => {
    btn.disabled = true;
});

// Generate Token
generateTokenBtn.addEventListener('click', async () => {
    try {
        const result = await makeRequest('/token', {});
        console.log('Token result:', result);
        currentToken = result.token;
        tokenDisplay.textContent = currentToken;
        document.querySelector('.token-display').classList.remove('hidden');

        // Enable other buttons
        [validateAddressBtn, resolveEnsBtn, checkContractBtn].forEach(btn => {
            btn.disabled = false;
        });

        showResult('Token generated successfully', false);
    } catch (error) {
        console.error('Token generation failed:', error);
        showResult(error.message, true);
    }
});

// Validate Address
validateAddressBtn.addEventListener('click', async () => {
    const address = addressInput.value.trim();
    if (!address) {
        showResult('Please enter an Ethereum address', true);
        return;
    }

    try {
        console.log('Validating address with token:', currentToken);
        const result = await makeRequest('/validate/' + address);
        console.log('Validation result:', result);

        showResult(
            result.isValid
                ? `Address ${address} is valid`
                : `Address ${address} is invalid`,
            !result.isValid
        );
    } catch (error) {
        console.error('Validation failed:', error);
        showResult(error.message, true);
    }
});

// Resolve ENS
resolveEnsBtn.addEventListener('click', async () => {
    const ens = ensInput.value.trim();
    if (!ens) {
        showResult('Please enter an ENS name', true);
        return;
    }

    try {
        const result = await makeRequest('/resolveEns/' + ens);
        showResult(`Resolved address: ${result.address}`, false);
    } catch (error) {
        showResult(error.message, true);
    }
});

// Check Contract
checkContractBtn.addEventListener('click', async () => {
    const address = contractInput.value.trim();
    if (!address) {
        showResult('Please enter an Ethereum address', true);
        return;
    }

    try {
        const result = await makeRequest('/isContract/' + address);
        console.log('Contract check result:', result); // Debug log

        // The server returns is_contract, but we're checking isContract
        const isContract = result.is_contract || result.isContract;

        showResult(
            isContract
                ? `Address ${address} is a contract`
                : `Address ${address} is not a contract`,
            false
        );
    } catch (error) {
        console.error('Contract check failed:', error);
        showResult(error.message, true);
    }
});

// Helper function to show results
function showResult(message, isError = false) {
    const resultDiv = document.createElement('div');
    resultDiv.className = `result ${isError ? 'error' : 'success'}`;
    resultDiv.textContent = message;

    const section = document.activeElement.closest('.tool-section');
    const existingResult = section.querySelector('.result');
    if (existingResult) {
        existingResult.remove();
    }
    section.appendChild(resultDiv);
} 