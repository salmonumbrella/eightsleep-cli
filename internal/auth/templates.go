package auth

const setupTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Eight Sleep CLI - Connect Your Pod</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #ffffff;
            --bg-secondary: #f8f9fa;
            --bg-input: #ffffff;
            --text-primary: #1a1a1a;
            --text-secondary: #666666;
            --text-muted: #999999;
            --border: #e5e5e5;
            --border-focus: #1a1a1a;
            --accent: #1a1a1a;
            --accent-hover: #333333;
            --success: #22c55e;
            --success-bg: #f0fdf4;
            --error: #ef4444;
            --error-bg: #fef2f2;
            --radius: 12px;
            --radius-sm: 8px;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        .hidden {
            display: none !important;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg);
            color: var(--text-primary);
            min-height: 100vh;
            line-height: 1.5;
            -webkit-font-smoothing: antialiased;
        }

        .container {
            max-width: 400px;
            margin: 0 auto;
            padding: 80px 24px 40px;
            animation: fadeIn 0.4s ease-out;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(8px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Header with Eight Sleep Logo */
        header {
            text-align: center;
            margin-bottom: 48px;
        }

        .logo {
            margin-bottom: 24px;
        }

        .logo svg {
            height: 40px;
            width: auto;
        }

        .subtitle {
            color: var(--text-secondary);
            font-size: 15px;
            font-weight: 400;
        }

        /* Form Card */
        .auth-card {
            background: var(--bg);
            border: 1px solid var(--border);
            border-radius: var(--radius);
            padding: 32px;
        }

        .card-title {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 8px;
        }

        .card-description {
            font-size: 14px;
            color: var(--text-secondary);
            margin-bottom: 24px;
        }

        /* Form */
        .form-group {
            margin-bottom: 20px;
        }

        .form-label {
            display: block;
            font-size: 14px;
            font-weight: 500;
            color: var(--text-primary);
            margin-bottom: 6px;
        }

        .form-input {
            width: 100%;
            padding: 12px 14px;
            background: var(--bg-input);
            border: 1px solid var(--border);
            border-radius: var(--radius-sm);
            color: var(--text-primary);
            font-family: inherit;
            font-size: 15px;
            transition: border-color 0.15s ease, box-shadow 0.15s ease;
        }

        .form-input:hover {
            border-color: #cccccc;
        }

        .form-input:focus {
            outline: none;
            border-color: var(--border-focus);
            box-shadow: 0 0 0 3px rgba(0, 0, 0, 0.05);
        }

        .form-input::placeholder {
            color: var(--text-muted);
        }

        .form-hint {
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            color: var(--text-muted);
            margin-top: 8px;
        }

        .form-hint svg {
            width: 14px;
            height: 14px;
            color: var(--success);
        }

        /* Buttons */
        .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            padding: 12px 20px;
            border-radius: var(--radius-sm);
            font-family: inherit;
            font-size: 15px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.15s ease;
            border: none;
            width: 100%;
        }

        .btn-primary {
            background: var(--accent);
            color: white;
        }

        .btn-primary:hover {
            background: var(--accent-hover);
        }

        .btn-primary:active {
            transform: scale(0.98);
        }

        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none !important;
        }

        .spinner {
            width: 16px;
            height: 16px;
            border: 2px solid transparent;
            border-top-color: currentColor;
            border-radius: 50%;
            animation: spin 0.7s linear infinite;
        }

        @keyframes spin {
            to { transform: rotate(360deg); }
        }

        /* Status */
        .status {
            padding: 12px 14px;
            border-radius: var(--radius-sm);
            font-size: 14px;
            margin-top: 16px;
            display: none;
            align-items: center;
            gap: 10px;
        }

        .status.visible {
            display: flex;
            animation: fadeIn 0.2s ease-out;
        }

        .status.success {
            background: var(--success-bg);
            color: #166534;
        }

        .status.error {
            background: var(--error-bg);
            color: #991b1b;
        }

        .status.loading {
            background: var(--bg-secondary);
            color: var(--text-secondary);
        }

        .status svg {
            width: 16px;
            height: 16px;
            flex-shrink: 0;
        }

        /* Info notice */
        .info-notice {
            background: var(--bg-secondary);
            border-radius: var(--radius-sm);
            padding: 14px;
            margin-bottom: 24px;
            font-size: 13px;
            color: var(--text-secondary);
            line-height: 1.5;
        }

        .info-notice strong {
            color: var(--text-primary);
            font-weight: 500;
        }

        /* Footer */
        footer {
            text-align: center;
            margin-top: 32px;
            padding-top: 24px;
        }

        .footer-links {
            display: flex;
            justify-content: center;
            gap: 24px;
            margin-bottom: 16px;
        }

        .footer-links a {
            display: inline-flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            color: var(--text-muted);
            text-decoration: none;
            transition: color 0.15s ease;
        }

        .footer-links a:hover {
            color: var(--text-primary);
        }

        .footer-links svg {
            width: 16px;
            height: 16px;
        }

        .footer-credit {
            font-size: 12px;
            color: var(--text-muted);
        }

        .footer-credit a {
            color: var(--text-secondary);
            text-decoration: none;
        }

        .footer-credit a:hover {
            color: var(--text-primary);
        }

        @media (max-width: 480px) {
            .container {
                padding: 60px 20px 32px;
            }

            .auth-card {
                padding: 24px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <div class="logo">
                <!-- Eight Sleep Logo -->
                <svg viewBox="0 0 124.3 50" fill="currentColor">
                    <path d="M24.5,39.3c0,1.7-1.4,3.1-3.2,3.1h-8.9c-1.8,0-3.2-1.4-3.2-3.1v-5.7c0-1.7,1.4-3.1,3.2-3.1h8.9c1.8,0,3.2,1.4,3.2,3.1V39.3z M9.2,10.7c0-1.7,1.4-3.1,3.2-3.1h8.9c1.8,0,3.2,1.4,3.2,3.1v4.1c0,1.7-1.4,3.1-3.2,3.1h-8.9c-1.8,0-3.2-1.4-3.2-3.1L9.2,10.7z M33.7,15.2V8.6c0-4.7-3.9-8.6-8.8-8.6H8.8C3.9,0,0,3.8,0,8.6v6.6c0,3.5,2.1,6.4,5.2,7.8v1.2C2.1,25.5,0,28.5,0,32v9.5C0,46.2,3.9,50,8.8,50h16.1c4.9,0,8.8-3.8,8.8-8.6V32c0-3.5-2.1-6.4-5.2-7.8V23C31.6,21.6,33.7,18.6,33.7,15.2z"/>
                    <path d="M76.7,46.4h12.7v-4h-8.4v-3.1h6.8v-4h-6.8v-3.1h8.4v-4H76.7V46.4z"/>
                    <path d="M102.7,12.9h-4.4V5.6h-4.4v18.2h4.4V17h4.4v6.9h4.4V5.6h-4.4V12.9z"/>
                    <path d="M93.9,46.4h12.7v-4h-8.4v-3.1h6.8v-4h-6.8v-3.1h8.4v-4H93.9V46.4z"/>
                    <path d="M118.3,36.1H115V32h3.3c1,0,1.8,0.3,1.8,1.9C120.1,35.4,119.5,36.1,118.3,36.1 M118.8,28.1h-8.2v18.2h4.3v-6.5h3.8c2.1,0,5.5-1.2,5.5-5.9C124.3,29,121,28.1,118.8,28.1"/>
                    <path d="M83,16.9h2.6c0,2.7-1,3.4-2.8,3.4c-2,0-2.7-1.1-2.7-3.4v-5.5c0-1.5,1.1-2.4,2.7-2.4c1.8,0,2.5,1.1,2.5,2.4h4.5c0-3.1-2.2-6.2-7-6.2c-4.9,0-7.1,2.4-7.1,6.2v5.5c0,3.1,1,7.1,7.1,7.1c5.8,0,7.1-3.4,7.1-7v-3.5H83L83,16.9z"/>
                    <path d="M52,35.3c-2.1-0.3-2.6-1.3-2.6-1.8c0-1.2,0.9-1.8,2.3-1.8c1.6,0,2.1,1,2.1,1.6h4.2c0-2.1-1.5-5.6-6.2-5.6c-5.3,0-6.7,3.4-6.7,5.9c0,2.9,2.4,4.8,6,5.3c1.2,0.2,2.9,0.5,2.9,2c0,1.2-0.7,1.9-2.3,1.9c-1,0-2.3-0.2-2.3-2.4H45c0,2.3,0.6,6.2,6.7,6.2c6,0,6.7-4.1,6.7-5.8C58.4,37.3,55.9,36,52,35.3"/>
                    <path d="M58,19.7h-8.4v-3.1h6.8v-4h-6.8V9.6H58v-4H45.3v18.2H58V19.7z"/>
                    <path d="M66.4,28.1h-4.3v18.2h10.7v-4.2h-6.3V28.1z"/>
                    <path d="M115,23.8h4.5v-14h4V5.6h-12.6v4.2h4.1V23.8z"/>
                    <path d="M65.1,19.9h-3v4h10.5v-4h-3.1V9.6h3v-4H62v4h3.1V19.9z"/>
                </svg>
            </div>
            <p class="subtitle">Control your Pod from the terminal</p>
        </header>

        <div class="auth-card">
            <h2 class="card-title">Connect Your Account</h2>
            <p class="card-description">Sign in with your Eight Sleep credentials</p>

            <div class="info-notice">
                Use your <strong>Eight Sleep app credentials</strong>. Your credentials are stored securely in your system keychain.
            </div>

            <form id="authForm">
                <div class="form-group">
                    <label class="form-label" for="email">Email</label>
                    <input
                        type="email"
                        id="email"
                        name="email"
                        class="form-input"
                        placeholder="you@example.com"
                        required
                        autocomplete="email"
                    >
                </div>

                <div class="form-group">
                    <label class="form-label" for="password">Password</label>
                    <input
                        type="password"
                        id="password"
                        name="password"
                        class="form-input"
                        placeholder="Your password"
                        required
                        autocomplete="current-password"
                    >
                    <p class="form-hint">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                            <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                        </svg>
                        Stored in your system keychain
                    </p>
                </div>

                <div id="status" class="status"></div>

                <button type="submit" id="submitBtn" class="btn btn-primary">
                    Sign In
                </button>
            </form>
        </div>

        <footer>
            <div class="footer-links">
                <a href="https://www.eightsleep.com" target="_blank">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                        <polyline points="15 3 21 3 21 9"/>
                        <line x1="10" y1="14" x2="21" y2="3"/>
                    </svg>
                    Eight Sleep
                </a>
                <a href="https://github.com/salmonumbrella/eightsleep-cli" target="_blank">
                    <svg viewBox="0 0 24 24" fill="currentColor">
                        <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                    GitHub
                </a>
            </div>
            <p class="footer-credit">
                Based on <a href="https://github.com/steipete/eightctl" target="_blank">eightctl</a> by Peter Steinberger
            </p>
        </footer>
    </div>

    <script>
        const csrfToken = '{{.CSRFToken}}';
        const form = document.getElementById('authForm');
        const emailInput = document.getElementById('email');
        const passwordInput = document.getElementById('password');
        const submitBtn = document.getElementById('submitBtn');
        const status = document.getElementById('status');

        function showStatus(message, type) {
            const icons = {
                success: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>',
                error: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>',
                loading: '<span class="spinner"></span>'
            };
            status.innerHTML = icons[type] + '<span>' + message + '</span>';
            status.className = 'status visible ' + type;
        }

        function hideStatus() {
            status.className = 'status';
        }

        function setLoading(loading) {
            if (loading) {
                submitBtn.disabled = true;
                submitBtn.innerHTML = '<span class="spinner"></span> Signing in...';
            } else {
                submitBtn.disabled = false;
                submitBtn.innerHTML = 'Sign In';
            }
        }

        form.addEventListener('submit', async (e) => {
            e.preventDefault();

            const email = emailInput.value.trim();
            const password = passwordInput.value;

            if (!email || !password) {
                showStatus('Please enter both email and password', 'error');
                return;
            }

            setLoading(true);
            showStatus('Authenticating...', 'loading');

            try {
                const response = await fetch('/submit', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-CSRF-Token': csrfToken
                    },
                    body: JSON.stringify({ email, password })
                });

                const data = await response.json();

                if (data.success) {
                    showStatus('Connected! Redirecting...', 'success');
                    setTimeout(() => {
                        window.location.href = '/success?email=' + encodeURIComponent(email);
                    }, 800);
                } else {
                    showStatus(data.error || 'Authentication failed', 'error');
                    setLoading(false);
                }
            } catch (err) {
                showStatus('Connection error: ' + err.message, 'error');
                setLoading(false);
            }
        });

        emailInput.focus();
    </script>
