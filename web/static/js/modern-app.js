// Modern WatchingCat Observability Platform
// Main Application JavaScript

// Global State
const state = {
    currentPage: 'dashboard',
    refreshInterval: null,
    charts: {},
    cart: [],
    loadGenStatus: 'stopped'
};

// Initialize App
document.addEventListener('DOMContentLoaded', () => {
    initializeNavigation();
    initializeDashboard();
    startAutoRefresh();

    console.log('🐱 WatchingCat Observability Platform initialized');
});

// ============================================
// Navigation
// ============================================

function initializeNavigation() {
    document.querySelectorAll('.nav-link').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const page = link.getAttribute('data-page');
            navigateTo(page);
        });
    });
}

function navigateTo(page) {
    // Update nav links
    document.querySelectorAll('.nav-link').forEach(link => {
        link.classList.remove('active');
        if (link.getAttribute('data-page') === page) {
            link.classList.add('active');
        }
    });

    // Update pages
    document.querySelectorAll('.page').forEach(p => {
        p.classList.remove('active');
    });
    document.getElementById(`page-${page}`).classList.add('active');

    state.currentPage = page;

    // Load page-specific data
    switch (page) {
        case 'services':
            loadServicesPage();
            break;
        case 'traces':
            loadTracesPage();
            break;
        case 'metrics':
            loadMetricsPage();
            break;
        case 'shop':
            loadShopPage();
            break;
    }
}

// ============================================
// Dashboard
// ============================================

function initializeDashboard() {
    updateSystemStatus();
    updateKeyMetrics();
    initializeCharts();
    loadServiceTopology();
    updateToolsStats();
}

async function updateSystemStatus() {
    try {
        const response = await fetch('/api/services');
        const services = await response.json();

        const allHealthy = services.every(s => s.healthy);
        const someHealthy = services.some(s => s.healthy);

        const indicator = document.getElementById('system-status');
        const statusText = document.getElementById('system-status-text');

        if (allHealthy) {
            indicator.className = 'fas fa-circle status-indicator healthy';
            statusText.textContent = 'All Systems Operational';
        } else if (someHealthy) {
            indicator.className = 'fas fa-circle status-indicator degraded';
            statusText.textContent = 'Degraded Performance';
        } else {
            indicator.className = 'fas fa-circle status-indicator down';
            statusText.textContent = 'System Down';
        }

        // Update uptime
        document.getElementById('system-uptime').textContent = '24h 15m';

    } catch (error) {
        console.error('Failed to update system status:', error);
    }
}

async function updateKeyMetrics() {
    try {
        const response = await fetch('/api/metrics');
        const metrics = await response.json();

        document.getElementById('metric-request-rate').textContent =
            (metrics.request_rate || 0).toFixed(1);

        const successRate = 100 - (metrics.error_rate || 0) * 100;
        document.getElementById('metric-success-rate').textContent =
            successRate.toFixed(1) + '%';

        document.getElementById('metric-latency').textContent =
            (metrics.avg_latency_ms || 0).toFixed(0);

        document.getElementById('metric-error-rate').textContent =
            ((metrics.error_rate || 0) * 100).toFixed(2) + '%';

    } catch (error) {
        console.error('Failed to update key metrics:', error);
    }
}

