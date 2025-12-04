// Dashboard JavaScript for Real-Time Updates

let updateInterval;

// Initialize dashboard on load
document.addEventListener('DOMContentLoaded', function () {
    initDashboard();
    startAutoRefresh();
});

function initDashboard() {
    updateServices();
    updateMetrics();
    updateLogs();
    updateLastUpdated();
}

function startAutoRefresh() {
    // Refresh every 5 seconds
    updateInterval = setInterval(() => {
        updateServices();
        updateMetrics();
        updateLogs();
        updateLastUpdated();
    }, 5000);
}

function stopAutoRefresh() {
    if (updateInterval) {
        clearInterval(updateInterval);
    }
}

function refreshDashboard() {
    const btn = event.target.closest('button');
    btn.disabled = true;
    btn.innerHTML = '<i class="fas fa-sync-alt fa-spin"></i> Refreshing...';

    initDashboard();

    setTimeout(() => {
        btn.disabled = false;
        btn.innerHTML = '<i class="fas fa-sync-alt"></i> Refresh';
    }, 1000);
}

// Update service status
async function updateServices() {
    try {
        const response = await fetch('/api/services');
        const services = await response.json();

        const grid = document.getElementById('services-grid');
        grid.innerHTML = '';

        services.forEach(service => {
            const card = document.createElement('div');
            card.className = `service-card ${service.healthy ? 'healthy' : 'unhealthy'}`;

            const icon = service.healthy
                ? '<i class="fas fa-check-circle"></i>'
                : '<i class="fas fa-times-circle"></i>';

            card.innerHTML = `
                <h3>${service.name}</h3>
                <div class="service-url">${service.url}</div>
                <div class="service-status">
                    ${icon} ${service.status}
                </div>
            `;

            grid.appendChild(card);
        });
    } catch (error) {
        console.error('Failed to update services:', error);
        showServiceError();
    }
}

function showServiceError() {
    const grid = document.getElementById('services-grid');
    grid.innerHTML = `
        <div class="service-card unhealthy">
            <h3>Error</h3>
            <div class="service-url">Failed to fetch service status</div>
            <div class="service-status">
                <i class="fas fa-exclamation-triangle"></i> Check services
            </div>
        </div>
    `;
}

// Update metrics
async function updateMetrics() {
    try {
        const response = await fetch('/api/metrics');
        const metrics = await response.json();

        document.getElementById('request-rate').textContent =
            metrics.request_rate ? metrics.request_rate.toFixed(1) : '--';

        document.getElementById('error-rate').textContent =
            metrics.error_rate ? (metrics.error_rate * 100).toFixed(2) : '--';

        document.getElementById('avg-latency').textContent =
            metrics.avg_latency_ms ? metrics.avg_latency_ms.toFixed(0) : '--';

        document.getElementById('p95-latency').textContent =
            metrics.p95_latency_ms ? metrics.p95_latency_ms.toFixed(0) : '--';
    } catch (error) {
        console.error('Failed to update metrics:', error);
    }
}

// Update logs
async function updateLogs() {
    try {
        const response = await fetch('/api/logs');
        const logs = await response.json();

        const container = document.getElementById('logs-container');
        container.innerHTML = '';

        logs.forEach(log => {
            const entry = document.createElement('div');
            entry.className = `log-entry ${log.level}`;

            const timestamp = new Date(log.timestamp).toLocaleTimeString();

            entry.innerHTML = `
                <span class="log-timestamp">${timestamp}</span>
                <span class="log-level ${log.level}">${log.level}</span>
                <span class="log-service">${log.service}</span>
                <span class="log-message">${log.message}</span>
                ${log.trace_id ? `<br><span class="log-timestamp">Trace ID: ${log.trace_id}</span>` : ''}
            `;

            container.appendChild(entry);
        });
    } catch (error) {
        console.error('Failed to update logs:', error);
    }
}

// Load generator controls
async function startLoadGen() {
    try {
        const response = await fetch('/api/loadgen/start', {
            method: 'POST'
        });
        const result = await response.json();

        if (result.status === 'started') {
            document.getElementById('loadgen-status').innerHTML =
                '<span style="color: #10b981; font-weight: 600;">Running</span>';
            showNotification('Load generator started successfully', 'success');
        }
    } catch (error) {
        console.error('Failed to start load generator:', error);
        showNotification('Failed to start load generator', 'error');
    }
}

async function stopLoadGen() {
    try {
        const response = await fetch('/api/loadgen/stop', {
            method: 'POST'
        });
        const result = await response.json();

        if (result.status === 'stopped') {
            document.getElementById('loadgen-status').innerHTML =
                '<span style="color: #64748b;">Stopped</span>';
            showNotification('Load generator stopped', 'info');
        }
    } catch (error) {
        console.error('Failed to stop load generator:', error);
        showNotification('Failed to stop load generator', 'error');
    }
}

// Update last updated timestamp
function updateLastUpdated() {
    const now = new Date();
    const timeString = now.toLocaleTimeString();
    document.getElementById('last-updated').textContent = timeString;
}

// Show notification
function showNotification(message, type) {
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 25px;
        border-radius: 8px;
        background: ${type === 'success' ? '#10b981' : type === 'error' ? '#ef4444' : '#4f46e5'};
        color: white;
        font-weight: 600;
        box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
        z-index: 1000;
        animation: slideIn 0.3s ease;
    `;

    document.body.appendChild(notification);

    // Remove after 3 seconds
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Add CSS animations
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(400px);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(400px);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    stopAutoRefresh();
});