</body>
</html>`

const successTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Connected - Eight Sleep CLI</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg: #ffffff;
            --bg-secondary: #f8f9fa;
            --bg-terminal: #1a1a1a;
            --text-primary: #1a1a1a;
            --text-secondary: #666666;
            --text-muted: #999999;
            --border: #e5e5e5;
            --success: #22c55e;
            --success-bg: #f0fdf4;
            --radius: 12px;
            --radius-sm: 8px;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: var(--bg);
            color: var(--text-primary);
            min-height: 100vh;
            line-height: 1.5;
            -webkit-font-smoothing: antialiased;
        }

        .container {
            max-width: 480px;
            margin: 0 auto;
            padding: 80px 24px 40px;
            text-align: center;
            animation: fadeIn 0.4s ease-out;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(8px); }
            to { opacity: 1; transform: translateY(0); }
        }

        /* Success icon */
        .success-icon {
            width: 72px;
            height: 72px;
            margin: 0 auto 24px;
            background: var(--success);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            animation: scaleIn 0.4s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
        }

        @keyframes scaleIn {
            from { transform: scale(0); opacity: 0; }
            to { transform: scale(1); opacity: 1; }
        }

        .success-icon svg {
            width: 36px;
            height: 36px;
            color: white;
        }

        h1 {
            font-size: 28px;
            font-weight: 600;
            margin-bottom: 8px;
        }

        .subtitle {
            color: var(--text-secondary);
            font-size: 15px;
            margin-bottom: 32px;
        }

        /* User badge */
        .user-badge {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            background: var(--bg-secondary);
            border: 1px solid var(--border);
            border-radius: 100px;
            padding: 8px 16px;
            font-size: 14px;
            color: var(--text-secondary);
            margin-bottom: 32px;
        }

        .user-badge .dot {
            width: 8px;
            height: 8px;
            background: var(--success);
            border-radius: 50%;
            animation: pulse 2s ease-in-out infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }

        /* Terminal card */
        .terminal {
            background: var(--bg-terminal);
            border-radius: var(--radius);
            overflow: hidden;
            text-align: left;
            box-shadow: 0 4px 24px rgba(0, 0, 0, 0.1);
        }

        .terminal-bar {
            background: #2d2d2d;
            padding: 12px 16px;
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .terminal-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
        }

        .terminal-dot.red { background: #ff5f57; }
        .terminal-dot.yellow { background: #febc2e; }
        .terminal-dot.green { background: #28c840; }

        .terminal-title {
            flex: 1;
            text-align: center;
            font-family: 'JetBrains Mono', monospace;
            font-size: 12px;
            color: #888888;
        }

        .terminal-body {
            padding: 20px;
        }

        .terminal-line {
            display: flex;
            align-items: center;
            gap: 8px;
            font-family: 'JetBrains Mono', monospace;
            font-size: 13px;
            margin-bottom: 10px;
        }

        .terminal-line:last-child {
            margin-bottom: 0;
        }

        .terminal-prompt {
            color: #22c55e;
            user-select: none;
        }

        .terminal-text {
            color: #ffffff;
        }

        .terminal-cursor {
            display: inline-block;
            width: 8px;
            height: 16px;
            background: #ffffff;
            animation: blink 1s step-end infinite;
            margin-left: 2px;
            vertical-align: middle;
        }

        @keyframes blink {
            0%, 50% { opacity: 1; }
            50.01%, 100% { opacity: 0; }
        }

        .terminal-output {
            color: #888888;
            padding-left: 18px;
            margin-top: -4px;
            margin-bottom: 10px;
            font-family: 'JetBrains Mono', monospace;
            font-size: 12px;
        }

        /* Message */
        .message {
            margin-top: 24px;
            padding: 16px;
            background: var(--success-bg);
            border-radius: var(--radius-sm);
            text-align: left;
        }

        .message-header {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 6px;
        }

        .message-icon {
            width: 20px;
            height: 20px;
            color: var(--success);
        }

        .message-title {
            font-weight: 500;
            font-size: 14px;
            color: #166534;
        }

        .message-text {
            font-size: 13px;
            color: #166534;
            padding-left: 30px;
        }

        .message-text code {
            font-family: 'JetBrains Mono', monospace;
            background: rgba(0, 0, 0, 0.05);
            padding: 2px 6px;
            border-radius: 4px;
            font-size: 12px;
        }

        /* Footer */
        .footer {
            margin-top: 24px;
            font-size: 13px;
            color: var(--text-muted);
        }

        @media (max-width: 480px) {
            .container {
                padding: 60px 20px 32px;
            }

            h1 {
                font-size: 24px;
            }

            .terminal-body {
                padding: 16px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="success-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12"/>
            </svg>
        </div>

        <h1>You're Connected</h1>
        <p class="subtitle">Eight Sleep CLI is ready to control your Pod</p>

        {{if .UserEmail}}
        <div class="user-badge">
            <span class="dot"></span>
            <span>{{.UserEmail}}</span>
        </div>
        {{end}}

        <div class="terminal">
            <div class="terminal-bar">
                <span class="terminal-dot red"></span>
                <span class="terminal-dot yellow"></span>
                <span class="terminal-dot green"></span>
                <span class="terminal-title">Terminal</span>
            </div>
            <div class="terminal-body">
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-text">eightsleep status</span>
                </div>
                <div class="terminal-output">Pod Pro 3 - Active, 68°F</div>
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-text">eightsleep temp --side left +2</span>
                </div>
                <div class="terminal-output">Left side temperature set to +2</div>
                <div class="terminal-line">
                    <span class="terminal-prompt">$</span>
                    <span class="terminal-cursor"></span>
                </div>
            </div>
        </div>

        <div class="message">
            <div class="message-header">
                <svg class="message-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="9 10 4 15 9 20"/>
                    <path d="M20 4v7a4 4 0 0 1-4 4H4"/>
                </svg>
                <span class="message-title">Return to your terminal</span>
            </div>
            <p class="message-text">You can close this window and start using the CLI. Try <code>eightsleep --help</code> to see all available commands.</p>
        </div>

        <p class="footer">This window will close automatically.</p>
    </div>

    <script>
        fetch('/complete', { method: 'POST' }).catch(() => {});
    </script>
</body>
</html>`