function initializeCharts() {
    // Request Volume Chart
    const reqCtx = document.getElementById('requestVolumeChart');
    if (reqCtx) {
        state.charts.requestVolume = new Chart(reqCtx, {
            type: 'line',
            data: {
                labels: generateTimeLabels(15),
                datasets: [{
                    label: 'Requests',
                    data: generateRandomData(15, 100, 200),
                    borderColor: 'rgb(99, 102, 241)',
                    backgroundColor: 'rgba(99, 102, 241, 0.1)',
                    tension: 0.4,
                    fill: true
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                aspectRatio: 2.5,
                plugins: {
                    legend: { display: false }
                },
                scales: {
                    y: { beginAtZero: true }
                }
            }
        });
    }

    // Latency Chart
    const latCtx = document.getElementById('latencyChart');
    if (latCtx) {
        state.charts.latency = new Chart(latCtx, {
            type: 'line',
            data: {
                labels: generateTimeLabels(15),
                datasets: [
                    {
                        label: 'P50',
                        data: generateRandomData(15, 100, 150),
                        borderColor: 'rgb(16, 185, 129)',
                        tension: 0.4
                    },
                    {
                        label: 'P95',
                        data: generateRandomData(15, 200, 300),
                        borderColor: 'rgb(245, 158, 11)',
                        tension: 0.4
                    },
                    {
                        label: 'P99',
                        data: generateRandomData(15, 350, 500),
                        borderColor: 'rgb(239, 68, 68)',
                        tension: 0.4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                aspectRatio: 2.5,
                plugins: {
                    legend: { display: false }
                },
                scales: {
                    y: { beginAtZero: true }
                }
            }
        });
    }
}

function loadServiceTopology() {
    const container = document.getElementById('service-topology');
    if (!container) return;

    // Clear container
    container.innerHTML = '';

    // Create SVG
    const width = container.clientWidth;
    const height = 500;

    const svg = d3.select('#service-topology')
        .append('svg')
        .attr('width', width)
        .attr('height', height);

    // Define nodes
    const nodes = [
        { id: 'load-gen', label: 'Load\nGenerator', x: 100, y: 250, type: 'client' },
        { id: 'frontend', label: 'Frontend', x: 300, y: 250, type: 'service' },
        { id: 'cart', label: 'Cart\nService', x: 500, y: 150, type: 'service' },
        { id: 'catalog', label: 'Product\nCatalog', x: 500, y: 250, type: 'service' },
        { id: 'checkout', label: 'Checkout\nService', x: 500, y: 350, type: 'service' },
        { id: 'collector', label: 'OTEL\nCollector', x: 700, y: 250, type: 'backend' },
        { id: 'jaeger', label: 'Jaeger', x: 900, y: 150, type: 'backend' },
        { id: 'prometheus', label: 'Prometheus', x: 900, y: 250, type: 'backend' },
        { id: 'elk', label: 'ELK Stack', x: 900, y: 350, type: 'backend' }
    ];

    // Define links
    const links = [
        { source: 'load-gen', target: 'frontend' },
        { source: 'frontend', target: 'cart' },
        { source: 'frontend', target: 'catalog' },
        { source: 'frontend', target: 'checkout' },
        { source: 'cart', target: 'collector' },
        { source: 'catalog', target: 'collector' },
        { source: 'checkout', target: 'collector' },
        { source: 'frontend', target: 'collector' },
        { source: 'collector', target: 'jaeger' },
        { source: 'collector', target: 'prometheus' },
        { source: 'collector', target: 'elk' }
    ];

    // Draw links
    links.forEach(link => {
        const source = nodes.find(n => n.id === link.source);
        const target = nodes.find(n => n.id === link.target);

        svg.append('line')
            .attr('x1', source.x)
            .attr('y1', source.y)
            .attr('x2', target.x)
            .attr('y2', target.y)
            .attr('stroke', '#cbd5e1')
            .attr('stroke-width', 2)
            .attr('opacity', 0.6);

        // Add animated dot
        svg.append('circle')
            .attr('r', 4)
            .attr('fill', '#6366f1')
            .append('animateMotion')
            .attr('dur', '3s')
            .attr('repeatCount', 'indefinite')
            .append('mpath')
            .attr('href', '#path-' + link.source + '-' + link.target);
    });

    // Draw nodes
    nodes.forEach(node => {
        const g = svg.append('g')
            .attr('transform', `translate(${node.x}, ${node.y})`)
            .attr('cursor', 'pointer')
            .on('mouseover', function () {
                d3.select(this).select('rect').attr('stroke-width', 3);
            })
            .on('mouseout', function () {
                d3.select(this).select('rect').attr('stroke-width', 2);
            });

        // Node color based on type
        let fillColor = '#f1f5f9';
        let strokeColor = '#cbd5e1';
        if (node.type === 'service') {
            fillColor = '#e0e7ff';
            strokeColor = '#6366f1';
        } else if (node.type === 'backend') {
            fillColor = '#fef3c7';
            strokeColor = '#f59e0b';
        }

        // Node rectangle
        g.append('rect')
            .attr('x', -50)
            .attr('y', -25)
            .attr('width', 100)
            .attr('height', 50)
            .attr('rx', 8)
            .attr('fill', fillColor)
            .attr('stroke', strokeColor)
            .attr('stroke-width', 2);

        // Node label
        g.append('text')
            .attr('text-anchor', 'middle')
            .attr('dominant-baseline', 'middle')
            .attr('font-size', '12px')
            .attr('font-weight', '600')
            .attr('fill', '#1e293b')
            .selectAll('tspan')
            .data(node.label.split('\n'))
            .enter()
            .append('tspan')
            .attr('x', 0)
            .attr('dy', (d, i) => i * 14 - 7)
            .text(d => d);

        // Status indicator
        g.append('circle')
            .attr('cx', 40)
            .attr('cy', -20)
            .attr('r', 5)
            .attr('fill', '#10b981')
            .attr('stroke', 'white')
            .attr('stroke-width', 2);
    });
}

async function updateToolsStats() {
    // Update tool statistics
    document.getElementById('jaeger-traces').textContent = '1,234';
    document.getElementById('prom-metrics').textContent = '542';
    document.getElementById('kibana-logs').textContent = '12.5K';
}

// ============================================
// Services Page
// ============================================

async function loadServicesPage() {
    try {
        const response = await fetch('/api/services');
        const services = await response.json();

        const grid = document.getElementById('services-grid');
        grid.innerHTML = '';

        services.forEach(service => {
            const card = createServiceCard(service);
            grid.appendChild(card);
        });
    } catch (error) {
        console.error('Failed to load services:', error);
    }
}

function createServiceCard(service) {
    const card = document.createElement('div');
    card.className = 'service-card-detailed';

    const badge = service.healthy ?
        '<span class="service-badge healthy">Healthy</span>' :
        '<span class="service-badge unhealthy">Unhealthy</span>';

    card.innerHTML = `
        <div class="service-card-header">
            <div>
                <div class="service-name">${service.name}</div>
                <div style="color: var(--text-secondary); font-size: 0.875rem;">${service.url}</div>
            </div>
            ${badge}
        </div>
        <div class="service-stats">
            <div class="service-stat">
                <div class="service-stat-value">${Math.floor(Math.random() * 200 + 100)}</div>
                <div class="service-stat-label">req/sec</div>
            </div>
            <div class="service-stat">
                <div class="service-stat-value">${Math.floor(Math.random() * 100 + 50)}ms</div>
                <div class="service-stat-label">latency</div>
            </div>
            <div class="service-stat">
                <div class="service-stat-value">${(Math.random() * 5).toFixed(2)}%</div>
                <div class="service-stat-label">error rate</div>
            </div>
        </div>
    `;

    return card;
}

// ============================================
// Traces Page
// ============================================

async function loadTracesPage() {
    const tracesList = document.getElementById('traces-list');
    tracesList.innerHTML = '';

    // Generate mock traces
    for (let i = 0; i < 10; i++) {
        const trace = {
            id: generateTraceId(),
            duration: Math.floor(Math.random() * 1000 + 100),
            services: ['frontend', 'cartservice', 'checkoutservice'],
            timestamp: new Date(Date.now() - Math.random() * 3600000)
        };

        const traceItem = createTraceItem(trace);
        tracesList.appendChild(traceItem);
    }
}

function createTraceItem(trace) {
    const item = document.createElement('div');
    item.className = 'trace-item';

    item.innerHTML = `
        <div class="trace-header">
            <span class="trace-id">${trace.id}</span>
            <span class="trace-duration">${trace.duration}ms</span>
        </div>
        <div style="color: var(--text-secondary); font-size: 0.875rem; margin-bottom: 0.75rem;">
            ${trace.timestamp.toLocaleString()}
        </div>
        <div class="trace-services">
            ${trace.services.map(s => `<span class="trace-service-badge">${s}</span>`).join('')}
        </div>
    `;

    item.addEventListener('click', () => {
        viewTraceDetails(trace);
    });

    return item;
}

function searchTraces() {
    const query = document.getElementById('trace-search').value;
    const service = document.getElementById('trace-service-filter').value;

    showToast(`Searching traces${service ? ` for ${service}` : ''}...`, 'info');
    loadTracesPage();
}

function viewTraceDetails(trace) {
    // Generate mock spans for the trace following OpenTelemetry structure
    const spans = generateMockSpans(trace);

    // Create trace viewer modal
    const modal = document.createElement('div');
    modal.className = 'trace-modal';
    modal.id = 'trace-modal';

    // Create backdrop
    const backdrop = document.createElement('div');
    backdrop.className = 'trace-modal-backdrop';
    backdrop.onclick = closeTraceModal;

    const content = document.createElement('div');
    content.className = 'trace-modal-content';
    content.onclick = (e) => e.stopPropagation(); // Prevent closing when clicking inside

    content.innerHTML = `
        <div class="trace-modal-header">
            <div>
                <h2><i class="fas fa-project-diagram"></i> Trace Details</h2>
                <div class="trace-id-display">Trace ID: ${trace.id}</div>
            </div>
            <button class="btn-close" onclick="closeTraceModal()">
                <i class="fas fa-times"></i>
            </button>
        </div>
        
        <div class="trace-summary">
            <div class="trace-summary-item">
                <span class="label">Duration:</span>
                <span class="value">${trace.duration}ms</span>
            </div>
            <div class="trace-summary-item">
                <span class="label">Services:</span>
                <span class="value">${trace.services.length}</span>
            </div>
            <div class="trace-summary-item">
                <span class="label">Spans:</span>
                <span class="value">${spans.length}</span>
            </div>
            <div class="trace-summary-item">
                <span class="label">Started:</span>
                <span class="value">${trace.timestamp.toLocaleTimeString()}</span>
            </div>
        </div>
        
        <div class="trace-timeline">
            <h3><i class="fas fa-stream"></i> Span Timeline</h3>
            <div class="spans-container">
                ${spans.map(span => createSpanView(span, trace.duration)).join('')}
            </div>
        </div>
        
        <div class="trace-actions">
            <button class="btn btn-outline" onclick="window.open('http://localhost:16686/trace/${trace.id}', '_blank')">
                <i class="fas fa-external-link-alt"></i> View in Jaeger
            </button>
            <button class="btn btn-primary" onclick="closeTraceModal()">
                Close
            </button>
        </div>
    `;

    modal.appendChild(backdrop);
    modal.appendChild(content);
    document.body.appendChild(modal);

    // Prevent body scroll
    document.body.style.overflow = 'hidden';

    // Animate in
    requestAnimationFrame(() => {
        modal.classList.add('active');
    });

    // ESC key handler
    const escHandler = (e) => {
        if (e.key === 'Escape') {
            closeTraceModal();
        }
    };
    document.addEventListener('keydown', escHandler);
    modal.escHandler = escHandler; // Store for cleanup
}

function generateMockSpans(trace) {
    const spans = [];
    let currentTime = 0;

    // Root span - Frontend request
    spans.push({
        spanId: generateSpanId(),
        operationName: 'HTTP GET /',
        service: 'frontend',
        startTime: currentTime,
        duration: trace.duration,
        tags: {
            'http.method': 'GET',
            'http.url': '/',
            'http.status_code': 200
        },
        level: 0
    });

    // Child spans
    const services = ['cartservice', 'productcatalog', 'checkoutservice'];
    services.forEach((service, index) => {
        const start = currentTime + (index * 30);
        const duration = Math.floor(trace.duration * 0.3);

        spans.push({
            spanId: generateSpanId(),
            operationName: `${service}.GetItems`,
            service: service,
            startTime: start,
            duration: duration,
            tags: {
                'rpc.service': service,
                'rpc.method': 'GetItems'
            },
            level: 1
        });
    });

    return spans;
}

function generateSpanId() {
    return Array.from({ length: 16 }, () =>
        Math.floor(Math.random() * 16).toString(16)
    ).join('');
}

function createSpanView(span, totalDuration) {
    const startPercent = (span.startTime / totalDuration) * 100;
    const widthPercent = Math.max((span.duration / totalDuration) * 100, 2); // Min 2% width for visibility
    const indent = span.level * 30;

    // Color based on service
    const colors = {
        'frontend': '#6366f1',
        'cartservice': '#10b981',
        'productcatalog': '#f59e0b',
        'checkoutservice': '#ef4444'
    };
    const color = colors[span.service] || '#64748b';

    // Show duration in bar if wide enough
    const showDurationInBar = widthPercent > 15;

    return `
        <div class="span-row" style="margin-left: ${indent}px;">
            <div class="span-info">
                <div class="span-service" style="color: ${color};">
                    <i class="fas fa-cube"></i> ${span.service}
                </div>
                <div class="span-operation">${span.operationName}</div>
                <div class="span-duration">${span.duration}ms</div>
            </div>
            <div class="span-bar-container" title="${span.service}: ${span.operationName} (${span.duration}ms)">
                <div class="span-bar" style="
                    left: ${startPercent}%;
                    width: ${widthPercent}%;
                    background: ${color};
                ">${showDurationInBar ? span.duration + 'ms' : ''}</div>
            </div>
            <div class="span-tags">
                ${Object.entries(span.tags).map(([key, value]) =>
        `<span class="span-tag"><strong>${key}:</strong> ${value}</span>`
    ).join('')}
            </div>
        </div>
    `;
}

function closeTraceModal() {
    const modal = document.getElementById('trace-modal');
    if (modal) {
        // Remove ESC key handler
        if (modal.escHandler) {
            document.removeEventListener('keydown', modal.escHandler);
        }

        // Restore body scroll
        document.body.style.overflow = '';

        // Animate out
        modal.classList.remove('active');
        setTimeout(() => {
            if (modal.parentNode) {
                modal.remove();
            }
        }, 300);
    }
}

// Make closeTraceModal globally accessible
window.closeTraceModal = closeTraceModal;

// ============================================
// Metrics Page
// ============================================

function loadMetricsPage() {
    // Initialize metric charts
    setTimeout(() => {
        initCPUChart();
        initMemoryChart();
        initNetworkChart();
    }, 100);
}

function initCPUChart() {
    const ctx = document.getElementById('cpuChart');
    if (!ctx || state.charts.cpu) return;

    state.charts.cpu = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['Frontend', 'Cart', 'Catalog', 'Checkout'],
            datasets: [{
                label: 'CPU Usage (%)',
                data: [45, 32, 28, 51],
                backgroundColor: [
                    'rgba(99, 102, 241, 0.8)',
                    'rgba(16, 185, 129, 0.8)',
                    'rgba(245, 158, 11, 0.8)',
                    'rgba(239, 68, 68, 0.8)'
                ]
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            aspectRatio: 2,
            scales: {
                y: { beginAtZero: true, max: 100 }
            }
        }
    });
}

function initMemoryChart() {
    const ctx = document.getElementById('memoryChart');
    if (!ctx || state.charts.memory) return;

    state.charts.memory = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['Frontend', 'Cart', 'Catalog', 'Checkout'],
            datasets: [{
                label: 'Memory (MB)',
                data: [256, 128, 192, 384],
                backgroundColor: [
                    'rgba(99, 102, 241, 0.8)',
                    'rgba(16, 185, 129, 0.8)',
                    'rgba(245, 158, 11, 0.8)',
                    'rgba(239, 68, 68, 0.8)'
                ]
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            aspectRatio: 2,
            scales: {
                y: { beginAtZero: true }
            }
        }
    });
}

function initNetworkChart() {
    const ctx = document.getElementById('networkChart');
    if (!ctx || state.charts.network) return;

    state.charts.network = new Chart(ctx, {
        type: 'line',
        data: {
            labels: generateTimeLabels(20),
            datasets: [{
                label: 'Throughput (MB/s)',
                data: generateRandomData(20, 10, 50),
                borderColor: 'rgb(99, 102, 241)',
                backgroundColor: 'rgba(99, 102, 241, 0.1)',
                tension: 0.4,
                fill: true
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            aspectRatio: 2.5,
            scales: {
                y: { beginAtZero: true }
            }
        }
    });
}

// ============================================
// Demo Shop
// ============================================

const products = [
    { id: 1, name: 'Observatory Telescope Pro', price: 2999.99, image: '🔭' },
    { id: 2, name: 'Star Chart Collection', price: 149.99, image: '🗺️' },
    { id: 3, name: 'Night Vision Binoculars', price: 599.99, image: '🔍' },
    { id: 4, name: 'Astronomy Guide Book', price: 49.99, image: '📚' },
    { id: 5, name: 'Portable Planetarium', price: 899.99, image: '🌍' },
    { id: 6, name: 'Space Photography Kit', price: 1299.99, image: '📷' }
];

function loadShopPage() {
    const grid = document.getElementById('product-grid');
    grid.innerHTML = '';

    products.forEach(product => {
        const card = createProductCard(product);
        grid.appendChild(card);
    });

    updateCart();
}

function createProductCard(product) {
    const card = document.createElement('div');
    card.className = 'product-card';

    card.innerHTML = `
        <div class="product-image" style="display: flex; align-items: center; justify-content: center; font-size: 4rem;">
            ${product.image}
        </div>
        <div class="product-info">
            <div class="product-name">${product.name}</div>
            <div class="product-price">$${product.price.toFixed(2)}</div>
            <button class="btn btn-primary btn-block" onclick="addToCart(${product.id})">
                <i class="fas fa-cart-plus"></i> Add to Cart
            </button>
        </div>
    `;

    return card;
}

function addToCart(productId) {
    const product = products.find(p => p.id === productId);
    if (product) {
        state.cart.push(product);
        updateCart();
        showToast(`Added ${product.name} to cart`, 'success');
    }
}

function updateCart() {
    const container = document.getElementById('cart-items');
    const totalEl = document.getElementById('cart-total');

    if (state.cart.length === 0) {
        container.innerHTML = '<div style="text-align: center; color: var(--text-secondary); padding: 2rem;">Cart is empty</div>';
        totalEl.textContent = '$0.00';
        return;
    }

    container.innerHTML = '';
    let total = 0;

    state.cart.forEach((item, index) => {
        total += item.price;

        const cartItem = document.createElement('div');
        cartItem.className = 'cart-item';
        cartItem.innerHTML = `
            <div>
                <div style="font-weight: 600; font-size: 0.875rem;">${item.name}</div>
                <div style="color: var(--text-secondary); font-size: 0.75rem;">$${item.price.toFixed(2)}</div>
            </div>
            <button class="btn btn-sm btn-danger" onclick="removeFromCart(${index})" style="padding: 0.25rem 0.5rem;">
                <i class="fas fa-times"></i>
            </button>
        `;

        container.appendChild(cartItem);
    });

    totalEl.textContent = `$${total.toFixed(2)}`;
}

function removeFromCart(index) {
    state.cart.splice(index, 1);
    updateCart();
}

async function checkout() {
    if (state.cart.length === 0) {
        showToast('Cart is empty', 'warning');
        return;
    }

    showToast('Processing checkout...', 'info');

    // Simulate checkout
    setTimeout(() => {
        showToast('Order placed successfully! Check traces in Jaeger', 'success');
        state.cart = [];
        updateCart();
    }, 1500);
}

// ============================================
// Load Generator
// ============================================

async function startLoadGen() {
    try {
        const response = await fetch('/api/loadgen/start', { method: 'POST' });
        const result = await response.json();

        if (result.status === 'started') {
            state.loadGenStatus = 'running';
            document.getElementById('loadgen-status').innerHTML =
                '<span class="status-badge running">Running</span>';
            showToast('Load generator started', 'success');
        }
    } catch (error) {
        showToast('Failed to start load generator', 'error');
    }
}

async function stopLoadGen() {
    try {
        const response = await fetch('/api/loadgen/stop', { method: 'POST' });
        const result = await response.json();

        if (result.status === 'stopped') {
            state.loadGenStatus = 'stopped';
            document.getElementById('loadgen-status').innerHTML =
                '<span class="status-badge stopped">Stopped</span>';
            showToast('Load generator stopped', 'info');
        }
    } catch (error) {
        showToast('Failed to stop load generator', 'error');
    }
}

// ============================================
// Auto Refresh
// ============================================

function startAutoRefresh() {
    state.refreshInterval = setInterval(() => {
        if (state.currentPage === 'dashboard') {
            updateSystemStatus();
            updateKeyMetrics();
            updateLastUpdate();
        }
    }, 5000);
}

function refreshAll() {
    showToast('Refreshing dashboard...', 'info');

    switch (state.currentPage) {
        case 'dashboard':
            initializeDashboard();
            break;
        case 'services':
            loadServicesPage();
            break;
        case 'traces':
            loadTracesPage();
            break;
        case 'metrics':
            loadMetricsPage();
            break;
        case 'shop':
            loadShopPage();
            break;
    }
}

function updateLastUpdate() {
    document.getElementById('last-update').textContent = new Date().toLocaleTimeString();
}

// ============================================
// Utilities
// ============================================

function generateTimeLabels(count) {
    const labels = [];
    const now = new Date();
    for (let i = count - 1; i >= 0; i--) {
        const time = new Date(now - i * 60000);
        labels.push(time.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }));
    }
    return labels;
}

function generateRandomData(count, min, max) {
    return Array.from({ length: count }, () =>
        Math.floor(Math.random() * (max - min + 1)) + min
    );
}

function generateTraceId() {
    return Array.from({ length: 32 }, () =>
        Math.floor(Math.random() * 16).toString(16)
    ).join('');
}

function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;

    const icon = type === 'success' ? 'check-circle' :
        type === 'error' ? 'exclamation-circle' :
            type === 'warning' ? 'exclamation-triangle' : 'info-circle';

    toast.innerHTML = `
        <i class="fas fa-${icon}"></i>
        <span>${message}</span>
    `;

    container.appendChild(toast);

    setTimeout(() => {
        toast.style.animation = 'slideIn 0.3s ease reverse';
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

function toggleTheme() {
    const html = document.documentElement;
    const current = html.getAttribute('data-theme') || 'light';
    const next = current === 'light' ? 'dark' : 'light';
    html.setAttribute('data-theme', next);
    showToast(`Switched to ${next} theme`, 'info');
}

function resetTopology() {
    loadServiceTopology();
    showToast('Topology view reset', 'info');
}

function updateCharts(timeRange) {
    showToast(`Updated charts to ${timeRange}`, 'info');
    // Regenerate chart data based on time range
    if (state.charts.requestVolume) {
        const count = timeRange === '5m' ? 5 : timeRange === '15m' ? 15 : timeRange === '1h' ? 60 : 24;
        state.charts.requestVolume.data.labels = generateTimeLabels(count);
        state.charts.requestVolume.data.datasets[0].data = generateRandomData(count, 100, 200);
        state.charts.requestVolume.update();
    }
}

// Update last update timestamp immediately
updateLastUpdate();

console.log('✅ WatchingCat loaded successfully');

